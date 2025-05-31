package main

import (
	"fmt"
	"net/http"
	"os"

	"io/fs"
	"path/filepath"

	"github.com/CQUT-handsomeboy/EZBiliMusic/database"
	"github.com/CQUT-handsomeboy/EZBiliMusic/downloader"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.senan.xyz/taglib"
)

// 下载歌曲请求
type SongDownloadRequest struct {
	Aid    int    `json:"aid"`
	BVid   string `json:"bvid"`
	Cid    int    `json:"cid"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
}

type SongDeleteRequest struct {
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

		ifExists, err := database.CheckMusicRecord(data.Title, data.Artist)

		if err != nil {
			fmt.Printf("database check error:%v\n", err)
			continue
		}

		if ifExists {
			fmt.Println("music already exists, skip download")
			continue
		}

		filename := fmt.Sprintf("%s_%d.m4a", data.BVid, data.Cid)

		fmt.Println("download start...")

		if err := downloader.DownloadAudio(data.Aid, data.Cid, data.BVid, data.Title, data.Artist); err != nil {
			fmt.Println("download failed...")
			continue
		} else {
			fmt.Println("download done...")
			// remove cache
			filepath.Walk(downloader.TempFilesRootPath, func(path string, info fs.FileInfo, err error) error {

				if filepath.Base(path) == ".gitkeep" {
					return nil
				}

				os.Remove(path)
				return nil
			})
		}

		err = database.InsertMuiscRecord(data.Title, data.Artist, filename)

		if err != nil {
			fmt.Printf("database insert error:%v\n", err)
			continue
		}

	}
}

func getMetadataWorker() {
	err := filepath.Walk(downloader.OutputFilesRootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || filepath.Ext(path) != ".m4a" {
			return nil
		}

		metadata, err := taglib.ReadTags(path)

		if err != nil {
			fmt.Printf("read tag error:%v\n", err)
			return err
		}

		title := metadata["TITLE"][0]
		artist := metadata["ARTIST"][0]

		exists, err := database.CheckMusicRecord(title, artist)

		if err != nil {
			fmt.Printf("database check error:%v\n", err)
			return err
		}

		if exists {
			fmt.Println("music already exists, skip insert")
			return nil
		}

		err = database.InsertMuiscRecord(title, artist, filepath.Base(path))

		if err != nil {
			fmt.Printf("database insert error:%v\n", err)
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Printf("file walk failed: %v\n", err)
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

	go getMetadataWorker()
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

	r.POST("/delete", func(ctx *gin.Context) {
		var req SongDeleteRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "request json format is valid"})
			return
		}

		exists, err := database.CheckMusicRecord(req.Title, req.Artist)

		if err != nil {
			fmt.Printf("database check error:%v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "database check error"})
			return
		}

		if !exists {
			fmt.Println("music not exists, skip delete")
			ctx.JSON(http.StatusBadRequest, gin.H{"info": "music not exists"})
			return
		}

		filename, err := database.DeleteMusicRecord(req.Title, req.Artist)

		if err != nil {
			fmt.Printf("database delete error:%v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "database delete error"})
			return
		}

		musicFilePath := filepath.Join(downloader.OutputFilesRootPath, filename)

		if _, err := os.Stat(musicFilePath); err != nil {
			fmt.Printf("file not exists:%v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "file not exists"})
			return
		}

		err = os.Remove(musicFilePath)

		if err != nil {
			fmt.Printf("file delete error:%v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "file delete error"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"info": "successfully delete music"})

	})

	r.Run(":8080")
}
