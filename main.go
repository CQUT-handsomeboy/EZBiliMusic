package main

import (
	"fmt"
	"net/http"
	"os"

	"os/exec"
	"path/filepath"

	"github.com/CQUT-handsomeboy/EZBiliMusic/downloader"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// 下载歌曲请求
type SongDownloadRequest struct {
	Aid    int    `json:"aid"`
	BVid   string `json:"bvid"`
	Cid    int    `json:"cid"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
}

func downloadSongWorker(ch chan SongDownloadRequest) {
	for data := range ch {
		path := filepath.Join(downloader.OutputFilesRootPath,
			fmt.Sprintf("%s_%d.m4a", data.BVid, data.Cid))
		if _, err := os.Stat(path); err == nil {
			// exists!
			fmt.Println("file already exists, skip download")
			continue
		}
		cmdArg := fmt.Sprintf("%s:%s", data.Artist, data.Title)
		cmd := exec.Command("skate", "get", cmdArg)
		if err := cmd.Run(); err == nil {
			// exists!
			fmt.Println("music already exists, skip download")
			continue
		}

		fmt.Println("download start...")
		if err := downloader.DownloadAudio(data.Aid, data.Cid, data.BVid, data.Title, data.Artist); err != nil {
			fmt.Println("download failed...")
			continue
		}
		cmd = exec.Command("skate", "set", cmdArg)

		if err := cmd.Run(); err != nil {
			fmt.Println("skate command failed...")
		}

		fmt.Println("download done...")
	}
}

type SongMetadataRequest struct {
	Url string `json:"url"`
}

func main() {
	ch := make(chan SongDownloadRequest, 100) // 缓存100个下载请求

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

	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/download", func(ctx *gin.Context) {
		var req SongDownloadRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "request json format is valid"})
			return
		}

		select {
		case ch <- req:
			ctx.JSON(http.StatusOK, gin.H{"info": "successfully parse request"})
			fmt.Println("successfully parse request")
		default:
			ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "server is busy, please try later"})
			fmt.Println("server is busy, please try later")
		}

	})

	r.POST("/metadata", func(ctx *gin.Context) {
		var req SongMetadataRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "request json format is valid"})
			return
		}
		html, err := downloader.GetHTMLContent(req.Url)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "url is invalid"})
			return
		}
		meta, err := downloader.ParseHTMLMeta(html)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "parse HTML meta error"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"info": "successfully get metadata", "meta": meta})
	})

	r.Run(":8080")
}
