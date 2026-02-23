package warg

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"slices"
	"strings"

	"github.com/mattn/go-shellwords"
	"github.com/reeflective/readline"
	"go.bbkane.com/warg/completion"
	"go.bbkane.com/warg/config"
	"go.bbkane.com/warg/path"
	"go.bbkane.com/warg/value"
	"go.bbkane.com/warg/value/scalar"
)

// AppOpt let's you customize the app. Most AppOpts panic if incorrectly called
type AppOpt func(*App)

// GlobalFlag adds an existing flag to a [Cmd]. It panics if a flag with the same name exists
func GlobalFlag(name string, value Flag) AppOpt {
	return func(com *App) {
		com.GlobalFlags.AddFlag(name, value)
	}
}

// GlobalFlagMap adds existing flags to a [Cmd]. It panics if a flag with the same name exists
func GlobalFlagMap(flagMap FlagMap) AppOpt {
	return func(com *App) {
		com.GlobalFlags.AddFlags(flagMap)
	}
}

// NewGlobalFlag adds a flag to the app. It panics if a flag with the same name exists
func NewGlobalFlag(name string, helpShort string, empty value.EmptyConstructor, opts ...FlagOpt) AppOpt {
	return GlobalFlag(name, NewFlag(helpShort, empty, opts...))

}

// ConfigFlag adds a global flag that will be used to read a config file. the flagMap must contain exactly one flag. When parsed, the config flag will be parsed before other flags, any config file found will be read, and any values found will be used to update other flags. This allows users to override flag defaults with values from a config file.
func ConfigFlag(reader config.NewReader, flagMap FlagMap) AppOpt {
	return func(app *App) {
		if len(flagMap) != 1 {
			panic(fmt.Sprintf("ConfigFlagMap must have exactly one flag, got %d", len(flagMap)))
		}
		app.NewConfigReader = reader
		app.ConfigFlagName = flagMap.SortedNames()[0]
		app.GlobalFlags.AddFlags(flagMap)
	}
}

// HelpFlag customizes your help  This option is only needed if you're also writing a custom help function. helpFlags be either `nil` to autogenerate or a flag map with one flat that with the followng properties:
//
//   - scalar string type
//   - choices that match the names in helpCommands
//   - default value set to one of the choices
//
// These properties are checked at runtime with app.Validate().
func HelpFlag(helpCmds CmdMap, helpFlags FlagMap) AppOpt {
	return func(a *App) {
		switch len(helpFlags) {
		case 0:
			helpFlags = DefaultHelpFlagMap("default", helpCmds.SortedNames())
		case 1:
			break
		default:
			panic(fmt.Sprintf("helpFlags must have 0 or 1 flags, got %d", len(helpFlags)))
		}

		a.HelpFlagName = helpFlags.SortedNames()[0]
		a.HelpCmds = helpCmds
		a.GlobalFlags.AddFlags(helpFlags)
	}
}

// SkipAll skips adding:
//   - the default completion commands (<app> completion)
//   - the default color flag map (<app> --color)
//   - the default version command map (<app> version)
//   - the default validation checks
//
// This is inteded for tests where you just want to assert against a minimal application
func SkipAll() AppOpt {
	return func(a *App) {
		a.SkipCompletionCmds = true
		a.SkipGlobalColorFlag = true
		a.SkipREPLCmd = true
		a.SkipValidation = true
		a.SkipVersionCmd = true
	}
}

// SkipCompletionCmds skips adding the default completion commands (<app> completion).
func SkipCompletionCmds() AppOpt {
	return func(a *App) {
		a.SkipCompletionCmds = true
	}
}

// SkipColorFlag skips adding the default color flag map (<app> --color).
func SkipGlobalColorFlag() AppOpt {
	return func(a *App) {
		a.SkipGlobalColorFlag = true
	}
}

// SkipValidation skips (most of) the app's internal consistency checks when the app is created.
// If used, make sure to call app.Validate() in a test!
func SkipValidation() AppOpt {
	return func(a *App) {
		a.SkipValidation = true
	}
}

// SkipVersionCmd skips adding the default version command (<app> version).
func SkipVersionCmd() AppOpt {
	return func(a *App) {
		a.SkipVersionCmd = true
	}
}

// SkipREPLCmd skips adding the default REPL command (<app> repl). NOTE: this command is not as polished as the rest of warg. I hope to improve it over time.
func SkipREPLCmd() AppOpt {
	return func(a *App) {
		a.SkipREPLCmd = true
	}
}

