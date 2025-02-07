package warg

import (
	"io"
	"text/template"
)

// writeCompletionScriptZsh prints a completions script to add to fpath
func writeCompletionScriptZsh(w io.Writer, appName string) {
	script := `#compdef {{.Name}}

date >> ~/_{{.Name}}_completion.log

local -a comp_values=()
local -a comp_descriptions=()
local -a comp_type

local state="expecting_type"
local line
# TODO: switch on comp_type to switch on the rest of the lines
while IFS= read -r line && [[ -n "$line" ]]; do
    if [[ "$state" == "expecting_type" ]]; then
        comp_type="$line"
        state="expecting_value"
    elif [[ "$state" == "expecting_value" ]]; then
        comp_values+=("$line")
        state="expecting_description"
    else
        comp_descriptions+=("$line")
        state="expecting_value"
    fi
done < <(${words[1]} --completion-bash "${(@)words[2,$CURRENT]}")

echo "$comp_values" >> ~/_{{.Name}}.log

compadd -d comp_descriptions -a comp_values
`
	scriptParams := struct {
		Name string
	}{
		Name: appName,
	}
	t := template.Must(template.New("template").Parse(script))

	err := t.Execute(w, scriptParams)
	if err != nil {
		panic("unexpected CompletionScriptZsh err " + err.Error())
	}
}

type CompletionType string

const CompletionType_ValueDescription CompletionType = "WARG_COMPLETION_TYPE_VALUE_DESCRIPTION"

type CompletionCandidate struct {
	Name        string
	Description string
}

type CompletionCandidates struct {
	Type   CompletionType
	Values []CompletionCandidate
}

// CompletionBashCandidates returns a list of possible completions for the current state of the parse
// TODO: descriptions lol
func (a *App) CompletionBashCandidates(args []string) CompletionCandidates {
	// TODO: maybe make an error state instead of returning an empty completion on error
	pr, err := a.parseArgs(args)
	if err != nil {
		return CompletionCandidates{
			Type:   CompletionType_ValueDescription,
			Values: []CompletionCandidate{},
		}
	}
	switch pr.State {
	case Parse_ExpectingSectionOrCommand:
		// TODO: this should produce CompletionCandidates natively instead of copying like this
		childrenNames := pr.CurrentSection.ChildrenNames()
		candidates := make([]CompletionCandidate, len(childrenNames))
		for i, name := range childrenNames {
			// TODO: get the description too
			candidates[i] = CompletionCandidate{Name: name, Description: ""}
		}
		return CompletionCandidates{
			Type:   CompletionType_ValueDescription,
			Values: candidates,
		}

	// case Parse_ExpectingFlagNameOrEnd:
	// 	// TODO: if a scalar flag has been passsed, don't suggest it again
	// 	cmdChildren := pr.CurrentCommand.ChildrenNames()
	// 	appFlagNames := a.GlobalFlags.SortedNames()
	// 	combined := slices.Concat(cmdChildren, appFlagNames)
	// 	candidates := make([]CompletionCandidate, len(combined))
	// 	for i, name := range combined {
	// 		candidates[i] = CompletionCandidate{Name: name}
	// 	}
	// 	return CompletionCandidates{
	// 		Type:   CompletionType_ValueDescription,
	// 		Values: candidates,
	// 	}

	// case Parse_ExpectingFlagValue:
	// 	err = a.resolveFlags(pr.CurrentCommand, pr.FlagValues)
	// 	if err != nil {
	// 		return CompletionCandidates{
	// 			Type:   CompletionType_ValueDescription,
	// 			Values: []CompletionCandidate{},
	// 		}
	// 	}

	// 	pf := pr.FlagValues.ToPassedFlags()
	// 	ctx := Context{
	// 		App:         a,
	// 		ParseResult: &pr,
	// 		Flags:       pf,
	// 	}
	// 	childrenNames := pr.CurrentFlag.ChildrenNames(ctx)
	// 	candidates := make([]CompletionCandidate, len(childrenNames))
	// 	for i, name := range childrenNames {
	// 		candidates[i] = CompletionCandidate{Name: name}
	// 	}
	// 	return CompletionCandidates{
	// 		Type:   CompletionType_ValueDescription,
	// 		Values: candidates,
	// 	}
	default:
		panic("unexpected state: " + pr.State)
	}

}
