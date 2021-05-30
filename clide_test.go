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
		passedFlagsWant   c.FlagMap
		wantErr           bool
	}{
		{
			name: "from main",

			app: c.App{
				Name: "app",
				RootCategory: c.NewCategory(
					c.AddCategoryFlag(
						"--af1",
						c.FlagValue{},
					),
					c.WithCategory(
						"cat1",
						c.WithCommand(
							"com1",
							c.AddCommandFlag(
								"--com1f1",
								c.FlagValue{},
							),
						),
					),
				),
			},

			args:              []string{"app", "cat1", "com1", "--com1f1", "flagarg"},
			passedCommandWant: []string{"cat1", "com1"},
			passedFlagsWant:   c.FlagMap{"--com1f1": c.FlagValue{Value: "flagarg"}},
			wantErr:           false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, got1, err := tt.app.RootCategory.Parse(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("RootCommand.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.passedCommandWant) {
				t.Errorf("RootCommand.Parse() got = %v, want %v", got, tt.passedCommandWant)
			}
			if !reflect.DeepEqual(got1, tt.passedFlagsWant) {
				t.Errorf("RootCommand.Parse() got1 = %v, want %v", got1, tt.passedFlagsWant)
			}
		})
	}
}