// FindVersion returns the version of the app. If the version is already set (eg. via a build flag), it returns that. Otherwise, it tries to read the go module version from the runtime info, or returns "unknown" if that fails. This is called automatically when passing an empty string to [New], but can also be used if you need the version independently of the app creation.
func FindVersion(version string) string {
	// if the version is already set (eg. via a build flag), return it
	if version != "" {
		return version
	}

	// If installed via `go install`, we'll be able to read runtime version info
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}
	// when run with `go run`, this will return "(devel)"
	return info.Main.Version
}

func CompletionsDirectories() CompletionsFunc {
	return func(cc CmdContext) (*completion.Candidates, error) {
		return &completion.Candidates{
			Type:   completion.Type_Directories,
			Values: nil,
		}, nil
	}
}

func CompletionsDirectoriesFiles() CompletionsFunc {
	return func(cc CmdContext) (*completion.Candidates, error) {
		return &completion.Candidates{
			Type:   completion.Type_DirectoriesFiles,
			Values: nil,
		}, nil
	}
}

func CompletionsNone() CompletionsFunc {
	return func(cc CmdContext) (*completion.Candidates, error) {
		return &completion.Candidates{
			Type:   completion.Type_None,
			Values: nil,
		}, nil
	}
}

func CompletionsValues(values []string) CompletionsFunc {
	var vals []completion.Candidate
	for _, v := range values {
		vals = append(vals, completion.Candidate{Name: v, Description: ""})
	}
	return func(ctx CmdContext) (*completion.Candidates, error) {

		return &completion.Candidates{
			Type:   completion.Type_Values,
			Values: vals,
		}, nil
	}
}

func CompletionsValuesDescriptions(values []completion.Candidate) CompletionsFunc {
	return func(ctx CmdContext) (*completion.Candidates, error) {
		return &completion.Candidates{
			Type:   completion.Type_ValuesDescriptions,
			Values: values,
		}, nil
	}
}

// New creates a warg app. name is used for help output only (though generally it should match the name of the compiled binary). version is the app version - if empty, warg will attempt to set it to the go module version, or "unknown" if that fails.
func New(name string, version string, rootSection Section, opts ...AppOpt) App {
	app := App{
		Name:                name,
		RootSection:         rootSection,
		ConfigFlagName:      "",
		NewConfigReader:     nil,
		HelpFlagName:        "",
		HelpCmds:            make(CmdMap),
		SkipCompletionCmds:  false,
		SkipGlobalColorFlag: false,
		SkipREPLCmd:         false,
		SkipValidation:      false,
		SkipVersionCmd:      false,
		Version:             version,
		GlobalFlags:         make(FlagMap),
	}
	for _, opt := range opts {
		opt(&app)
	}

	if app.HelpFlagName == "" {
		HelpFlag(
			DefaultHelpCmdMap(),
			DefaultHelpFlagMap("default", DefaultHelpCmdMap().SortedNames()),
		)(&app)
	}

	app.Version = FindVersion(app.Version)

	if !app.SkipGlobalColorFlag {
		GlobalFlagMap(FlagMap{
			"--color": NewFlag(
				"Use ANSI colors",
				scalar.String(
					scalar.Choices("true", "false", "auto"),
					scalar.Default("auto"),
				),
				EnvVars("WARG_COLOR"),
			),
		})(&app)
	}

	if !app.SkipCompletionCmds {
		NewSubSection(
			"completion",
			"Print shell completion scripts",
			NewSubCmd(
				"zsh",
				"Print zsh completion script",
				func(ctx CmdContext) error {
					completion.ZshCompletionScriptWrite(ctx.Stdout, app.Name)
					return nil
				},
			),
		)(&app.RootSection)
	}

	if !app.SkipVersionCmd {
		NewSubCmd(
			"version",
			"Print version",
			func(ctx CmdContext) error {
				fmt.Fprintln(ctx.Stdout, ctx.App.Version)
				return nil
			},
		)(&app.RootSection)
	}

	if !app.SkipREPLCmd {
		NewSubCmd(
			"repl",
			"Start a REPL to interactively run commands",
			replCmdAction,
		)(&app.RootSection)
	}

	// validate or not and return
	if app.SkipValidation {
		return app
	}

	err := app.Validate()
	if err != nil {
		panic(err)
	}
	return app
}

