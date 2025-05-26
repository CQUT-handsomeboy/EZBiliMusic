package downloader

import (
	"encoding/json"
	"fmt"
	"go.senan.xyz/taglib"
	"os"
	"testing"
)

var urls []string = []string{
	"https://www.bilibili.com/video/BV1q4FTeBEMj", // 周杰伦歌曲 分P
	"https://www.bilibili.com/video/BV1t14y1W7jE", // 薛之谦《无数》 单曲
	"https://www.bilibili.com/video/BV1YmFPe4EnY", // 陈奕迅歌曲 分P
}

func TestDownloadAudio(t *testing.T) {
	url := urls[1]
	bodyStr, err := GetHTMLContent(url)

	if err != nil {
		fmt.Println("Failed to get HTML content")
		return
	}

	htmlmeta, err := ParseHTMLMeta(bodyStr)

	if err != nil {
		fmt.Println("Failed to parse HTML meta")
		return
	}

	var aid int
	var bvid string
	var cid int
	var outputPathStem string

	switch len(htmlmeta.VideoData.Pages) {
	case 0:
		fmt.Println("Failed to get page")
		return
	default:
		aid = htmlmeta.Aid
		bvid = htmlmeta.BVid
		cid = htmlmeta.VideoData.Pages[0].Cid
		outputPathStem = fmt.Sprintf("%s_%d", bvid, cid)
	}

	api := "https://api.bilibili.com/x/player/playurl?"
	params := fmt.Sprintf(
		"avid=%d&cid=%d&bvid=%s&qn=%d&type=&otype=json&fourk=1&fnver=0&fnval=2000",
		aid, cid, bvid, 127,
	)
	api += params

	var jsonContentP *PlayURLResponse
	jsonContentP, err = GetJSONContent(api)
	if err != nil {
		fmt.Println("Failed to get JSON content")
		return
	}

	jsonContent := *jsonContentP
	audioArray := jsonContent.Data.Dash.Audio

	if len(audioArray) == 0 {
		fmt.Println("Failed to get audio")
		return
	}

	maxBandwith := 0
	var bestAudio DashItem

	for _, audio := range audioArray {
		if audio.Bandwidth > maxBandwith {
			maxBandwith = audio.Bandwidth
			bestAudio = audio
		}
	}

	err = DownloadPerChunkM4a(bestAudio.BaseURL, outputPathStem, "标题", "艺术家")

	if err != nil {
		fmt.Println("Failed to download audio")
		return
	}

	fmt.Println("Download audio successfully")

}

func TestGetMetaAndDownloadAudio(t *testing.T) {
	url := urls[1]
	meta, err := GetVideoHTMLMeta(url)

	if err != nil {
		fmt.Println("Failed to get video meta")
		return
	}

	aid := meta.Aid
	bvid := meta.BVid
	cid := meta.VideoData.Pages[0].Cid

	err = DownloadAudio(aid, cid, bvid, "标题", "艺术家")

	if err != nil {
		fmt.Println("Failed to download audio")
		return
	}

	fmt.Println("Successfully downloaded audio,check output folder")
}

func TestGetVideoMeta(t *testing.T) {
	url := urls[2]
	bodyStr, err := GetHTMLContent(url)

	if err != nil {
		fmt.Println("Failed to get HTML content")
		return
	}

	htmlmeta, err := ParseHTMLMeta(bodyStr)

	if err != nil {
		fmt.Println("Failed to parse HTML meta")
		return
	}

	jsonData, err := json.MarshalIndent(htmlmeta, "", "  ")

	if err != nil {
		fmt.Println("Failed to marshal JSON")
		return
	}

	err = os.WriteFile("test_output.json", jsonData, 0644)

	if err != nil {
		fmt.Println("Failed to write JSON to file")
		return
	}

	fmt.Println("JSON metadata written to file")

}

func TestTempfileIsExist(t *testing.T) {
	tempFilePath := "../hello.txt"
	var file *os.File
	var err error
	if file, err = os.OpenFile(tempFilePath, os.O_APPEND|os.O_WRONLY, 0644); err != nil {
		fmt.Println("no exists!")
		if file, err = os.Create(tempFilePath); err != nil {
			fmt.Println("create error!")
			return
		}
	}
	fmt.Println("successfully!")
	defer file.Close()
}

func TestGetSizeOfResources(t *testing.T) {
	size, err := GetResourceSize("https://qiang.xiaoyin.link/d/xiaoyin/Download/geek.zip?sign=B9t-JcEmFW4L0EhE3k9a2bzP8j2iUl2rcoQMDJgSoLU=:0") // geek installer from xiaoyin

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(*size)
	}
}

func TestEditMusicMeta(t *testing.T) {
	filePath := "../output/BV1YmFPe4EnY_28086437741.m4a"

	err := taglib.WriteTags(filePath, map[string][]string{
		taglib.Title:  {"孤独患者"},
		taglib.Artist: {"陈奕迅"},
	}, 1)

	if err != nil {
		fmt.Println("写入标签时出错:", err)
	} else {
		fmt.Println("标签写入成功")
	}
}
