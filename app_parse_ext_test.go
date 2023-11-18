package warg_test

// external tests - import warg like it's an external package

import (
	"path/filepath"
	"testing"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/config"
	"go.bbkane.com/warg/config/jsonreader"
	"go.bbkane.com/warg/config/yamlreader"
	"go.bbkane.com/warg/flag"
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
func testDataFilePath(testName string, subTestName string, fileName string) string {
	return filepath.Join("testdata", testName, subTestName, fileName)
}

func TestApp_Parse(t *testing.T) {

	tests := []struct {
		name                     string
		app                      warg.App
		args                     []string
		lookup                   warg.LookupFunc
		expectedPassedPath       []string
		expectedPassedFlagValues command.PassedFlags
		expectedErr              bool
	}{
		{
			name: "fromMain",
			app: warg.New(
				"newAppName",
				section.New(
					"help for test",
					section.Flag(
						flag.Name("--af1"),
						"flag help",
						scalar.Int(),
					),
					section.Section(
						"cat1",
						"help for cat1",
						section.Command(
							"com1",
							"help for com1",
							command.DoNothing,
							command.Flag(
								"--com1f1",
								"flag help",
								scalar.Int(
									scalar.Default(10),
								),
							),
						),
					),
				),
				warg.SkipValidation(),
			),

			args:                     []string{"app", "cat1", "com1", "--com1f1", "1"},
			lookup:                   warg.LookupMap(nil),
			expectedPassedPath:       []string{"cat1", "com1"},
			expectedPassedFlagValues: command.PassedFlags{"--com1f1": int(1), "--help": "default"},
			expectedErr:              false,
		},
		{
			name: "noSection",
			app: warg.New(
				"newAppName",
				section.New(
					"help for test",
					section.Flag(
						"--af1",
						"flag help",
						scalar.Int(),
					),
					section.Command("com", "command for validation", command.DoNothing),
				),
				warg.SkipValidation(),
			),
			args:                     []string{"app"},
			lookup:                   warg.LookupMap(nil),
			expectedPassedPath:       nil,
			expectedPassedFlagValues: map[string]interface{}{"--help": "default"},
			expectedErr:              false,
		},
		{
			name: "flagDefault",
			app: warg.New(
				"newAppName",
				section.New(
					"help for test",
					section.Command(
						"com",
						"com help",
						command.DoNothing,
						command.Flag(
							"--flag",
							"flag help",
							scalar.String(
								scalar.Default("hi"),
							),
						),
					),
				),
				warg.SkipValidation(),
			),
			args:                     []string{"test", "com"},
			lookup:                   warg.LookupMap(nil),
			expectedPassedPath:       []string{"com"},
			expectedPassedFlagValues: command.PassedFlags{"--flag": "hi", "--help": "default"},
			expectedErr:              false,
		},
		{
			name: "extraFlag",
			app: warg.New(
				"newAppName",
				section.New(
					"help for test",
					section.Command(
						"com",
						"com help",
						command.DoNothing,
						command.Flag(
							"--flag",
							"flag help",
							scalar.String(
								scalar.Default("hi"),
							),
						),
					),
				),
				warg.SkipValidation(),
			),
			args:                     []string{"test", "com", "--unexpected", "value"},
			lookup:                   warg.LookupMap(nil),
			expectedPassedPath:       nil,
			expectedPassedFlagValues: nil,
			expectedErr:              true,
		},
		{
			name: "configFlag",
			app: warg.New(
				"newAppName",
				section.New(
					"help for test",
					section.Flag(
						"--key",
						"a key",
						scalar.String(
							scalar.Default("defaultkeyval"),
						),
						flag.ConfigPath("key"),
					),
					section.Command("print", "print key value", command.DoNothing),
				),
				warg.ConfigFlag(
					"--config",
					[]scalar.ScalarOpt[string]{scalar.Default("defaultconfigval")},
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
			lookup:             warg.LookupMap(nil),
			expectedPassedPath: []string{"print"},
			expectedPassedFlagValues: command.PassedFlags{
				"--key":    "mapkeyval",
				"--config": "passedconfigval",
				"--help":   "default",
			},
			expectedErr: false,
		},
		{
			name: "sectionFlag",
			app: warg.New(
				"newAppName",
				section.New(
					"help for test",
					section.Flag(
						"--sflag",
						"help for --sflag",
						scalar.String(
							scalar.Default("sflagval"),
						),
					),
					section.Command(
						"com",
						"help for com",
						command.DoNothing,
					),
				),
				warg.SkipValidation(),
			),
			args:               []string{"test", "com"},
			lookup:             warg.LookupMap(nil),
			expectedPassedPath: []string{"com"},
			expectedPassedFlagValues: command.PassedFlags{
				"--sflag": "sflagval",
				"--help":  "default",
			},
			expectedErr: false,
		},
		{
			name: "simpleJSONConfig",
			app: warg.New(
				"newAppName",
				section.New("help for test",
					section.Flag(
						"--val",
						"flag help",
						scalar.String(),
						flag.ConfigPath("params.val"),
					),
					section.Command(
						"com",
						"help for com",
						command.DoNothing,
					),
				),
				warg.ConfigFlag(
					"--config",
					[]scalar.ScalarOpt[string]{
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
			lookup:             warg.LookupMap(nil),
			expectedPassedPath: []string{"com"},
			expectedPassedFlagValues: command.PassedFlags{
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
				"newAppName",
				section.New("help for test",
					// TODO: looks like I don't have a float value type yet...
					// section.Flag(
					// 	"--floatval",
					// 	"flag help",
					// 	scalar.Float(),
					// 	flag.ConfigPath("params.floatval"),
					// ),
					section.Flag(
						"--intval",
						"flag help",
						scalar.Int(),
						flag.ConfigPath("params.intval"),
					),
					section.Command(
						"com",
						"help for com",
						command.DoNothing,
					),
				),
				warg.ConfigFlag(
					"--config",
					[]scalar.ScalarOpt[string]{
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
			lookup:             warg.LookupMap(nil),
			expectedPassedPath: []string{"com"},
			expectedPassedFlagValues: command.PassedFlags{
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
				"newAppName",
				section.New(
					"help for test",
					section.Flag(
						"--subreddits",
						"the subreddits",
						slice.String(),
						flag.ConfigPath("subreddits[].name"),
					),
					section.Command("print", "print key value", command.DoNothing),
				),
				warg.ConfigFlag(
					"--config",
					[]scalar.ScalarOpt[string]{
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
			lookup:             warg.LookupMap(nil),
			expectedPassedPath: []string{"print"},
			expectedPassedFlagValues: command.PassedFlags{
				"--subreddits": []string{"earthporn", "wallpapers"},
				"--config":     testDataFilePath(t.Name(), "configSlice", "config_slice.json"),
				"--help":       "default",
			},
			expectedErr: false,
		},
		{
			name: "envvar",
			app: warg.New(
				"newAppName",
				section.New(
					"help for test",
					section.Flag(
						"--flag",
						"help for --flag",
						scalar.String(),
						flag.EnvVars("notthere", "there", "alsothere"),
					),
					section.Command(
						"test",
						"blah",
						command.DoNothing,
					),
				),
				warg.SkipValidation(),
			),
			args: []string{t.Name(), "test"},
			lookup: warg.LookupMap(
				map[string]string{
					"there":     "there",
					"alsothere": "alsothere",
				},
			),
			expectedPassedPath: []string{"test"},
			expectedPassedFlagValues: command.PassedFlags{
				"--flag": "there",
				"--help": "default",
			},
			expectedErr: false,
		},
		{
			name: "requiredFlag",
			app: warg.New(
				"newAppName",
				section.New(
					"help for test",
					section.Flag(
						"--flag",
						"help for --flag",
						scalar.String(),
						flag.Required(),
					),
					section.Command(
						"test",
						"blah",
						command.DoNothing,
					),
				),
				warg.SkipValidation(),
			),
			args: []string{t.Name(), "test"},
			lookup: warg.LookupMap(
				nil,
			),
			expectedPassedPath:       []string{"test"},
			expectedPassedFlagValues: command.PassedFlags{},
			expectedErr:              true,
		},
		{
			name: "flagAlias",
			app: warg.New(
				"newAppName",
				section.New(
					"help for section",
					section.Flag(
						"--flag",
						"help for --flag",
						scalar.String(),
						flag.Alias("-f"),
					),
					section.Command(
						"test",
						"help for test",
						command.DoNothing,
					),
				),
				warg.SkipValidation(),
			),
			args:                     []string{t.Name(), "test", "-f", "val"},
			lookup:                   warg.LookupMap(nil),
			expectedPassedPath:       []string{"test"},
			expectedPassedFlagValues: command.PassedFlags{"--flag": "val", "--help": "default"},
			expectedErr:              false,
		},
		{
			name: "flagAliasWithList",
			app: warg.New(
				"newAppName",
				section.New(
					"help for section",
					section.Command(
						"test",
						"help for test",
						command.DoNothing,
						command.Flag(
							"--flag",
							"help for --flag",
							slice.String(),
							flag.Alias("-f"),
						),
					),
				),
				warg.SkipValidation(),
			),
			args:                     []string{t.Name(), "test", "-f", "1", "--flag", "2", "-f", "3", "--flag", "4"},
			lookup:                   warg.LookupMap(nil),
			expectedPassedPath:       []string{"test"},
			expectedPassedFlagValues: command.PassedFlags{"--flag": []string{"1", "2", "3", "4"}, "--help": "default"},
			expectedErr:              false,
		},
		{
			name: "badHelp",
			app: warg.New(
				"newAppName",
				section.New(
					"help for section",
					section.Command(
						"test",
						"help for test",
						command.DoNothing,
					),
				),
				warg.SkipValidation(),
			),
			args:                     []string{t.Name(), "test", "-h", "badhelpval"},
			lookup:                   warg.LookupMap(nil),
			expectedPassedPath:       nil,
			expectedPassedFlagValues: nil,
			expectedErr:              true,
		},
		{
			name: "addSectionFlags",
			app: func() warg.App {
				fm := flag.FlagMap{
					"--flag1": flag.New("--flag1 value", scalar.String()),
					"--flag2": flag.New("--flag1 value", scalar.String()),
				}
				app := warg.New(
					"newAppName",
					section.New(
						"help for section",
						section.ExistingFlags(fm),
						section.Command(
							"test",
							"help for test",
							command.DoNothing,
						),
					),
					warg.SkipValidation(),
				)
				return app
			}(),

			args:                     []string{t.Name(), "test", "--flag1", "val1"},
			lookup:                   warg.LookupMap(nil),
			expectedPassedPath:       []string{"test"},
			expectedPassedFlagValues: command.PassedFlags{"--flag1": "val1", "--help": "default"},
			expectedErr:              false,
		},
		{
			name: "addCommandFlags",
			app: func() warg.App {
				fm := flag.FlagMap{
					"--flag1": flag.New("--flag1 value", scalar.String()),
					"--flag2": flag.New("--flag1 value", scalar.String()),
				}
				app := warg.New(
					"newAppName",
					section.New(
						"help for section",
						section.Command(
							"test",
							"help for test",
							command.DoNothing,
							command.ExistingFlags(fm),
						),
					),
					warg.SkipValidation(),
				)
				return app
			}(),

			args:                     []string{t.Name(), "test", "--flag1", "val1"},
			lookup:                   warg.LookupMap(nil),
			expectedPassedPath:       []string{"test"},
			expectedPassedFlagValues: command.PassedFlags{"--flag1": "val1", "--help": "default"},
			expectedErr:              false,
		},
		{
			name: "JSONConfigStringSlice",
			app: warg.New(
				"newAppName",
				section.New("help for test",
					section.Flag(
						"--val",
						"flag help",
						slice.String(),
						flag.ConfigPath("val"),
					),
					section.Command(
						"com",
						"help for com",
						command.DoNothing,
					),
				),
				warg.ConfigFlag(
					"--config",
					[]scalar.ScalarOpt[string]{
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
			lookup:             warg.LookupMap(nil),
			expectedPassedPath: []string{"com"},
			expectedPassedFlagValues: command.PassedFlags{
				"--config": testDataFilePath(t.Name(), "JSONConfigStringSlice", "config.json"),
				"--val":    []string{"from", "config"},
				"--help":   "default",
			},
			expectedErr: false,
		},
		{
			name: "YAMLConfigStringSlice",
			app: warg.New(
				"newAppName",
				section.New("help for test",
					section.Flag(
						"--val",
						"flag help",
						slice.String(),
						flag.ConfigPath("val"),
					),
					section.Command(
						"com",
						"help for com",
						command.DoNothing,
					),
				),
				warg.ConfigFlag(
					"--config",
					[]scalar.ScalarOpt[string]{
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
			lookup:             warg.LookupMap(nil),
			expectedPassedPath: []string{"com"},
			expectedPassedFlagValues: command.PassedFlags{
				"--config": testDataFilePath(t.Name(), "YAMLConfigStringSlice", "config.yaml"),
				"--val":    []string{"from", "config"},
				"--help":   "default",
			},
			expectedErr: false,
		},
		// Note: Will need to update this if https://github.com/bbkane/warg/issues/36 gets implemented
		{
			name: "invalidFlagsErrorEvenForHelp",
			app: warg.New(
				"invalidFlagsErrorEvenForHelp",
				section.New(
					section.HelpShort("A virtual assistant"),
					section.Command(
						"present",
						"Formally present a guest (guests are never introduced, always presented).",
						command.DoNothing,
						command.Flag(
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
			lookup:                   warg.LookupMap(map[string]string{"USER": "bbkane"}),
			expectedPassedPath:       []string{"present"},
			expectedPassedFlagValues: command.PassedFlags{"--help": "default"},
			expectedErr:              true,
		},
		{
			name: "JSONConfigMap",
			app: warg.New(
				"newAppName",
				section.New("help for test",
					section.Flag(
						"--val",
						"flag help",
						dict.Int(),
						flag.ConfigPath("val"),
					),
					section.Command(
						"com",
						"help for com",
						command.DoNothing,
					),
				),
				warg.ConfigFlag(
					"--config",
					[]scalar.ScalarOpt[string]{
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
			lookup:             warg.LookupMap(nil),
			expectedPassedPath: []string{"com"},
			expectedPassedFlagValues: command.PassedFlags{
				"--config": testDataFilePath(t.Name(), "JSONConfigMap", "config.json"),
				"--val":    map[string]int{"a": 1},
				"--help":   "default",
			},
			expectedErr: false,
		},
		{
			name: "dictUpdate",
			app: warg.New(
				"newAppName",
				section.New(
					"help for test",
					section.Flag(
						flag.Name("--flag"),
						"flag help",
						dict.Bool(),
					),

					section.Command(
						"com1",
						"help for com1",
						command.DoNothing,
					),
				),
				warg.SkipValidation(),
			),

			args:                     []string{"app", "com1", "--flag", "true=true", "--flag", "false=false"},
			lookup:                   warg.LookupMap(nil),
			expectedPassedPath:       []string{"com1"},
			expectedPassedFlagValues: command.PassedFlags{"--flag": map[string]bool{"true": true, "false": false}, "--help": "default"},
			expectedErr:              false,
		},
		{
			name: "unsetSentinelScalarSuccess",
			app: warg.New(
				"newAppName",
				section.New(
					"help for test",
					section.Flag(
						"--flag",
						"help for --flag",
						scalar.String(scalar.Default("default")),
						flag.UnsetSentinel("UNSET"),
					),
					section.Command(
						"test",
						"help for test",
						command.DoNothing,
					),
				),
				warg.SkipValidation(),
			),
			args:               []string{t.Name(), "test", "--flag", "UNSET"},
			lookup:             warg.LookupMap(nil),
			expectedPassedPath: []string{"test"},
			expectedPassedFlagValues: command.PassedFlags{
				"--help": "default",
			},
			expectedErr: false,
		},
		{
			name: "unsetSentinelScalarError",
			app: warg.New(
				"newAppName",
				section.New(
					"help for test",
					section.Flag(
						"--flag",
						"help for --flag",
						scalar.String(scalar.Default("default")),
						flag.UnsetSentinel("UNSET"),
					),
					section.Command(
						"test",
						"help for test",
						command.DoNothing,
					),
				),
				warg.SkipValidation(),
			),
			args:                     []string{t.Name(), "test", "--flag", "UNSET", "--flag", "justsayno"},
			lookup:                   warg.LookupMap(nil),
			expectedPassedPath:       []string{"test"},
			expectedPassedFlagValues: nil,
			expectedErr:              true,
		},
		{
			name: "unsetSentinelSlice",
			app: warg.New(
				"newAppName",
				section.New(
					"help for test",
					section.Flag(
						"--flag",
						"help for --flag",
						slice.String(slice.Default([]string{"default"})),
						flag.UnsetSentinel("UNSET"),
					),
					section.Command(
						"test",
						"help for test",
						command.DoNothing,
					),
				),
				warg.SkipValidation(),
			),
			args:               []string{t.Name(), "test", "--flag", "a", "--flag", "UNSET", "--flag", "b", "--flag", "c"},
			lookup:             warg.LookupMap(nil),
			expectedPassedPath: []string{"test"},
			expectedPassedFlagValues: command.PassedFlags{
				"--help": "default",
				"--flag": []string{"b", "c"},
			},
			expectedErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := tt.app.Validate()
			require.Nil(t, err)

			actualPR, actualErr := tt.app.Parse(tt.args, tt.lookup)

			if tt.expectedErr {
				require.NotNil(t, actualErr)
				return
			} else {
				require.Nil(t, actualErr)
			}
			require.Equal(t, tt.expectedPassedPath, actualPR.Path)
			require.Equal(t, tt.expectedPassedFlagValues, actualPR.Context.Flags)
		})
	}
}

func TestContextVersion(t *testing.T) {
	app := warg.New(
		"appName",
		section.New(
			"test",
			section.Command("version", "Print version", command.PrintVersion),
		),
		warg.OverrideVersion("customversion"),
	)
	err := app.Validate()
	require.Nil(t, err)

	actualPR, err := app.Parse([]string{"appName"}, warg.LookupMap(nil))
	require.Nil(t, err)

	expectedVersion := "customversion"
	require.Equal(t, expectedVersion, actualPR.Context.Version)
}
