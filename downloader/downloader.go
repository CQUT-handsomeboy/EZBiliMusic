package downloader

import (
	"fmt"
)

func GetVideoHTMLMeta(url string) (*HTMLMeta, error) {
	bodyStr, err := GetHTMLContent(url)

	if err != nil {
		fmt.Println("Failed to get HTML content")
		return nil, err
	}

	return ParseHTMLMeta(bodyStr)
}

func DownloadAudio(aid int, cid int, bvid string, title string, artist string) error {
	outputPathStem := fmt.Sprintf("%s_%d", bvid, cid)

	api := "https://api.bilibili.com/x/player/playurl?"
	params := fmt.Sprintf(
		"avid=%d&cid=%d&bvid=%s&qn=%d&type=&otype=json&fourk=1&fnver=0&fnval=2000",
		aid, cid, bvid, 127,
	)
	api += params

	var jsonContentP *PlayURLResponse
	jsonContentP, err := GetJSONContent(api)

	if err != nil {
		return fmt.Errorf("failed to get JSON content")
	}

	jsonContent := *jsonContentP
	audioArray := jsonContent.Data.Dash.Audio

	if len(audioArray) == 0 {
		return fmt.Errorf("failed to get JSON content")
	}

	maxBandwith := 0
	var bestAudio DashItem

	for _, audio := range audioArray {
		if audio.Bandwidth > maxBandwith {
			maxBandwith = audio.Bandwidth
			bestAudio = audio
		}
	}

	pathPointer, err := DownloadPerChunkM4a(bestAudio.BaseURL, outputPathStem)

	if err != nil {
		fmt.Println("failed to download audio")
		return err
	}

	fmt.Println("download audio successfully")

	err = AddMusicTag(*pathPointer, title, artist)

	if err != nil {
		fmt.Print("failed to add music tag:")
		fmt.Printf("%v\n", err)
		return err
	}

	return nil

}
