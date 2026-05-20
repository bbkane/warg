package warg

import (
	"os"

	"go.bbkane.com/warg/metadata"
)

// ParseWithMetadata attaches custom metadata to the parse context,
// accessible later via [CmdContext].ParseMetadata. Useful for injecting test mocks.
func ParseWithMetadata(md metadata.Metadata) ParseOpt {
	return func(poh *ParseOpts) {
		poh.ParseMetadata = md
	}
}

// func ParseWithArgs(args []string) ParseOpt {
// 	return func(poh *ParseOpts) {
// 		poh.Args = args
// 	}
// }

// ParseWithLookupEnv overrides the environment variable lookup function used during parsing.
// Defaults to [os.LookupEnv]. Use [LookupMap] to supply a static map for tests.
func ParseWithLookupEnv(lookup LookupEnv) ParseOpt {
	return func(poh *ParseOpts) {
		poh.LookupEnv = lookup
	}
}

// ParseWithStderr overrides the stderr file passed to [CmdContext].
// The caller is responsible for closing the file.
func ParseWithStderr(stderr *os.File) ParseOpt {
	return func(poh *ParseOpts) {
		poh.Stderr = stderr
	}
}

// ParseWithStdin overrides the stdin file passed to [CmdContext].
// The caller is responsible for closing the file.
func ParseWithStdin(stdin *os.File) ParseOpt {
	return func(poh *ParseOpts) {
		poh.Stdin = stdin
	}
}

// ParseWithStdout overrides the stdout file passed to [CmdContext].
// The caller is responsible for closing the file.
func ParseWithStdout(stdout *os.File) ParseOpt {
	return func(poh *ParseOpts) {
		poh.Stdout = stdout
	}
}
