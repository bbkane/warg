package main

import (
	"testing"
)

func TestBuildApp(t *testing.T) {
	app := buildApp()

	if err := app.Validate(); err != nil {
		t.Fatal(err)
	}
}
