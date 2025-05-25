package downloader

import (
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

// const bilibiliAPI = "https://api.bilibili.com/x/player/playurl?"

var fakeHeader map[string]string = map[string]string{
	"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.81 Safari/537.36",
	"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
	"Accept-Charset":  "UTF-8,*;q=0.5",
	"Accept-Encoding": "gzip,deflate,sdch",
	"Accept-Language": "en-US,en;q=0.8",
}

const defaultRefer = "https://www.bilibili.com/"

var OutputFilesRootPath string
var TempFilesRootPath string

func GetHTMLContent(url string) (string, error) {

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println("Failed to create request for html")
		return "", err
	}

	req.Header.Set("Referer", defaultRefer)

	for k, v := range fakeHeader {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	var reader io.ReadCloser

	switch res.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(res.Body)
	case "deflate":
		reader = flate.NewReader(res.Body)
	default:
		reader = res.Body
	}
	defer reader.Close() // nolint

	body, err := io.ReadAll(reader)
	if err != nil && err != io.EOF {
		return "", err
	}

	return string(body), nil
}

func GetJSONContent(url string) (*PlayURLResponse, error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println("Failed to create request for json")
		return nil, err
	}

	req.Header.Set("Referer", defaultRefer)

	for k, v := range fakeHeader {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var reader io.ReadCloser

	switch res.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(res.Body)
	case "deflate":
		reader = flate.NewReader(res.Body)
	default:
		reader = res.Body
	}
	defer reader.Close() // nolint

	var jsonString string

	if body, err := io.ReadAll(reader); err != nil && err != io.EOF {
		return nil, err
	} else {
		jsonString = string(body)
	}

	// 解析JSON

	var data PlayURLResponse

	err = json.Unmarshal([]byte(jsonString), &data)

	if err != nil {
		return nil, err
	}

	return &data, nil

}

func GetResourceSize(url string) (*int64, error) {
	req, err := http.NewRequest("HEAD", url, nil)

	if err != nil {
		fmt.Println("Failed to create request for html")
		return nil, err
	}

	req.Header.Set("Referer", defaultRefer)

	for k, v := range fakeHeader {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	contentLength := res.Header.Get("Content-Length")

	if contentLength == "" {
		return nil, errors.New("Content-Length is empty")
	}

	size, err := strconv.ParseInt(contentLength, 10, 64)

	if err != nil {
		return nil, err
	}

	return &size, nil

}

func ParseHTMLMeta(html string) (*HTMLMeta, error) {
	var data HTMLMeta

	multiPageDataString := MatchOneOf(
		html, `window.__INITIAL_STATE__=(.+?);\(function`,
	)

	if multiPageDataString == nil {
		return &data, errors.New("this page has no metadata")
	}

	if err := json.Unmarshal([]byte(multiPageDataString[1]), &data); err != nil {
		return nil, err
	}

	return &data, nil

}

func MatchOneOf(text string, patterns ...string) []string {
	var (
		re    *regexp.Regexp
		value []string
	)
	for _, pattern := range patterns {
		// (?flags): set flags within current group; non-capturing
		// s: let . match \n (default false)
		// https://github.com/google/re2/wiki/Syntax
		re = regexp.MustCompile(pattern)
		value = re.FindStringSubmatch(text)
		if len(value) > 0 {
			return value
		}
	}
	return nil
}

func DownloadPerChunkM4a(url string, outputPathStem string) (err error) {
	tempFilePath := filepath.Join(TempFilesRootPath, outputPathStem+".download")
	outputFilePath := filepath.Join(OutputFilesRootPath, outputPathStem+".m4a")

	var tempFile *os.File

	tempFile, err = os.OpenFile(tempFilePath, os.O_APPEND|os.O_WRONLY, 0644) // 注意是追加模式

	if err != nil {
		tempFile, err = os.Create(tempFilePath)

		if err != nil {
			fmt.Println("temp file not exists and create error!")
			return
		}
	}

	tempFileInfo, _ := os.Stat(tempFilePath)
	tempFileSize := tempFileInfo.Size()

	headers := map[string]string{
		"Referer": defaultRefer,
	}

	defer func() {
		tempFile.Close()
		if err == nil {
			os.Rename(tempFilePath, outputFilePath)
		}
	}()

	var chunkSize int64 = 1024 * 1024 * 1024 // 1MB
	var remainSize int64

	if remainSizeP, err := GetResourceSize(url); err != nil {
		fmt.Println("try to get resource size failed, try \"GET\" method instead")
		return err
	} else {
		remainSize = *remainSizeP
	}

	remainSize -= tempFileSize
	totalChunkNums := remainSize / chunkSize

	if remainSize%chunkSize != 0 {
		totalChunkNums++
	}

	start := tempFileSize

	var i int64 = 1
	var end int64

	for ; i <= totalChunkNums; i++ {
		end = start + chunkSize - 1
		headers["Range"] = fmt.Sprintf("bytes=%d-%d", start, end)

		req, err := http.NewRequest("GET", url, nil)

		if err != nil {
			fmt.Println("error during init download request")
			return err
		}

		for key, value := range headers {
			req.Header.Set(key, value)
		}

		client := &http.Client{}

		res, err := client.Do(req)

		if err != nil {
			fmt.Println("error during downloading")
			return err
		}

		defer res.Body.Close()

		io.Copy(tempFile, res.Body)
	}

	return nil
}
