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

	"github.com/stretchr/testify/assert"
)

func TestApp_Parse(t *testing.T) {

	tests := []struct {
		name                 string
		app                  w.App
		args                 []string
		passedPathWant       []string
		passedFlagValuesWant f.FlagValues
		wantErr              bool
	}{
		{
			name: "from main",
			app: w.New(
				"test",
				"v0.0.0",
				s.NewSection(
					"help for test",
					s.WithFlag(
						"--af1",
						"flag help",
						v.IntEmpty,
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
								v.IntEmpty,
								f.Default("10"),
							),
						),
					),
				),
			),

			args:                 []string{"app", "cat1", "com1", "--com1f1", "1"},
			passedPathWant:       []string{"cat1", "com1"},
			passedFlagValuesWant: f.FlagValues{"--com1f1": int(1)},
			wantErr:              false,
		},
		{
			name: "no section",
			app: w.New(
				"test",
				"v0.0.0",
				s.NewSection(
					"help for test",
					s.WithFlag(
						"--af1",
						"flag help",
						v.IntEmpty,
					),
				),
			),

			args:                 []string{"app"},
			passedPathWant:       nil,
			passedFlagValuesWant: nil,
			wantErr:              false,
		},
		{
			name: "flag default",
			app: w.New(
				"test",
				"v0.0.0",
				s.NewSection(
					"help for test",
					s.WithCommand(
						"com",
						"com help",
						c.DoNothing,
						c.WithFlag(
							"--flag",
							"flag help",
							v.StringEmpty,
							f.Default("hi"),
						),
					),
				),
			),
			args:                 []string{"test", "com"},
			passedPathWant:       []string{"com"},
			passedFlagValuesWant: f.FlagValues{"--flag": "hi"},
			wantErr:              false,
		},
		{
			name: "extra flag",
			app: w.New(
				"test",
				"v0.0.0",
				s.NewSection(
					"help for test",
					s.WithCommand(
						"com",
						"com help",
						c.DoNothing,
						c.WithFlag(
							"--flag",
							"flag help",
							v.StringEmpty,
							f.Default("hi"),
						),
					),
				),
			),
			args:                 []string{"test", "com", "--unexpected", "value"},
			passedPathWant:       nil,
			passedFlagValuesWant: nil,
			wantErr:              true,
		},
		{
			name: "config_flag",
			app: w.New(
				"test",
				"v0.0.0",
				s.NewSection(
					"help for test",
					s.WithFlag(
						"--key",
						"a key",
						v.StringEmpty,
						f.ConfigPath("key", v.StringFromInterface),
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
			args:           []string{"test", "print", "--config", "passedconfigval"},
			passedPathWant: []string{"print"},
			passedFlagValuesWant: f.FlagValues{
				"--key":    "mapkeyval",
				"--config": "passedconfigval",
			},
			wantErr: false,
		},
		{
			name: "section flag",
			app: w.New(
				"test",
				"v0.0.0",
				s.NewSection(
					"help for test",
					s.WithFlag(
						"--sflag",
						"help for --sflag",
						v.StringEmpty,
						f.Default("sflagval"),
					),
					s.WithCommand(
						"com",
						"help for com",
						c.DoNothing,
					),
				),
			),
			args:           []string{"test", "com"},
			passedPathWant: []string{"com"},
			passedFlagValuesWant: f.FlagValues{
				"--sflag": "sflagval",
			},
			wantErr: false,
		},
		{
			name: "simple JSON config",
			app: w.New(
				"test",
				"v0.0.0",
				s.NewSection("help for test",
					s.WithFlag(
						"--val",
						"flag help",
						v.StringEmpty,
						f.ConfigPath("params.val", v.StringFromInterface),
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

			args:           []string{"app", "com"},
			passedPathWant: []string{"com"},
			passedFlagValuesWant: f.FlagValues{
				"--config": "testdata/simple_json_config.json",
				"--val":    "hi",
			},
			wantErr: false,
		},
		{
			name: "config_slice",
			app: w.New(
				"test",
				"v0.0.0",
				s.NewSection(
					"help for test",
					s.WithFlag(
						"--subreddits",
						"the subreddits",
						v.StringSliceEmpty,
						// TODO: gonna need something here
						f.ConfigPath("subreddits[].name", v.StringFromInterface),
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
			args:           []string{"test", "print"},
			passedPathWant: []string{"print"},
			passedFlagValuesWant: f.FlagValues{
				"--subreddits": []string{"earthporn", "wallpapers"},
				"--config":     "testdata/config_slice.json",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			pr, err := tt.app.Parse(tt.args)

			// return early if there's an error
			// don't want to deref a null pr
			if (err != nil) && tt.wantErr {
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("RootCommand.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.passedPathWant, pr.PasssedPath)
			assert.Equal(t, tt.passedFlagValuesWant, pr.PassedFlags)
		})
	}
}
