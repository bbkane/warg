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
		lookup                   f.LookupFunc
		expectedPassedPath       []string
		expectedPassedFlagValues f.PassedFlags
		expectedErr              bool
	}{
		{
			name: "from main",
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
			lookup:                   warg.DictLookup(nil),
			expectedPassedPath:       []string{"cat1", "com1"},
			expectedPassedFlagValues: f.PassedFlags{"--com1f1": int(1)},
			expectedErr:              false,
		},
		{
			name: "no section",
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
			lookup:                   warg.DictLookup(nil),
			expectedPassedPath:       nil,
			expectedPassedFlagValues: nil,
			expectedErr:              false,
		},
		{
			name: "flag default",
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
			lookup:                   warg.DictLookup(nil),
			expectedPassedPath:       []string{"com"},
			expectedPassedFlagValues: f.PassedFlags{"--flag": "hi"},
			expectedErr:              false,
		},
		{
			name: "extra_flag",
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
			lookup:                   warg.DictLookup(nil),
			expectedPassedPath:       nil,
			expectedPassedFlagValues: nil,
			expectedErr:              true,
		},
		{
			name: "config_flag",
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
			lookup:             warg.DictLookup(nil),
			expectedPassedPath: []string{"print"},
			expectedPassedFlagValues: f.PassedFlags{
				"--key":    "mapkeyval",
				"--config": "passedconfigval",
			},
			expectedErr: false,
		},
		{
			name: "section flag",
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
			lookup:             warg.DictLookup(nil),
			expectedPassedPath: []string{"com"},
			expectedPassedFlagValues: f.PassedFlags{
				"--sflag": "sflagval",
			},
			expectedErr: false,
		},
		{
			name: "simple JSON config",
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
			lookup:             warg.DictLookup(nil),
			expectedPassedPath: []string{"com"},
			expectedPassedFlagValues: f.PassedFlags{
				"--config": "testdata/simple_json_config.json",
				"--val":    "hi",
			},
			expectedErr: false,
		},
		{
			name: "config_slice",
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
			lookup:             warg.DictLookup(nil),
			expectedPassedPath: []string{"print"},
			expectedPassedFlagValues: f.PassedFlags{
				"--subreddits": []string{"earthporn", "wallpapers"},
				"--config":     "testdata/config_slice.json",
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
