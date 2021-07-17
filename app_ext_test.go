package warg_test

import (
	"testing"

	w "github.com/bbkane/warg"
	c "github.com/bbkane/warg/command"
	f "github.com/bbkane/warg/flag"
	s "github.com/bbkane/warg/section"
	v "github.com/bbkane/warg/value"
	"github.com/stretchr/testify/assert"
)

func TestApp_Parse(t *testing.T) {

	tests := []struct {
		name             string
		app              w.App
		args             []string
		passedPathWant   []string
		passedValuesWant v.ValueMap
		wantErr          bool
	}{
		{
			name: "from main",
			app: w.New("test", "v0.0.0",
				s.NewSection("help for test",
					s.WithFlag("--af1", "flag help", v.IntEmpty()),
					s.WithSection("cat1", "help for cat1",
						s.WithCommand("com1", "help for com1", c.DoNothing,
							c.WithFlag("--com1f1", "flag help", v.IntEmpty(),
								f.Default(v.IntNew(10)),
							),
						),
					),
				),
			),

			args:             []string{"app", "cat1", "com1", "--com1f1", "1"},
			passedPathWant:   []string{"cat1", "com1"},
			passedValuesWant: v.ValueMap{"--com1f1": v.IntNew(1)},
			wantErr:          false,
		},
		{
			name: "no category",
			app: w.New("test", "v0.0.0",
				s.NewSection("help for test",
					s.WithFlag("--af1", "flag help", v.IntEmpty()),
				),
			),

			args:             []string{"app"},
			passedPathWant:   nil,
			passedValuesWant: nil,
			wantErr:          false,
		},
		{
			name: "flag default",
			app: w.New("test", "v0.0.0",
				s.NewSection(
					"help for test",
					s.WithCommand("com", "com help", c.DoNothing,
						c.WithFlag("--flag", "flag help", v.StringEmpty(),
							f.Default(v.StringNew("hi")),
						),
					),
				),
			),
			args:             []string{"test", "com"},
			passedPathWant:   []string{"com"},
			passedValuesWant: v.ValueMap{"--flag": v.StringNew("hi")},
			wantErr:          false,
		},
		{
			name: "extra flag",
			app: w.New("test", "v0.0.0",
				s.NewSection(
					"help for test",
					s.WithCommand("com", "com help", c.DoNothing,
						c.WithFlag("--flag", "flag help", v.StringEmpty(),
							f.Default(v.StringNew("hi")),
						),
					),
				),
			),
			args:             []string{"test", "com", "--unexpected", "value"},
			passedPathWant:   nil,
			passedValuesWant: nil,
			wantErr:          true,
		},
		{
			name: "config flag",
			app: w.New("test", "v0.0.0",
				s.NewSection(
					"help for test",
					s.WithFlag("--key", "a key", v.StringEmpty(),
						f.ConfigPath("key", v.StringFromInterface),
						f.Default(v.StringNew("defaultkeyval")),
					),
					s.WithCommand("print", "print key value", c.DoNothing),
				),
				w.ConfigFlag(
					"--config",
					// dummy function just to get me a map
					func(s string) (w.ConfigMap, error) {
						return w.ConfigMap{
							"configName": s,
							"key":        "mapkeyval",
						}, nil
					},
					"config flag",
					f.Default(v.StringNew("defaultconfigval")),
				),
			),
			args:           []string{"test", "print", "--config", "passedconfigval"},
			passedPathWant: []string{"print"},
			passedValuesWant: v.ValueMap{
				"--key":    v.StringNew("mapkeyval"),
				"--config": v.StringNew("passedconfigval"),
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
						v.StringEmpty(),
						f.Default(v.StringNew("sflagval")),
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
			passedValuesWant: v.ValueMap{
				"--sflag": v.StringNew("sflagval"),
			},
			wantErr: false,
		},
		{
			name: "simple JSON config",
			app: w.New(
				"test",
				"v0.0.0",
				s.NewSection("help for test",
					s.WithFlag("--val", "flag help", v.StringEmpty(),
						f.ConfigPath("params.val", v.StringFromInterface),
					),
					s.WithCommand("com", "help for com", c.DoNothing),
				),
				w.ConfigFlag(
					"--config",
					w.JSONUnmarshaller,
					"path to config",
					// TODO: make this test work by following the config cases
					// in the README
					f.Default(v.StringNew("testdata/simple_json_config.json")),
				),
			),

			args:           []string{"app", "com"},
			passedPathWant: []string{"com"},
			passedValuesWant: v.ValueMap{
				"--config": v.StringNew("testdata/simple_json_config.json"),
				"--val":    v.StringNew("hi"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// pr, err := tt.app.RootCategory.Parse(tt.args)
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
			assert.Equal(t, tt.passedValuesWant, pr.PassedFlags)
		})
	}
}
