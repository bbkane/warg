package warg

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"slices"
	"strconv"
	"strings"

	"github.com/mattn/go-shellwords"
	"github.com/reeflective/readline"
	"go.bbkane.com/warg/colerr"
	"go.bbkane.com/warg/completion"
	"go.bbkane.com/warg/config"
	"go.bbkane.com/warg/path"
	"go.bbkane.com/warg/value"
	"go.bbkane.com/warg/value/contained"
	"go.bbkane.com/warg/value/scalar"
)

// AppOpt is a functional option for configuring an [App] during creation via [New].
// Most AppOpts panic if given invalid input (e.g., duplicate flag names).
type AppOpt func(*App)

// GlobalFlag registers an existing [Flag] as a global flag available to all commands.
// It panics if a flag with the same name already exists.
func GlobalFlag(name string, value Flag) AppOpt {
	return func(com *App) {
		com.GlobalFlags.AddFlag(name, value)
	}
}

// GlobalFlagMap registers multiple existing flags as global flags available to all commands.
// It panics if any flag name already exists.
func GlobalFlagMap(flagMap FlagMap) AppOpt {
	return func(com *App) {
		com.GlobalFlags.AddFlags(flagMap)
	}
}

// NewGlobalFlag creates and registers a new global flag available to all commands.
// It panics if a flag with the same name already exists.
func NewGlobalFlag(name string, helpShort string, empty value.EmptyConstructor, opts ...FlagOpt) AppOpt {
	return GlobalFlag(name, NewFlag(helpShort, empty, opts...))

}

// ConfigFlag enables config file support by adding a global flag whose value is the path
// to a config file. The flagMap must contain exactly one flag (panics otherwise).
// During parsing, this flag is resolved first; the referenced config file is then read,
// and its values are used to fill other unset flags before environment variables or defaults.
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

// HelpFlag customizes the help system by providing custom help command implementations
// and an optional flag map. Only needed if writing custom help output.
// helpFlags may be nil (to auto-generate) or a [FlagMap] with exactly one flag that:
//
//   - is a scalar string type
//   - has choices matching the keys of helpCmds
//   - has a default value that is one of those choices
//
// These constraints are validated at runtime by [App.Validate].
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

// SkipAll disables all automatically added features:
//   - completion commands (<app> completion)
//   - --color global flag
//   - --term-width global flag
//   - repl command
//   - version command
//   - startup validation
//
// Intended for tests that need a minimal application without auto-generated commands/flags.
func SkipAll() AppOpt {
	return func(a *App) {
		a.SkipCompletionCmds = true
		a.SkipGlobalColorFlag = true
		a.SkipGlobalTermWidthFlag = true
		a.SkipREPLCmd = true
		a.SkipValidation = true
		a.SkipVersionCmd = true
	}
}

// SkipCompletionCmds disables the auto-generated "completion" section and its subcommands.
func SkipCompletionCmds() AppOpt {
	return func(a *App) {
		a.SkipCompletionCmds = true
	}
}

// SkipGlobalColorFlag disables the auto-generated --color global flag.
func SkipGlobalColorFlag() AppOpt {
	return func(a *App) {
		a.SkipGlobalColorFlag = true
	}
}

// SkipGlobalTermWidthFlag disables the auto-generated --term-width global flag.
func SkipGlobalTermWidthFlag() AppOpt {
	return func(a *App) {
		a.SkipGlobalTermWidthFlag = true
	}
}

// SkipValidation disables the app's internal consistency checks during [New].
// If used, call [App.Validate] in a test to catch configuration errors.
func SkipValidation() AppOpt {
	return func(a *App) {
		a.SkipValidation = true
	}
}

// SkipVersionCmd disables the auto-generated "version" command.
func SkipVersionCmd() AppOpt {
	return func(a *App) {
		a.SkipVersionCmd = true
	}
}

// SkipREPLCmd disables the auto-generated "repl" command.
// NOTE: the REPL command is experimental and may change in future versions.
func SkipREPLCmd() AppOpt {
	return func(a *App) {
		a.SkipREPLCmd = true
	}
}

