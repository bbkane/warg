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
	"go.bbkane.com/warg/value"

	"github.com/alecthomas/assert"
)

// NOTE: this is is a bit of a hack to mock out a configreader
// NOTE: see https://karthikkaranth.me/blog/functions-implementing-interfaces-in-go/
// for how to use ConfigReaderFunc in tests
type ConfigReaderFunc func(path string) (config.SearchResult, error)

func (f ConfigReaderFunc) Search(path string) (config.SearchResult, error) {
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
						value.Int,
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
								value.Int,
								flag.Default("10"),
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
						value.Int,
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
							value.String,
							flag.Default("hi"),
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
							value.String,
							flag.Default("hi"),
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
						value.String,
						flag.ConfigPath("key"),
						flag.Default("defaultkeyval"),
					),
					section.Command("print", "print key value", command.DoNothing),
				),
				warg.ConfigFlag(
					"--config",
					func(_ string) (config.Reader, error) {
						var cr ConfigReaderFunc = func(path string) (config.SearchResult, error) {
							if path == "key" {
								return config.SearchResult{
									IFace:        "mapkeyval",
									Exists:       true,
									IsAggregated: false,
								}, nil
							}
							return config.SearchResult{}, nil
						}

						return cr, nil
					},
					"config flag",
					flag.Default("defaultconfigval"),
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
						value.String,
						flag.Default("sflagval"),
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
						value.String,
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
					jsonreader.New,
					"path to config",
					flag.Default(testDataFilePath(t.Name(), "simpleJSONConfig", "simple_json_config.json")),
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
			name: "configSlice",
			app: warg.New(
				"newAppName",
				section.New(
					"help for test",
					section.Flag(
						"--subreddits",
						"the subreddits",
						value.StringSlice,
						flag.ConfigPath("subreddits[].name"),
					),
					section.Command("print", "print key value", command.DoNothing),
				),
				warg.ConfigFlag(
					"--config",
					jsonreader.New,
					"config flag",
					flag.Default(testDataFilePath(t.Name(), "configSlice", "config_slice.json")),
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
						value.String,
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
						value.String,
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
						value.String,
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
							value.StringSlice,
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
					"--flag1": flag.New("--flag1 value", value.String),
					"--flag2": flag.New("--flag1 value", value.String),
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
					"--flag1": flag.New("--flag1 value", value.String),
					"--flag2": flag.New("--flag1 value", value.String),
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
						value.StringSlice,
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
					jsonreader.New,
					"path to config",
					flag.Default(testDataFilePath(t.Name(), "JSONConfigStringSlice", "config.json")),
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
						value.StringSlice,
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
					yamlreader.New,
					"path to config",
					flag.Default(testDataFilePath(t.Name(), "YAMLConfigStringSlice", "config.yaml")),
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := tt.app.Validate()
			assert.Nil(t, err)

			actualPR, actualErr := tt.app.Parse(tt.args, tt.lookup)

			if tt.expectedErr {
				assert.NotNil(t, actualErr)
				return
			} else {
				assert.Nil(t, actualErr)
			}

			// TODO: I wish I'd made this compare a parse result with a parse result instead of field by field
			assert.Equal(t, tt.expectedPassedPath, actualPR.Path)
			assert.Equal(t, tt.expectedPassedFlagValues, actualPR.Context.Flags)
		})
	}
}
