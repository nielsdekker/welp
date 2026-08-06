package requests

import (
	"mime"
	"strings"
)

type ContentType string

const (
	ContentTypeCSS      = "text/css"
	ContentTypeHTML     = "text/html"
	ContentTypeMarkdown = "text/markdown"
	ContentTypeText     = "text/plain"
	ContentTypeBinary   = "application/octet-stream"
	ContentTypeJSON     = "application/json"
	ContentTypeJSONLD   = "application/ld+json"
	ContentTypePDF      = "application/pdf"
	ContentTypeXML      = "application/xml"
	ContentTypeZip      = "application/zip"
	ContentTypeAudio    = "audio/"
	ContentTypeFont     = "font/"
	ContentTypeIMG      = "image/"
	ContentTypeVideo    = "video/"
	ContentTypeUnknown  = "unknown/"
)

func ParseContentType(headerValue string) ContentType {
	mediatype, _, err := mime.ParseMediaType(headerValue)

	if err != nil {
		return ContentTypeUnknown
	}

	for _, ct := range []ContentType{
		ContentTypeCSS,
		ContentTypeHTML,
		ContentTypeMarkdown,
		ContentTypeText,
		ContentTypeBinary,
		ContentTypeJSON,
		ContentTypeJSONLD,
		ContentTypePDF,
		ContentTypeXML,
		ContentTypeZip,
		ContentTypeAudio,
		ContentTypeFont,
		ContentTypeIMG,
		ContentTypeVideo,
	} {
		if strings.HasPrefix(mediatype, string(ct)) {
			return ct
		}
	}

	return ContentTypeUnknown
}

func MatchContentType(value string) []ContentType {
	matches := []ContentType{}
	for _, ct := range []ContentType{
		ContentTypeCSS,
		ContentTypeHTML,
		ContentTypeMarkdown,
		ContentTypeText,
		ContentTypeBinary,
		ContentTypeJSON,
		ContentTypeJSONLD,
		ContentTypePDF,
		ContentTypeXML,
		ContentTypeZip,
		ContentTypeAudio,
		ContentTypeFont,
		ContentTypeIMG,
		ContentTypeVideo,
	} {
		if strings.Contains(string(ct), value) {
			matches = append(matches, ct)
		}
	}

	return matches
}
