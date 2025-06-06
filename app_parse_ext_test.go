package warg_test

// external tests - import warg like it's an external package

import (
	"context"
	"path/filepath"
	"slices"
	"testing"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/config"
	"go.bbkane.com/warg/config/jsonreader"
	"go.bbkane.com/warg/config/yamlreader"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/parseopt"
	"go.bbkane.com/warg/path"
	"go.bbkane.com/warg/section"
	"go.bbkane.com/warg/value/dict"
	"go.bbkane.com/warg/value/scalar"
	"go.bbkane.com/warg/value/slice"
	"go.bbkane.com/warg/wargcore"

	"github.com/stretchr/testify/require"
)

// NOTE: this is is a bit of a hack to mock out a configreader
// NOTE: see https://karthikkaranth.me/blog/functions-implementing-interfaces-in-go/
// for how to use ConfigReaderFunc in tests
type ConfigReaderFunc func(path string) (*config.SearchResult, error)

func (f ConfigReaderFunc) Search(path string) (*config.SearchResult, error) {
	return f(path)
}

// testDataFilePath generates a path to a file needed for a test
func testDataFilePath(testName string, subTestName string, fileName string) path.Path {
	return path.New(filepath.Join("testdata", testName, subTestName, fileName))
}

func TestApp_Parse(t *testing.T) {

	tests := []struct {
		name                     string
		app                      wargcore.App
		args                     []string
		lookup                   wargcore.LookupEnv
		expectedPassedPath       []string
		expectedPassedFlagValues wargcore.PassedFlags
		expectedErr              bool
	}{

		{
			name: "envvar",
			app: warg.New(
				"newAppName", "v1.0.0",
				section.New(
					"help for test",
					section.NewCommand(
						"test",
						"blah",
						command.DoNothing,
						command.NewFlag(
							"--flag",
							"help for --flag",
							scalar.String(),
							flag.EnvVars("notthere", "there", "alsothere"),
						),
					),
				),
				warg.SkipValidation(),
			),
			args: []string{t.Name(), "test"},
			lookup: wargcore.LookupMap(
				map[string]string{
					"there":     "there",
					"alsothere": "alsothere",
				},
			),
			expectedPassedPath: []string{"test"},
			expectedPassedFlagValues: wargcore.PassedFlags{
				"--flag": "there",
				"--help": "default",
			},
			expectedErr: false,
		},
		{
			name: "addCommandFlags",
			app: func() wargcore.App {
				fm := wargcore.FlagMap{
					"--flag1": flag.New("--flag1 value", scalar.String()),
					"--flag2": flag.New("--flag1 value", scalar.String()),
				}
				app := warg.New(
					"newAppName", "v1.0.0",
					section.New(
						"help for section",
						section.NewCommand(
							"test",
							"help for test",
							command.DoNothing,
							command.FlagMap(fm),
						),
					),
					warg.SkipValidation(),
				)
				return app
			}(),

			args:                     []string{t.Name(), "test", "--flag1", "val1"},
			lookup:                   wargcore.LookupMap(nil),
			expectedPassedPath:       []string{"test"},
			expectedPassedFlagValues: wargcore.PassedFlags{"--flag1": "val1", "--help": "default"},
			expectedErr:              false,
		},

		// Note: Will need to update this if https://github.com/bbkane/warg/issues/36 gets implemented
		{
			name: "invalidFlagsErrorEvenForHelp",
			app: warg.New(
				"newAppName", "v1.0.0",
				section.New(
					string("A virtual assistant"),
					section.NewCommand(
						"present",
						"Formally present a guest (guests are never introduced, always presented).",
						command.DoNothing,
						command.NewFlag(
							"--name",
							"Guest to address.",
							scalar.String(scalar.Choices("bob")),
							flag.Alias("-n"),
							flag.EnvVars("BUTLER_PRESENT_NAME", "USER"),
							flag.Required(),
						),
					),
				),
				warg.SkipValidation(),
			),

			args:                     []string{"app", "present", "-h"},
			lookup:                   wargcore.LookupMap(map[string]string{"USER": "bbkane"}),
			expectedPassedPath:       []string{"present"},
			expectedPassedFlagValues: wargcore.PassedFlags{"--help": "default"},
			expectedErr:              true,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			err := tt.app.Validate()
			require.Nil(t, err)

			actualPR, actualErr := tt.app.Parse(parseopt.Args(tt.args), parseopt.LookupEnv(tt.lookup))

			if tt.expectedErr {
				require.NotNil(t, actualErr)
				return
			} else {
				require.Nil(t, actualErr)
			}
			actualPath := slices.Clone(actualPR.Context.ParseState.SectionPath)
			if actualPR.Context.ParseState.CurrentCommandName != "" {
				actualPath = append(actualPath, actualPR.Context.ParseState.CurrentCommandName)
			}

			require.Equal(t, tt.expectedPassedPath, actualPath)
			require.Equal(t, tt.expectedPassedFlagValues, actualPR.Context.Flags)
		})
	}
}

