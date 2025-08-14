package warg

import (
	"context"
	"os"
)

func ParseContext(ctx context.Context) ParseOpt {
	return func(poh *ParseOpts) {
		poh.Context = ctx
	}
}

func Args(args []string) ParseOpt {
	return func(poh *ParseOpts) {
		poh.Args = args
	}
}

func ParseLookupEnv(lookup LookupEnv) ParseOpt {
	return func(poh *ParseOpts) {
		poh.LookupEnv = lookup
	}
}

func Stderr(stderr *os.File) ParseOpt {
	return func(poh *ParseOpts) {
		poh.Stderr = stderr
	}
}

func Stdout(stdout *os.File) ParseOpt {
	return func(poh *ParseOpts) {
		poh.Stdout = stdout
	}
}
