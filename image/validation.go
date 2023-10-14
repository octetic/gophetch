package image

// IsValidImage checks if the given URL is a valid image.
func IsValidImage(url string) bool {
	contentType, err := ContentTypeFromURL(url)
	if err != nil {
		return false
	}
	return IsValidImageContentType(contentType)
}

// IsValidFavicon checks if the given URL is a valid favicon.
func IsValidFavicon(url string) bool {
	contentType, err := ContentTypeFromURL(url)
	if err != nil {
		return false
	}
	return IsValidFaviconContentType(contentType)
}
