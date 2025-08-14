package testapp

import (
	"go.bbkane.com/warg"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/completion"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/section"
	"go.bbkane.com/warg/value/scalar"
	"go.bbkane.com/warg/wargcore"
)

func BuildApp() *wargcore.App {
	app := warg.New(
		"testappcmd",
		"v1.0.0",
		section.NewSection(
			"root section help",
			section.NewChildCmd(
				"command1",
				"command1 help",
				command.DoNothing,
				command.NewChildFlag(
					"--flag1",
					"flag1 help",
					scalar.String(
						scalar.Choices("alpha", "beta", "gamma"),
					),
				),
			),
			section.NewChildCmd(
				"manual",
				"commands with flags using all completion types for manual testing",
				command.DoNothing,
				command.NewChildFlag(
					"--dirs",
					"dirs completion",
					scalar.Path(),
					flag.FlagCompletions(warg.CompletionsDirectories),
				),
				command.NewChildFlag(
					"--dirs-files",
					"dirs/files completion",
					scalar.Path(),
					flag.FlagCompletions(func(ctx wargcore.Context) (*completion.Candidates, error) {
						return &completion.Candidates{
							Type:   completion.Type_DirectoriesFiles,
							Values: nil,
						}, nil
					}),
				),
				command.NewChildFlag(
					"--none",
					"no completion",
					scalar.String(),
					flag.FlagCompletions(warg.CompletionsNone),
				),
				command.NewChildFlag(
					"--values",
					"values completion",
					scalar.String(),
					flag.FlagCompletions(warg.CompletionsValues([]string{"alpha", "beta"})),
				),
				command.NewChildFlag(
					"--values-descriptions",
					"values completion with descriptions",
					scalar.String(),
					flag.FlagCompletions(warg.CompletionsValuesDescriptions([]completion.Candidate{
						{Name: "gamma", Description: "gamma description"},
						{Name: "delta", Description: "delta description"},
					})),
				),
			),
			section.NewChildSection(
				"section1",
				"section1 help",
				section.NewChildCmd(
					"command2",
					"command2 help",
					command.DoNothing,
					command.NewChildFlag(
						"--bool",
						"bool completion is special cased to return true/false",
						scalar.Bool(),
					),
					command.NewChildFlag(
						"--flag2",
						"flag2 help",
						scalar.String(),
						flag.FlagCompletions(func(ctx wargcore.Context) (*completion.Candidates, error) {
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
		warg.SkipGlobalColorFlag(),
		warg.SkipVersionCommand(),
	)
	return &app
}
