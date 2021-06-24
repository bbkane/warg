///usr/bin/true; exec /usr/bin/env go run "$0" .
package warg_test

import (
	"reflect"
	"testing"

	a "github.com/bbkane/warg/app"
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
			app: a.NewApp(
				a.AppRootCategory(
					s.WithCategoryFlag(
						"--af1",
						v.NewEmptyIntValue(),
					),
					s.WithCategory(
						"cat1",
						s.WithCommand(
							"com1",
							c.WithCommandFlag(
								"--com1f1",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// pr, err := tt.app.RootCategory.Parse(tt.args)
			pr, err := tt.app.Parse(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("RootCommand.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(pr.PassedCmd, tt.passedCommandWant) {
				t.Errorf("RootCommand.Parse() got = %v, want %v", pr.PassedCmd, tt.passedCommandWant)
			}
			if !reflect.DeepEqual(pr.PassedFlags, tt.passedValuesWant) {
				t.Errorf("RootCommand.Parse() got1 = %v, want %v", pr.PassedFlags, tt.passedValuesWant)
			}
		})
	}
}
