package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Progress is used to track the progress of a file upload.
// It implements the io.Writer interface so it can be passed
// to an io.TeeReader()
type Progress struct {
	TotalSize int64
	BytesRead int64
}

// Write is used to satisfy the io.Writer interface.
// Instead of writing somewhere, it simply aggregates
// the total bytes on each read
func (pr *Progress) Write(p []byte) (n int, err error) {
	time.Sleep(time.Millisecond * 500)
	n, err = len(p), nil
	pr.BytesRead += int64(n)
	pr.Print()
	return
}

// Print displays the current progress of the file upload
// each time Write is called
func (pr *Progress) Print() {
	if pr.BytesRead == pr.TotalSize {
		fmt.Println("DONE!")
		return
	}

	fmt.Printf("Progress File upload in progress: %.2f%%\n", float32(pr.BytesRead)/float32(pr.TotalSize)*100)
}

func main() {

	url := "http://127.0.0.1:4500/upload"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	name := "C:\\Users\\王昕成\\Downloads\\Feishu-win32_ia32-5.31.6-signed.exe"
	file, _ := os.Open(name)
	defer file.Close()
	part1, _ := writer.CreateFormFile("file", filepath.Base(name))
	_, _ = io.Copy(part1, file)

	// progress file
	pr := &Progress{
		TotalSize: int64(payload.Len()),
	}

	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, io.TeeReader(payload, pr))

	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