// An App contains your defined sections, commands, and flags
// Create a new App with New()
type App struct {
	// Config
	ConfigFlagName  string
	NewConfigReader config.NewReader

	// Help
	HelpFlagName string
	HelpCmds     CmdMap

	GlobalFlags         FlagMap
	Name                string
	RootSection         Section
	SkipGlobalColorFlag bool
	SkipCompletionCmds  bool
	SkipValidation      bool
	SkipVersionCmd      bool
	SkipREPLCmd         bool
	Version             string
}

func (app *App) Run(args []string, opts ...ParseOpt) {
	pr, err := app.Parse(args, opts...)
	if err != nil {
		// https://unix.stackexchange.com/a/254747
		fmt.Fprintln(os.Stderr, err)
		os.Exit(64)
	}
	err = pr.Action(pr.Context)
	if err != nil {
		fmt.Fprintln(pr.Context.Stderr, err)
		os.Exit(1)
	}
}

// MustRun runs the app.
// Any flag parsing errors will be printed to stderr and os.Exit(64) (EX_USAGE) will be called.
// Any errors on an Action will be printed to stderr and os.Exit(1) will be called.
func (app *App) MustRun(opts ...ParseOpt) {
	if len(os.Args) >= 3 && os.Args[1] == "--completion-zsh" {
		// app --completion-zsh <args> . Note that <args> must be something, even if it's the empty string

		// parseOpts.Args looks like: <exe> --completion-zsh <args>... <partialOrEmptyString>
		// the partial or empty string is passed to us from the completion script. Empty if the user just typed space and pressed tab, partial if the user pressed tab after typing part of something. zsh will filter that for us
		// so we need to remove the first two args and the last arg leaving <args>...
		args := os.Args[2 : len(os.Args)-1]

		candidates, err := app.Completions(args, opts...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		completion.ZshCompletionsWrite(os.Stdout, candidates)

	} else {
		app.Run(os.Args[1:], opts...)
	}
}

// Look up keys (meant for environment variable parsing) - fulfillable with os.LookupEnv or warg.LookupMap(map)
type LookupEnv func(key string) (string, bool)

// LookupMap loooks up keys from a provided map. Useful to mock os.LookupEnv when parsing
func LookupMap(m map[string]string) LookupEnv {
	return func(key string) (string, bool) {
		val, exists := m[key]
		return val, exists
	}
}

// validateFlags checks that global and command flag names and aliases start with "-" and are unique.
// It does not need to check the following scenarios:
//
//   - global flag names don't collide with global flag names (app will panic when adding the second global flag) - TOOD: ensure there's a test for this
//   - command flag names in the same command don't collide with each other (app will panic when adding the second command flag) TODO: ensure there's a test for this
//   - command flag names/aliases don't collide with command flag names/aliases in other commands (since only one command will be run, this is not a problem)
func validateFlags(
	globalFlags FlagMap,
	comFlags FlagMap,
) error {
	nameCount := make(map[string]int)
	for name, fl := range globalFlags {
		nameCount[name]++
		if fl.Alias != "" {
			nameCount[fl.Alias]++
		}
	}
	for name, fl := range comFlags {
		nameCount[name]++
		if fl.Alias != "" {
			nameCount[fl.Alias]++
		}
	}
	var errs []error
	for name, count := range nameCount {
		if !strings.HasPrefix(string(name), "-") {
			errs = append(errs, fmt.Errorf("flag and alias names must start with '-': %#v", name))
		}
		if count > 1 {
			errs = append(errs, fmt.Errorf("flag or alias name exists %d times: %v", count, name))
		}
	}
	return errors.Join(errs...)
}

// Validate checks app for creation errors. It checks:
//
//   - the help flag is the right type
//   - Sections and commands don't start with "-" (needed for parsing)
//   - Flag names and aliases do start with "-" (needed for parsing)
//   - Flag names and aliases don't collide
func (app *App) Validate() error {

	// validate --help flag
	if app.HelpFlagName == "" {
		return fmt.Errorf("HelpFlagName must be set")
	}
	helpFlag, exists := app.GlobalFlags[app.HelpFlagName]
	if !exists {
		return fmt.Errorf("HelpFlagName not found in GlobalFlags: %v", app.HelpFlagName)
	}
	helpFlagValEmpty, ok := helpFlag.EmptyValueConstructor().(value.ScalarValue)
	if !ok {
		return fmt.Errorf("HelpFlagName must be a scalar: %v", app.HelpFlagName)
	}
	if _, ok := helpFlagValEmpty.Get().(string); !ok {
		return fmt.Errorf("HelpFlagName must be a string: %v", app.HelpFlagName)
	}
	if !helpFlagValEmpty.HasDefault() {
		return fmt.Errorf("HelpFlagName must have a default value: %v", app.HelpFlagName)
	}
	if !slices.Equal(helpFlagValEmpty.Choices(), app.HelpCmds.SortedNames()) {
		return fmt.Errorf("HelpFlagName choices must match HelpCmds: %v", app.HelpFlagName)
	}
	if !slices.Contains(helpFlagValEmpty.Choices(), helpFlagValEmpty.DefaultString()) {
		return fmt.Errorf("HelpFlagName default value (%v) must be in choices (%v): %v", helpFlagValEmpty.DefaultString(), helpFlagValEmpty.Choices(), app.HelpFlagName)
	}

	// validate --config flag
	if app.ConfigFlagName != "" {
		if app.NewConfigReader == nil {
			return fmt.Errorf("ConfigFlagName must have a NewConfigReader: %v", app.ConfigFlagName)
		}
		configFlag, exists := app.GlobalFlags[app.ConfigFlagName]
		if !exists {
			return fmt.Errorf("ConfigFlagName not found in GlobalFlags: %v", app.ConfigFlagName)
		}
		configFlagValEmpty, ok := configFlag.EmptyValueConstructor().(value.ScalarValue)
		if !ok {
			return fmt.Errorf("ConfigFlagName must be a scalar: %v", app.ConfigFlagName)
		}
		if _, ok := configFlagValEmpty.Get().(path.Path); !ok {
			return fmt.Errorf("ConfigFlagName must be a path: %v", app.ConfigFlagName)
		}
	}

	// TODO: check that the default value is in the choices and the choices match app help mappings and that the flag is a scalar

	// NOTE: we need to be able to validate before we parse, and we may not know the app name
	// till after prsing so set the root path to "root"
	rootPath := []string{string(app.Name)}
	it := app.RootSection.breadthFirst(rootPath)

	for it.HasNext() {
		flatSec := it.Next()

		// Sections don't start with "-"
		secName := flatSec.Path[len(flatSec.Path)-1]
		if strings.HasPrefix(string(secName), "-") {
			return fmt.Errorf("section names must not start with '-': %#v", secName)
		}

		// Sections must not be leaf nodes
		if flatSec.Sec.Sections.Empty() && flatSec.Sec.Cmds.Empty() {
			return fmt.Errorf("sections must have either child sections or child commands: %#v", secName)
		}

		{
			// child section names should not clash with child command names
			nameCount := make(map[string]int)
			for name := range flatSec.Sec.Cmds {
				nameCount[string(name)]++
			}
			for name := range flatSec.Sec.Sections {
				nameCount[string(name)]++
			}
			errs := []error{}
			for name, count := range nameCount {
				if count > 1 {
					errs = append(errs, fmt.Errorf("command and section name clash: %s", name))
				}
			}
			err := errors.Join(errs...)
			if err != nil {
				return fmt.Errorf("name collision: %w", err)
			}
		}

		for name, com := range flatSec.Sec.Cmds {

			// Commands must not start wtih "-"
			if strings.HasPrefix(string(name), "-") {
				return fmt.Errorf("command names must not start with '-': %#v", name)
			}

			err := validateFlags(app.GlobalFlags, com.Flags)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// CompletionsFunc is a function that returns completion candidates for a flag. See warg.Completions[Type] for convenience functions to make this
type CompletionsFunc func(CmdContext) (*completion.Candidates, error)

// Completions generates completions from a list of args, by looking at the app structure starting at the root section. Args must start with a section, command in root section or a global flag and must not end with a partial string
func (a *App) Completions(args []string, opts ...ParseOpt) (*completion.Candidates, error) {
	parseOpts := NewParseOpts(opts...)

	// I could to a full parse here, but that would be slower and more prone to failure than just parsing the args - we don't need a lot of info to complete section/command names
	parseState, err := a.parseArgs(args)
	if err != nil {
		return nil, fmt.Errorf("unexpected parseArgs err: %w", err)
	}

	// special case if help is passed
	if parseState.HelpPassed {
		// if the value of the flag has been passed, don't suggest anything
		if parseState.FlagValues[a.HelpFlagName].UpdatedBy() == value.UpdatedByFlag {
			return &completion.Candidates{
				Type:   completion.Type_None,
				Values: nil,
			}, nil
		}

		// otherwise suggest the help commands as the values of the help flag
		res := &completion.Candidates{
			Type:   completion.Type_Values,
			Values: []completion.Candidate{},
		}
		for _, name := range a.HelpCmds.SortedNames() {
			res.Values = append(res.Values, completion.Candidate{
				Name:        string(name),
				Description: "",
			})
		}
		return res, nil
	}

	if parseState.ParseArgState == ParseArgState_WantSectionOrCmd {
		s := parseState.CurrentSection
		ret := completion.Candidates{
			Type:   completion.Type_ValuesDescriptions,
			Values: []completion.Candidate{},
		}
		for _, name := range s.Cmds.SortedNames() {
			ret.Values = append(ret.Values, completion.Candidate{
				Name:        string(name),
				Description: string(s.Cmds[name].HelpShort),
			})
		}
		for _, name := range s.Sections.SortedNames() {
			ret.Values = append(ret.Values, completion.Candidate{
				Name:        string(name),
				Description: string(s.Sections[name].HelpShort),
			})
		}
		ret.Values = append(ret.Values, completion.Candidate{
			Name:        a.HelpFlagName,
			Description: a.GlobalFlags[a.HelpFlagName].HelpShort,
		})
		return &ret, nil
	}

	// Finish the parse!
	err = a.resolveFlags(parseState.CurrentCmd, parseState.FlagValues, parseOpts.LookupEnv, parseState.UnsetFlagNames)
	if err != nil {
		return nil, fmt.Errorf("unexpected resolveFlags err: %w", err)
	}
	cmdContext := CmdContext{
		App:           a,
		ParseMetadata: parseOpts.ParseMetadata,
		Flags:         parseState.FlagValues.ToPassedFlags(),
		ForwardedArgs: parseState.CurrentCmdForwardedArgs, // should always be nil during completions as completions occur at the end
		ParseState:    &parseState,
		Stderr:        parseOpts.Stderr,
		Stdin:         parseOpts.Stdin,
		Stdout:        parseOpts.Stdout,
	}

	switch parseState.ParseArgState {
	case ParseArgState_WantFlagNameOrEnd:
		return cmdCompletions(cmdContext)
	case ParseArgState_WantFlagValue:
		return parseState.CurrentFlag.Completions(cmdContext)
	case ParseArgState_WantSectionOrCmd:
		panic("unreachable state: ExpectingArg_SectionOrCommand")
	default:
		return nil, fmt.Errorf("unexpected ParseState: %v", parseState.ParseArgState)
	}
}

func cmdCompletions(cmdCtx CmdContext) (*completion.Candidates, error) {
	// FZF (or maybe zsh) auto-sorts by alphabetical order, so no need to get fancy with the following ideas
	//  - if the flag is required and is not set, suggest it first
	//  - suggest command flags before global flags
	//  - let the flags define rank or priority for completion order
	candidates := &completion.Candidates{
		Type:   completion.Type_ValuesDescriptions,
		Values: []completion.Candidate{},
	}
	// command flags
	for _, name := range cmdCtx.ParseState.CurrentCmd.Flags.SortedNames() {
		// scalar flags set by passed arg can't be appended to or overridden, so don't suggest them
		val, isScalar := cmdCtx.ParseState.FlagValues[name].(value.ScalarValue)
		if isScalar && val.UpdatedBy() == value.UpdatedByFlag {
			continue
		}
		var valStr string
		// TODO: does it matter if valstring is a large list?
		if cmdCtx.ParseState.FlagValues[name].UpdatedBy() != value.UpdatedByUnset {
			valStr = fmt.Sprint(cmdCtx.ParseState.FlagValues[name].Get())
			valStr = strings.ReplaceAll(valStr, "\n", " ")
			valStr = " (" + valStr + ")"
		}

		candidates.Values = append(candidates.Values, completion.Candidate{
			Name:        string(name),
			Description: string(cmdCtx.ParseState.CurrentCmd.Flags[name].HelpShort) + valStr,
		})
	}
	// global flags
	for _, name := range cmdCtx.App.GlobalFlags.SortedNames() {
		candidates.Values = append(candidates.Values, completion.Candidate{
			Name:        string(name),
			Description: string(cmdCtx.App.GlobalFlags[name].HelpShort),
		})
	}

	// AllowForwardedArgs
	if cmdCtx.ParseState.CurrentCmd.AllowForwardedArgs {
		candidates.Values = append(candidates.Values, completion.Candidate{
			Name:        "--",
			Description: "Indicates the end of flag parsing and the beginning of forwarded args",
		})
	}

	return candidates, nil
}

// Examples to get an intution of this:
//
// butler >>> hi
// line: 2, cursor: 2
// hi
// --^
// func debugCompleter(line []rune, cursor int) readline.Completions {
// 	spaces := make([]rune, cursor+1)
// 	for i := 0; i < cursor+1; i++ {
// 		spaces[i] = '-'
// 	}
// 	spaces[cursor] = '^'
// 	stats := fmt.Sprintf("line: %d, cursor: %d", len(line), cursor)
// 	msg := strings.Join([]string{stats, string(line), string(spaces)}, "\n")
// 	return readline.CompleteMessage(msg)
// }

func replCmdAction(cmdCtx CmdContext) error {
	rl := readline.NewShell()
	err := rl.Config.Set("menu-complete-display-prefix", true)
	if err != nil {
		return fmt.Errorf("could not set readline config: %w", err)
	}
	rl.Prompt.Primary(func() string {
		return cmdCtx.App.Name + " >>> "
	})

	rl.Completer = func(line []rune, cursor int) readline.Completions {

		// don't care about stuff after the cursor.
		truncatedLine := line[:cursor]

		lineStr := string(truncatedLine)
		// fmt.Fprintf(os.Stderr, "lineStr: %s\n", lineStr)

		words, err := shellwords.Parse(lineStr)
		if err != nil {
			err = fmt.Errorf("could not parse args for completion: args: %v, %w", words, err)
			return readline.CompleteMessage(err.Error())
		}

		// Completions are accepted with " ", so if we don't end in " ", it's a partial word. We can't handle partial words, so chop it off
		if len(lineStr) != 0 && !strings.HasSuffix(lineStr, " ") {
			words = words[:len(words)-1]
		}
		// fmt.Fprintf(os.Stderr, "words: %#v\n", words)

		// TODO: should I copy parseOpts from cmdCtx?
		candidates, err := cmdCtx.App.Completions(
			words,
		)
		if err != nil {
			err = fmt.Errorf("could not get completions: args: %v, %w", words, err)
			return readline.CompleteMessage(err.Error())
		}

		//nolint:exhaustive  // the default handles the cases we don't support yet
		switch candidates.Type {
		case completion.Type_ValuesDescriptions:
			vals := make([]string, 0, len(candidates.Values)*2)
			for _, c := range candidates.Values {
				vals = append(vals, c.Name)
				vals = append(vals, c.Description)
			}
			return readline.CompleteValuesDescribed(vals...)

		case completion.Type_Values:
			vals := make([]string, 0, len(candidates.Values))
			for _, c := range candidates.Values {
				vals = append(vals, c.Name)
			}
			return readline.CompleteValues(vals...)

		default:
			return readline.CompleteMessage("TODO: support type: " + string(candidates.Type))
		}

	}

	// fmt.Printf("Completing for line: %q at cursor position %d\n", string(line), cursor)
	// completions := readline.CompleteValuesDescribed("hi", "Hi loooooooooonnnnnnggggg description", "there", "there description")
	// completions.Usage("Use tab to complete")
	// return completions
	for {
		line, err := rl.Readline()
		switch err {
		case readline.ErrInterrupt, io.EOF:
			return nil
		}
		if err != nil {
			return fmt.Errorf("could not read line: %w", err)
		}
		words, err := shellwords.Parse(line)
		if err != nil {
			fmt.Fprintf(cmdCtx.Stderr, "could not parse line: %v\n", err)
			continue
		}
		pr, err := cmdCtx.App.Parse(words)
		if err != nil {
			fmt.Fprintf(cmdCtx.Stderr, "could not parse args: %v\n", err)
			continue
		}
		err = pr.Action(pr.Context)
		if err != nil {
			fmt.Fprintf(cmdCtx.Stderr, "error running command: %v\n", err)
			continue
		}
	}
}
