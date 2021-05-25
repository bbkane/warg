///usr/bin/true; exec /usr/bin/env go run "$0" .
package main

import (
	"reflect"
	"testing"
)

func TestApp_Parse(t *testing.T) {

	type args struct {
		args []string
	}
	tests := []struct {
		name              string
		app               App
		args              []string
		passedCommandWant []string
		passedFlagsWant   FlagMap
		wantErr           bool
	}{
		{
			name: "from main",

			app: App{
				Name: "app",
				RootCategory: NewCategory(
					AddCategoryFlag(
						"--af1",
						FlagValue{},
					),
					WithCategory(
						"cat1",
						WithCommand(
							"com1",
							AddCommandFlag(
								"--com1f1",
								FlagValue{},
							),
						),
					),
				),
			},

			args:              []string{"app", "cat1", "com1", "--com1f1", "flagarg"},
			passedCommandWant: []string{"cat1", "com1"},
			passedFlagsWant:   FlagMap{"--com1f1": FlagValue{Value: "flagarg"}},
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
