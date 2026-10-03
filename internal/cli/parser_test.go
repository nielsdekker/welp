package cli_test

import (
	"net/url"
	"testing"

	"github.com/nielsdekker/welp/internal/_tests/asserts"
	"github.com/nielsdekker/welp/internal/cli"
)

func TestParseValidOptions(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected cli.Options
	}{
		{"Help short argument", []string{"-h"}, genOpt(func(o *cli.Options) { o.ShowHelp = true })},
		{"Help long argument", []string{"--help"}, genOpt(func(o *cli.Options) { o.ShowHelp = true })},
		{"Url short argument", []string{"-uhttp://localhost.test"}, genOpt(func(o *cli.Options) {
			u, _ := url.Parse("http://localhost.test")
			o.Target = u
		})},
		{"Url long argument", []string{"--url", "http://localhost.test"}, genOpt(func(o *cli.Options) {
			u, _ := url.Parse("http://localhost.test")
			o.Target = u
		})},
		{"Prefix argument", []string{"-pfoo", "-p", "bar", "--prefix", "baz", "hello", "world"}, genOpt(func(o *cli.Options) {
			o.Prefixes = map[string]struct{}{
				"foo":   struct{}{},
				"bar":   struct{}{},
				"baz":   struct{}{},
				"hello": struct{}{},
				"world": struct{}{},
			}
		})},
		{"Threads short argument", []string{"-t100"}, genOpt(func(o *cli.Options) { o.ConcurrentRequests = 100 })},
		{"Threads long argument", []string{"--thread", "100"}, genOpt(func(o *cli.Options) { o.ConcurrentRequests = 100 })},
		{"Max-depth short argument", []string{"-d100"}, genOpt(func(o *cli.Options) { o.SearchDepth = 100 })},
		{"Max-depth long argument", []string{"--max-depth", "100"}, genOpt(func(o *cli.Options) { o.SearchDepth = 100 })},
		{"Insecure short argument", []string{"-k"}, genOpt(func(o *cli.Options) { o.SSLIgnore = true })},
		{"Insecure long argument", []string{"--insecure"}, genOpt(func(o *cli.Options) { o.SSLIgnore = true })},
		{"Filter codes", []string{"-fc400", "--filter-code", "401", "-fc", "402", "403"}, genOpt(func(o *cli.Options) {
			o.FilterCodes = []int{400, 401, 402, 403}
		})},
		{"Filter types", []string{"-fthtml/", "--filter-type", "img/", "-ft", "plain", "xml"}, genOpt(func(o *cli.Options) {
			o.FilterContentType = []string{"html/", "img/", "plain", "xml"}
		})},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// For now ignore the errors
			actual, _ := cli.Parse(tt.args)
			compOpt(t, tt.expected, actual)
		})
	}
}

func genOpt(overrides func(*cli.Options)) cli.Options {
	opt := &cli.Options{
		// Defaults
		SSLIgnore:          false,
		Prefixes:           map[string]struct{}{},
		ConcurrentRequests: 10,
		SearchDepth:        25,
	}

	overrides(opt)

	opt.ApplyUnsetDefaults()
	return *opt
}

func compOpt(t *testing.T, expected cli.Options, actual cli.Options) {
	asserts.Eq(t, expected.ShowHelp, actual.ShowHelp)
	if expected.Target != nil {
		if actual.Target == nil {
			t.Errorf("Expected %s but was nil", expected.Target.String())
		} else {
			asserts.Eq(t, expected.Target.String(), actual.Target.String())
		}
	}
	asserts.KeysEq(t, expected.Prefixes, actual.Prefixes)
	asserts.Eq(t, expected.ConcurrentRequests, actual.ConcurrentRequests)
	asserts.Eq(t, expected.SearchDepth, actual.SearchDepth)
	asserts.Eq(t, expected.SSLIgnore, actual.SSLIgnore)
	asserts.SliceEq(t, expected.FilterCodes, actual.FilterCodes)
	asserts.SliceEq(t, expected.FilterContentType, actual.FilterContentType)
}
