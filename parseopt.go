package warg

import (
	"context"
	"os"
)

func ParseWithContext(ctx context.Context) ParseOpt {
	return func(poh *ParseOpts) {
		poh.Context = ctx
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

func ParseWithStdout(stdout *os.File) ParseOpt {
	return func(poh *ParseOpts) {
		poh.Stdout = stdout
	}
}
