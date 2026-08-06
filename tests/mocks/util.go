package mocks

import (
	"bytes"
	"io"
	"net/http"
)

func createResponse(
	request *http.Request,
	statuscode int,
	body string,
	contentType string,
) *http.Response {
	return &http.Response{
		StatusCode: statuscode,
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Header: http.Header{
			"Content-Type": []string{contentType},
		},
		Request: request,
	}
}
