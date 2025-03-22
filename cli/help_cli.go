package cli

type CommandHelp func(cur *Command, helpInfo HelpInfo) Action
type SectionHelp func(cur *SectionT, helpInfo HelpInfo) Action

// HelpInfo lists common information available to a help function
type HelpInfo struct {

	// AvailableFlags for the current section or commmand, including inherted flags from parent sections.
	// All flags are Resolved if possible (i.e., flag.SetBy != "")
	AvailableFlags FlagMap
	// RootSection of the app. Especially useful for printing all sections and commands
	RootSection SectionT
}

func HelpToCommand(commandHelp CommandHelp, secHelp SectionHelp) Command {
	return Command{ //nolint:exhaustruct  // This help is never used since this is a generated command
		Action: func(cmdCtx Context) error {
			// build ftar.AvailableFlags - it's a map of string to flag for the app globals + current command. Don't forget to set each flag.IsCommandFlag and Value for now..
			// TODO:
			ftarAllowedFlags := make(FlagMap)
			for flagName, fl := range cmdCtx.App.GlobalFlags {
				fl.Value = cmdCtx.ParseResult.FlagValues[flagName]
				fl.IsCommandFlag = false
				ftarAllowedFlags.AddFlag(flagName, fl)
			}

			// If we're in Parse_ExpectingSectionOrCommand, we haven't received a command
			if cmdCtx.ParseResult.State != Parse_ExpectingSectionOrCommand {
				for flagName, fl := range cmdCtx.ParseResult.CurrentCommand.Flags {
					fl.Value = cmdCtx.ParseResult.FlagValues[flagName]
					fl.IsCommandFlag = true
					ftarAllowedFlags.AddFlag(flagName, fl)
				}
			}

			hi := HelpInfo{
				AvailableFlags: ftarAllowedFlags,
				RootSection:    cmdCtx.App.RootSection,
			}
			com := cmdCtx.ParseResult.CurrentCommand
			if com != nil {
				return commandHelp(com, hi)(cmdCtx)
			} else {
				return secHelp(cmdCtx.ParseResult.CurrentSection, hi)(cmdCtx)
			}
		},
	}

}
