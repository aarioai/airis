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
	Duration Second `json:"duration"`
}

type DocSrc struct {
	FileSrc
}

type ImgSrc struct {
	FileSrc

	CropUrl UrlPattern `json:"crop_url"` // e.g.  https://xxx/img.jpg?width={width:int}&height={height:int}
	Width   int        `json:"width"`
	Height  int        `json:"height"`
	Allowed [][2]int   `json:"allowed"` // allowed [width, height][]
}

type VideoSrc struct {
	FileSrc
	Preview  *ImgSrc  `json:"preview"`
	Duration Second   `json:"duration"`
	Width    int      `json:"width"`
	Height   int      `json:"height"`
	Allowed  [][2]int `json:"allowed"` // allowed [width, height][]
}
