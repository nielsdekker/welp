package welp

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/nielsdekker/welp/internal/requests"
)

type Options struct {
	Target             *url.URL
	ConcurrentRequests int
	FilterCodes        []int
	FilterContentType  []requests.ContentType
	ShowHelp           bool
}

func ParseOptions() (Options, error) {
	opt := Options{
		ConcurrentRequests: 10,
		ShowHelp:           false,
	}

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]

		switch arg {
		// URL argument
		case "-u":
			fallthrough
		case "--url":
			targetUrl, _ := url.Parse(os.Args[i+1])
			opt.Target = targetUrl
			i++
		case "-t":
			fallthrough
		case "--threads":
			c, _ := strconv.ParseInt(os.Args[i+1], 10, 32)
			opt.ConcurrentRequests = int(c)
			i++
		case "-h":
			fallthrough
		case "--help":
			opt.ShowHelp = true
		case "-fc":
			fallthrough
		case "--filter-code":
			c, _ := strconv.ParseInt(os.Args[i+1], 10, 32)
			opt.FilterCodes = append(opt.FilterCodes, int(c))
			i++
		case "-ft":
			fallthrough
		case "--filter-type":
			opt.FilterContentType = append(opt.FilterContentType, requests.MatchContentType(os.Args[i+1]))
			i++
		}
	}

	if opt.Target == nil || opt.Target.String() == "" {
		return opt, fmt.Errorf("-u No target URL specified")
	}
	if !strings.HasPrefix(opt.Target.Scheme, "http") {
		return opt, fmt.Errorf("-u %s:// is not a valid scheme", opt.Target.Scheme)
	}
	if opt.ConcurrentRequests <= 0 {
		return opt, fmt.Errorf("-t Can not be zero or smaller, is %d", opt.ConcurrentRequests)
	}

	return opt, nil
}

func (o Options) PrintHelp() {
	for _, s := range []string{
		"WELP",
		"",
		"Usage:",
		fmt.Sprintf("  %-24s%s", "-h, --help", "Shows this help command"),
		fmt.Sprintf("  %-24s%s", "-u, --url", "The target url to start crawling from"),
		fmt.Sprintf("  %-24s%s", "-t, --threads", "Number of concurrent requests"),
		fmt.Sprintf("  %-24s%s", "-fc, --filter-code", "Filters out URL results with the given status code, multiple arguments can be passed"),
		fmt.Sprintf("  %-24s%s", "-ft, --filter-type", "Filters out URL results for the given content type, multiple arguments can be passed"),
	} {
		fmt.Println(s)
	}
}
