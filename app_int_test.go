package warg

// internal tests - part of the warg package

import (
	"testing"
)

func TestGatherArgs(t *testing.T) {
	_, err := gatherArgs(
		[]string{"app", "cat", "--flag", "value", "--flag", "value2"},
		[]string{},
	)
	if err != nil {
		t.Fatal(err)
	}
	// fmt.Print(res)
}
