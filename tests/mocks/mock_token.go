package mocks

import (
	"context"
	"net/http"

	"github.com/nielsdekker/welp/src/requests"
)

type tokenInJsMock struct{}

func NewTokenInJSMock() requests.Pool {
	return tokenInJsMock{}
}

func (t tokenInJsMock) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	index := `
		<html>
			<head></head>
			<body>
				<script type="module" src="scripts/main.js"></script>
			</body>
		</html>
	`
	jsMain := `
		import t from "./tokens.js"

		func main() {
			console.log("t: " + tokens)
		}
	`
	jsToken := `
		const GH="ghp_123abc"
		const LLM="sk-123-456"
		const JWT="ey123.eyabc.def"
		const API="zILEpsOAxrvFnMOOxZTMkVrItcyZw6jPpCHHolXFnsaiy5/OgMSywrjGlMW4zLHNhsqLyLrIsD8kxpY="
	`

	switch req.URL.Path {
	case "/scripts/main.js":
		return createResponse(req, 200, jsMain, "application/javascript"), nil
	case "/scripts/tokens.js":
		return createResponse(req, 200, jsToken, "application/javascript"), nil
	case "":
		fallthrough
	case "/":
		return createResponse(req, 200, index, "text/html"), nil
	default:
		return createResponse(req, 404, "Not found", "text/plain"), nil
	}
}
