package main

import (
	"fmt"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/value/scalar"
	"go.bbkane.com/warg/value/slice"
)

func app() *warg.App {
	app := warg.New(
		"flaggroup",
		"v1.0.0",
		warg.NewSection(
			"A sample app to demonstrate flag groups",
			warg.NewSubCmd(
				"deploy",
				"Deploy the application.",
				deploy,
				warg.NewCmdFlag(
					"--name",
					"Deployment name.",
					scalar.String(),
					warg.Required(),
				),
				warg.NewCmdFlag(
					"--verbose",
					"Enable verbose logging.",
					scalar.Bool(scalar.Default(false)),
				),
				warg.NewCmdFlag(
					"--db-host",
					"Database host.",
					scalar.String(scalar.Default("localhost")),
					warg.FlagGroup("Database"),
				),
				warg.NewCmdFlag(
					"--db-port",
					"Database port.",
					scalar.Int(scalar.Default(5432)),
					warg.FlagGroup("Database"),
				),
				warg.NewCmdFlag(
					"--db-name",
					"Database name.",
					scalar.String(),
					warg.Required(),
					warg.FlagGroup("Database"),
				),
				warg.NewCmdFlag(
					"--redis-url",
					"Redis connection URL.",
					scalar.String(scalar.Default("redis://localhost:6379")),
					warg.FlagGroup("Cache"),
				),
				warg.NewCmdFlag(
					"--cache-ttl",
					"Cache TTL in seconds.",
					scalar.Int(scalar.Default(300)),
					warg.FlagGroup("Cache"),
				),
				warg.NewCmdFlag(
					"--replicas",
					"Number of replicas.",
					scalar.Int(scalar.Default(1)),
					warg.FlagGroup("Scaling"),
				),
				warg.NewCmdFlag(
					"--regions",
					"Regions to deploy to.",
					slice.String(),
					warg.FlagGroup("Scaling"),
				),
			),
		),
	)
	return &app
}

func deploy(ctx warg.CmdContext) error {
	name := ctx.Flags["--name"].(string)
	dbHost := ctx.Flags["--db-host"].(string)
	dbPort := ctx.Flags["--db-port"].(int)
	dbName := ctx.Flags["--db-name"].(string)

	fmt.Fprintf(ctx.Stdout, "Deploying %s\n", name)
	fmt.Fprintf(ctx.Stdout, "Database: %s:%d/%s\n", dbHost, dbPort, dbName)

	if redisURL, ok := ctx.Flags["--redis-url"]; ok {
		fmt.Fprintf(ctx.Stdout, "Redis: %s\n", redisURL)
	}

	if regions, ok := ctx.Flags["--regions"]; ok {
		fmt.Fprintf(ctx.Stdout, "Regions: %v\n", regions)
	}

	return nil
}

func main() {
	app := app()
	app.MustRun()
}