// TestApp_Parse_rootSection tests the Parse method, when only a root section is needed (i.e., no special app opts, names, versions, or LookupMaps
func TestApp_Parse_rootSection(t *testing.T) {
	tests := []struct {
		name                     string
		rootSection              wargcore.Section
		args                     []string
		expectedPassedPath       []string
		expectedPassedFlagValues wargcore.PassedFlags
		expectedErr              bool
	}{
		{
			name: "fromMain",
			rootSection: section.New(
				"help for test",
				section.NewSection(
					"cat1",
					"help for cat1",
					section.NewCommand(
						"com1",
						"help for com1",
						command.DoNothing,
						command.NewFlag(
							"--com1f1",
							"flag help",
							scalar.Int(
								scalar.Default(10),
							),
						),
					),
				),
			),
			args:                     []string{"app", "cat1", "com1", "--com1f1", "1"},
			expectedPassedPath:       []string{"cat1", "com1"},
			expectedPassedFlagValues: wargcore.PassedFlags{"--com1f1": int(1), "--help": "default"},
			expectedErr:              false,
		},
		{
			name: "noSection",
			rootSection: section.New(
				"help for test",
				section.NewCommand("com", "command for validation", command.DoNothing),
			),

			args:                     []string{"app"},
			expectedPassedPath:       nil,
			expectedPassedFlagValues: map[string]interface{}{"--help": "default"},
			expectedErr:              false,
		},
		{
			name: "flagDefault",
			rootSection: section.New(
				"help for test",
				section.NewCommand(
					"com",
					"com help",
					command.DoNothing,
					command.NewFlag(
						"--flag",
						"flag help",
						scalar.String(
							scalar.Default("hi"),
						),
					),
				),
			),
			args:                     []string{"test", "com"},
			expectedPassedPath:       []string{"com"},
			expectedPassedFlagValues: wargcore.PassedFlags{"--flag": "hi", "--help": "default"},
			expectedErr:              false,
		},
		{
			name: "extraFlag",
			rootSection: section.New(
				"help for test",
				section.NewCommand(
					"com",
					"com help",
					command.DoNothing,
					command.NewFlag(
						"--flag",
						"flag help",
						scalar.String(
							scalar.Default("hi"),
						),
					),
				),
			),
			args:                     []string{"test", "com", "--unexpected", "value"},
			expectedPassedPath:       nil,
			expectedPassedFlagValues: nil,
			expectedErr:              true,
		},
		{
			name: "requiredFlag",
			rootSection: section.New(
				"help for test",
				section.NewCommand(
					"test",
					"blah",
					command.DoNothing,
					command.NewFlag(
						"--flag",
						"help for --flag",
						scalar.String(),
						flag.Required(),
					),
				),
			),
			args:                     []string{t.Name(), "test"},
			expectedPassedPath:       []string{"test"},
			expectedPassedFlagValues: wargcore.PassedFlags{},
			expectedErr:              true,
		},
		{
			name: "flagAlias",
			rootSection: section.New(
				"help for section",
				section.NewCommand(
					"test",
					"help for test",
					command.DoNothing,
					command.NewFlag(
						"--flag",
						"help for --flag",
						scalar.String(),
						flag.Alias("-f"),
					),
				),
			),
			args:                     []string{t.Name(), "test", "-f", "val"},
			expectedPassedPath:       []string{"test"},
			expectedPassedFlagValues: wargcore.PassedFlags{"--flag": "val", "--help": "default"},
			expectedErr:              false,
		},
		{
			name: "flagAliasWithList",
			rootSection: section.New(
				"help for section",
				section.NewCommand(
					"test",
					"help for test",
					command.DoNothing,
					command.NewFlag(
						"--flag",
						"help for --flag",
						slice.String(),
						flag.Alias("-f"),
					),
				),
			),
			args:                     []string{t.Name(), "test", "-f", "1", "--flag", "2", "-f", "3", "--flag", "4"},
			expectedPassedPath:       []string{"test"},
			expectedPassedFlagValues: wargcore.PassedFlags{"--flag": []string{"1", "2", "3", "4"}, "--help": "default"},
			expectedErr:              false,
		},
		{
			name: "badHelp",
			rootSection: section.New(
				"help for section",
				section.NewCommand(
					"test",
					"help for test",
					command.DoNothing,
				),
			),
			args:                     []string{t.Name(), "test", "-h", "badhelpval"},
			expectedPassedPath:       nil,
			expectedPassedFlagValues: nil,
			expectedErr:              true,
		},
		{
			name: "dictUpdate",
			rootSection: section.New(
				"help for test",
				section.NewCommand(
					"com1",
					"help for com1",
					command.DoNothing,
					command.NewFlag(
						string("--flag"),
						"flag help",
						dict.Bool(),
					),
				),
			),

			args:                     []string{"app", "com1", "--flag", "true=true", "--flag", "false=false"},
			expectedPassedPath:       []string{"com1"},
			expectedPassedFlagValues: wargcore.PassedFlags{"--flag": map[string]bool{"true": true, "false": false}, "--help": "default"},
			expectedErr:              false,
		},
		{
			name: "passAbsentSection",
			rootSection: section.New(
				"help for test",
				section.NewCommand(
					"com",
					"help for com",
					command.DoNothing,
				),
			),

			args:                     []string{"app", "badSectionName"},
			expectedPassedPath:       []string{"com1"},
			expectedPassedFlagValues: wargcore.PassedFlags{"--help": "default"},
			expectedErr:              true,
		},
		{
			name: "scalarFlagPassedTwice",
			rootSection: section.New(
				"help for test",
				section.NewCommand(
					"com",
					"help for com1",
					command.DoNothing,
					command.NewFlag(
						"--flag",
						"flag help",
						scalar.Int(),
					),
				),
			),

			args:                     []string{"app", "com", "--flag", "1", "--flag", "2"},
			expectedPassedPath:       []string{"com"},
			expectedPassedFlagValues: wargcore.PassedFlags{"--flag": int(1), "--help": "default"},
			expectedErr:              true,
		},
		{
			name: "passedFlagBeforeCommand",
			rootSection: section.New(
				"help for test",
				section.NewCommand(
					"com",
					"help for com",
					command.DoNothing,
					command.NewFlag(
						"--flag",
						"flag help",
						scalar.Int(),
					),
				),
			),

			args:                     []string{"app", "--flag", "1", "com"},
			expectedPassedPath:       []string{"com"},
			expectedPassedFlagValues: wargcore.PassedFlags{"--flag": int(1), "--help": "default"},
			expectedErr:              true,
		},
		{
			name: "existingSectionsExistingCommands",
			rootSection: section.New(
				"help for test",
				section.SectionMap(
					wargcore.SectionMap{
						"section": section.New(
							"help for section",
							section.CommandMap(
								wargcore.CommandMap{
									"command": command.New(
										"help for command",
										command.DoNothing,
										command.NewFlag(
											"--flag",
											"flag help",
											scalar.Int(),
										),
									),
								},
							),
						),
					},
				),
			),
			args:               []string{"app", "section", "command", "--flag", "1"},
			expectedPassedPath: []string{"section", "command"},
			expectedPassedFlagValues: wargcore.PassedFlags{
				"--flag": 1,
				"--help": "default",
			},
			expectedErr: false,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			app := warg.New(
				"newAppName", "v1.0.0",
				tt.rootSection,
				warg.SkipValidation(),
			)

			err := app.Validate()
			require.Nil(t, err)

			actualPR, actualErr := app.Parse(parseopt.Args(tt.args), parseopt.LookupEnv(wargcore.LookupMap(nil)))

			if tt.expectedErr {
				require.Error(t, actualErr)
				return
			} else {
				require.NoError(t, actualErr)
			}
			actualPath := slices.Clone(actualPR.Context.ParseState.SectionPath)
			if actualPR.Context.ParseState.CurrentCommandName != "" {
				actualPath = append(actualPath, actualPR.Context.ParseState.CurrentCommandName)
			}
			require.Equal(t, tt.expectedPassedPath, actualPath)
			require.Equal(t, tt.expectedPassedFlagValues, actualPR.Context.Flags)
		})
	}
}

