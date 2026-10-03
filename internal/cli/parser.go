package cli

import (
	_ "embed"
	"errors"
	"fmt"
	"strings"
)

func Parse(args []string) (Options, error) {
	// Follow the next rules:
	// - One `-` is new arg, greedy match until we match a valid option
	//   - Any characters after that are options for that arg
	//   - Any args after that not starting with a `-` are other options for
	//     that arg
	// - Two `--` is a full name arg. Match completely
	//   - Args after that not starting with a `-` are options for that arg.
	//
	// Following these rules the following is possible:
	//
	// ```bash
	// # Targets localhost
	// welp -uhttp://localhost
	//
	// # Also targets localhost
	// welp --url http://localhost
	//
	// # This is invalid (full name arg without space)
	// welp --urlhttp://localhost
	// ```

	foundErrors := []error{}
	parsedArgs := map[string][]string{}
	lastArg := ""

	// Parse the args
	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			// Full arg
			lastArg = arg
		} else if strings.HasPrefix(arg, "-") {
			// Single arg, try to get the corresponding full arg
			remainingValue := ""
			found := false

			for k, v := range mapping {

				if strings.HasPrefix(arg, k) {
					lastArg = v
					found = true
					remainingValue = arg[len(k):]
					break
				}
			}

			if !found {
				foundErrors = append(foundErrors, fmt.Errorf("%s is not a valid argument", arg))
			} else if remainingValue != "" {
				// This was a argument like `-fc404`, so make sure the 404 is
				// added as an option
				parsedArgs[lastArg] = append(parsedArgs[lastArg], remainingValue)
			}
		} else {
			// An option
			if lastArg != "" {
				parsedArgs[lastArg] = append(parsedArgs[lastArg], arg)
			} else {
				foundErrors = append(foundErrors, fmt.Errorf("%s is not a valid argument", arg))
			}
		}

		// Make sure the options is set if
		if _, ok := parsedArgs[lastArg]; !ok {
			parsedArgs[lastArg] = []string{}
		}
	}

	// Apply the args
	opt := &Options{
		// Defaults
		SSLIgnore:          false,
		Prefixes:           map[string]struct{}{},
		ConcurrentRequests: 10,
		SearchDepth:        25,
	}

	for k, v := range parsedArgs {
		handler, ok := argParsers[k]
		if !ok {
			foundErrors = append(foundErrors, fmt.Errorf("%s is not a valid argument", k))
			continue
		}

		// Also handle error
		if err := handler(opt, v); err != nil {
			foundErrors = append(foundErrors, err)
		}
	}

	foundErrors = append(foundErrors, opt.IsValid()...)
	opt.ApplyUnsetDefaults()

	return *opt, errors.Join(foundErrors...)
}
