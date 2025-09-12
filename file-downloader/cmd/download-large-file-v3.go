package main

import (
	"context"
	"dobrainya/file-downloader/internal/wp"
	"encoding/binary"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type IPool interface {
	Create()
	Handle(r *Range)
	Wait()
	Stats()
}

type Range struct {
	Url        string
	Path2save  string
	Offset     int
	Start      int64
	End        int64
	Errors     int
	Downloaded bool
}

var processedRanges int = 0
var locker *sync.Mutex = &sync.Mutex{}
var totalTimeout int = 3600
var chunkSize int64 = 400000000 //50Mb

func main() {
	// downloadFile(
	// 	"https://fias-file.nalog.ru/downloads/2025.09.05/gar_delta_xml.zip",
	// 	"gar_delta.zip",
	// )

	downloadFile(
		"https://fias-file.nalog.ru/downloads/2025.09.05/gar_xml.zip",
		"gar_full.zip",
		chunkSize,
	)
}

func downloadRangeWorkerFn(workerId int, r *Range) {
	for {
		if r.Errors == 3 {
			fmt.Println("Range %d-%d: download attempt exeeded", r.Start, r.End)
			return
		}

		if r.Errors > 0 {
			time.Sleep(time.Duration(5) * time.Duration(r.Errors) * time.Second)
		}

		err := downloadRange(r)

		if err != nil {
			fmt.Println(err.Error())
			r.Errors += 1
			continue
		}

		processedRanges++
		return
	}
}

func downloadFile(url string, filePath string, partSize int64) {
	resp, err := http.Head(url)

	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Cannot download file. Got %d HTTP status code from server\n", resp.StatusCode)
		os.Exit(1)
	}

	cl := resp.Header.Get("Content-Length")

	if cl == "" {
		panic(fmt.Errorf("Content is empty"))
	}

	numProc := runtime.NumCPU()
	fileSize, err := strconv.Atoi(cl)

	if err != nil {
		panic(err)
	}

	fileSize64 := int64(fileSize)

	rangesCount := int(fileSize64 / partSize)

	if mod := fileSize64 % partSize; mod != 0 {
		rangesCount += 1
	}

	fmt.Printf("Jobs count: %d \nFile size: %d\nPart Size: %d\n", rangesCount, fileSize64, partSize)

	ranges := make([]*Range, 0, rangesCount)
	var startRange int64 = 0
	var endRange int64 = 0

	for i := range rangesCount {
		if i > 0 {
			startRange = endRange + 1
		}

		endRange = partSize * (int64(i) + int64(1))

		if endRange > fileSize64 {
			endRange = fileSize64
		}

		ranges = append(ranges, &Range{url, filePath, int(i), startRange, endRange, 0, false})
	}

	go func() {
		for {
			if processedRanges >= rangesCount {
				return
			}

			p := float64(processedRanges) / float64(rangesCount) * 100

			fmt.Printf("\rProgress: %.2f", p)

			time.Sleep(5 * time.Second)
		}
	}()

	//===================worker pool================================
	wp.InitWorkers(numProc)
	var pool IPool = wp.New(downloadRangeWorkerFn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(totalTimeout)*time.Second)
	defer cancel()
l:
	for {
		select {
		case <-ctx.Done():
			break l
		default:
		}

		pool.Create()

		for _, r := range ranges {
			pool.Handle(r)
		}

		pool.Wait()
	}

	//pool.Stats()
	fmt.Printf("File %s successfully download to\n", filePath)
}

func downloadRange(ranges *Range) error {
	req, err := http.NewRequest("GET", ranges.Url, nil)

	if err != nil {
		return err
	}

	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", ranges.Start, ranges.End))

	client := &http.Client{
		Timeout: 1800 * time.Second,
	}

	fileResp, err := client.Do(req)
	defer fileResp.Body.Close()

	if err != nil {
		return err
	}

	contentLength, err := strconv.Atoi(fileResp.Header.Get("Content-Length"))

	if err != nil {
		return err
	}

	buffer := make([]byte, contentLength)

	err = binary.Read(fileResp.Body, binary.LittleEndian, buffer)

	if err != nil {
		return err
	}

	//fmt.Printf("\r%d: Download by range %d-%d\n", ranges.Offset, ranges.Start, ranges.End)
	// fmt.Println("Content-Range:", fileResp.Header.Get("Content-Range"))
	// fmt.Println("Content-Length:", contentLength)

	file, err := os.OpenFile(ranges.Path2save, os.O_WRONLY|os.O_CREATE, 0644)
	defer file.Close()

	if err != nil {
		fmt.Println(err)
		return nil
	}

	file.Seek(ranges.Start, 0)
	a, err := file.Write(buffer)

	_ = a

	if err != nil {
		return err
	}

	return nil
}
