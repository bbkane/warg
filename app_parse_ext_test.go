package warg_test

// external tests - import warg like it's an external package

import (
	"context"
	"path/filepath"
	"testing"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/cli"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/config"
	"go.bbkane.com/warg/config/jsonreader"
	"go.bbkane.com/warg/config/yamlreader"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/path"
	"go.bbkane.com/warg/section"
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
		app                      cli.App
		args                     []string
		lookup                   cli.LookupFunc
		expectedPassedPath       []string
		expectedPassedFlagValues cli.PassedFlags
		expectedErr              bool
	}{

		{
			name: "envvar",
			app: warg.NewApp(
				"newAppName", "v1.0.0",
				section.NewSectionT(
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
			lookup: cli.LookupMap(
				map[string]string{
					"there":     "there",
					"alsothere": "alsothere",
				},
			),
			expectedPassedPath: []string{"test"},
			expectedPassedFlagValues: cli.PassedFlags{
				"--flag": "there",
				"--help": "default",
			},
			expectedErr: false,
		},
		{
			name: "addCommandFlags",
			app: func() cli.App {
				fm := cli.FlagMap{
					"--flag1": flag.NewFlag("--flag1 value", scalar.String()),
					"--flag2": flag.NewFlag("--flag1 value", scalar.String()),
				}
				app := warg.NewApp(
					"newAppName", "v1.0.0",
					section.NewSectionT(
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
			lookup:                   cli.LookupMap(nil),
			expectedPassedPath:       []string{"test"},
			expectedPassedFlagValues: cli.PassedFlags{"--flag1": "val1", "--help": "default"},
			expectedErr:              false,
		},

		// Note: Will need to update this if https://github.com/bbkane/warg/issues/36 gets implemented
		{
			name: "invalidFlagsErrorEvenForHelp",
			app: warg.NewApp(
				"newAppName", "v1.0.0",
				section.NewSectionT(
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
			lookup:                   cli.LookupMap(map[string]string{"USER": "bbkane"}),
			expectedPassedPath:       []string{"present"},
			expectedPassedFlagValues: cli.PassedFlags{"--help": "default"},
			expectedErr:              true,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			err := tt.app.Validate()
			require.Nil(t, err)

			actualPR, actualErr := tt.app.Parse(cli.OverrideArgs(tt.args), cli.OverrideLookupFunc(tt.lookup))

			if tt.expectedErr {
				require.NotNil(t, actualErr)
				return
			} else {
				require.Nil(t, actualErr)
			}
			require.Equal(t, tt.expectedPassedPath, actualPR.Context.Path)
			require.Equal(t, tt.expectedPassedFlagValues, actualPR.Context.Flags)
		})
	}
}

// TestApp_Parse_rootSection tests the Parse method, when only a root section is needed (i.e., no special app opts, names, versions, or LookupMaps
func TestApp_Parse_rootSection(t *testing.T) {
	tests := []struct {
		name                     string
		rootSection              cli.SectionT
		args                     []string
		expectedPassedPath       []string
		expectedPassedFlagValues cli.PassedFlags
		expectedErr              bool
	}{
		{
			name: "fromMain",
			rootSection: section.NewSectionT(
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
			expectedPassedFlagValues: cli.PassedFlags{"--com1f1": int(1), "--help": "default"},
			expectedErr:              false,
		},
		{
			name: "noSection",
			rootSection: section.NewSectionT(
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
			rootSection: section.NewSectionT(
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
			expectedPassedFlagValues: cli.PassedFlags{"--flag": "hi", "--help": "default"},
			expectedErr:              false,
		},
		{
			name: "extraFlag",
			rootSection: section.NewSectionT(
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
			rootSection: section.NewSectionT(
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
			expectedPassedFlagValues: cli.PassedFlags{},
			expectedErr:              true,
		},
		{
			name: "flagAlias",
			rootSection: section.NewSectionT(
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
			expectedPassedFlagValues: cli.PassedFlags{"--flag": "val", "--help": "default"},
			expectedErr:              false,
		},
		{
			name: "flagAliasWithList",
			rootSection: section.NewSectionT(
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
			expectedPassedFlagValues: cli.PassedFlags{"--flag": []string{"1", "2", "3", "4"}, "--help": "default"},
			expectedErr:              false,
		},
		{
			name: "badHelp",
			rootSection: section.NewSectionT(
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
			rootSection: section.NewSectionT(
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
			expectedPassedFlagValues: cli.PassedFlags{"--flag": map[string]bool{"true": true, "false": false}, "--help": "default"},
			expectedErr:              false,
		},
		{
			name: "passAbsentSection",
			rootSection: section.NewSectionT(
				"help for test",
				section.NewCommand(
					"com",
					"help for com",
					command.DoNothing,
				),
			),

			args:                     []string{"app", "badSectionName"},
			expectedPassedPath:       []string{"com1"},
			expectedPassedFlagValues: cli.PassedFlags{"--help": "default"},
			expectedErr:              true,
		},
		{
			name: "scalarFlagPassedTwice",
			rootSection: section.NewSectionT(
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
			expectedPassedFlagValues: cli.PassedFlags{"--flag": int(1), "--help": "default"},
			expectedErr:              true,
		},
		{
			name: "passedFlagBeforeCommand",
			rootSection: section.NewSectionT(
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
			expectedPassedFlagValues: cli.PassedFlags{"--flag": int(1), "--help": "default"},
			expectedErr:              true,
		},
		{
			name: "existingSectionsExistingCommands",
			rootSection: section.NewSectionT(
				"help for test",
				section.SectionMap(
					cli.SectionMapT{
						"section": section.NewSectionT(
							"help for section",
							section.CommandMap(
								cli.CommandMap{
									"command": command.NewCommand(
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
			expectedPassedFlagValues: cli.PassedFlags{
				"--flag": 1,
				"--help": "default",
			},
			expectedErr: false,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			app := warg.NewApp(
				"newAppName", "v1.0.0",
				tt.rootSection,
				warg.SkipValidation(),
			)

			err := app.Validate()
			require.Nil(t, err)

			actualPR, actualErr := app.Parse(cli.OverrideArgs(tt.args), cli.OverrideLookupFunc(cli.LookupMap(nil)))

			if tt.expectedErr {
				require.Error(t, actualErr)
				return
			} else {
				require.NoError(t, actualErr)
			}
			require.Equal(t, tt.expectedPassedPath, actualPR.Context.Path)
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
		expectedPassedFlagValues cli.PassedFlags
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
			expectedPassedFlagValues: cli.PassedFlags{
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
			expectedPassedFlagValues: cli.PassedFlags{"--flag": "setAfter", "--help": "default"},
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
			expectedPassedFlagValues: cli.PassedFlags{
				"--help": "default",
				"--flag": []string{"b", "c"},
			},
			expectedErr: false,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			app := warg.NewApp(
				"newAppName", "v1.0.0",
				section.NewSectionT(
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

			actualPR, actualErr := app.Parse(cli.OverrideArgs(tt.args), cli.OverrideLookupFunc(cli.LookupMap(nil)))

			if tt.expectedErr {
				require.Error(t, actualErr)
				return
			} else {
				require.NoError(t, actualErr)
			}
			require.Equal(t, tt.expectedPassedPath, actualPR.Context.Path)
			require.Equal(t, tt.expectedPassedFlagValues, actualPR.Context.Flags)
		})
	}
}

func TestApp_Parse_config(t *testing.T) {
	tests := []struct {
		name                     string
		app                      cli.App
		args                     []string
		lookup                   cli.LookupFunc
		expectedPassedPath       []string
		expectedPassedFlagValues cli.PassedFlags
		expectedErr              bool
	}{
		{
			name: "configFlag",
			app: warg.NewApp(
				"newAppName", "v1.0.0",
				section.NewSectionT(
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
					"--config",
					[]scalar.ScalarOpt[path.Path]{scalar.Default(path.New("defaultconfigval"))},
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
					"config flag",
				),
				warg.SkipValidation(),
			),
			args:               []string{"test", "print", "--config", "passedconfigval"},
			lookup:             cli.LookupMap(nil),
			expectedPassedPath: []string{"print"},
			expectedPassedFlagValues: cli.PassedFlags{
				"--key":    "mapkeyval",
				"--config": path.New("passedconfigval"),
				"--help":   "default",
			},
			expectedErr: false,
		},
		{
			name: "simpleJSONConfig",
			app: warg.NewApp(
				"newAppName", "v1.0.0",
				section.NewSectionT("help for test",
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
					"--config",
					[]scalar.ScalarOpt[path.Path]{
						scalar.Default(
							testDataFilePath(t.Name(), "simpleJSONConfig", "simple_json_config.json"),
						),
					},
					jsonreader.New,
					"path to config",
				),
				warg.SkipValidation(),
			),

			args:               []string{"app", "com"},
			lookup:             cli.LookupMap(nil),
			expectedPassedPath: []string{"com"},
			expectedPassedFlagValues: cli.PassedFlags{
				"--config": testDataFilePath(t.Name(), "simpleJSONConfig", "simple_json_config.json"),
				"--val":    "hi",
				"--help":   "default",
			},
			expectedErr: false,
		},
		{
			// JSON needs some help to support number decoding
			name: "numJSONConfig",
			app: warg.NewApp(
				"newAppName", "v1.0.0",
				section.NewSectionT("help for test",
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
					"--config",
					[]scalar.ScalarOpt[path.Path]{
						scalar.Default(
							testDataFilePath(t.Name(), "numJSONConfig", "config.json"),
						),
					},
					jsonreader.New,
					"path to config",
				),
				warg.SkipValidation(),
			),

			args:               []string{"app", "com"},
			lookup:             cli.LookupMap(nil),
			expectedPassedPath: []string{"com"},
			expectedPassedFlagValues: cli.PassedFlags{
				"--config": testDataFilePath(t.Name(), "numJSONConfig", "config.json"),
				// "--floatval": 1.5,
				"--intval": 1,
				"--help":   "default",
			},
			expectedErr: false,
		},
		{
			name: "configSlice",
			app: warg.NewApp(
				"newAppName", "v1.0.0",
				section.NewSectionT(
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
					"--config",
					[]scalar.ScalarOpt[path.Path]{
						scalar.Default(
							testDataFilePath(t.Name(), "configSlice", "config_slice.json"),
						),
					},
					jsonreader.New,
					"config flag",
				),
				warg.SkipValidation(),
			),
			args:               []string{"test", "print"},
			lookup:             cli.LookupMap(nil),
			expectedPassedPath: []string{"print"},
			expectedPassedFlagValues: cli.PassedFlags{
				"--subreddits": []string{"earthporn", "wallpapers"},
				"--config":     testDataFilePath(t.Name(), "configSlice", "config_slice.json"),
				"--help":       "default",
			},
			expectedErr: false,
		},
		{
			name: "JSONConfigStringSlice",
			app: warg.NewApp(
				"newAppName", "v1.0.0",
				section.NewSectionT("help for test",
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
					"--config",
					[]scalar.ScalarOpt[path.Path]{
						scalar.Default(
							testDataFilePath(t.Name(), "JSONConfigStringSlice", "config.json"),
						),
					},
					jsonreader.New,
					"path to config",
				),
				warg.SkipValidation(),
			),
			args:               []string{"app", "com"},
			lookup:             cli.LookupMap(nil),
			expectedPassedPath: []string{"com"},
			expectedPassedFlagValues: cli.PassedFlags{
				"--config": testDataFilePath(t.Name(), "JSONConfigStringSlice", "config.json"),
				"--val":    []string{"from", "config"},
				"--help":   "default",
			},
			expectedErr: false,
		},
		{
			name: "YAMLConfigStringSlice",
			app: warg.NewApp(
				"newAppName", "v1.0.0",
				section.NewSectionT("help for test",
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
					"--config",
					[]scalar.ScalarOpt[path.Path]{
						scalar.Default(
							testDataFilePath(t.Name(), "YAMLConfigStringSlice", "config.yaml"),
						),
					},
					yamlreader.New,
					"path to config",
				),
				warg.SkipValidation(),
			),

			args:               []string{"app", "com"},
			lookup:             cli.LookupMap(nil),
			expectedPassedPath: []string{"com"},
			expectedPassedFlagValues: cli.PassedFlags{
				"--config": testDataFilePath(t.Name(), "YAMLConfigStringSlice", "config.yaml"),
				"--val":    []string{"from", "config"},
				"--help":   "default",
			},
			expectedErr: false,
		},
		{
			name: "JSONConfigMap",
			app: warg.NewApp(
				"newAppName", "v1.0.0",
				section.NewSectionT(
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
					"--config",
					[]scalar.ScalarOpt[path.Path]{
						scalar.Default(
							testDataFilePath(t.Name(), "JSONConfigMap", "config.json"),
						),
					},
					jsonreader.New,
					"path to config",
				),
				warg.SkipValidation(),
			),

			args:               []string{"app", "com"},
			lookup:             cli.LookupMap(nil),
			expectedPassedPath: []string{"com"},
			expectedPassedFlagValues: cli.PassedFlags{
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

			actualPR, actualErr := tt.app.Parse(cli.OverrideArgs(tt.args), cli.OverrideLookupFunc(tt.lookup))

			if tt.expectedErr {
				require.NotNil(t, actualErr)
				return
			} else {
				require.Nil(t, actualErr)
			}
			require.Equal(t, tt.expectedPassedPath, actualPR.Context.Path)
			require.Equal(t, tt.expectedPassedFlagValues, actualPR.Context.Flags)
		})
	}
}

// This is the same as TestApp_Parse, but that's too long for a single test
func TestApp_Parse_GlobalFlag(t *testing.T) {
	tests := []struct {
		name                     string
		app                      cli.App
		args                     []string
		lookup                   cli.LookupFunc
		expectedPassedPath       []string
		expectedPassedFlagValues cli.PassedFlags
		expectedErr              bool
	}{
		{
			name: "globalFlag",
			app: warg.NewApp(
				"newAppName", "v1.0.0",
				section.NewSectionT(
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
			lookup:                   cli.LookupMap(nil),
			expectedPassedPath:       []string{"com"},
			expectedPassedFlagValues: cli.PassedFlags{"--global": "globalval", "--help": "default"},
			expectedErr:              false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test validation in TestApp_Validate
			err := tt.app.Validate()
			require.Nil(t, err)

			actualPR, actualErr := tt.app.Parse(cli.OverrideArgs(tt.args), cli.OverrideLookupFunc(tt.lookup))

			if tt.expectedErr {
				require.NotNil(t, actualErr)
				return
			} else {
				require.Nil(t, actualErr)
			}
			require.Equal(t, tt.expectedPassedPath, actualPR.Context.Path)
			require.Equal(t, tt.expectedPassedFlagValues, actualPR.Context.Flags)
		})
	}
}
func TestCustomVersion(t *testing.T) {
	expectedVersion := "customversion"

	app := warg.NewApp(
		"appName",
		expectedVersion,
		section.NewSectionT(
			"test",
			section.NewCommand("version", "Print version", command.DoNothing),
		),
	)
	err := app.Validate()
	require.Nil(t, err)

	actualPR, err := app.Parse(
		cli.OverrideArgs([]string{"appName"}),
		cli.OverrideLookupFunc(cli.LookupMap(nil)),
	)
	require.Nil(t, err)

	require.Equal(t, expectedVersion, actualPR.Context.Version)
}

func TestContextContainsValue(t *testing.T) {
	app := warg.NewApp(
		"appName",
		"v1.0.0",
		section.NewSectionT(
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
		cli.OverrideArgs([]string{"appName"}),
		cli.OverrideLookupFunc(cli.LookupMap(nil)),
		cli.AddContext(ctx),
	)
	require.Nil(t, err)

	require.Equal(t, expectedValue, actualPR.Context.Context.Value(contextKey{}).(string))
}

func TestAppFlagToAddr(t *testing.T) {
	require := require.New(t)
	var flagVal string
	expectedFlagVal := "flag value"
	app := warg.NewApp(
		"appName",
		"v1.0.0",
		section.NewSectionT(
			"test",
			section.NewCommand(
				"command",
				"Test Command",
				func(ctx cli.Context) error {
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

	pr, err := app.Parse(cli.OverrideArgs([]string{"appName", "command", "--flag", "flag value"}))
	require.NoError(err)
	err = pr.Action(pr.Context)
	require.NoError(err)
	require.Equal(expectedFlagVal, flagVal)
}
