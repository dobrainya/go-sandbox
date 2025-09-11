package main

import (
	"encoding/binary"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
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
		1024,
	)

}

func downloadFile(url string, filePath string, partSize int) {
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

	jobsCount := fileSize / partSize

	if mod := fileSize % partSize; mod != 0 {
		jobsCount += 1
	}

	fmt.Println("Jobs count: ", jobsCount, "File size: ", fileSize)

	jobsCh := make(chan *Range, jobsCount)

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)

	defer file.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	startRange := 0
	endRange := 0
	processedJobs := 0

	go func() {
		for i := range jobsCount {
			if i > 1 {
				startRange = endRange + 1
			}

			endRange = partSize * (i + 1)

			if endRange > fileSize {
				endRange = fileSize
			}

			jobsCh <- &Range{i, startRange, endRange}
		}

		close(jobsCh)
	}()

	go func() {
		for {
			if processedJobs >= jobsCount {
				return
			}

			p := float64(processedJobs) / float64(jobsCount) * 100

			fmt.Printf("\rProgress: %f", p)

			time.Sleep(1 * time.Second)
		}
	}()

	wg := &sync.WaitGroup{}
	locker := &sync.Mutex{}
	for range numProc {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case ranges, ok := <-jobsCh:
					if !ok {
						return
					}
					err := downloadRange(url, ranges, file)

					if err != nil {
						jobsCh <- ranges
						continue
					}

					locker.Lock()
					processedJobs++
					locker.Unlock()
				}
			}
		}()
	}

	wg.Wait()

	fmt.Printf("File %s successfully download to\n", filePath)
}

func downloadRange(url string, ranges *Range, file *os.File) error {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return err
	}

	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", ranges.Start, ranges.End))

	client := &http.Client{
		Timeout: 15 * time.Second,
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

	file.Seek(int64(ranges.Start), 0)
	a, err := file.Write(buffer)

	_ = a

	if err != nil {
		return err
	}

	return nil
}
