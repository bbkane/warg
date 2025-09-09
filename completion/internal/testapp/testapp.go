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
			warg.NewSubCmd(
				"command1",
				"command1 help",
				warg.Unimplemented(),
				warg.NewCmdFlag(
					"--flag1",
					"flag1 help",
					scalar.String(
						scalar.Choices("alpha", "beta", "gamma"),
					),
				),
			),
			warg.NewSubCmd(
				"manual",
				"commands with flags using all completion types for manual testing",
				warg.Unimplemented(),
				warg.NewCmdFlag(
					"--dirs",
					"dirs completion",
					scalar.Path(),
					warg.FlagCompletions(warg.CompletionsDirectories()),
				),
				warg.NewCmdFlag(
					"--dirs-files",
					"dirs/files completion",
					scalar.Path(),
					warg.FlagCompletions(func(ctx warg.CmdContext) (*completion.Candidates, error) {
						return &completion.Candidates{
							Type:   completion.Type_DirectoriesFiles,
							Values: nil,
						}, nil
					}),
				),
				warg.NewCmdFlag(
					"--none",
					"no completion",
					scalar.String(),
					warg.FlagCompletions(warg.CompletionsNone()),
				),
				warg.NewCmdFlag(
					"--values",
					"values completion",
					scalar.String(),
					warg.FlagCompletions(warg.CompletionsValues([]string{"alpha", "beta"})),
				),
				warg.NewCmdFlag(
					"--values-descriptions",
					"values completion with descriptions",
					scalar.String(),
					warg.FlagCompletions(warg.CompletionsValuesDescriptions([]completion.Candidate{
						{Name: "gamma", Description: "gamma description"},
						{Name: "delta", Description: "delta description"},
					})),
				),
			),
			warg.NewSubSection(
				"section1",
				"section1 help",
				warg.NewSubCmd(
					"command2",
					"command2 help",
					warg.Unimplemented(),
					warg.NewCmdFlag(
						"--bool",
						"bool completion is special cased to return true/false",
						scalar.Bool(),
					),
					warg.NewCmdFlag(
						"--flag2",
						"flag2 help",
						scalar.String(),
						warg.FlagCompletions(func(ctx warg.CmdContext) (*completion.Candidates, error) {
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
		warg.SkipVersionCmd(),
		warg.HelpFlag(
			warg.CmdMap{
				"default":  warg.NewCmd("", warg.Unimplemented()),
				"detailed": warg.NewCmd("", warg.Unimplemented()),
			},
			warg.DefaultHelpFlagMap("default", []string{"default", "detailed"}),
		),
	)
	return &app
}