// FindVersion determines the application version. If version is non-empty, it is returned as-is.
// Otherwise, it reads the Go module version from [debug.ReadBuildInfo], returning "(devel)" when
// run via "go run" or "unknown" if build info is unavailable.
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

// CompletionsDirectories returns a [CompletionsFunc] that suggests only directory names.
func CompletionsDirectories() CompletionsFunc {
	return func(cc CmdContext) (*completion.Candidates, error) {
		return &completion.Candidates{
			Type:   completion.Type_Directories,
			Values: nil,
		}, nil
	}
}

// CompletionsDirectoriesFiles returns a [CompletionsFunc] that suggests both files and directories.
func CompletionsDirectoriesFiles() CompletionsFunc {
	return func(cc CmdContext) (*completion.Candidates, error) {
		return &completion.Candidates{
			Type:   completion.Type_DirectoriesFiles,
			Values: nil,
		}, nil
	}
}

// CompletionsNone returns a [CompletionsFunc] that provides no completion suggestions.
func CompletionsNone() CompletionsFunc {
	return func(cc CmdContext) (*completion.Candidates, error) {
		return &completion.Candidates{
			Type:   completion.Type_None,
			Values: nil,
		}, nil
	}
}

// CompletionsValues returns a [CompletionsFunc] that suggests the given fixed string values.
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

// CompletionsValuesDescriptions returns a [CompletionsFunc] that suggests the given values
// with descriptions shown alongside each candidate.
func CompletionsValuesDescriptions(values []completion.Candidate) CompletionsFunc {
	return func(ctx CmdContext) (*completion.Candidates, error) {
		return &completion.Candidates{
			Type:   completion.Type_ValuesDescriptions,
			Values: values,
		}, nil
	}
}

