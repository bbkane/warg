///usr/bin/true; exec /usr/bin/env go run "$0" .
package warg_test

import (
	"reflect"
	"testing"

	w "github.com/bbkane/warg"
)

func TestApp_Parse(t *testing.T) {

	tests := []struct {
		name              string
		app               w.App
		args              []string
		passedCommandWant []string
		passedValuesWant  w.ValueMap
		wantErr           bool
	}{
		{
			name: "from main",
			app: w.NewApp(
				w.AppRootCategory(
					w.WithCategoryFlag(
						"--af1",
						w.NewIntValue(0),
					),
					w.WithCategory(
						"cat1",
						w.WithCommand(
							"com1",
							w.WithCommandFlag(
								"--com1f1",
								w.NewIntValue(0),
							),
						),
					),
				),
			),

			args:              []string{"app", "cat1", "com1", "--com1f1", "1"},
			passedCommandWant: []string{"cat1", "com1"},
			passedValuesWant:  w.ValueMap{"--com1f1": w.NewIntValue(1)},
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
