package cli

import (
	"fmt"
	"net/url"
	"strconv"
)

const ARG_HELP = "--help"
const ARG_URL = "--url"
const ARG_PREFIX = "--prefix"
const ARG_THREADS = "--threads"
const ARG_MAX_DEPTH = "--max-depth"
const ARG_INSECURE = "--insecure"
const ARG_FILTER_CODE = "--filter-code"
const ARG_FILTER_TYPE = "--filter-type"

// Mapping between shorthand to longhand form
var mapping = map[string]string{
	"-h":  ARG_HELP,
	"-u":  ARG_URL,
	"-p":  ARG_PREFIX,
	"-t":  ARG_THREADS,
	"-d":  ARG_MAX_DEPTH,
	"-k":  ARG_INSECURE,
	"-fc": ARG_FILTER_CODE,
	"-ft": ARG_FILTER_TYPE,
}

var argParsers = map[string]func(*Options, []string) error{
	ARG_HELP: func(o *Options, s []string) error {
		if len(s) > 0 {
			return fmt.Errorf("No values expected for help argument")
		}
		o.ShowHelp = true
		return nil
	},
	ARG_URL: func(o *Options, s []string) error {
		if len(s) != 1 {
			return fmt.Errorf("Only one value expected for url argument")
		}
		u, err := url.Parse(s[0])
		o.Target = u
		return err
	},
	ARG_PREFIX: func(o *Options, s []string) error {
		if len(s) == 0 {
			return fmt.Errorf("Expected one or more values for prefix argument")
		}

		o.Prefixes = asSet(s)
		return nil
	},
	ARG_THREADS: func(o *Options, s []string) error {
		if len(s) != 1 {
			return fmt.Errorf("Expected one value for threads argument")
		}

		i, err := asIntSlice(s)
		if err != nil {
			return err
		} else {
			o.ConcurrentRequests = i[0]
			return nil
		}
	},
	ARG_MAX_DEPTH: func(o *Options, s []string) error {
		if len(s) != 1 {
			return fmt.Errorf("Expected one value for threads argument")
		}

		i, err := asIntSlice(s)
		if err != nil {
			return err
		} else {
			o.SearchDepth = i[0]
			return nil
		}
	},
	ARG_INSECURE: func(o *Options, s []string) error {
		if len(s) > 0 {
			return fmt.Errorf("No values expected for insecure argument")
		}

		o.SSLIgnore = true
		return nil
	},
	ARG_FILTER_CODE: func(o *Options, s []string) error {
		if len(s) == 0 {
			return fmt.Errorf("Expected one or more values for the filter-code argument")
		}

		i, err := asIntSlice(s)
		o.FilterCodes = i
		return err
	},
	ARG_FILTER_TYPE: func(o *Options, s []string) error {
		if len(s) == 0 {
			return fmt.Errorf("Expected one or more values for the filter-type argument")
		}

		o.FilterContentType = s
		return nil
	},
}

func asSet(s []string) map[string]struct{} {
	m := map[string]struct{}{}
	for _, pre := range s {
		m[pre] = struct{}{}
	}
	return m
}

func asIntSlice(s []string) ([]int, error) {
	i := []int{}
	for _, raw := range s {
		fc, err := strconv.ParseInt(raw, 10, 32)
		if err != nil {
			return i, err
		}
		i = append(i, int(fc))
	}

	return i, nil
}
