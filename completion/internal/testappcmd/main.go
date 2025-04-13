package main

import "go.bbkane.com/warg/completion/internal/testapp"

func main() {
	app := testapp.BuildApp()
	if err := app.Validate(); err != nil {
		panic(err)
	}
	app.MustRun()
}
