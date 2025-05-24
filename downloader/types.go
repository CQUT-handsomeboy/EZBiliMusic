package downloader

type episode struct {
	Aid   int    `json:"aid"`
	Cid   int    `json:"cid"`
	Title string `json:"title"`
	BVid  string `json:"bvid"`
}

type multiEpisodeData struct {
	Seasionid int       `json:"season_id"`
	Episodes  []episode `json:"episodes"`
}

type videoPagesData struct {
	Cid  int    `json:"cid"`
	Part string `json:"part"`
	Page int    `json:"page"`
}

type multiPageVideoData struct {
	Title string           `json:"title"`
	Pages []videoPagesData `json:"pages"`
}

type HTMLMeta struct {
	Aid       int                `json:"aid"`
	BVid      string             `json:"bvid"`
	Sections  []multiEpisodeData `json:"sections"`
	VideoData multiPageVideoData `json:"videoData"`
}

// API JSON
type PlayURLResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	TTL     int          `json:"ttl"`
	Data    PlayURLData  `json:"data"`
}

type PlayURLData struct {
	AcceptDescription []string         `json:"accept_description"`
	AcceptFormat      string           `json:"accept_format"`
	AcceptQuality     []int            `json:"accept_quality"`
	Dash              DashStream       `json:"dash"`
	Format            string           `json:"format"`
	From              string           `json:"from"`
	HighFormat        interface{}      `json:"high_format"` // 可能是null
	LastPlayCid       int              `json:"last_play_cid"`
	LastPlayTime      int              `json:"last_play_time"`
	Message           string           `json:"message"`
	PlayConf          PlayConfig       `json:"play_conf"`
	Quality           int              `json:"quality"`
	Result            string           `json:"result"`
	SeekParam         string           `json:"seek_param"`
	SeekType          string           `json:"seek_type"`
	SupportFormats    []SupportFormat  `json:"support_formats"`
	Timelength        int              `json:"timelength"`
	VideoCodecid      int              `json:"video_codecid"`
	ViewInfo          interface{}      `json:"view_info"` // 可能是null
}

type DashStream struct {
	Audio         []DashItem `json:"audio"`
	Dolby         DolbyInfo `json:"dolby"`
	Duration      int       `json:"duration"`
	Flac          interface{} `json:"flac"` // 可能是null
	MinBufferTime float64   `json:"min_buffer_time"`
	Video         []DashItem `json:"video"`
}

type DashItem struct {
	SegmentBase   SegmentBase `json:"SegmentBase"`
	BackupUrl     []string    `json:"backupUrl"`
	BackupURL     []string    `json:"backup_url"`
	Bandwidth     int         `json:"bandwidth"`
	BaseUrl       string      `json:"baseUrl"`
	BaseURL       string      `json:"base_url"`
	Codecid       int         `json:"codecid"`
	Codecs        string      `json:"codecs"`
	FrameRate     string      `json:"frameRate"`
	FrameRateAlt  string      `json:"frame_rate"`
	Height        int         `json:"height"`
	ID            int         `json:"id"`
	MimeType      string      `json:"mimeType"`
	MimeTypeAlt   string      `json:"mime_type"`
	SAR           string      `json:"sar"`
	SegmentBaseAlt SegmentBase `json:"segment_base"`
	StartWithSap  int         `json:"startWithSap"`
	StartWithSapAlt int       `json:"start_with_sap"`
	Width         int         `json:"width"`
}

type SegmentBase struct {
	Initialization string `json:"Initialization"`
	IndexRange     string `json:"indexRange"`
}

type DolbyInfo struct {
	Audio interface{} `json:"audio"` // 可能是null
	Type  int         `json:"type"`
}

type PlayConfig struct {
	IsNewDescription bool `json:"is_new_description"`
}

type SupportFormat struct {
	Codecs         []string `json:"codecs"`
	DisplayDesc    string   `json:"display_desc"`
	Format         string   `json:"format"`
	NewDescription string   `json:"new_description"`
	Quality        int      `json:"quality"`
	Superscript    string   `json:"superscript"`
}