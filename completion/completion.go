package completion

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
)

// ZshCompletionScript holds the embedded zsh completion script template.
//
//go:embed completion_script.zsh
var ZshCompletionScript []byte

// ZshCompletionScriptWrite writes the zsh completion script for appName to w.
func ZshCompletionScriptWrite(w io.Writer, appName string) {
	script := bytes.ReplaceAll(ZshCompletionScript, []byte("WARG_COMPLETION_APPNAME"), []byte(appName))
	_, err := w.Write(script)
	if err != nil {
		panic("unexpected CompletionScriptZsh err " + err.Error())
	}
}

// BashCompletionScript holds the embedded bash completion script template.
//
//go:embed completion_script.bash
var BashCompletionScript []byte

// BashCompletionScriptWrite writes the bash completion script for appName to w.
func BashCompletionScriptWrite(w io.Writer, appName string) {
	script := bytes.ReplaceAll(BashCompletionScript, []byte("WARG_COMPLETION_APPNAME"), []byte(appName))
	_, err := w.Write(script)
	if err != nil {
		panic("unexpected CompletionScriptBash err " + err.Error())
	}
}

// FishCompletionScript holds the embedded fish completion script template.
//
//go:embed completion_script.fish
var FishCompletionScript []byte

// FishCompletionScriptWrite writes the fish completion script for appName to w.
func FishCompletionScriptWrite(w io.Writer, appName string) {
	script := bytes.ReplaceAll(FishCompletionScript, []byte("WARG_COMPLETION_APPNAME"), []byte(appName))
	_, err := w.Write(script)
	if err != nil {
		panic("unexpected CompletionScriptFish err " + err.Error())
	}
}

// FishCompletionsWrite writes completion candidates in fish shell format to w.
func FishCompletionsWrite(w io.Writer, c *Candidates) {
	fmt.Fprintln(w, c.Type)
	switch c.Type {
	case Type_Directories, Type_DirectoriesFiles, Type_None:
		// nothing else needed
		return
	case Type_Values:
		for _, v := range c.Values {
			fmt.Fprintln(w, v.Name)
		}
	case Type_ValuesDescriptions:
		// fish uses tab-separated name\tdescription format
		for _, v := range c.Values {
			if v.Description != "" {
				fmt.Fprintf(w, "%s\t%s\n", v.Name, v.Description)
			} else {
				fmt.Fprintln(w, v.Name)
			}
		}
	default:
		panic("unexpected completion type: " + string(c.Type))
	}
}

// BashCompletionsWrite writes completion candidates in bash shell format to w.
func BashCompletionsWrite(w io.Writer, c *Candidates) {
	fmt.Fprintln(w, c.Type)
	switch c.Type {
	case Type_Directories, Type_DirectoriesFiles, Type_None:
		// nothing else needed
		return
	case Type_Values, Type_ValuesDescriptions:
		// bash doesn't support descriptions, so just output names for both types
		for _, v := range c.Values {
			fmt.Fprintln(w, v.Name)
		}
	default:
		panic("unexpected completion type: " + string(c.Type))
	}
}

// ZshCompletionsWrite writes completion candidates in zsh shell format to w.
func ZshCompletionsWrite(w io.Writer, c *Candidates) {
	fmt.Fprintln(w, c.Type)
	switch c.Type {
	case Type_Directories, Type_DirectoriesFiles, Type_None:
		// nothing else needed
		return
	case Type_Values:
		for _, v := range c.Values {
			fmt.Fprintln(w, v.Name)
		}
	case Type_ValuesDescriptions:
		for _, v := range c.Values {
			fmt.Fprintln(w, v.Name)
			if v.Description != "" {
				fmt.Fprintln(w, v.Name+" - "+v.Description)
			} else {
				fmt.Fprintln(w, v.Name)
			}

		}
	default:
		panic("unexpected completion type: " + string(c.Type))
	}
}

// Type indicates what kind of completion the shell should perform.
type Type string

const (
	Type_Directories        Type = "COMPLETION_TYPE_DIRECTORIES"
	Type_DirectoriesFiles   Type = "COMPLETION_TYPE_DIRECTORIES_FILES"
	Type_None               Type = "COMPLETION_TYPE_NONE"
	Type_Values             Type = "COMPLETION_TYPE_VALUES"
	Type_ValuesDescriptions Type = "COMPLETION_TYPE_VALUES_DESCRIPTIONS"
)

// Candidate is a single tab-completion suggestion with an optional description.
type Candidate struct {
	Name        string
	Description string
}

// Candidates holds the completion type and a list of candidate values
// to present to the user during tab completion.
type Candidates struct {
	Type   Type
	Values []Candidate
}
