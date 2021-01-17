package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type WriteCounter struct {
	TotalCount uint
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.TotalCount += uint(n)
	return n, nil
}

func main() {
	var URL string
	fmt.Print("URL: ")
	if _, err := fmt.Scanln(&URL); err != nil {
		println("Error when try to read URL")
		os.Exit(1)
	}
	fileName := URL[strings.LastIndex(URL, "/")+1 : len(URL)]

	downloadFile(URL, fileName)
}

func downloadFile(URL string, fileName string){
	fmt.Println("Connecting to site...")
	res, err := http.Get(URL)
	if err != nil || res.StatusCode != http.StatusOK {
		println("Failed to get file from", URL)
		os.Exit(2)
	}
	defer res.Body.Close()

	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error when try to create file:", err.Error())
		os.Exit(3)
	}
	defer file.Close()

	writeCounter := WriteCounter{0}
	reader := io.TeeReader(res.Body, &writeCounter)
	end := false
	go func() {
		for !end {
			fmt.Println("Downloaded ", writeCounter.TotalCount / 1024, "Kb")
			time.Sleep(time.Second)
		}
	}()
	_, err = io.Copy(file, reader)
	if err != nil {
		fmt.Println("Error when write to the file:")
		os.Exit(4)
	}
	end = true

	fmt.Println("Total size:", writeCounter.TotalCount / 1024, "Kb")
	fmt.Println("File is downloaded as ", fileName)
}