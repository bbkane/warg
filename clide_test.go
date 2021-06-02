///usr/bin/true; exec /usr/bin/env go run "$0" .
package clide_test

import (
	"reflect"
	"testing"

	c "github.com/bbkane/clide"
)

func TestApp_Parse(t *testing.T) {

	tests := []struct {
		name              string
		app               c.App
		args              []string
		passedCommandWant []string
		passedValuesWant  c.ValueMap
		wantErr           bool
	}{
		{
			name: "from main",

			app: c.App{
				Name: "app",
				RootCategory: c.NewCategory(
					c.AddCategoryFlag(
						"--af1",
						c.Flag{
							Value: c.NewIntValue(0),
						},
					),
					c.WithCategory(
						"cat1",
						c.WithCommand(
							"com1",
							c.AddCommandFlag(
								"--com1f1",
								c.Flag{
									Value: c.NewIntValue(0),
								},
							),
						),
					),
				),
			},

			args:              []string{"app", "cat1", "com1", "--com1f1", "1"},
			passedCommandWant: []string{"cat1", "com1"},
			passedValuesWant:  c.ValueMap{"--com1f1": c.NewIntValue(1)},
			wantErr:           false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			pr, err := tt.app.RootCategory.Parse(tt.args)
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
