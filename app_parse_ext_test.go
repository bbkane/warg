package warg_test

// external tests - import warg like it's an external package

import (
	"context"
	"path/filepath"
	"slices"
	"testing"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/config"
	"go.bbkane.com/warg/config/jsonreader"
	"go.bbkane.com/warg/config/yamlreader"
	"go.bbkane.com/warg/path"
	"go.bbkane.com/warg/value/dict"
	"go.bbkane.com/warg/value/scalar"
	"go.bbkane.com/warg/value/slice"

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
		app                      warg.App
		args                     []string
		lookup                   warg.LookupEnv
		expectedPassedPath       []string
		expectedPassedFlagValues warg.PassedFlags
		expectedErr              bool
	}{

		{
			name: "envvar",
			app: warg.New(
				"newAppName", "v1.0.0",
				warg.NewSection(
					"help for test",
					warg.NewSubCmd(
						"test",
						"blah",
						warg.Unimplemented(),
						warg.NewCmdFlag(
							"--flag",
							"help for --flag",
							scalar.String(),
							warg.EnvVars("notthere", "there", "alsothere"),
						),
					),
				),

				warg.SkipAll(),
			),
			args: []string{t.Name(), "test"},
			lookup: warg.LookupMap(
				map[string]string{
					"there":     "there",
					"alsothere": "alsothere",
				},
			),
			expectedPassedPath: []string{"test"},
			expectedPassedFlagValues: warg.PassedFlags{
				"--flag": "there",
				"--help": "default",
			},
			expectedErr: false,
		},
		{
			name: "addCommandFlags",
			app: func() warg.App {
				fm := warg.FlagMap{
					"--flag1": warg.NewFlag("--flag1 value", scalar.String()),
					"--flag2": warg.NewFlag("--flag1 value", scalar.String()),
				}
				app := warg.New(
					"newAppName", "v1.0.0",
					warg.NewSection(
						"help for section",
						warg.NewSubCmd(
							"test",
							"help for test",
							warg.Unimplemented(),
							warg.CmdFlagMap(fm),
						),
					),

					warg.SkipAll(),
				)
				return app
			}(),

			args:                     []string{t.Name(), "test", "--flag1", "val1"},
			lookup:                   warg.LookupMap(nil),
			expectedPassedPath:       []string{"test"},
			expectedPassedFlagValues: warg.PassedFlags{"--flag1": "val1", "--help": "default"},
			expectedErr:              false,
		},

		// Note: Will need to update this if https://github.com/bbkane/warg/issues/36 gets implemented
		{
			name: "invalidFlagsErrorEvenForHelp",
			app: warg.New(
				"newAppName", "v1.0.0",
				warg.NewSection(
					string("A virtual assistant"),
					warg.NewSubCmd(
						"present",
						"Formally present a guest (guests are never introduced, always presented).",
						warg.Unimplemented(),
						warg.NewCmdFlag(
							"--name",
							"Guest to address.",
							scalar.String(scalar.Choices("bob")),
							warg.Alias("-n"),
							warg.EnvVars("BUTLER_PRESENT_NAME", "USER"),
							warg.Required(),
						),
					),
				),
				warg.SkipAll(),
			),

			args:                     []string{"app", "present", "-h"},
			lookup:                   warg.LookupMap(map[string]string{"USER": "bbkane"}),
			expectedPassedPath:       []string{"present"},
			expectedPassedFlagValues: warg.PassedFlags{"--help": "default"},
			expectedErr:              true,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			err := tt.app.Validate()
			require.Nil(t, err)

			actualPR, actualErr := tt.app.Parse(warg.ParseWithArgs(tt.args), warg.ParseWithLookupEnv(tt.lookup))

			if tt.expectedErr {
				require.NotNil(t, actualErr)
				return
			} else {
				require.Nil(t, actualErr)
			}
			actualPath := slices.Clone(actualPR.Context.ParseState.SectionPath)
			if actualPR.Context.ParseState.CurrentCmdName != "" {
				actualPath = append(actualPath, actualPR.Context.ParseState.CurrentCmdName)
			}

			require.Equal(t, tt.expectedPassedPath, actualPath)
			require.Equal(t, tt.expectedPassedFlagValues, actualPR.Context.Flags)
		})
	}
}

