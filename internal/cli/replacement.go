package cli

import (
	"fmt"
	"strings"
)

type TextReplacement struct {
	glob           string
	parts          []string
	replaceWith    string
	startsWithGlob bool
	endsWithGlob   bool
}

func NewReplacement(
	glob string,
	replaceWith string,
) TextReplacement {
	parts := []string{}
	for s := range strings.SplitSeq(glob, "*") {
		if len(s) > 0 {
			parts = append(parts, s)
		}
	}

	return TextReplacement{
		glob:           glob,
		parts:          parts,
		replaceWith:    replaceWith,
		startsWithGlob: strings.HasPrefix(glob, "*"),
		endsWithGlob:   strings.HasSuffix(glob, "*"),
	}
}

// Applies the text replacement to the given string value, if no replacement is
// necessary the given string is returned.
func (t TextReplacement) Apply(input string) string {
	if len(t.glob) == 0 {
		return input
	}

	i := 0
	partI := 0
	firstMatchStart := -1

	for i < len(input) {
		currentPart := t.parts[partI]
		if input[i] != currentPart[0] {
			// No match for the part
			i++
			continue
		}

		// Try to match the entire part
		n := 0
		for n < len(currentPart) && i+n < len(input) {
			if currentPart[n] != input[i+n] {
				break
			}

			n++
		}

		if n == len(currentPart) {
			// Matched the entire glob
			if partI == 0 {
				firstMatchStart = i
			}

			partI++
			if partI == len(t.parts) {
				// Everything matched
				inputBeforeMatch := input[0:firstMatchStart]
				inputAfterMatch := input[i+n:]

				input = ""
				if t.startsWithGlob {
					inputBeforeMatch = ""
				}
				if t.endsWithGlob {
					inputAfterMatch = ""
				}

				input = fmt.Sprintf("%s%s%s",
					inputBeforeMatch,
					t.replaceWith,
					inputAfterMatch,
				)

				i = len(input) - len(inputAfterMatch)
				partI = 0
			}
		} else {
			// Check multiple so can skip some chars
			i += n
		}
	}

	return input
}
