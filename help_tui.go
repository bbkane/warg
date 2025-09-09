package warg

import (
	"errors"
	"fmt"
)

func tuiHelp() Cmd {
	action := func(cmdCtx CmdContext) error {
		if cmdCtx.ParseState.ParseArgState != ParseArgState_WantFlagNameOrEnd {
			return fmt.Errorf("unexpected parse state: %v", cmdCtx.ParseState.ParseArgState)
		}

		return errors.New("TODO")
	}
	return NewCmd("", action)
}
