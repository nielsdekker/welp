package cli

// Mapping between shorthand to longhand form
var mapping = map[string]string{
	"-h":  "--help",
	"-u":  "--url",
	"-m":  "--module",
	"-p":  "--prefix",
	"-o":  "--output",
	"-t":  "--threads",
	"-d":  "--max-depth",
	"-k":  "--insecure",
	"-fc": "--filter-code",
	"-ft": "--filter-type",
}
