package testapp

import (
	"go.bbkane.com/warg"
	"go.bbkane.com/warg/completion"
	"go.bbkane.com/warg/value/scalar"
)

func BuildApp() *warg.App {
	app := warg.New(
		"testappcmd",
		"v1.0.0",
		warg.NewSection(
			"root section help",
			warg.NewChildCmd(
				"command1",
				"command1 help",
				warg.DoNothing,
				warg.NewChildFlag(
					"--flag1",
					"flag1 help",
					scalar.String(
						scalar.Choices("alpha", "beta", "gamma"),
					),
				),
			),
			warg.NewChildCmd(
				"manual",
				"commands with flags using all completion types for manual testing",
				warg.DoNothing,
				warg.NewChildFlag(
					"--dirs",
					"dirs completion",
					scalar.Path(),
					warg.FlagCompletions(warg.CompletionsDirectories),
				),
				warg.NewChildFlag(
					"--dirs-files",
					"dirs/files completion",
					scalar.Path(),
					warg.FlagCompletions(func(ctx warg.Context) (*completion.Candidates, error) {
						return &completion.Candidates{
							Type:   completion.Type_DirectoriesFiles,
							Values: nil,
						}, nil
					}),
				),
				warg.NewChildFlag(
					"--none",
					"no completion",
					scalar.String(),
					warg.FlagCompletions(warg.CompletionsNone),
				),
				warg.NewChildFlag(
					"--values",
					"values completion",
					scalar.String(),
					warg.FlagCompletions(warg.CompletionsValues([]string{"alpha", "beta"})),
				),
				warg.NewChildFlag(
					"--values-descriptions",
					"values completion with descriptions",
					scalar.String(),
					warg.FlagCompletions(warg.CompletionsValuesDescriptions([]completion.Candidate{
						{Name: "gamma", Description: "gamma description"},
						{Name: "delta", Description: "delta description"},
					})),
				),
			),
			warg.NewChildSection(
				"section1",
				"section1 help",
				warg.NewChildCmd(
					"command2",
					"command2 help",
					warg.DoNothing,
					warg.NewChildFlag(
						"--bool",
						"bool completion is special cased to return true/false",
						scalar.Bool(),
					),
					warg.NewChildFlag(
						"--flag2",
						"flag2 help",
						scalar.String(),
						warg.FlagCompletions(func(ctx warg.Context) (*completion.Candidates, error) {
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
