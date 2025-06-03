package atype

type VideoPattern struct {
}

type ImagePattern struct {
	Height    int    `json:"height"`
	Width     int    `json:"width"`
	Quality   int    `json:"quality"`
	MaxWidth  int    `json:"max_width"`
	MaxHeight int    `json:"max_height"`
	Watermark string `json:"watermark"`
}

type FileSrc struct {
	Provider        int        `json:"provider"`
	URL             URL        `json:"url"`         // e.g. https://xxx/img.jpg
	UrlPattern      UrlPattern `json:"url_pattern"` // e.g. https://xxx/img.jpg?maxwidth={max_width:int}
	AlterUrlPattern UrlPattern `json:"alter_url_pattern"`
	BaseURL         URL        `json:"base_url"`
	Path            Path       `json:"path"`
	Filetype        FileType   `json:"filetype"`
	Size            int        `json:"size"`
	Info            string     `json:"info"`
	Checksum        string     `json:"checksum"`
	Jsonkey         string     `json:"jsonkey"`
}

type AudioSrc struct {
	FileSrc
	Bitrate    int    `json:"bitrate"` // bit per second
	Duration   Second `json:"duration"`
	SampleRate int    `json:"sample_rate"` // HZ
}

type DocSrc struct {
	FileSrc
}

type ImgSrc struct {
	FileSrc
	Allowed        [][2]int   `json:"allowed"`          // allowed [width, height][]
	CropUrlPattern UrlPattern `json:"crop_url_pattern"` // e.g.  https://xxx/img.jpg?width={width:int}&height={height:int}
	Height         int        `json:"height"`
	Width          int        `json:"width"`
}

type VideoSrc struct {
	FileSrc
	Allowed    [][2]int `json:"allowed"` // allowed [width, height][]
	Bitrate    int      `json:"bitrate"`
	Codec      string   `json:"codec"`
	Duration   Second   `json:"duration"`
	Framerate  int      `json:"framerate"`
	Height     int      `json:"height"`
	Preview    *ImgSrc  `json:"preview"`
	SampleRate int      `json:"sample_rate"`
	Width      int      `json:"width"`
}
