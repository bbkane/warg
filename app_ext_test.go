package warg_test

import (
	"testing"

	a "github.com/bbkane/warg"
	c "github.com/bbkane/warg/command"
	f "github.com/bbkane/warg/flag"
	s "github.com/bbkane/warg/section"
	v "github.com/bbkane/warg/value"
	"github.com/stretchr/testify/assert"
)

func TestApp_Parse(t *testing.T) {

	tests := []struct {
		name             string
		app              a.App
		args             []string
		passedPathWant   []string
		passedValuesWant v.ValueMap
		wantErr          bool
	}{
		{
			name: "from main",
			app: a.New("test", "v0.0.0",
				a.WithRootSection("help for test",
					s.WithFlag("--af1", "flag help", v.NewEmptyIntValue()),
					s.WithSection("cat1", "help for cat1",
						s.WithCommand("com1", "help for com1", c.DoNothing,
							c.WithFlag("--com1f1", "flag help", v.NewEmptyIntValue(),
								f.WithDefault(v.NewIntValue(10)),
							),
						),
					),
				),
			),

			args:             []string{"app", "cat1", "com1", "--com1f1", "1"},
			passedPathWant:   []string{"cat1", "com1"},
			passedValuesWant: v.ValueMap{"--com1f1": v.NewIntValue(1)},
			wantErr:          false,
		},
		{
			name: "no category",
			app: a.New("test", "v0.0.0",
				a.WithRootSection("help for test",
					s.WithFlag("--af1", "flag help", v.NewEmptyIntValue()),
				),
			),

			args:             []string{"app"},
			passedPathWant:   nil,
			passedValuesWant: map[string]v.Value{},
			wantErr:          false,
		},
		{
			name: "flag default",
			app: a.New("test", "v0.0.0",
				a.WithRootSection(
					"help for test",
					s.WithCommand("com", "com help", c.DoNothing,
						c.WithFlag("--flag", "flag help", v.NewEmptyStringValue(),
							f.WithDefault(v.NewStringValue("hi")),
						),
					),
				),
			),
			args:             []string{"test", "com"},
			passedPathWant:   []string{"com"},
			passedValuesWant: v.ValueMap{"--flag": v.NewStringValue("hi")},
			wantErr:          false,
		},
		{
			name: "extra flag",
			app: a.New("test", "v0.0.0",
				a.WithRootSection(
					"help for test",
					s.WithCommand("com", "com help", c.DoNothing,
						c.WithFlag("--flag", "flag help", v.NewEmptyStringValue(),
							f.WithDefault(v.NewStringValue("hi")),
						),
					),
				),
			),
			args:             []string{"test", "com", "--unexpected", "value"},
			passedPathWant:   nil,
			passedValuesWant: nil,
			wantErr:          true,
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
