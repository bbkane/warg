package cli

import (
	"context"
	"os"
)

// ParseResult holds the result of parsing the command line.
type ParseResult struct {
	Context Context
	// Action holds the passed command's action to execute.
	Action Action
}

type ParseOptHolder struct {
	Args []string

	// Context for unstructured data. Useful for setting up mocks for tests (i.e., pass in in memory database and use it if it's here in the context)
	Context context.Context

	LookupFunc LookupFunc

	// Stderr will be passed to command.Context for user commands to print to.
	// This file is never closed by warg, so if setting to something other than stderr/stdout,
	// remember to close the file after running the command.
	// Useful for saving output for tests. Defaults to os.Stderr if not passed
	Stderr *os.File

	// Stdout will be passed to command.Context for user commands to print to.
	// This file is never closed by warg, so if setting to something other than stderr/stdout,
	// remember to close the file after running the command.
	// Useful for saving output for tests. Defaults to os.Stdout if not passed
	Stdout *os.File
}

type ParseOpt func(*ParseOptHolder)

func AddContext(ctx context.Context) ParseOpt {
	return func(poh *ParseOptHolder) {
		poh.Context = ctx
	}
}

func OverrideArgs(args []string) ParseOpt {
	return func(poh *ParseOptHolder) {
		poh.Args = args
	}
}

func OverrideLookupFunc(lookup LookupFunc) ParseOpt {
	return func(poh *ParseOptHolder) {
		poh.LookupFunc = lookup
	}
}

func OverrideStderr(stderr *os.File) ParseOpt {
	return func(poh *ParseOptHolder) {
		poh.Stderr = stderr
	}
}

func OverrideStdout(stdout *os.File) ParseOpt {
	return func(poh *ParseOptHolder) {
		poh.Stdout = stdout
	}
}

func NewParseOptHolder(opts ...ParseOpt) ParseOptHolder {
	parseOptHolder := ParseOptHolder{
		Context:    nil,
		Args:       nil,
		LookupFunc: nil,
		Stderr:     nil,
		Stdout:     nil,
	}

	for _, opt := range opts {
		opt(&parseOptHolder)
	}

	if parseOptHolder.Args == nil {
		OverrideArgs(os.Args)(&parseOptHolder)
	}

	if parseOptHolder.Context == nil {
		AddContext(context.Background())(&parseOptHolder)
	}

	if parseOptHolder.LookupFunc == nil {
		OverrideLookupFunc(os.LookupEnv)(&parseOptHolder)
	}

	if parseOptHolder.Stderr == nil {
		OverrideStderr(os.Stderr)(&parseOptHolder)
	}

	if parseOptHolder.Stdout == nil {
		OverrideStdout(os.Stdout)(&parseOptHolder)
	}

	return parseOptHolder
}

// Parse parses the args, but does not execute anything.
func (app *App) Parse(opts ...ParseOpt) (*ParseResult, error) {

	parseOptHolder := NewParseOptHolder(opts...)
	return app.parseWithOptHolder2(parseOptHolder)
}
