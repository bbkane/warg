package completion

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
)

//go:embed completion_script.zsh
var ZshCompletionScript []byte

func ZshCompletionScriptWrite(w io.Writer, appName string) {
	script := bytes.ReplaceAll(ZshCompletionScript, []byte("WARG_COMPLETION_APPNAME"), []byte(appName))
	_, err := w.Write(script)
	if err != nil {
		panic("unexpected CompletionScriptZsh err " + err.Error())
	}
}

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

type Type string

const (
	Type_Directories        Type = "COMPLETION_TYPE_DIRECTORIES"
	Type_DirectoriesFiles   Type = "COMPLETION_TYPE_DIRECTORIES_FILES"
	Type_None               Type = "COMPLETION_TYPE_NONE"
	Type_Values             Type = "COMPLETION_TYPE_VALUES"
	Type_ValuesDescriptions Type = "COMPLETION_TYPE_VALUES_DESCRIPTIONS"
)

type Candidate struct {
	Name        string
	Description string
}

type Candidates struct {
	Type   Type
	Values []Candidate
}
