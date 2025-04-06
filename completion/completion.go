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
	Type_ValueDescription Type = "COMPLETION_TYPE_VALUE_DESCRIPTION"
	Type_DirectoriesFiles Type = "COMPLETION_TYPE_DIRECTORIES_FILES"
)

type Candidate struct {
	Name        string
	Description string
}

type Candidates struct {
	Type   Type
	Values []Candidate
}
