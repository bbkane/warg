package testapp

import (
	"go.bbkane.com/warg"
	"go.bbkane.com/warg/completion"
	"go.bbkane.com/warg/value/scalar"
	"go.bbkane.com/warg/wargcore"
)

func BuildApp() *wargcore.App {
	app := warg.New(
		"testappcmd",
		"v1.0.0",
		wargcore.NewSection(
			"root section help",
			wargcore.NewChildCmd(
				"command1",
				"command1 help",
				wargcore.DoNothing,
				wargcore.NewChildFlag(
					"--flag1",
					"flag1 help",
					scalar.String(
						scalar.Choices("alpha", "beta", "gamma"),
					),
				),
			),
			wargcore.NewChildCmd(
				"manual",
				"commands with flags using all completion types for manual testing",
				wargcore.DoNothing,
				wargcore.NewChildFlag(
					"--dirs",
					"dirs completion",
					scalar.Path(),
					wargcore.FlagCompletions(warg.CompletionsDirectories),
				),
				wargcore.NewChildFlag(
					"--dirs-files",
					"dirs/files completion",
					scalar.Path(),
					wargcore.FlagCompletions(func(ctx wargcore.Context) (*completion.Candidates, error) {
						return &completion.Candidates{
							Type:   completion.Type_DirectoriesFiles,
							Values: nil,
						}, nil
					}),
				),
				wargcore.NewChildFlag(
					"--none",
					"no completion",
					scalar.String(),
					wargcore.FlagCompletions(warg.CompletionsNone),
				),
				wargcore.NewChildFlag(
					"--values",
					"values completion",
					scalar.String(),
					wargcore.FlagCompletions(warg.CompletionsValues([]string{"alpha", "beta"})),
				),
				wargcore.NewChildFlag(
					"--values-descriptions",
					"values completion with descriptions",
					scalar.String(),
					wargcore.FlagCompletions(warg.CompletionsValuesDescriptions([]completion.Candidate{
						{Name: "gamma", Description: "gamma description"},
						{Name: "delta", Description: "delta description"},
					})),
				),
			),
			wargcore.NewChildSection(
				"section1",
				"section1 help",
				wargcore.NewChildCmd(
					"command2",
					"command2 help",
					wargcore.DoNothing,
					wargcore.NewChildFlag(
						"--bool",
						"bool completion is special cased to return true/false",
						scalar.Bool(),
					),
					wargcore.NewChildFlag(
						"--flag2",
						"flag2 help",
						scalar.String(),
						wargcore.FlagCompletions(func(ctx wargcore.Context) (*completion.Candidates, error) {
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