// TestApp_Parse_rootSection tests the Parse method, when only a root section is needed (i.e., no special app opts, names, versions, or LookupMaps
func TestApp_Parse_rootSection(t *testing.T) {
	//exhaustruct:ignore  // in tests I like to only set what I don't expect
	tests := []struct {
		name                     string
		rootSection              warg.Section
		args                     []string
		expectedPassedPath       []string
		expectedPassedFlagValues warg.PassedFlags
		expectedForwardedArgs    []string
		expectedErr              bool
	}{
		{
			name: "minimal",
			rootSection: warg.NewSection(
				"help for test",
				warg.NewSubCmd(
					"com",
					"com help",
					warg.Unimplemented(),
				),
			),
			args:                     []string{"test", "com"},
			expectedPassedPath:       []string{"com"},
			expectedPassedFlagValues: warg.PassedFlags{"--help": "default"},
		},
		{
			name: "fromMain",
			rootSection: warg.NewSection(
				"help for test",
				warg.NewSubSection(
					"cat1",
					"help for cat1",
					warg.NewSubCmd(
						"com1",
						"help for com1",
						warg.Unimplemented(),
						warg.NewCmdFlag(
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
			expectedPassedFlagValues: warg.PassedFlags{"--com1f1": int(1), "--help": "default"},
		},
		{
			name: "noSection",
			rootSection: warg.NewSection(
				"help for test",
				warg.NewSubCmd("com", "command for validation", warg.Unimplemented()),
			),

			args:                     []string{"app"},
			expectedPassedPath:       nil,
			expectedPassedFlagValues: map[string]interface{}{"--help": "default"},
		},
		{
			name: "flagDefault",
			rootSection: warg.NewSection(
				"help for test",
				warg.NewSubCmd(
					"com",
					"com help",
					warg.Unimplemented(),
					warg.NewCmdFlag(
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
			expectedPassedFlagValues: warg.PassedFlags{"--flag": "hi", "--help": "default"},
		},
		{
			name: "extraFlag",
			rootSection: warg.NewSection(
				"help for test",
				warg.NewSubCmd(
					"com",
					"com help",
					warg.Unimplemented(),
					warg.NewCmdFlag(
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
			rootSection: warg.NewSection(
				"help for test",
				warg.NewSubCmd(
					"test",
					"blah",
					warg.Unimplemented(),
					warg.NewCmdFlag(
						"--flag",
						"help for --flag",
						scalar.String(),
						warg.Required(),
					),
				),
			),
			args:                     []string{t.Name(), "test"},
			expectedPassedPath:       []string{"test"},
			expectedPassedFlagValues: warg.PassedFlags{},
			expectedErr:              true,
		},
		{
			name: "flagAlias",
			rootSection: warg.NewSection(
				"help for section",
				warg.NewSubCmd(
					"test",
					"help for test",
					warg.Unimplemented(),
					warg.NewCmdFlag(
						"--flag",
						"help for --flag",
						scalar.String(),
						warg.Alias("-f"),
					),
				),
			),
			args:                     []string{t.Name(), "test", "-f", "val"},
			expectedPassedPath:       []string{"test"},
			expectedPassedFlagValues: warg.PassedFlags{"--flag": "val", "--help": "default"},
		},
		{
			name: "flagAliasWithList",
			rootSection: warg.NewSection(
				"help for section",
				warg.NewSubCmd(
					"test",
					"help for test",
					warg.Unimplemented(),
					warg.NewCmdFlag(
						"--flag",
						"help for --flag",
						slice.String(),
						warg.Alias("-f"),
					),
				),
			),
			args:                     []string{t.Name(), "test", "-f", "1", "--flag", "2", "-f", "3", "--flag", "4"},
			expectedPassedPath:       []string{"test"},
			expectedPassedFlagValues: warg.PassedFlags{"--flag": []string{"1", "2", "3", "4"}, "--help": "default"},
			expectedErr:              false,
		},
		{
			name: "badHelp",
			rootSection: warg.NewSection(
				"help for section",
				warg.NewSubCmd(
					"test",
					"help for test",
					warg.Unimplemented(),
				),
			),
			args:                     []string{t.Name(), "test", "-h", "badhelpval"},
			expectedPassedPath:       nil,
			expectedPassedFlagValues: nil,
			expectedErr:              true,
		},
		{
			name: "dictUpdate",
			rootSection: warg.NewSection(
				"help for test",
				warg.NewSubCmd(
					"com1",
					"help for com1",
					warg.Unimplemented(),
					warg.NewCmdFlag(
						string("--flag"),
						"flag help",
						dict.Bool(),
					),
				),
			),

			args:                     []string{"app", "com1", "--flag", "true=true", "--flag", "false=false"},
			expectedPassedPath:       []string{"com1"},
			expectedPassedFlagValues: warg.PassedFlags{"--flag": map[string]bool{"true": true, "false": false}, "--help": "default"},
		},
		{
			name: "passAbsentSection",
			rootSection: warg.NewSection(
				"help for test",
				warg.NewSubCmd(
					"com",
					"help for com",
					warg.Unimplemented(),
				),
			),

			args:                     []string{"app", "badSectionName"},
			expectedPassedPath:       []string{"com1"},
			expectedPassedFlagValues: warg.PassedFlags{"--help": "default"},
			expectedErr:              true,
		},
		{
			name: "scalarFlagPassedTwice",
			rootSection: warg.NewSection(
				"help for test",
				warg.NewSubCmd(
					"com",
					"help for com1",
					warg.Unimplemented(),
					warg.NewCmdFlag(
						"--flag",
						"flag help",
						scalar.Int(),
					),
				),
			),

			args:                     []string{"app", "com", "--flag", "1", "--flag", "2"},
			expectedPassedPath:       []string{"com"},
			expectedPassedFlagValues: warg.PassedFlags{"--flag": int(1), "--help": "default"},
			expectedErr:              true,
		},
		{
			name: "passedFlagBeforeCommand",
			rootSection: warg.NewSection(
				"help for test",
				warg.NewSubCmd(
					"com",
					"help for com",
					warg.Unimplemented(),
					warg.NewCmdFlag(
						"--flag",
						"flag help",
						scalar.Int(),
					),
				),
			),

			args:                     []string{"app", "--flag", "1", "com"},
			expectedPassedPath:       []string{"com"},
			expectedPassedFlagValues: warg.PassedFlags{"--flag": int(1), "--help": "default"},
			expectedErr:              true,
		},
		{
			name: "existingSectionsExistingCommands",
			rootSection: warg.NewSection(
				"help for test",
				warg.SubSectionMap(
					warg.SectionMap{
						"section": warg.NewSection(
							"help for section",
							warg.SubCmdMap(
								warg.CmdMap{
									"command": warg.NewCmd(
										"help for command",
										warg.Unimplemented(),
										warg.NewCmdFlag(
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
			expectedPassedFlagValues: warg.PassedFlags{
				"--flag": 1,
				"--help": "default",
			},
		},
		{
			name: "flagWithEmptyValue",
			rootSection: warg.NewSection(
				"help for test",
				warg.NewSubCmd(
					"command",
					"help for command",
					warg.Unimplemented(),
					warg.NewCmdFlag(
						"--flag",
						"flag help",
						scalar.String(),
					),
				),
			),
			args:               []string{"app", "command", "--flag", ""},
			expectedPassedPath: []string{"command"},
			expectedPassedFlagValues: warg.PassedFlags{
				"--flag": "",
				"--help": "default",
			},
		},
		{
			name: "allowForwardedArgs",
			rootSection: warg.NewSection(
				"help for test",
				warg.NewSubCmd(
					"com",
					"com help",
					warg.Unimplemented(),
					warg.AllowForwardedArgs(),
				),
			),
			args:                     []string{"test", "com", "--", "arg1", "arg2"},
			expectedPassedPath:       []string{"com"},
			expectedPassedFlagValues: warg.PassedFlags{"--help": "default"},
			expectedForwardedArgs:    []string{"arg1", "arg2"},
			expectedErr:              false,
		},
		{
			name: "forwardedArgsButNoPassedArgs",
			rootSection: warg.NewSection(
				"help for test",
				warg.NewSubCmd(
					"com",
					"com help",
					warg.Unimplemented(),
					warg.AllowForwardedArgs(),
				),
			),
			args:                     []string{"test", "com", "--"},
			expectedPassedPath:       []string{"com"},
			expectedPassedFlagValues: warg.PassedFlags{"--help": "default"},
			expectedErr:              true,
		},
		{
			name: "forwardedArgsNotAllowedButPassedAnyway",
			rootSection: warg.NewSection(
				"help for test",
				warg.NewSubCmd(
					"com",
					"com help",
					warg.Unimplemented(),
				),
			),
			args:                     []string{"test", "com", "--", "arg1", "arg2"},
			expectedPassedPath:       []string{"com"},
			expectedPassedFlagValues: warg.PassedFlags{"--help": "default"},
			expectedErr:              true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			app := warg.New(
				"newAppName", "v1.0.0",
				tt.rootSection,

				warg.SkipAll(),
			)

			err := app.Validate()
			require.Nil(t, err)

			actualPR, actualErr := app.Parse(warg.ParseWithArgs(tt.args), warg.ParseWithLookupEnv(warg.LookupMap(nil)))

			if tt.expectedErr {
				require.Error(t, actualErr)
				return
			} else {
				require.NoError(t, actualErr)
			}
			actualPath := slices.Clone(actualPR.Context.ParseState.SectionPath)
			if actualPR.Context.ParseState.CurrentCmdName != "" {
				actualPath = append(actualPath, actualPR.Context.ParseState.CurrentCmdName)
			}
			require.Equal(t, tt.expectedPassedPath, actualPath)
			require.Equal(t, tt.expectedPassedFlagValues, actualPR.Context.Flags)
		})
	}
}

func TestApp_Parse_unsetSetinel(t *testing.T) {
	tests := []struct {
		name                     string
		flagDef                  warg.CmdOpt
		args                     []string
		expectedPassedPath       []string
		expectedPassedFlagValues warg.PassedFlags
		expectedErr              bool
	}{
		{
			name: "unsetSentinelScalarSuccess",
			flagDef: warg.NewCmdFlag(
				"--flag",
				"help for --flag",
				scalar.String(scalar.Default("default")),
				warg.UnsetSentinel("UNSET"),
			),
			args:               []string{t.Name(), "test", "--flag", "UNSET"},
			expectedPassedPath: []string{"test"},
			expectedPassedFlagValues: warg.PassedFlags{
				"--help": "default",
			},
			expectedErr: false,
		},
		{
			name: "unsetSentinelScalarUpdate",
			flagDef: warg.NewCmdFlag(
				"--flag",
				"help for --flag",
				scalar.String(scalar.Default("default")),
				warg.UnsetSentinel("UNSET"),
			),
			args:                     []string{t.Name(), "test", "--flag", "UNSET", "--flag", "setAfter"},
			expectedPassedPath:       []string{"test"},
			expectedPassedFlagValues: warg.PassedFlags{"--flag": "setAfter", "--help": "default"},
			expectedErr:              false,
		},
		{
			name: "unsetSentinelSlice",
			flagDef: warg.NewCmdFlag(
				"--flag",
				"help for --flag",
				slice.String(slice.Default([]string{"default"})),
				warg.UnsetSentinel("UNSET"),
			),
			args:               []string{t.Name(), "test", "--flag", "a", "--flag", "UNSET", "--flag", "b", "--flag", "c"},
			expectedPassedPath: []string{"test"},
			expectedPassedFlagValues: warg.PassedFlags{
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
				warg.NewSection(
					"help for test",
					warg.NewSubCmd(
						"test",
						"help for test",
						warg.Unimplemented(),
						tt.flagDef,
					),
				),

				warg.SkipAll(),
			)

			err := app.Validate()
			require.Nil(t, err)

			actualPR, actualErr := app.Parse(warg.ParseWithArgs(tt.args), warg.ParseWithLookupEnv(warg.LookupMap(nil)))

			if tt.expectedErr {
				require.Error(t, actualErr)
				return
			} else {
				require.NoError(t, actualErr)
			}
			actualPath := slices.Clone(actualPR.Context.ParseState.SectionPath)
			if actualPR.Context.ParseState.CurrentCmdName != "" {
				actualPath = append(actualPath, actualPR.Context.ParseState.CurrentCmdName)
			}
			require.Equal(t, tt.expectedPassedPath, actualPath)
			require.Equal(t, tt.expectedPassedFlagValues, actualPR.Context.Flags)
		})
	}
}

func TestApp_Parse_config(t *testing.T) {
	tests := []struct {
		name                     string
		app                      warg.App
		args                     []string
		lookup                   warg.LookupEnv
		expectedPassedPath       []string
		expectedPassedFlagValues warg.PassedFlags
		expectedErr              bool
	}{
		{
			name: "configFlag",
			app: warg.New(
				"newAppName", "v1.0.0",
				warg.NewSection(
					"help for test",
					warg.NewSubCmd(
						"print",
						"print key value",
						warg.Unimplemented(),
						warg.NewCmdFlag(
							"--key",
							"a key",
							scalar.String(
								scalar.Default("defaultkeyval"),
							),
							warg.ConfigPath("key"),
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
					warg.FlagMap{
						"--config": warg.NewFlag(
							"Path to config file",
							scalar.Path(
								scalar.Default(path.New("defaultconfigval")),
							),
						),
					},
				),

				warg.SkipAll(),
			),
			args:               []string{"test", "print", "--config", "passedconfigval"},
			lookup:             warg.LookupMap(nil),
			expectedPassedPath: []string{"print"},
			expectedPassedFlagValues: warg.PassedFlags{
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
				warg.NewSection("help for test",
					warg.NewSubCmd(
						"com",
						"help for com",
						warg.Unimplemented(),
						warg.NewCmdFlag(
							"--val",
							"flag help",
							scalar.String(),
							warg.ConfigPath("params.val"),
						),
					),
				),
				warg.ConfigFlag(
					jsonreader.New,
					warg.FlagMap{
						"--config": warg.NewFlag(
							"path to config",
							scalar.Path(
								scalar.Default(
									testDataFilePath(t.Name(), "simpleJSONConfig", "simple_json_config.json"),
								),
							),
							warg.Alias("-c"),
						),
					},
				),
				warg.SkipAll(),
			),

			args:               []string{"app", "com"},
			lookup:             warg.LookupMap(nil),
			expectedPassedPath: []string{"com"},
			expectedPassedFlagValues: warg.PassedFlags{
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
				warg.NewSection("help for test",
					warg.NewSubCmd(
						"com",
						"help for com",
						warg.Unimplemented(),
						warg.NewCmdFlag(
							"--intval",
							"flag help",
							scalar.Int(),
							warg.ConfigPath("params.intval"),
						),
					),
				),
				warg.ConfigFlag(
					jsonreader.New,
					warg.FlagMap{
						"--config": warg.NewFlag(
							"path to config",
							scalar.Path(
								scalar.Default(
									testDataFilePath(t.Name(), "numJSONConfig", "config.json"),
								),
							),
							warg.Alias("-c"),
						),
					},
				),
				warg.SkipAll(),
			),

			args:               []string{"app", "com"},
			lookup:             warg.LookupMap(nil),
			expectedPassedPath: []string{"com"},
			expectedPassedFlagValues: warg.PassedFlags{
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
				warg.NewSection(
					"help for test",
					warg.NewSubCmd(
						"print",
						"print key value",
						warg.Unimplemented(),
						warg.NewCmdFlag(
							"--subreddits",
							"the subreddits",
							slice.String(),
							warg.ConfigPath("subreddits[].name"),
						),
					),
				),
				warg.ConfigFlag(
					jsonreader.New,
					warg.FlagMap{
						"--config": warg.NewFlag(
							"path to config",
							scalar.Path(
								scalar.Default(
									testDataFilePath(t.Name(), "configSlice", "config_slice.json"),
								),
							),
							warg.Alias("-c"),
						),
					},
				),
				warg.SkipAll(),
			),
			args:               []string{"test", "print"},
			lookup:             warg.LookupMap(nil),
			expectedPassedPath: []string{"print"},
			expectedPassedFlagValues: warg.PassedFlags{
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
				warg.NewSection("help for test",
					warg.NewSubCmd(
						"com",
						"help for com",
						warg.Unimplemented(),
						warg.NewCmdFlag(
							"--val",
							"flag help",
							slice.String(),
							warg.ConfigPath("val"),
						),
					),
				),
				warg.ConfigFlag(
					jsonreader.New,
					warg.FlagMap{
						"--config": warg.NewFlag(
							"path to config",
							scalar.Path(
								scalar.Default(
									testDataFilePath(t.Name(), "JSONConfigStringSlice", "config.json"),
								),
							),
							warg.Alias("-c"),
						),
					},
				),

				warg.SkipAll(),
			),
			args:               []string{"app", "com"},
			lookup:             warg.LookupMap(nil),
			expectedPassedPath: []string{"com"},
			expectedPassedFlagValues: warg.PassedFlags{
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
				warg.NewSection("help for test",
					warg.NewSubCmd(
						"com",
						"help for com",
						warg.Unimplemented(),
						warg.NewCmdFlag(
							"--val",
							"flag help",
							slice.String(),
							warg.ConfigPath("val"),
						),
					),
				),
				warg.ConfigFlag(
					yamlreader.New,
					warg.FlagMap{
						"--config": warg.NewFlag(
							"path to config",
							scalar.Path(
								scalar.Default(
									testDataFilePath(t.Name(), "YAMLConfigStringSlice", "config.yaml"),
								),
							),
							warg.Alias("-c"),
						),
					},
				),

				warg.SkipAll(),
			),

			args:               []string{"app", "com"},
			lookup:             warg.LookupMap(nil),
			expectedPassedPath: []string{"com"},
			expectedPassedFlagValues: warg.PassedFlags{
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
				warg.NewSection(
					"help for test",
					warg.NewSubCmd(
						"com",
						"help for com",
						warg.Unimplemented(),
						warg.NewCmdFlag(
							"--val",
							"flag help",
							dict.Int(),
							warg.ConfigPath("val"),
						),
					),
				),
				warg.ConfigFlag(
					jsonreader.New,
					warg.FlagMap{
						"--config": warg.NewFlag(
							"path to config",
							scalar.Path(
								scalar.Default(
									testDataFilePath(t.Name(), "JSONConfigMap", "config.json"),
								),
							),
							warg.Alias("-c"),
						),
					},
				),

				warg.SkipAll(),
			),

			args:               []string{"app", "com"},
			lookup:             warg.LookupMap(nil),
			expectedPassedPath: []string{"com"},
			expectedPassedFlagValues: warg.PassedFlags{
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

			actualPR, actualErr := tt.app.Parse(warg.ParseWithArgs(tt.args), warg.ParseWithLookupEnv(tt.lookup))

			if tt.expectedErr {
				require.NotNil(t, actualErr)
				return
			} else {
				require.Nil(t, actualErr)
			}
			actualPath := slices.Clone(actualPR.Context.ParseState.SectionPath)
			if actualPR.Context.ParseState.CurrentCmdName != "" {
				actualPath = append(actualPath, actualPR.Context.ParseState.CurrentCmdName)
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
		app                      warg.App
		args                     []string
		lookup                   warg.LookupEnv
		expectedPassedPath       []string
		expectedPassedFlagValues warg.PassedFlags
		expectedErr              bool
	}{
		{
			name: "globalFlag",
			app: warg.New(
				"newAppName", "v1.0.0",
				warg.NewSection(
					"help for test",
					warg.NewSubCmd(
						"com",
						"help for com",
						warg.Unimplemented(),
					),
				),
				warg.SkipAll(),
				warg.NewGlobalFlag(
					"--global",
					"global flag",
					scalar.String(),
				),
			),

			args:                     []string{"app", "com", "--global", "globalval"},
			lookup:                   warg.LookupMap(nil),
			expectedPassedPath:       []string{"com"},
			expectedPassedFlagValues: warg.PassedFlags{"--global": "globalval", "--help": "default"},
			expectedErr:              false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test validation in TestApp_Validate
			err := tt.app.Validate()
			require.Nil(t, err)

			actualPR, actualErr := tt.app.Parse(warg.ParseWithArgs(tt.args), warg.ParseWithLookupEnv(tt.lookup))

			if tt.expectedErr {
				require.NotNil(t, actualErr)
				return
			} else {
				require.Nil(t, actualErr)
			}
			actualPath := slices.Clone(actualPR.Context.ParseState.SectionPath)
			if actualPR.Context.ParseState.CurrentCmdName != "" {
				actualPath = append(actualPath, actualPR.Context.ParseState.CurrentCmdName)
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
		warg.NewSection(
			"test",
			warg.NewSubCmd("version", "Print version", warg.Unimplemented()),
		),
		warg.SkipVersionCmd(),
	)
	err := app.Validate()
	require.Nil(t, err)

	actualPR, err := app.Parse(
		warg.ParseWithArgs([]string{"appName"}),
		warg.ParseWithLookupEnv(warg.LookupMap(nil)),
	)
	require.Nil(t, err)

	require.Equal(t, expectedVersion, actualPR.Context.App.Version)
}

func TestContextContainsValue(t *testing.T) {
	app := warg.New(
		"appName",
		"v1.0.0",
		warg.NewSection(
			"test",
			warg.NewSubCmd("dummycommand", "Do nothing", warg.Unimplemented()),
		),
	)
	err := app.Validate()
	require.Nil(t, err)

	type contextKey struct{}
	expectedValue := "value"

	ctx := context.WithValue(context.Background(), contextKey{}, expectedValue)
	actualPR, err := app.Parse(
		warg.ParseWithArgs([]string{"appName"}),
		warg.ParseWithLookupEnv(warg.LookupMap(nil)),
		warg.ParseWithContext(ctx),
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
		warg.NewSection(
			"test",
			warg.NewSubCmd(
				"command",
				"Test Command",
				func(ctx warg.CmdContext) error {
					require.Equal(expectedFlagVal, flagVal)
					return nil
				},
				warg.NewCmdFlag(
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

	pr, err := app.Parse(warg.ParseWithArgs([]string{"appName", "command", "--flag", "flag value"}))
	require.NoError(err)
	err = pr.Action(pr.Context)
	require.NoError(err)
	require.Equal(expectedFlagVal, flagVal)
}

func TestSet(t *testing.T) {
	require := require.New(t)
	s := warg.NewSet[string]()

	require.False(s.Contains("a"))
	s.Add("a")
	require.True(s.Contains("a"))
	s.Delete("a")
	require.False(s.Contains("a"))
}
