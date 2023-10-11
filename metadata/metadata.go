package metadata

// Metadata is the struct that encapsulates the extracted metadata.
type Metadata struct {
	Audio            Audio    `json:"audio"`
	Author           string   `json:"author"`
	HTML             string   `json:"html"`
	CanonicalURL     string   `json:"canonical_url"`
	Date             string   `json:"date"`
	Description      string   `json:"description"`
	FaviconURL       string   `json:"favicon_url"`
	FeedURLs         []string `json:"feed_url"`
	LeadImageURL     string   `json:"lead_image_url"`
	LeadImageInMeta  bool     `json:"lead_image_in_meta"`
	Image            Image    `json:"image"`
	Lang             string   `json:"lang"`
	Logo             Image    `json:"logo"`
	Meta             Meta     `json:"meta"`
	Publisher        string   `json:"publisher"`
	Title            string   `json:"title"`
	URL              string   `json:"url"`
	Video            Video    `json:"video"`
	ReadableText     string   `json:"readable_text"`
	ReadableHTML     string   `json:"readable_html"`
	ReadableExcerpt  string   `json:"readable_excerpt"`
	ReadableImage    string   `json:"readable_image"`
	ReadableLang     string   `json:"readable_lang"`
	ReadableTitle    string   `json:"readable_title"`
	ReadableByline   string   `json:"readable_byline"`
	ReadableSiteName string   `json:"readable_site_name"`
	SiteName         string   `json:"site_name"`
	Dynamic          map[string]any
}

// Meta is the struct that encapsulates the extracted metadata from <meta> tags
type Meta struct {
	Charset  string `json:"charset"`
	Viewport string `json:"viewport"`
	// ... Other meta properties
}

// Image is the struct that encapsulates the extracted metadata from <img> tags
type Image struct {
	URL        string `json:"url"`
	Type       string `json:"type"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Size       int    `json:"size"`
	SizePretty string `json:"size_pretty"`
}

// Video is the struct that encapsulates the extracted metadata from <video> tags
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

// Audio is the struct that encapsulates the extracted metadata from <audio> tags
type Audio struct {
	URL            string  `json:"url"`
	Type           string  `json:"type"`
	Duration       float64 `json:"duration"`
	DurationPretty string  `json:"duration_pretty"`
	Size           int     `json:"size"`
	SizePretty     string  `json:"size_pretty"`
}
