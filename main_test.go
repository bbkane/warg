///usr/bin/true; exec /usr/bin/env go run "$0" .
package main

import (
	"reflect"
	"testing"
)

func TestRootCommand_Parse(t *testing.T) {

	type args struct {
		args []string
	}
	tests := []struct {
		name              string
		command           App
		args              []string
		passedCommandWant []string
		passedFlagsWant   FlagMap
		wantErr           bool
	}{
		{
			name: "from main",
			command: App{
				Name: "rc",
				Flags: FlagMap{
					"--rcf1": Flag{},
				},
				Commands: CommandMap{},
				Categories: CategoryMap{
					"sc1": Category{
						Flags: FlagMap{},
						Commands: CommandMap{
							"lc1": Command{
								Flags: FlagMap{
									"--lc1f1": Flag{},
								},
							},
						},
					},
				},
			},
			args:              []string{"rc", "sc1", "lc1", "--lc1f1", "flagarg"},
			passedCommandWant: []string{"sc1", "lc1"},
			passedFlagsWant:   FlagMap{"--lc1f1": Flag{Value: "flagarg"}},
			wantErr:           false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, got1, err := tt.command.Parse(tt.args)
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