func TestApp_Parse_unsetSetinel(t *testing.T) {
	tests := []struct {
		name                     string
		flagDef                  command.CommandOpt
		args                     []string
		expectedPassedPath       []string
		expectedPassedFlagValues wargcore.PassedFlags
		expectedErr              bool
	}{
		{
			name: "unsetSentinelScalarSuccess",
			flagDef: command.NewFlag(
				"--flag",
				"help for --flag",
				scalar.String(scalar.Default("default")),
				flag.UnsetSentinel("UNSET"),
			),
			args:               []string{t.Name(), "test", "--flag", "UNSET"},
			expectedPassedPath: []string{"test"},
			expectedPassedFlagValues: wargcore.PassedFlags{
				"--help": "default",
			},
			expectedErr: false,
		},
		{
			name: "unsetSentinelScalarUpdate",
			flagDef: command.NewFlag(
				"--flag",
				"help for --flag",
				scalar.String(scalar.Default("default")),
				flag.UnsetSentinel("UNSET"),
			),
			args:                     []string{t.Name(), "test", "--flag", "UNSET", "--flag", "setAfter"},
			expectedPassedPath:       []string{"test"},
			expectedPassedFlagValues: wargcore.PassedFlags{"--flag": "setAfter", "--help": "default"},
			expectedErr:              false,
		},
		{
			name: "unsetSentinelSlice",
			flagDef: command.NewFlag(
				"--flag",
				"help for --flag",
				slice.String(slice.Default([]string{"default"})),
				flag.UnsetSentinel("UNSET"),
			),
			args:               []string{t.Name(), "test", "--flag", "a", "--flag", "UNSET", "--flag", "b", "--flag", "c"},
			expectedPassedPath: []string{"test"},
			expectedPassedFlagValues: wargcore.PassedFlags{
				"--help": "default",
				"--flag": []string{"b", "c"},
			},
			expectedErr: false,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			app := warg.New(
				"newAppName", "v1.0.0",
				section.New(
					"help for test",
					section.NewCommand(
						"test",
						"help for test",
						command.DoNothing,
						tt.flagDef,
					),
				),
				warg.SkipValidation(),
			)

			err := app.Validate()
			require.Nil(t, err)

			actualPR, actualErr := app.Parse(parseopt.Args(tt.args), parseopt.LookupEnv(wargcore.LookupMap(nil)))

			if tt.expectedErr {
				require.Error(t, actualErr)
				return
			} else {
				require.NoError(t, actualErr)
			}
			actualPath := slices.Clone(actualPR.Context.ParseState.SectionPath)
			if actualPR.Context.ParseState.CurrentCommandName != "" {
				actualPath = append(actualPath, actualPR.Context.ParseState.CurrentCommandName)
			}
			require.Equal(t, tt.expectedPassedPath, actualPath)
			require.Equal(t, tt.expectedPassedFlagValues, actualPR.Context.Flags)
		})
	}
}

