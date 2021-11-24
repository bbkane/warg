package warg_test

// external tests - import warg like it's an external package

import (
	"testing"

	"github.com/bbkane/warg"
	c "github.com/bbkane/warg/command"
	"github.com/bbkane/warg/configreader"
	"github.com/bbkane/warg/configreader/jsonreader"
	f "github.com/bbkane/warg/flag"
	s "github.com/bbkane/warg/section"
	v "github.com/bbkane/warg/value"

	"github.com/stretchr/testify/require"
)

// NOTE: this is is a bit of a hack to mock out a configreader
// NOTE: see https://karthikkaranth.me/blog/functions-implementing-interfaces-in-go/
// for how to use ConfigReaderFunc in tests
type ConfigReaderFunc func(path string) (configreader.ConfigSearchResult, error)

func (f ConfigReaderFunc) Search(path string) (configreader.ConfigSearchResult, error) {
	return f(path)
}

func TestApp_Parse(t *testing.T) {

	tests := []struct {
		name                     string
		app                      warg.App
		args                     []string
		lookup                   warg.LookupFunc
		expectedPassedPath       []string
		expectedPassedFlagValues f.PassedFlags
		expectedErr              bool
	}{
		{
			name: "fromMain",
			app: warg.New(
				"test",
				s.New(
					"help for test",
					s.WithFlag(
						"--af1",
						"flag help",
						v.Int,
					),
					s.WithSection(
						"cat1",
						"help for cat1",
						s.WithCommand(
							"com1",
							"help for com1",
							c.DoNothing,
							c.WithFlag(
								"--com1f1",
								"flag help",
								v.Int,
								f.Default("10"),
							),
						),
					),
				),
			),

			args:                     []string{"app", "cat1", "com1", "--com1f1", "1"},
			lookup:                   warg.LookupMap(nil),
			expectedPassedPath:       []string{"cat1", "com1"},
			expectedPassedFlagValues: f.PassedFlags{"--com1f1": int(1), "--help": "default"},
			expectedErr:              false,
		},
		{
			name: "noSection",
			app: warg.New(
				"test",
				s.New(
					"help for test",
					s.WithFlag(
						"--af1",
						"flag help",
						v.Int,
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
				s.New(
					"help for test",
					s.WithCommand(
						"com",
						"com help",
						c.DoNothing,
						c.WithFlag(
							"--flag",
							"flag help",
							v.String,
							f.Default("hi"),
						),
					),
				),
			),
			args:                     []string{"test", "com"},
			lookup:                   warg.LookupMap(nil),
			expectedPassedPath:       []string{"com"},
			expectedPassedFlagValues: f.PassedFlags{"--flag": "hi", "--help": "default"},
			expectedErr:              false,
		},
		{
			name: "extraFlag",
			app: warg.New(
				"test",
				s.New(
					"help for test",
					s.WithCommand(
						"com",
						"com help",
						c.DoNothing,
						c.WithFlag(
							"--flag",
							"flag help",
							v.String,
							f.Default("hi"),
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
				s.New(
					"help for test",
					s.WithFlag(
						"--key",
						"a key",
						v.String,
						f.ConfigPath("key"),
						f.Default("defaultkeyval"),
					),
					s.WithCommand("print", "print key value", c.DoNothing),
				),
				warg.ConfigFlag(
					"--config",
					func(_ string) (configreader.ConfigReader, error) {
						var cr ConfigReaderFunc = func(path string) (configreader.ConfigSearchResult, error) {
							if path == "key" {
								return configreader.ConfigSearchResult{
									IFace:        "mapkeyval",
									Exists:       true,
									IsAggregated: false,
								}, nil
							}
							return configreader.ConfigSearchResult{}, nil
						}

						return cr, nil
					},
					"config flag",
					f.Default("defaultconfigval"),
				),
			),
			args:               []string{"test", "print", "--config", "passedconfigval"},
			lookup:             warg.LookupMap(nil),
			expectedPassedPath: []string{"print"},
			expectedPassedFlagValues: f.PassedFlags{
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
				s.New(
					"help for test",
					s.WithFlag(
						"--sflag",
						"help for --sflag",
						v.String,
						f.Default("sflagval"),
					),
					s.WithCommand(
						"com",
						"help for com",
						c.DoNothing,
					),
				),
			),
			args:               []string{"test", "com"},
			lookup:             warg.LookupMap(nil),
			expectedPassedPath: []string{"com"},
			expectedPassedFlagValues: f.PassedFlags{
				"--sflag": "sflagval",
				"--help":  "default",
			},
			expectedErr: false,
		},
		{
			name: "simpleJSONConfig",
			app: warg.New(
				"test",
				s.New("help for test",
					s.WithFlag(
						"--val",
						"flag help",
						v.String,
						f.ConfigPath("params.val"),
					),
					s.WithCommand(
						"com",
						"help for com",
						c.DoNothing,
					),
				),
				warg.ConfigFlag(
					"--config",
					jsonreader.NewJSONConfigReader,
					"path to config",
					// TODO: make this test work by following the config cases
					// in the README
					f.Default("testdata/simple_json_config.json"),
				),
			),

			args:               []string{"app", "com"},
			lookup:             warg.LookupMap(nil),
			expectedPassedPath: []string{"com"},
			expectedPassedFlagValues: f.PassedFlags{
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
				s.New(
					"help for test",
					s.WithFlag(
						"--subreddits",
						"the subreddits",
						v.StringSlice,
						f.ConfigPath("subreddits[].name"),
					),
					s.WithCommand("print", "print key value", c.DoNothing),
				),
				warg.ConfigFlag(
					"--config",
					jsonreader.NewJSONConfigReader,
					"config flag",
					f.Default("testdata/config_slice.json"),
				),
			),
			args:               []string{"test", "print"},
			lookup:             warg.LookupMap(nil),
			expectedPassedPath: []string{"print"},
			expectedPassedFlagValues: f.PassedFlags{
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
				s.New(
					"help for test",
					s.WithFlag(
						"--flag",
						"help for --flag",
						v.String,
						f.EnvVars("notthere", "there", "alsothere"),
					),
					s.WithCommand(
						"test",
						"blah",
						c.DoNothing,
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
			expectedPassedFlagValues: f.PassedFlags{
				"--flag": "there",
				"--help": "default",
			},
			expectedErr: false,
		},
		{
			name: "requiredFlag",
			app: warg.New(
				t.Name(),
				s.New(
					"help for test",
					s.WithFlag(
						"--flag",
						"help for --flag",
						v.String,
						f.Required(),
					),
					s.WithCommand(
						"test",
						"blah",
						c.DoNothing,
					),
				),
			),
			args: []string{t.Name(), "test"},
			lookup: warg.LookupMap(
				nil,
			),
			expectedPassedPath:       []string{"test"},
			expectedPassedFlagValues: f.PassedFlags{},
			expectedErr:              true,
		},
		{
			name: "flagAlias",
			app: warg.New(
				t.Name(),
				s.New(
					"help for section",
					s.WithFlag(
						"--flag",
						"help for --flag",
						v.String,
						f.Alias("-f"),
					),
					s.WithCommand(
						"test",
						"help for test",
						c.DoNothing,
					),
				),
			),
			args:                     []string{t.Name(), "test", "-f", "val"},
			lookup:                   warg.LookupMap(nil),
			expectedPassedPath:       []string{"test"},
			expectedPassedFlagValues: f.PassedFlags{"--flag": "val", "--help": "default"},
			expectedErr:              false,
		},
		{
			name: "flagAliasWithList",
			app: warg.New(
				t.Name(),
				s.New(
					"help for section",
					s.WithCommand(
						"test",
						"help for test",
						c.DoNothing,
						c.WithFlag(
							"--flag",
							"help for --flag",
							v.StringSlice,
							f.Alias("-f"),
						),
					),
				),
			),
			args:                     []string{t.Name(), "test", "-f", "1", "--flag", "2", "-f", "3", "--flag", "4"},
			lookup:                   warg.LookupMap(nil),
			expectedPassedPath:       []string{"test"},
			expectedPassedFlagValues: f.PassedFlags{"--flag": []string{"1", "2", "3", "4"}, "--help": "default"},
			expectedErr:              false,
		},
		{
			name: "badHelp",
			app: warg.New(
				t.Name(),
				s.New(
					"help for section",
					s.WithCommand(
						"test",
						"help for test",
						c.DoNothing,
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
				fm := map[string]f.Flag{
					"--flag1": f.New("--flag1 value", v.String),
					"--flag2": f.New("--flag1 value", v.String),
				}
				app := warg.New(
					t.Name(),
					s.New(
						"help for section",
						s.AddFlags(fm),
						s.WithCommand(
							"test",
							"help for test",
							c.DoNothing,
						),
					),
				)
				return app
			}(),

			args:                     []string{t.Name(), "test", "--flag1", "val1"},
			lookup:                   warg.LookupMap(nil),
			expectedPassedPath:       []string{"test"},
			expectedPassedFlagValues: f.PassedFlags{"--flag1": "val1", "--help": "default"},
			expectedErr:              false,
		},
		{
			name: "addCommandFlags",
			app: func() warg.App {
				fm := map[string]f.Flag{
					"--flag1": f.New("--flag1 value", v.String),
					"--flag2": f.New("--flag1 value", v.String),
				}
				app := warg.New(
					t.Name(),
					s.New(
						"help for section",
						s.WithCommand(
							"test",
							"help for test",
							c.DoNothing,
							c.AddFlags(fm),
						),
					),
				)
				return app
			}(),

			args:                     []string{t.Name(), "test", "--flag1", "val1"},
			lookup:                   warg.LookupMap(nil),
			expectedPassedPath:       []string{"test"},
			expectedPassedFlagValues: f.PassedFlags{"--flag1": "val1", "--help": "default"},
			expectedErr:              false,
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