// New creates a new CLI application. name is used for help output and should match
// the compiled binary name. version is displayed by the auto-generated "version" command;
// if empty, warg attempts to read it from the Go module build info.
// [New] panics if [App.Validate] fails (disable with [SkipValidation]).
func New(name string, version string, rootSection Section, opts ...AppOpt) App {
	app := App{
		Name:                    name,
		RootSection:             rootSection,
		ConfigFlagName:          "",
		NewConfigReader:         nil,
		HelpFlagName:            "",
		HelpCmds:                make(CmdMap),
		SkipCompletionCmds:      false,
		SkipGlobalColorFlag:     false,
		SkipGlobalTermWidthFlag: false,
		SkipREPLCmd:             false,
		SkipValidation:          false,
		SkipVersionCmd:          false,
		Version:                 version,
		GlobalFlags:             make(FlagMap),
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

	if !app.SkipGlobalTermWidthFlag {
		GlobalFlagMap(FlagMap{
			"--term-width": NewFlag(
				"Terminal width. Should be a positive integer, \"auto\", or \"infinite\". If \"auto\" the app will attempt to detect terminal width and fall back to 120 if detection fails",
				scalar.New(
					contained.TypeInfo[string]{
						Description: "string",
						FromIFace: func(iFace interface{}) (string, error) {
							s, ok := iFace.(string)
							if !ok {
								return "", contained.ErrIncompatibleInterface
							}
							return parseTermWidth(s)
						},
						FromString: parseTermWidth,
						FromZero:   contained.FromZero[string],
						Equals:     contained.Equals[string],
					},
					scalar.Default("auto"),
				),
				EnvVars("WARG_TERM_WIDTH"),
				FlagCompletions(CompletionsValues([]string{"auto", "infinite"})),
			),
		})(&app)
	}

	if !app.SkipCompletionCmds {
		NewSubSection(
			"completion",
			"Print shell completion scripts",
			NewSubCmd(
				"bash",
				"Print bash completion script",
				func(ctx CmdContext) error {
					completion.BashCompletionScriptWrite(ctx.Stdout, app.Name)
					return nil
				},
			),
			NewSubCmd(
				"zsh",
				"Print zsh completion script",
				func(ctx CmdContext) error {
					completion.ZshCompletionScriptWrite(ctx.Stdout, app.Name)
					return nil
				},
			),
			NewSubCmd(
				"fish",
				"Print fish completion script",
				func(ctx CmdContext) error {
					completion.FishCompletionScriptWrite(ctx.Stdout, app.Name)
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

// App is the top-level container for a warg CLI application, holding sections,
// commands, flags, and configuration. Create with [New].
type App struct {
	// Config
	ConfigFlagName  string
	NewConfigReader config.NewReader

	// Help
	HelpFlagName string
	HelpCmds     CmdMap

	GlobalFlags             FlagMap
	Name                    string
	RootSection             Section
	SkipGlobalColorFlag     bool
	SkipGlobalTermWidthFlag bool
	SkipCompletionCmds      bool
	SkipValidation          bool
	SkipVersionCmd          bool
	SkipREPLCmd             bool
	Version                 string
}

func parseTermWidth(s string) (string, error) {
	if s == "auto" || s == "infinite" {
		return s, nil
	}

	parsed, err := strconv.Atoi(s)
	if err != nil {
		return "", colerr.NewWrappedf(nil, "Expected a positive integer, \"infinite\", or \"auto\", got %s", fmt.Sprintf("%q", s))
	}
	if parsed <= 0 {
		return "", colerr.NewWrappedf(nil, "Expected a positive integer, \"infinite\", or \"auto\", got %s", fmt.Sprintf("%q", s))
	}

	return s, nil
}

// MustRunWithArgs parses and executes the app with the given args (without the program name).
// Parse errors are printed to stderr with exit code 64 (EX_USAGE).
// Action errors are printed to stderr with exit code 1.
func (app *App) MustRunWithArgs(args []string, opts ...ParseOpt) {
	// TODO: make colors optional!
	pr, err := app.Parse(args, opts...)
	if err != nil {
		// TODO: right now passing nil since if there's a parsing error we don't have passed flags. I'd like to also gate this on an env var in that case.
		disable := os.Getenv("WARG_COLOR") == "false"
		style, _ := conditionallyEnableStyle(disable, nil, os.Stderr)
		colerr.Stacktrace(os.Stderr, &style, err)
		// https://unix.stackexchange.com/a/254747
		os.Exit(64)
	}
	err = pr.Action(pr.Context)
	if err != nil {
		// note that this is user code, so let's not impose styles
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// MustRun parses os.Args and executes the matched command, or produces shell completions
// if a --completion-{bash,zsh,fish} flag is detected. Parse errors exit with code 64 (EX_USAGE);
// action errors exit with code 1.
func (app *App) MustRun(opts ...ParseOpt) {
	if len(os.Args) >= 3 && os.Args[1] == "--completion-bash" {
		// app --completion-bash <args> . Note that <args> must be something, even if it's the empty string

		// parseOpts.Args looks like: <exe> --completion-bash <args>... <partialOrEmptyString>
		// the partial or empty string is passed to us from the completion script. Empty if the user just typed space and pressed tab, partial if the user pressed tab after typing part of something. bash will filter that for us via compgen
		// so we need to remove the first two args and the last arg leaving <args>...
		args := os.Args[2 : len(os.Args)-1]
		partiallyTypedArg := os.Args[len(os.Args)-1]

		candidates, err := app.Complete(args, partiallyTypedArg, opts...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		completion.BashCompletionsWrite(os.Stdout, candidates)

	} else if len(os.Args) >= 3 && os.Args[1] == "--completion-zsh" {
		// app --completion-zsh <args> . Note that <args> must be something, even if it's the empty string

		// parseOpts.Args looks like: <exe> --completion-zsh <args>... <partialOrEmptyString>
		// the partial or empty string is passed to us from the completion script. Empty if the user just typed space and pressed tab, partial if the user pressed tab after typing part of something. zsh will filter that for us
		// so we need to remove the first two args and the last arg leaving <args>...
		args := os.Args[2 : len(os.Args)-1]
		partiallyTypedArg := os.Args[len(os.Args)-1]

		candidates, err := app.Complete(args, partiallyTypedArg, opts...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		completion.ZshCompletionsWrite(os.Stdout, candidates)

	} else if len(os.Args) >= 3 && os.Args[1] == "--completion-fish" {
		// app --completion-fish <args> . Note that <args> must be something, even if it's the empty string

		// parseOpts.Args looks like: <exe> --completion-fish <args>... <partialOrEmptyString>
		// the partial or empty string is passed to us from the completion script. Empty if the user just typed space and pressed tab, partial if the user pressed tab after typing part of something. fish will filter that for us
		// so we need to remove the first two args and the last arg leaving <args>...
		args := os.Args[2 : len(os.Args)-1]
		partiallyTypedArg := os.Args[len(os.Args)-1]

		candidates, err := app.Complete(args, partiallyTypedArg, opts...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		completion.FishCompletionsWrite(os.Stdout, candidates)

	} else {
		app.MustRunWithArgs(os.Args[1:], opts...)
	}
}

// LookupEnv is a function for resolving environment variable names to values.
// Satisfiable by [os.LookupEnv] or [LookupMap].
type LookupEnv func(key string) (string, bool)

// LookupMap returns a [LookupEnv] backed by a static map. Useful for mocking
// environment variables in tests.
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
			errs = append(errs, colerr.NewWrappedf(nil, "Flag and alias names must start with '-': %s", fmt.Sprintf("%#v", name)))
		}
		if count > 1 {
			errs = append(errs, colerr.NewWrappedf(nil, "Flag or alias name exists %s times: %s", fmt.Sprintf("%d", count), fmt.Sprintf("%v", name)))
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
		return errors.New("HelpFlagName must be set")
	}
	helpFlag, exists := app.GlobalFlags[app.HelpFlagName]
	if !exists {
		return colerr.NewWrappedf(nil, "HelpFlagName not found in GlobalFlags: %s", fmt.Sprintf("%v", app.HelpFlagName))
	}
	helpFlagValEmpty, ok := helpFlag.EmptyValueConstructor().(value.ScalarValue)
	if !ok {
		return colerr.NewWrappedf(nil, "HelpFlagName must be a scalar: %s", fmt.Sprintf("%v", app.HelpFlagName))
	}
	if _, ok := helpFlagValEmpty.Get().(string); !ok {
		return colerr.NewWrappedf(nil, "HelpFlagName must be a string: %s", fmt.Sprintf("%v", app.HelpFlagName))
	}
	if !helpFlagValEmpty.HasDefault() {
		return colerr.NewWrappedf(nil, "HelpFlagName must have a default value: %s", fmt.Sprintf("%v", app.HelpFlagName))
	}
	if !slices.Equal(helpFlagValEmpty.Choices(), app.HelpCmds.SortedNames()) {
		return colerr.NewWrappedf(nil, "HelpFlagName choices must match HelpCmds: %s", fmt.Sprintf("%v", app.HelpFlagName))
	}
	if !slices.Contains(helpFlagValEmpty.Choices(), helpFlagValEmpty.DefaultString()) {
		return colerr.NewWrappedf(nil, "HelpFlagName default value (%s) must be in choices (%s): %s", fmt.Sprintf("%v", helpFlagValEmpty.DefaultString()), fmt.Sprintf("%v", helpFlagValEmpty.Choices()), fmt.Sprintf("%v", app.HelpFlagName))
	}

	// validate --config flag
	if app.ConfigFlagName != "" {
		if app.NewConfigReader == nil {
			return colerr.NewWrappedf(nil, "ConfigFlagName must have a NewConfigReader: %s", fmt.Sprintf("%v", app.ConfigFlagName))
		}
		configFlag, exists := app.GlobalFlags[app.ConfigFlagName]
		if !exists {
			return colerr.NewWrappedf(nil, "ConfigFlagName not found in GlobalFlags: %s", fmt.Sprintf("%v", app.ConfigFlagName))
		}
		configFlagValEmpty, ok := configFlag.EmptyValueConstructor().(value.ScalarValue)
		if !ok {
			return colerr.NewWrappedf(nil, "ConfigFlagName must be a scalar: %s", fmt.Sprintf("%v", app.ConfigFlagName))
		}
		if _, ok := configFlagValEmpty.Get().(path.Path); !ok {
			return colerr.NewWrappedf(nil, "ConfigFlagName must be a path: %s", fmt.Sprintf("%v", app.ConfigFlagName))
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
			return colerr.NewWrappedf(nil, "Section names must not start with '-': %s", fmt.Sprintf("%#v", secName))
		}

		// Sections must not be leaf nodes
		if flatSec.Sec.Sections.Empty() && flatSec.Sec.Cmds.Empty() {
			return colerr.NewWrappedf(nil, "Sections must have either child sections or child commands: %s", fmt.Sprintf("%#v", secName))
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
					errs = append(errs, colerr.NewWrappedf(nil, "Command and section name clash: %s", name))
				}
			}
			err := errors.Join(errs...)
			if err != nil {
				return colerr.NewWrapped(err, "Name collision")
			}
		}

		for name, com := range flatSec.Sec.Cmds {

			// Commands must not start wtih "-"
			if strings.HasPrefix(string(name), "-") {
				return colerr.NewWrappedf(nil, "Command names must not start with '-': %s", fmt.Sprintf("%#v", name))
			}

			err := validateFlags(app.GlobalFlags, com.Flags)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// CompletionsFunc generates shell completion candidates for a flag value.
// Use the Completions* convenience functions (e.g., [CompletionsValues]) to construct one.
type CompletionsFunc func(CmdContext) (*completion.Candidates, error)

// Complete generates tab-completion candidates for the given args and partially typed argument.
// args should not include the program name. partiallyTypedArg is the in-progress token
// the user is typing (may be empty). This is called internally by [App.MustRun] to support
// shell completion scripts.
func (app *App) Complete(args []string, partiallyTypedArg string, opts ...ParseOpt) (*completion.Candidates, error) {
	parseOpts := NewParseOpts(opts...)

	// I could to a full parse here, but that would be slower and more prone to failure than just parsing the args - we don't need a lot of info to complete section/command names
	parseState, err := app.parseArgs(args)
	if err != nil {
		return nil, colerr.NewWrapped(err, "Unexpected parseArgs err")
	}

	// special case if help is passed
	if parseState.HelpPassed {
		// if the value of the flag has been passed, don't suggest anything
		if parseState.FlagValues[app.HelpFlagName].UpdatedBy() == value.UpdatedByFlag {
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
		for _, name := range app.HelpCmds.SortedNames() {
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
			Name:        app.HelpFlagName,
			Description: app.GlobalFlags[app.HelpFlagName].HelpShort,
		})
		return &ret, nil
	}

	// Finish the parse!
	err = app.resolveFlags(parseState.CurrentCmd, parseState.FlagValues, parseOpts.LookupEnv, parseState.UnsetFlagNames)
	if err != nil {
		return nil, colerr.NewWrapped(err, "Unexpected resolveFlags err")
	}
	cmdContext := CmdContext{
		App:           app,
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
		return nil, colerr.NewWrappedf(nil, "Unexpected ParseState: %s", fmt.Sprintf("%v", parseState.ParseArgState))
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
		return colerr.NewWrapped(err, "Could not set readline config")
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
			err = colerr.NewWrappedf(err, "Could not parse args for completion: args: %s", fmt.Sprintf("%v", words))
			return readline.CompleteMessage(err.Error())
		}

		// Completions are accepted with " ", so if we don't end in " ", it's a partial word.
		var partiallyTypedArg string
		if len(lineStr) != 0 && !strings.HasSuffix(lineStr, " ") {
			words = words[:len(words)-1]
			// TODO: this makes it panic. I need to extract this line, cursor -> args, partiallyTypedArg logic into a separate function and test it. For now, just don't don't try to get the partiallyTypedArg
			// partiallyTypedArg = words[len(words)-1]
		}
		// fmt.Fprintf(os.Stderr, "words: %#v\n", words)

		// TODO: should I copy parseOpts from cmdCtx?
		candidates, err := cmdCtx.App.Complete(words, partiallyTypedArg)
		if err != nil {
			err = colerr.NewWrappedf(err, "Could not get completions: args: %s", fmt.Sprintf("%v", words))
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
			return colerr.NewWrapped(err, "Could not read line")
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
