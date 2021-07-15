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
					s.WithFlag("--af1", "flag help", v.IntValueEmpty()),
					s.WithSection("cat1", "help for cat1",
						s.WithCommand("com1", "help for com1", c.DoNothing,
							c.WithFlag("--com1f1", "flag help", v.IntValueEmpty(),
								f.Default(v.IntValueNew(10)),
							),
						),
					),
				),
			),

			args:             []string{"app", "cat1", "com1", "--com1f1", "1"},
			passedPathWant:   []string{"cat1", "com1"},
			passedValuesWant: v.ValueMap{"--com1f1": v.IntValueNew(1)},
			wantErr:          false,
		},
		{
			name: "no category",
			app: w.New("test", "v0.0.0",
				s.NewSection("help for test",
					s.WithFlag("--af1", "flag help", v.IntValueEmpty()),
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
						c.WithFlag("--flag", "flag help", v.StringValueEmpty(),
							f.Default(v.StringValueNew("hi")),
						),
					),
				),
			),
			args:             []string{"test", "com"},
			passedPathWant:   []string{"com"},
			passedValuesWant: v.ValueMap{"--flag": v.StringValueNew("hi")},
			wantErr:          false,
		},
		{
			name: "extra flag",
			app: w.New("test", "v0.0.0",
				s.NewSection(
					"help for test",
					s.WithCommand("com", "com help", c.DoNothing,
						c.WithFlag("--flag", "flag help", v.StringValueEmpty(),
							f.Default(v.StringValueNew("hi")),
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
					s.WithFlag("--key", "a key", v.StringValueEmpty(),
						f.ConfigPath("key", v.StringValueFromInterface),
						f.Default(v.StringValueNew("defaultkeyval")),
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
					f.Default(v.StringValueNew("defaultconfigval")),
				),
			),
			args:           []string{"test", "print", "--config", "passedconfigval"},
			passedPathWant: []string{"print"},
			passedValuesWant: v.ValueMap{
				"--key":    v.StringValueNew("mapkeyval"),
				"--config": v.StringValueNew("passedconfigval"),
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
						v.StringValueEmpty(),
						f.Default(v.StringValueNew("sflagval")),
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
				"--sflag": v.StringValueNew("sflagval"),
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
