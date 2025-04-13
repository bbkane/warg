package parseopt

import (
	"context"
	"os"

	"go.bbkane.com/warg/cli"
)

func Context(ctx context.Context) cli.ParseOpt {
	return func(poh *cli.ParseOpts) {
		poh.Context = ctx
	}
}

func Args(args []string) cli.ParseOpt {
	return func(poh *cli.ParseOpts) {
		poh.Args = args
	}
}

func LookupEnv(lookup cli.LookupEnv) cli.ParseOpt {
	return func(poh *cli.ParseOpts) {
		poh.LookupEnv = lookup
	}
}

func Stderr(stderr *os.File) cli.ParseOpt {
	return func(poh *cli.ParseOpts) {
		poh.Stderr = stderr
	}
}

func Stdout(stdout *os.File) cli.ParseOpt {
	return func(poh *cli.ParseOpts) {
		poh.Stdout = stdout
	}
}
