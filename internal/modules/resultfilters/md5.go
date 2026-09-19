package resultfilters

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/nielsdekker/welp/internal/welp"
)

type md5Filter struct {
	cache map[string]struct{}
}

var _ welp.ResultFilterModule = &md5Filter{}

func NewMD5Filter() *md5Filter {
	return &md5Filter{
		cache: make(map[string]struct{}),
	}
}

func (m *md5Filter) ShouldCrawl(result welp.CrawlResult) bool {
	md5sum := md5.New()
	md5Value := hex.EncodeToString(md5sum.Sum(result.Raw))

	if _, ok := m.cache[md5Value]; ok {
		return false
	} else {
		// Make sure to update the cache with this value
		m.cache[md5Value] = struct{}{}
		return true
	}
}
