package mocks

import (
	"net/http"
	"strings"
)

type spaMock struct{}

// Goal of these mocks is to simulate a single page application that returns the
// same index.html for each request
func newSpaMock() mock {
	return spaMock{}
}

func (s spaMock) handle(path string) (*http.Response, bool) {
	spaIndex := `
<html>
	<head>
		<link rel="stylesheet" type="text/css" href="style/default.css">
	</head>
	<body>
		Welcome to the main SPA page :)
	</body>
</html>
	`
	spaCss := `
html {
	color: red;
}
	`

	if !strings.HasPrefix(path, "/spa/") {
		return nil, false
	}

	switch path {
	case "/spa/style/default.css":
		return createResponse(200, spaCss, "text/html"), true
	default:
		return createResponse(200, spaIndex, "text/html"), true
	}
}
