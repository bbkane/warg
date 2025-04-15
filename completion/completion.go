package completion

import (
	"bytes"
	_ "embed"
	"io"
)

//go:embed completion_script.zsh
var CompletionScriptZsh []byte

func WriteCompletionScriptZsh(w io.Writer, appName string) {
	script := bytes.ReplaceAll(CompletionScriptZsh, []byte("WARG_COMPLETION_APPNAME"), []byte(appName))
	_, err := w.Write(script)
	if err != nil {
		panic("unexpected CompletionScriptZsh err " + err.Error())
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
