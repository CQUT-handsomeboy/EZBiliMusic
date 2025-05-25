package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/CQUT-handsomeboy/EZBiliMusic/downloader"
	"path/filepath"
)

// 下载歌曲请求
type SongDownloadRequest struct {
	Aid  int    `json:"aid"`
	BVid string `json:"bvid"`
	Cid  int    `json:"cid"`
}

func downloadSongWorker(ch chan SongDownloadRequest) {
	for data := range ch {
		path := filepath.Join(downloader.OutputFilesRootPath,
			fmt.Sprintf("%s_%d.m4a", data.BVid, data.Cid))
		if _, err := os.Stat(path); err == nil {
			fmt.Println("file already exists, skip download")
			continue
		}
		fmt.Println("download start...")
		downloader.DownloadAudio(data.Aid, data.Cid, data.BVid)
		fmt.Println("download done...")
	}
}

type SongMetadataRequest struct {
	Url string `json:"url"`
}

func main() {
	ch := make(chan SongDownloadRequest)

	exePath, _ := os.Executable()
	downloader.OutputFilesRootPath = filepath.Join(filepath.Dir(exePath), "output")
	downloader.TempFilesRootPath = filepath.Join(filepath.Dir(exePath), "temp")

	if err := os.MkdirAll(downloader.OutputFilesRootPath, 0755); err != nil {
		fmt.Printf("mkdir fail: %v\n", err)
		return
	}

	if err := os.MkdirAll(downloader.TempFilesRootPath, 0755); err != nil {
		fmt.Printf("mkdir fail: %v\n", err)
		return
	}

	fmt.Printf("your music is in %s\n", downloader.OutputFilesRootPath)

	go downloadSongWorker(ch)

	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			io.WriteString(w, "Method Not Allowed")
			return
		}

		var req SongDownloadRequest
		err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "Bad Request")
			return
		}

		fmt.Println("successfully parse request")

		select {
		case ch <- req:
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, "order received and processing...")
		default:
			w.WriteHeader(http.StatusTooManyRequests)
			io.WriteString(w, "server is busy, please try later")
			return
		}

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "order received and processing...")

	})

	http.HandleFunc("/metadata", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			io.WriteString(w, "Method Not Allowed")
			return
		}

		var req SongMetadataRequest
		err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "Bad Request")
			return
		}

		fmt.Println("successfully parse request")

		html, err := downloader.GetHTMLContent(req.Url)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "url is invalid")
			return
		}

		meta, err := downloader.ParseHTMLMeta(html)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, "parse HTML meta error")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(meta)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, "encode JSON error")
			return
		}

	})

	fmt.Println("Hello,server is running...")

	err := http.ListenAndServe(":8080", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		return
	}
}