func TestApp_Parse_config(t *testing.T) {
	tests := []struct {
		name                     string
		app                      wargcore.App
		args                     []string
		lookup                   wargcore.LookupEnv
		expectedPassedPath       []string
		expectedPassedFlagValues wargcore.PassedFlags
		expectedErr              bool
	}{
		{
			name: "configFlag",
			app: warg.New(
				"newAppName", "v1.0.0",
				section.New(
					"help for test",
					section.NewCommand(
						"print",
						"print key value",
						command.DoNothing,
						command.NewFlag(
							"--key",
							"a key",
							scalar.String(
								scalar.Default("defaultkeyval"),
							),
							flag.ConfigPath("key"),
						),
					),
				),
				warg.ConfigFlag(
					func(_ string) (config.Reader, error) {
						var cr ConfigReaderFunc = func(path string) (*config.SearchResult, error) {
							if path == "key" {

								return &config.SearchResult{
									IFace:        "mapkeyval",
									IsAggregated: false,
								}, nil
							}

							return nil, nil
						}

						return cr, nil
					},
					wargcore.FlagMap{
						"--config": flag.New(
							"Path to config file",
							scalar.Path(
								scalar.Default(path.New("defaultconfigval")),
							),
						),
					},
				),
				warg.SkipValidation(),
			),
			args:               []string{"test", "print", "--config", "passedconfigval"},
			lookup:             wargcore.LookupMap(nil),
			expectedPassedPath: []string{"print"},
			expectedPassedFlagValues: wargcore.PassedFlags{
				"--key":    "mapkeyval",
				"--config": path.New("passedconfigval"),
				"--help":   "default",
			},
			expectedErr: false,
		},
		{
			name: "simpleJSONConfig",
			app: warg.New(
				"newAppName", "v1.0.0",
				section.New("help for test",
					section.NewCommand(
						"com",
						"help for com",
						command.DoNothing,
						command.NewFlag(
							"--val",
							"flag help",
							scalar.String(),
							flag.ConfigPath("params.val"),
						),
					),
				),
				warg.ConfigFlag(
					jsonreader.New,
					wargcore.FlagMap{
						"--config": flag.New(
							"path to config",
							scalar.Path(
								scalar.Default(
									testDataFilePath(t.Name(), "simpleJSONConfig", "simple_json_config.json"),
								),
							),
							flag.Alias("-c"),
						),
					},
				),
				warg.SkipValidation(),
			),

			args:               []string{"app", "com"},
			lookup:             wargcore.LookupMap(nil),
			expectedPassedPath: []string{"com"},
			expectedPassedFlagValues: wargcore.PassedFlags{
				"--config": testDataFilePath(t.Name(), "simpleJSONConfig", "simple_json_config.json"),
				"--val":    "hi",
				"--help":   "default",
			},
			expectedErr: false,
		},
		{
			// JSON needs some help to support number decoding
			name: "numJSONConfig",
			app: warg.New(
				"newAppName", "v1.0.0",
				section.New("help for test",
					section.NewCommand(
						"com",
						"help for com",
						command.DoNothing,
						command.NewFlag(
							"--intval",
							"flag help",
							scalar.Int(),
							flag.ConfigPath("params.intval"),
						),
					),
				),
				warg.ConfigFlag(
					jsonreader.New,
					wargcore.FlagMap{
						"--config": flag.New(
							"path to config",
							scalar.Path(
								scalar.Default(
									testDataFilePath(t.Name(), "numJSONConfig", "config.json"),
								),
							),
							flag.Alias("-c"),
						),
					},
				),
				warg.SkipValidation(),
			),

			args:               []string{"app", "com"},
			lookup:             wargcore.LookupMap(nil),
			expectedPassedPath: []string{"com"},
			expectedPassedFlagValues: wargcore.PassedFlags{
				"--config": testDataFilePath(t.Name(), "numJSONConfig", "config.json"),
				// "--floatval": 1.5,
				"--intval": 1,
				"--help":   "default",
			},
			expectedErr: false,
		},
		{
			name: "configSlice",
			app: warg.New(
				"newAppName", "v1.0.0",
				section.New(
					"help for test",
					section.NewCommand(
						"print",
						"print key value",
						command.DoNothing,
						command.NewFlag(
							"--subreddits",
							"the subreddits",
							slice.String(),
							flag.ConfigPath("subreddits[].name"),
						),
					),
				),
				warg.ConfigFlag(
					jsonreader.New,
					wargcore.FlagMap{
						"--config": flag.New(
							"path to config",
							scalar.Path(
								scalar.Default(
									testDataFilePath(t.Name(), "configSlice", "config_slice.json"),
								),
							),
							flag.Alias("-c"),
						),
					},
				),
				warg.SkipValidation(),
			),
			args:               []string{"test", "print"},
			lookup:             wargcore.LookupMap(nil),
			expectedPassedPath: []string{"print"},
			expectedPassedFlagValues: wargcore.PassedFlags{
				"--subreddits": []string{"earthporn", "wallpapers"},
				"--config":     testDataFilePath(t.Name(), "configSlice", "config_slice.json"),
				"--help":       "default",
			},
			expectedErr: false,
		},
		{
			name: "JSONConfigStringSlice",
			app: warg.New(
				"newAppName", "v1.0.0",
				section.New("help for test",
					section.NewCommand(
						"com",
						"help for com",
						command.DoNothing,
						command.NewFlag(
							"--val",
							"flag help",
							slice.String(),
							flag.ConfigPath("val"),
						),
					),
				),
				warg.ConfigFlag(
					jsonreader.New,
					wargcore.FlagMap{
						"--config": flag.New(
							"path to config",
							scalar.Path(
								scalar.Default(
									testDataFilePath(t.Name(), "JSONConfigStringSlice", "config.json"),
								),
							),
							flag.Alias("-c"),
						),
					},
				),
				warg.SkipValidation(),
			),
			args:               []string{"app", "com"},
			lookup:             wargcore.LookupMap(nil),
			expectedPassedPath: []string{"com"},
			expectedPassedFlagValues: wargcore.PassedFlags{
				"--config": testDataFilePath(t.Name(), "JSONConfigStringSlice", "config.json"),
				"--val":    []string{"from", "config"},
				"--help":   "default",
			},
			expectedErr: false,
		},
		{
			name: "YAMLConfigStringSlice",
			app: warg.New(
				"newAppName", "v1.0.0",
				section.New("help for test",
					section.NewCommand(
						"com",
						"help for com",
						command.DoNothing,
						command.NewFlag(
							"--val",
							"flag help",
							slice.String(),
							flag.ConfigPath("val"),
						),
					),
				),
				warg.ConfigFlag(
					yamlreader.New,
					wargcore.FlagMap{
						"--config": flag.New(
							"path to config",
							scalar.Path(
								scalar.Default(
									testDataFilePath(t.Name(), "YAMLConfigStringSlice", "config.yaml"),
								),
							),
							flag.Alias("-c"),
						),
					},
				),
				warg.SkipValidation(),
			),

			args:               []string{"app", "com"},
			lookup:             wargcore.LookupMap(nil),
			expectedPassedPath: []string{"com"},
			expectedPassedFlagValues: wargcore.PassedFlags{
				"--config": testDataFilePath(t.Name(), "YAMLConfigStringSlice", "config.yaml"),
				"--val":    []string{"from", "config"},
				"--help":   "default",
			},
			expectedErr: false,
		},
		{
			name: "JSONConfigMap",
			app: warg.New(
				"newAppName", "v1.0.0",
				section.New(
					"help for test",
					section.NewCommand(
						"com",
						"help for com",
						command.DoNothing,
						command.NewFlag(
							"--val",
							"flag help",
							dict.Int(),
							flag.ConfigPath("val"),
						),
					),
				),
				warg.ConfigFlag(
					jsonreader.New,
					wargcore.FlagMap{
						"--config": flag.New(
							"path to config",
							scalar.Path(
								scalar.Default(
									testDataFilePath(t.Name(), "JSONConfigMap", "config.json"),
								),
							),
							flag.Alias("-c"),
						),
					},
				),
				warg.SkipValidation(),
			),

			args:               []string{"app", "com"},
			lookup:             wargcore.LookupMap(nil),
			expectedPassedPath: []string{"com"},
			expectedPassedFlagValues: wargcore.PassedFlags{
				"--config": testDataFilePath(t.Name(), "JSONConfigMap", "config.json"),
				"--val":    map[string]int{"a": 1},
				"--help":   "default",
			},
			expectedErr: false,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			err := tt.app.Validate()
			require.Nil(t, err)

			actualPR, actualErr := tt.app.Parse(parseopt.Args(tt.args), parseopt.LookupEnv(tt.lookup))

			if tt.expectedErr {
				require.NotNil(t, actualErr)
				return
			} else {
				require.Nil(t, actualErr)
			}
			actualPath := slices.Clone(actualPR.Context.ParseState.SectionPath)
			if actualPR.Context.ParseState.CurrentCommandName != "" {
				actualPath = append(actualPath, actualPR.Context.ParseState.CurrentCommandName)
			}
			require.Equal(t, tt.expectedPassedPath, actualPath)
			require.Equal(t, tt.expectedPassedFlagValues, actualPR.Context.Flags)
		})
	}
}

