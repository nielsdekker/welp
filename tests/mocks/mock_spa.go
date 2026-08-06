package mocks

import (
	"net/http"

	"github.com/nielsdekker/welp/src/requests"
)

type spaMock struct{}

// Goal of these mocks is to simulate a single page application that returns the
// same index.html for each request
func NewSpaMock() requests.Pool {
	return spaMock{}
}

func (s spaMock) Do(req *http.Request) (*http.Response, error) {
	spaIndex := `
		<html>
			<head>
				<link rel="stylesheet" type="text/css" href="style/default.css">
			</head>
			<body>
				<nav>
					<a href="/other">other page</a>
				</nav>

				<main>
				Welcome to the this single page application
				</main>
			</body>
		</html>
	`
	spaCss := `
		html {
			color: red;
		}
	`

	switch req.URL.Path {
	case "/style/default.css":
		return createResponse(200, spaCss, "text/html"), nil
	default:
		return createResponse(200, spaIndex, "text/html"), nil
	}
}
