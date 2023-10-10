package metadata

type Metadata struct {
	Audio       Audio  `json:"audio"`
	Author      string `json:"author"`
	HTML        string `json:"html"`
	Date        string `json:"date"`
	Description string `json:"description"`
	Image       Image  `json:"image"`
	Lang        string `json:"lang"`
	Logo        Image  `json:"logo"`
	Meta        Meta   `json:"meta"`
	Publisher   string `json:"publisher"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Video       Video  `json:"video"`
	Dynamic     map[string]any
}

type Meta struct {
	Charset  string `json:"charset"`
	Viewport string `json:"viewport"`
	// ... Other meta properties
}

type Image struct {
	URL        string `json:"url"`
	Type       string `json:"type"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Size       int    `json:"size"`
	SizePretty string `json:"size_pretty"`
}

type Video struct {
	URL            string  `json:"url"`
	Type           string  `json:"type"`
	Duration       float64 `json:"duration"`
	DurationPretty string  `json:"duration_pretty"`
	Width          int     `json:"width"`
	Height         int     `json:"height"`
	Size           int     `json:"size"`
	SizePretty     string  `json:"size_pretty"`
}

type Audio struct {
	URL            string  `json:"url"`
	Type           string  `json:"type"`
	Duration       float64 `json:"duration"`
	DurationPretty string  `json:"duration_pretty"`
	Size           int     `json:"size"`
	SizePretty     string  `json:"size_pretty"`
}
