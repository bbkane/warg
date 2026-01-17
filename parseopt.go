package warg

import (
	"os"

	"go.bbkane.com/warg/metadata"
)

func ParseWithMetadata(md metadata.Metadata) ParseOpt {
	return func(poh *ParseOpts) {
		poh.ParseMetadata = md
	}
}

func ParseWithArgs(args []string) ParseOpt {
	return func(poh *ParseOpts) {
		poh.Args = args
	}
}

func ParseWithLookupEnv(lookup LookupEnv) ParseOpt {
	return func(poh *ParseOpts) {
		poh.LookupEnv = lookup
	}
}

func ParseWithStderr(stderr *os.File) ParseOpt {
	return func(poh *ParseOpts) {
		poh.Stderr = stderr
	}
}

func ParseWithStdin(stdin *os.File) ParseOpt {
	return func(poh *ParseOpts) {
		poh.Stdin = stdin
	}
}

func ParseWithStdout(stdout *os.File) ParseOpt {
	return func(poh *ParseOpts) {
		poh.Stdout = stdout
	}
}
