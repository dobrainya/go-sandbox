package main

import (
	"encoding/binary"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
)

type Range struct {
	Offset int
	Start  int
	End    int
}

func main() {
	// downloadFile(
	// 	"https://fias-file.nalog.ru/downloads/2025.09.05/gar_delta_xml.zip",
	// 	"gar_delta.zip",
	// )

	downloadFile(
		"https://fias-file.nalog.ru/downloads/2025.09.05/gar_xml.zip",
		"gar_full.zip",
	)

}

func downloadFile(url string, filePath string) {
	resp, err := http.Head(url)

	if err != nil {
		panic(err)
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
	workersCnt := numProc

	if mod := fileSize % numProc; mod != 0 {
		workersCnt += 1
	}

	bachSize := fileSize / numProc
	rangeCh := make(chan *Range, numProc)

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)

	defer file.Close()

	if err != nil {
		panic(err)
	}

	startRange := 0
	endRange := 0

	go func() {
		for i := range workersCnt {
			if i > 1 {
				startRange = endRange + 1
			}

			endRange = bachSize * (i + 1)

			if endRange > fileSize {
				endRange = fileSize
			}

			rangeCh <- &Range{i, startRange, endRange}
		}

		close(rangeCh)
	}()
	wg := &sync.WaitGroup{}
	for range workersCnt {
		ranges := <-rangeCh

		wg.Add(1)
		go downloadRange(url, wg, ranges, file)
	}

	wg.Wait()

	fmt.Printf("File %s successfully download to %s\n", url, filePath)
}

func downloadRange(url string, wg *sync.WaitGroup, ranges *Range, file *os.File) {
	defer wg.Done()
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		panic(err)
	}

	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", ranges.Start, ranges.End))

	client := &http.Client{}

	fileResp, err := client.Do(req)
	defer fileResp.Body.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	contentLength, err := strconv.Atoi(fileResp.Header.Get("Content-Length"))

	if err != nil {
		fmt.Println(err)
		return
	}

	buffer := make([]byte, contentLength)

	err = binary.Read(fileResp.Body, binary.LittleEndian, buffer)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%d: Download by range %d-%d\n", ranges.Offset, ranges.Start, ranges.End)
	// fmt.Println("Content-Range:", fileResp.Header.Get("Content-Range"))
	// fmt.Println("Content-Length:", contentLength)

	file.Seek(int64(ranges.Start), 0)
	a, err := file.Write(buffer)
	_ = a
	if err != nil {
		fmt.Println(err)
		return
	}
}
