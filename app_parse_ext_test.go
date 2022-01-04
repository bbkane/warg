package warg_test

// external tests - import warg like it's an external package

import (
	"path/filepath"
	"testing"

	"github.com/bbkane/warg"
	"github.com/bbkane/warg/command"
	"github.com/bbkane/warg/config"
	"github.com/bbkane/warg/config/jsonreader"
	"github.com/bbkane/warg/config/yamlreader"
	"github.com/bbkane/warg/flag"
	"github.com/bbkane/warg/section"
	"github.com/bbkane/warg/value"

	"github.com/stretchr/testify/require"
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
		expectedPassedFlagValues flag.PassedFlags
		expectedErr              bool
	}{
		{
			name: "fromMain",
			app: warg.New(
				"test",
				section.New(
					"help for test",
					section.Flag(
						"--af1",
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
			),

			args:                     []string{"app", "cat1", "com1", "--com1f1", "1"},
			lookup:                   warg.LookupMap(nil),
			expectedPassedPath:       []string{"cat1", "com1"},
			expectedPassedFlagValues: flag.PassedFlags{"--com1f1": int(1), "--help": "default"},
			expectedErr:              false,
		},
		{
			name: "noSection",
			app: warg.New(
				"test",
				section.New(
					"help for test",
					section.Flag(
						"--af1",
						"flag help",
						value.Int,
					),
				),
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
				"test",
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
			),
			args:                     []string{"test", "com"},
			lookup:                   warg.LookupMap(nil),
			expectedPassedPath:       []string{"com"},
			expectedPassedFlagValues: flag.PassedFlags{"--flag": "hi", "--help": "default"},
			expectedErr:              false,
		},
		{
			name: "extraFlag",
			app: warg.New(
				"test",
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
				"test",
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
			),
			args:               []string{"test", "print", "--config", "passedconfigval"},
			lookup:             warg.LookupMap(nil),
			expectedPassedPath: []string{"print"},
			expectedPassedFlagValues: flag.PassedFlags{
				"--key":    "mapkeyval",
				"--config": "passedconfigval",
				"--help":   "default",
			},
			expectedErr: false,
		},
		{
			name: "sectionFlag",
			app: warg.New(
				"test",
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
			),
			args:               []string{"test", "com"},
			lookup:             warg.LookupMap(nil),
			expectedPassedPath: []string{"com"},
			expectedPassedFlagValues: flag.PassedFlags{
				"--sflag": "sflagval",
				"--help":  "default",
			},
			expectedErr: false,
		},
		{
			name: "simpleJSONConfig",
			app: warg.New(
				"test",
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
					// TODO: make this test work by following the config cases
					// in the README
					flag.Default("testdata/simple_json_config.json"),
				),
			),

			args:               []string{"app", "com"},
			lookup:             warg.LookupMap(nil),
			expectedPassedPath: []string{"com"},
			expectedPassedFlagValues: flag.PassedFlags{
				"--config": "testdata/simple_json_config.json",
				"--val":    "hi",
				"--help":   "default",
			},
			expectedErr: false,
		},
		{
			name: "configSlice",
			app: warg.New(
				"test",
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
					flag.Default("testdata/config_slice.json"),
				),
			),
			args:               []string{"test", "print"},
			lookup:             warg.LookupMap(nil),
			expectedPassedPath: []string{"print"},
			expectedPassedFlagValues: flag.PassedFlags{
				"--subreddits": []string{"earthporn", "wallpapers"},
				"--config":     "testdata/config_slice.json",
				"--help":       "default",
			},
			expectedErr: false,
		},
		{
			name: "envvar",
			app: warg.New(
				t.Name(),
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
			),
			args: []string{t.Name(), "test"},
			lookup: warg.LookupMap(
				map[string]string{
					"there":     "there",
					"alsothere": "alsothere",
				},
			),
			expectedPassedPath: []string{"test"},
			expectedPassedFlagValues: flag.PassedFlags{
				"--flag": "there",
				"--help": "default",
			},
			expectedErr: false,
		},
		{
			name: "requiredFlag",
			app: warg.New(
				t.Name(),
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
			),
			args: []string{t.Name(), "test"},
			lookup: warg.LookupMap(
				nil,
			),
			expectedPassedPath:       []string{"test"},
			expectedPassedFlagValues: flag.PassedFlags{},
			expectedErr:              true,
		},
		{
			name: "flagAlias",
			app: warg.New(
				t.Name(),
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
			),
			args:                     []string{t.Name(), "test", "-f", "val"},
			lookup:                   warg.LookupMap(nil),
			expectedPassedPath:       []string{"test"},
			expectedPassedFlagValues: flag.PassedFlags{"--flag": "val", "--help": "default"},
			expectedErr:              false,
		},
		{
			name: "flagAliasWithList",
			app: warg.New(
				t.Name(),
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
			),
			args:                     []string{t.Name(), "test", "-f", "1", "--flag", "2", "-f", "3", "--flag", "4"},
			lookup:                   warg.LookupMap(nil),
			expectedPassedPath:       []string{"test"},
			expectedPassedFlagValues: flag.PassedFlags{"--flag": []string{"1", "2", "3", "4"}, "--help": "default"},
			expectedErr:              false,
		},
		{
			name: "badHelp",
			app: warg.New(
				t.Name(),
				section.New(
					"help for section",
					section.Command(
						"test",
						"help for test",
						command.DoNothing,
					),
				),
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
				fm := map[string]flag.Flag{
					"--flag1": flag.New("--flag1 value", value.String),
					"--flag2": flag.New("--flag1 value", value.String),
				}
				app := warg.New(
					t.Name(),
					section.New(
						"help for section",
						section.ExistingFlags(fm),
						section.Command(
							"test",
							"help for test",
							command.DoNothing,
						),
					),
				)
				return app
			}(),

			args:                     []string{t.Name(), "test", "--flag1", "val1"},
			lookup:                   warg.LookupMap(nil),
			expectedPassedPath:       []string{"test"},
			expectedPassedFlagValues: flag.PassedFlags{"--flag1": "val1", "--help": "default"},
			expectedErr:              false,
		},
		{
			name: "addCommandFlags",
			app: func() warg.App {
				fm := map[string]flag.Flag{
					"--flag1": flag.New("--flag1 value", value.String),
					"--flag2": flag.New("--flag1 value", value.String),
				}
				app := warg.New(
					t.Name(),
					section.New(
						"help for section",
						section.Command(
							"test",
							"help for test",
							command.DoNothing,
							command.ExistingFlags(fm),
						),
					),
				)
				return app
			}(),

			args:                     []string{t.Name(), "test", "--flag1", "val1"},
			lookup:                   warg.LookupMap(nil),
			expectedPassedPath:       []string{"test"},
			expectedPassedFlagValues: flag.PassedFlags{"--flag1": "val1", "--help": "default"},
			expectedErr:              false,
		},
		{
			name: "JSONConfigStringSlice",
			app: warg.New(
				"test",
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
			),

			args:               []string{"app", "com"},
			lookup:             warg.LookupMap(nil),
			expectedPassedPath: []string{"com"},
			expectedPassedFlagValues: flag.PassedFlags{
				"--config": testDataFilePath(t.Name(), "JSONConfigStringSlice", "config.json"),
				"--val":    []string{"from", "config"},
				"--help":   "default",
			},
			expectedErr: false,
		},
		{
			name: "YAMLConfigStringSlice",
			app: warg.New(
				"test",
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
			),

			args:               []string{"app", "com"},
			lookup:             warg.LookupMap(nil),
			expectedPassedPath: []string{"com"},
			expectedPassedFlagValues: flag.PassedFlags{
				"--config": testDataFilePath(t.Name(), "YAMLConfigStringSlice", "config.yaml"),
				"--val":    []string{"from", "config"},
				"--help":   "default",
			},
			expectedErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			actualPR, actualErr := tt.app.Parse(tt.args, tt.lookup)

			if tt.expectedErr {
				require.NotNil(t, actualErr)
				return
			} else {
				require.Nil(t, actualErr)
			}

			require.Equal(t, tt.expectedPassedPath, actualPR.Path)
			require.Equal(t, tt.expectedPassedFlagValues, actualPR.PassedFlags)
		})
	}
}
