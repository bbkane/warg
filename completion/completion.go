package completion

import (
	"io"
	"text/template"
)

// WriteCompletionScriptZsh prints a completions script to add to fpath
func WriteCompletionScriptZsh(w io.Writer, appName string) {
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
done < <(${words[1]} --completion-zsh "${(@)words[2,$CURRENT]}")

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
