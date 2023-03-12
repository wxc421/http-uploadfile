package main

import (
	"context"
	"fmt"
	"github.com/freshman-tech/file-upload-starter-files/iox"
	"io"
	"io/ioutil"
	"log"
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

func (pr *Progress) Close() error {
	return nil
}

// Write is used to satisfy the io.Writer interface.
// Instead of writing somewhere, it simply aggregates
// the total bytes on each read
func (pr *Progress) Write(p []byte) (n int, err error) {
	time.Sleep(time.Millisecond * 200)
	n, err = len(p), nil
	pr.BytesRead += int64(n)
	pr.Print()
	return
}

// Print displays the current progress of the file upload
// each time Write is called
func (pr *Progress) Print() {
	if pr.BytesRead == pr.TotalSize {
		fmt.Printf("Progress File upload in progress: %.2f%%\n", float32(pr.BytesRead)/float32(pr.TotalSize)*100)
		fmt.Println("DONE!")
		return
	}
	fmt.Printf("Progress File upload in progress: %.2f%%\n", float32(pr.BytesRead)/float32(pr.TotalSize)*100)
}

func main() {
	defer func() {
		time.Sleep(time.Second * 20)
	}()
	log.SetFlags(log.Ldate | log.Lshortfile)

	url := "http://127.0.0.1:4500/upload"
	method := "POST"
	r, w := io.Pipe()
	m := multipart.NewWriter(w)

	name := "C:\\Users\\王昕成\\Downloads\\XunLeiWebSetup11.4.1.2030dl.exe"
	fileMeta, _ := os.Stat(name)

	p := &Progress{
		TotalSize: fileMeta.Size(),
	}
	teeReader := iox.TeeReader(r, p)

	go func() {
		defer func() {
			log.Println("w is close")
			log.Println("m is close")
		}()
		defer w.Close()
		file, _ := os.Open(name)
		defer file.Close()
		part, err := m.CreateFormFile("file", filepath.Base(name))
		if err == nil {
			_, err = io.Copy(part, file)
			if err == nil {
				m.Close()
			}
			return
		}

	}()

	client := &http.Client{}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	req, err := http.NewRequestWithContext(ctx, method, url, teeReader)
	go func() {
		time.Sleep(time.Second * 10)
		cancel()
		log.Println("cancel success")
	}()
	//req, err := http.NewRequest(method, url, teeReader)

	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("Content-Type", m.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(body))

}
