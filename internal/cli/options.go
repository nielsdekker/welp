package cli

import (
	_ "embed"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

//go:embed usage.md
var usage string

type Options struct {
	Target             *url.URL
	ConcurrentRequests int
	ShowHelp           bool
	SearchDepth        int
	Modules            map[string]struct{}
	Prefixes           map[string]struct{}
	OutputFile         string
	FilterCodes        []int
	FilterContentType  []string
	TextReplacements   []TextReplacement
	TextMinLength      int
	TextMaxLength      int
}

func Parse() (Options, error) {
	opt := Options{
		Modules:            map[string]struct{}{},
		Prefixes:           map[string]struct{}{},
		ConcurrentRequests: 10,
		SearchDepth:        25,
		TextMinLength:      4,
		TextMaxLength:      200,
	}

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]

		if longhand, ok := mapping[arg]; ok {
			// Arg was given in shorthand form, use longhand form for matching
			// the correct options.
			arg = longhand
		}

		switch arg {
		case "--help":
			opt.ShowHelp = true
		case "--url":
			i++
			if u, err := urlArg(i); err == nil {
				opt.Target = u
			}
		case "--module":
			i++
			if m, err := stringArg(i); err == nil {
				opt.Modules[m] = struct{}{}
			}
		case "--prefix":
			i++
			if p, err := stringArg(i); err == nil {
				opt.Prefixes[p] = struct{}{}
			}
		case "--output":
			i++
			if s, err := stringArg(i); err == nil {
				opt.OutputFile = s
			}
		case "--threads":
			i++
			if c, err := intArg(i); err == nil {
				opt.ConcurrentRequests = c
			}
		case "--depth":
			i++
			if c, err := intArg(i); err == nil {
				opt.SearchDepth = c
			}
		case "--filter-code":
			i++
			if c, err := intArg(i); err == nil {
				opt.FilterCodes = append(opt.FilterCodes, c)
			}
		case "--filter-type":
			i++
			if t, err := stringArg(i); err == nil {
				opt.FilterContentType = append(opt.FilterContentType, t)
			}
		case "--text-min-length":
			i++
			if l, err := intArg(i); err == nil {
				opt.TextMinLength = l
			}
		case "--text-max-length":
			i++
			if l, err := intArg(i); err == nil {
				opt.TextMaxLength = l
			}
		case "--text-replace":
			g, errG := stringArg(i + 1)
			r, errR := stringArg(i + 2)
			i += 2

			if errG == nil && errR == nil {
				opt.TextReplacements = append(opt.TextReplacements, NewReplacement(g, r))
			}
		}
	}

	return opt, opt.validate()
}

func intArg(index int) (int, error) {
	if index >= len(os.Args) {
		return 0, fmt.Errorf("Out of bounds")
	}

	c, err := strconv.ParseInt(os.Args[index], 10, 32)
	if err != nil {
		return 0, err
	} else {
		return int(c), nil
	}
}

func urlArg(index int) (*url.URL, error) {
	if index >= len(os.Args) {
		return nil, fmt.Errorf("Out of bounds")
	}

	return url.Parse(os.Args[index])
}

func stringArg(index int) (string, error) {
	if index >= len(os.Args) {
		return "", fmt.Errorf("Out of bounds")
	}

	return os.Args[index], nil
}

func (o *Options) validate() error {
	if o.Target == nil || o.Target.String() == "" {
		return fmt.Errorf("-u No target URL specified")
	}
	if !strings.HasPrefix(o.Target.Scheme, "http") {
		return fmt.Errorf("-u %s:// is not a valid scheme", o.Target.Scheme)
	}
	if o.ConcurrentRequests <= 0 {
		return fmt.Errorf("-t Can not be zero or smaller, is %d", o.ConcurrentRequests)
	}

	if len(o.FilterCodes) == 0 {
		o.FilterCodes = []int{404}
	}

	return nil
}

func (o Options) Usage() string {
	return usage
}
