package welp_test

import (
	"context"
	"testing"

	"github.com/nielsdekker/welp/src/module"
	"github.com/nielsdekker/welp/src/welp"
	"github.com/nielsdekker/welp/tests/mocks"
)

func TestSpa(t *testing.T) {
	mockPool := mocks.NewMockPool()
	w := welp.New(
		mustParse("http://localhost/spa/"),
		mockPool,
		[]module.Module{},
	)

	w.StartCrawl(context.Background())

	urlResult := w.CrawledURLs()

	eq(t, len(urlResult), 4)
	hasPath(t, urlResult, "/spa/")
	hasPath(t, urlResult, "/spa/style/default.css")
	hasPath(t, urlResult, "/spa/stylesheet")
	hasPath(t, urlResult, "/spa/text/css")
}
