package tree

import (
	"fmt"
	"strings"
	"unicode/utf8"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_javascript "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
)

func StringValues(raw []byte, contentType string) map[string]struct{} {
	switch contentType {
	case "application/javascript":
		return parseJS(raw)
	default:
		return parseUnknown(raw)
	}
}

func parseJS(raw []byte) map[string]struct{} {
	p := tree_sitter.NewParser()
	lang_js := tree_sitter.NewLanguage(tree_sitter_javascript.Language())
	p.SetLanguage(lang_js)

	tree := p.Parse(raw, nil)
	defer tree.Close()

	q, err := tree_sitter.NewQuery(lang_js, `[(string) (template_string)] @str`)

	if err != nil {
		fmt.Printf("err: Unable to parse %v", err)
		return make(map[string]struct{})
	}

	cursor := tree_sitter.NewQueryCursor()
	defer cursor.Close()

	retValue := make(map[string]struct{})
	matches := cursor.Matches(q, tree.RootNode(), raw)

	for {
		currentMatch := matches.Next()
		if currentMatch == nil {
			break
		}

		for _, c := range currentMatch.Captures {
			// This still contains the `"` or ``` values so remove those. Should
			// be possible in the treesitter query but then template strings get
			// multiple values/are split which is not something we want.
			strValue := strings.TrimFunc(c.Node.Utf8Text(raw), func(r rune) bool {
				return r == '\'' || r == '"' || r == '`'
			})
			retValue[strValue] = struct{}{}
		}
	}

	return retValue
}

func parseUnknown(raw []byte) map[string]struct{} {
	result := make(map[string]struct{})

	quoteIndices := map[byte]int{
		'\'': -1,
		'"':  -1,
		'`':  -1,
	}

	for rawIndex, b := range raw {
		if quoteIndex, ok := quoteIndices[b]; ok {
			// This is a quote/string character so parse it
			if quoteIndex >= 0 {
				bytes := raw[quoteIndex+1 : rawIndex]
				if utf8.Valid(bytes) {
					foundValue := strings.TrimSpace(string(bytes))
					result[foundValue] = struct{}{}
				}
				quoteIndices[b] = -1
			} else {
				quoteIndices[b] = rawIndex
			}
		}
	}

	return result
}
