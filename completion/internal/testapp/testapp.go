package testapp

import (
	"go.bbkane.com/warg"
	"go.bbkane.com/warg/cli"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/completion"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/section"
	"go.bbkane.com/warg/value/scalar"
)

func BuildApp() *cli.App {
	app := warg.New(
		"newAppName",
		"v1.0.0",
		section.New(
			"root section help",
			section.NewCommand(
				"command1",
				"command1 help",
				command.DoNothing,
				command.NewFlag(
					"--flag1",
					"flag1 help",
					scalar.String(
						scalar.Choices("alpha", "beta", "gamma"),
					),
				),
			),
			// TODO: actually test this
			section.NewCommand(
				"filecompletioncommands",
				"filecompletioncommands help",
				command.DoNothing,
				command.NewFlag(
					"--fileflag",
					"fileflag help",
					scalar.Path(),
					flag.CompletionCandidate(func(ctx cli.Context) (*completion.Candidates, error) {
						return &completion.Candidates{
							Type:   completion.Type_DirectoriesFiles,
							Values: nil,
						}, nil
					},
					),
				),
			),
			section.NewSection(
				"section1",
				"section1 help",
				section.NewCommand(
					"command2",
					"command2 help",
					command.DoNothing,
					command.NewFlag(
						"--flag2",
						"flag2 help",
						scalar.String(),
						flag.CompletionCandidate(func(ctx cli.Context) (*completion.Candidates, error) {
							if ctx.Flags["--globalFlag"].(string) == "nondefault" {
								return &completion.Candidates{
									Type: completion.Type_ValueDescription,
									Values: []completion.Candidate{
										{
											Name:        "nondefault",
											Description: "nondefault completion",
										},
									},
								}, nil
							}
							return &completion.Candidates{
								Type: completion.Type_ValueDescription,
								Values: []completion.Candidate{
									{
										Name:        "default",
										Description: "default completion",
									},
								},
							}, nil
						}),
					),
				),
			),
		),
		warg.NewGlobalFlag(
			"--globalFlag",
			"globalFlag help",
			scalar.String(
				scalar.Default("default"),
			),
		),
	)
	return &app
}