// This is the same as TestApp_Parse, but that's too long for a single test
func TestApp_Parse_GlobalFlag(t *testing.T) {
	tests := []struct {
		name                     string
		app                      wargcore.App
		args                     []string
		lookup                   wargcore.LookupEnv
		expectedPassedPath       []string
		expectedPassedFlagValues wargcore.PassedFlags
		expectedErr              bool
	}{
		{
			name: "globalFlag",
			app: warg.New(
				"newAppName", "v1.0.0",
				section.New(
					"help for test",
					section.NewCommand(
						"com",
						"help for com",
						command.DoNothing,
					),
				),
				warg.SkipValidation(),
				warg.NewGlobalFlag(
					"--global",
					"global flag",
					scalar.String(),
				),
			),

			args:                     []string{"app", "com", "--global", "globalval"},
			lookup:                   wargcore.LookupMap(nil),
			expectedPassedPath:       []string{"com"},
			expectedPassedFlagValues: wargcore.PassedFlags{"--global": "globalval", "--help": "default"},
			expectedErr:              false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test validation in TestApp_Validate
			err := tt.app.Validate()
			require.Nil(t, err)

			actualPR, actualErr := tt.app.Parse(parseopt.Args(tt.args), parseopt.LookupEnv(tt.lookup))

			if tt.expectedErr {
				require.NotNil(t, actualErr)
				return
			} else {
				require.Nil(t, actualErr)
			}
			actualPath := slices.Clone(actualPR.Context.ParseState.SectionPath)
			if actualPR.Context.ParseState.CurrentCommandName != "" {
				actualPath = append(actualPath, actualPR.Context.ParseState.CurrentCommandName)
			}
			require.Equal(t, tt.expectedPassedPath, actualPath)
			require.Equal(t, tt.expectedPassedFlagValues, actualPR.Context.Flags)
		})
	}
}
func TestCustomVersion(t *testing.T) {
	expectedVersion := "customversion"

	app := warg.New(
		"appName",
		expectedVersion,
		section.New(
			"test",
			section.NewCommand("version", "Print version", command.DoNothing),
		),
	)
	err := app.Validate()
	require.Nil(t, err)

	actualPR, err := app.Parse(
		parseopt.Args([]string{"appName"}),
		parseopt.LookupEnv(wargcore.LookupMap(nil)),
	)
	require.Nil(t, err)

	require.Equal(t, expectedVersion, actualPR.Context.App.Version)
}

