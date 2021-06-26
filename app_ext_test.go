package warg_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	a "github.com/bbkane/warg"
	c "github.com/bbkane/warg/command"
	f "github.com/bbkane/warg/flag"
	s "github.com/bbkane/warg/section"
	v "github.com/bbkane/warg/value"
)

func TestApp_Parse(t *testing.T) {

	tests := []struct {
		name              string
		app               a.App
		args              []string
		passedCommandWant []string
		passedValuesWant  v.ValueMap
		wantErr           bool
	}{
		{
			name: "from main",
			app: a.New(
				"test",
				"v0.0.0",
				a.WithRootSection(
					"help for test",
					s.WithFlag(
						"--af1",
						"flag help",
						v.NewEmptyIntValue(),
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
								v.NewEmptyIntValue(),
								f.WithDefault(v.NewIntValue(10)),
							),
						),
					),
				),
			),

			args:              []string{"app", "cat1", "com1", "--com1f1", "1"},
			passedCommandWant: []string{"cat1", "com1"},
			passedValuesWant:  v.ValueMap{"--com1f1": v.NewIntValue(1)},
			wantErr:           false,
		},
		{
			name: "no category",
			app: a.New(
				"test",
				"v0.0.0",
				a.WithRootSection(
					"help for test",
					s.WithFlag(
						"--af1",
						"flag help",
						v.NewEmptyIntValue(),
					),
				),
			),

			args:              []string{"app"},
			passedCommandWant: nil,
			passedValuesWant:  map[string]v.Value{},
			wantErr:           false,
		},
		{
			name: "flag default",
			app: a.New(
				"test",
				"v0.0.0",
				a.WithRootSection(
					"help for test",
					s.WithCommand(
						"com",
						"com help",
						c.DoNothing,
						c.WithFlag(
							"--flag",
							"flag help",
							v.NewEmptyStringValue(),
							f.WithDefault(v.NewStringValue("hi")),
						),
					),
				),
			),
			args:              []string{"test", "com"},
			passedCommandWant: []string{"com"},
			passedValuesWant:  v.ValueMap{"--flag": v.NewStringValue("hi")},
			wantErr:           false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// pr, err := tt.app.RootCategory.Parse(tt.args)
			pr, err := tt.app.Parse(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("RootCommand.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, pr.PassedCmd, tt.passedCommandWant)
			assert.Equal(t, pr.PassedFlags, tt.passedValuesWant)
		})
	}
}
