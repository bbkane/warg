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
		"testappcmd",
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
				"manual",
				"commands with flags using all completion types for manual testing",
				command.DoNothing,
				command.NewFlag(
					"--dirs",
					"dirs completion",
					scalar.Path(),
					flag.CompletionCandidate(func(ctx cli.Context) (*completion.Candidates, error) {
						return &completion.Candidates{
							Type:   completion.Type_Directories,
							Values: nil,
						}, nil
					}),
				),
				command.NewFlag(
					"--dirs-files",
					"dirs/files completion",
					scalar.Path(),
					flag.CompletionCandidate(func(ctx cli.Context) (*completion.Candidates, error) {
						return &completion.Candidates{
							Type:   completion.Type_DirectoriesFiles,
							Values: nil,
						}, nil
					}),
				),
				command.NewFlag(
					"--none",
					"no completion",
					scalar.String(),
					flag.CompletionCandidate(func(ctx cli.Context) (*completion.Candidates, error) {
						return &completion.Candidates{
							Type:   completion.Type_None,
							Values: nil,
						}, nil
					}),
				),
				command.NewFlag(
					"--values",
					"values completion",
					scalar.String(),
					flag.CompletionCandidate(func(ctx cli.Context) (*completion.Candidates, error) {
						return &completion.Candidates{
							Type: completion.Type_Values,
							Values: []completion.Candidate{
								{Name: "alpha", Description: "THIS SHOULDN'T SHOW"},
								{Name: "beta", Description: "THIS SHOULDN'T SHOW"},
							},
						}, nil
					}),
				),
				command.NewFlag(
					"--values-descriptions",
					"values completion with descriptions",
					scalar.String(),
					flag.CompletionCandidate(func(ctx cli.Context) (*completion.Candidates, error) {
						return &completion.Candidates{
							Type: completion.Type_ValuesDescriptions,
							Values: []completion.Candidate{
								{Name: "gamma", Description: "THIS SHOULD SHOW"},
								{Name: "delta", Description: "THIS SHOULD SHOW"},
							},
						}, nil
					}),
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
									Type: completion.Type_ValuesDescriptions,
									Values: []completion.Candidate{
										{
											Name:        "nondefault",
											Description: "nondefault completion",
										},
									},
								}, nil
							}
							return &completion.Candidates{
								Type: completion.Type_ValuesDescriptions,
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
