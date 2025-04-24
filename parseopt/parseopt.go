package parseopt

import (
	"context"
	"os"

	"go.bbkane.com/warg/wargcore"
)

func Context(ctx context.Context) wargcore.ParseOpt {
	return func(poh *wargcore.ParseOpts) {
		poh.Context = ctx
	}
}

func Args(args []string) wargcore.ParseOpt {
	return func(poh *wargcore.ParseOpts) {
		poh.Args = args
	}
}

func LookupEnv(lookup wargcore.LookupEnv) wargcore.ParseOpt {
	return func(poh *wargcore.ParseOpts) {
		poh.LookupEnv = lookup
	}
}

func Stderr(stderr *os.File) wargcore.ParseOpt {
	return func(poh *wargcore.ParseOpts) {
		poh.Stderr = stderr
	}
}

func Stdout(stdout *os.File) wargcore.ParseOpt {
	return func(poh *wargcore.ParseOpts) {
		poh.Stdout = stdout
	}
}
