package warg_test

// external tests - import warg like it's an external package

import (
	"testing"

	w "github.com/bbkane/warg"
	c "github.com/bbkane/warg/command"
	"github.com/bbkane/warg/configreader"
	"github.com/bbkane/warg/configreader/jsonreader"
	f "github.com/bbkane/warg/flag"
	s "github.com/bbkane/warg/section"
	v "github.com/bbkane/warg/value"

	"github.com/stretchr/testify/require"
)

func TestApp_Parse(t *testing.T) {

	tests := []struct {
		name                     string
		app                      w.App
		args                     []string
		expectedPassedPath       []string
		expectedPassedFlagValues f.FlagValues
		expectedErr              bool
	}{
		{
			name: "from main",
			app: w.New(
				"test",
				s.NewSection(
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
			expectedPassedPath:       []string{"cat1", "com1"},
			expectedPassedFlagValues: f.FlagValues{"--com1f1": int(1)},
			expectedErr:              false,
		},
		{
			name: "no section",
			app: w.New(
				"test",
				s.NewSection(
					"help for test",
					s.WithFlag(
						"--af1",
						"flag help",
						v.Int,
					),
				),
			),

			args:                     []string{"app"},
			expectedPassedPath:       nil,
			expectedPassedFlagValues: nil,
			expectedErr:              false,
		},
		{
			name: "flag default",
			app: w.New(
				"test",
				s.NewSection(
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
			expectedPassedPath:       []string{"com"},
			expectedPassedFlagValues: f.FlagValues{"--flag": "hi"},
			expectedErr:              false,
		},
		{
			name: "extra_flag",
			app: w.New(
				"test",
				s.NewSection(
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
			expectedPassedPath:       nil,
			expectedPassedFlagValues: nil,
			expectedErr:              true,
		},
		{
			name: "config_flag",
			app: w.New(
				"test",
				s.NewSection(
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
				w.ConfigFlag(
					"--config",
					func(_ string) (configreader.ConfigReader, error) {

						var cr configreader.ConfigReaderFunc = func(path string) (configreader.ConfigSearchResult, error) {
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
			expectedPassedPath: []string{"print"},
			expectedPassedFlagValues: f.FlagValues{
				"--key":    "mapkeyval",
				"--config": "passedconfigval",
			},
			expectedErr: false,
		},
		{
			name: "section flag",
			app: w.New(
				"test",
				s.NewSection(
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
			expectedPassedPath: []string{"com"},
			expectedPassedFlagValues: f.FlagValues{
				"--sflag": "sflagval",
			},
			expectedErr: false,
		},
		{
			name: "simple JSON config",
			app: w.New(
				"test",
				s.NewSection("help for test",
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
				w.ConfigFlag(
					"--config",
					jsonreader.NewJSONConfigReader,
					"path to config",
					// TODO: make this test work by following the config cases
					// in the README
					f.Default("testdata/simple_json_config.json"),
				),
			),

			args:               []string{"app", "com"},
			expectedPassedPath: []string{"com"},
			expectedPassedFlagValues: f.FlagValues{
				"--config": "testdata/simple_json_config.json",
				"--val":    "hi",
			},
			expectedErr: false,
		},
		{
			name: "config_slice",
			app: w.New(
				"test",
				s.NewSection(
					"help for test",
					s.WithFlag(
						"--subreddits",
						"the subreddits",
						v.StringSlice,
						f.ConfigPath("subreddits[].name"),
					),
					s.WithCommand("print", "print key value", c.DoNothing),
				),
				w.ConfigFlag(
					"--config",
					jsonreader.NewJSONConfigReader,
					"config flag",
					f.Default("testdata/config_slice.json"),
				),
			),
			args:               []string{"test", "print"},
			expectedPassedPath: []string{"print"},
			expectedPassedFlagValues: f.FlagValues{
				"--subreddits": []string{"earthporn", "wallpapers"},
				"--config":     "testdata/config_slice.json",
			},
			expectedErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			actualPR, actualErr := tt.app.Parse(tt.args)

			if tt.expectedErr {
				require.NotNil(t, actualErr)
				return
			} else {
				require.Nil(t, actualErr)
			}

			require.Equal(t, tt.expectedPassedPath, actualPR.PasssedPath)
			require.Equal(t, tt.expectedPassedFlagValues, actualPR.PassedFlags)
		})
	}
}
