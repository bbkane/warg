package clide

import (
	"fmt"
	"testing"
)

func TestGatherArgs(t *testing.T) {
	res, err := gatherArgs([]string{"app", "cat", "--flag", "value", "--flag", "value2"})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Print(res)
}