func TestContextContainsValue(t *testing.T) {
	app := warg.New(
		"appName",
		"v1.0.0",
		section.New(
			"test",
			section.NewCommand("version", "Print version", command.DoNothing),
		),
	)
	err := app.Validate()
	require.Nil(t, err)

	type contextKey struct{}
	expectedValue := "value"

	ctx := context.WithValue(context.Background(), contextKey{}, expectedValue)
	actualPR, err := app.Parse(
		parseopt.Args([]string{"appName"}),
		parseopt.LookupEnv(wargcore.LookupMap(nil)),
		parseopt.Context(ctx),
	)
	require.Nil(t, err)

	require.Equal(t, expectedValue, actualPR.Context.Context.Value(contextKey{}).(string))
}

func TestAppFlagToAddr(t *testing.T) {
	require := require.New(t)
	var flagVal string
	expectedFlagVal := "flag value"
	app := warg.New(
		"appName",
		"v1.0.0",
		section.New(
			"test",
			section.NewCommand(
				"command",
				"Test Command",
				func(ctx wargcore.Context) error {
					require.Equal(expectedFlagVal, flagVal)
					return nil
				},
				command.NewFlag(
					"--flag",
					"Flag for test",
					scalar.String(
						scalar.PointerTo(&flagVal),
					),
				),
			),
		),
	)
	err := app.Validate()
	require.NoError(err)

	pr, err := app.Parse(parseopt.Args([]string{"appName", "command", "--flag", "flag value"}))
	require.NoError(err)
	err = pr.Action(pr.Context)
	require.NoError(err)
	require.Equal(expectedFlagVal, flagVal)
}
