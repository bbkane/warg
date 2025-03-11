package cli

type CommandHelp func(cur *Command, helpInfo HelpInfo) Action
type SectionHelp func(cur *SectionT, helpInfo HelpInfo) Action

// HelpFlagMapping adds a new option to your --help flag
type HelpFlagMapping struct {
	Name        string
	CommandHelp CommandHelp
	SectionHelp SectionHelp
}

// HelpInfo lists common information available to a help function
type HelpInfo struct {

	// AvailableFlags for the current section or commmand, including inherted flags from parent sections.
	// All flags are Resolved if possible (i.e., flag.SetBy != "")
	AvailableFlags FlagMap
	// RootSection of the app. Especially useful for printing all sections and commands
	RootSection SectionT
}
