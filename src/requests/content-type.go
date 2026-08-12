package requests

import (
	"mime"
	"strings"
)

type ContentType string

const (
	ContentTypeAudio      = "audio/"
	ContentTypeBinary     = "application/octet-stream"
	ContentTypeCSS        = "text/css"
	ContentTypeFont       = "font/"
	ContentTypeHTML       = "text/html"
	ContentTypeIMG        = "image/"
	ContentTypeJSON       = "application/json"
	ContentTypeJSONLD     = "application/ld+json"
	ContentTypeJavaScript = "application/javascript"
	ContentTypeMarkdown   = "text/markdown"
	ContentTypePDF        = "application/pdf"
	ContentTypeText       = "text/plain"
	ContentTypeUnknown    = "unknown/"
	ContentTypeVideo      = "video/"
	ContentTypeXML        = "application/xml"
	ContentTypeZip        = "application/zip"
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
		ContentTypeJavaScript,
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
		ContentTypeJavaScript,
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
