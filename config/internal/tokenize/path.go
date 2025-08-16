package tokenize

import (
	"fmt"
	"strings"
)

// TODO: make these their own type
const TokenTypeKey = "tokenTypeKey"
const TokenTypeSlice = "tokenTypeSlice"

type Token struct {
	Text string
	Type string
}

// Tokenize a configpath into a list of tokens
func Tokenize(path string) ([]Token, error) {
	// TODO: make this better :) - I'm only checking for [] at the end of strings and not looking for escaped dots or anything
	pathElements := strings.Split(path, ".")
	lenPathElemennts := len(pathElements)

	var tokens []Token
	for i, el := range pathElements {
		if strings.HasSuffix(el, "[]") {
			if i != lenPathElemennts-2 {
				return nil, fmt.Errorf("[] is only allowed as an element before the last element: path: %#v", path)
			}
			tokens = append(tokens, Token{Text: el[:len(el)-2], Type: TokenTypeKey})
			tokens = append(tokens, Token{Text: "[]", Type: TokenTypeSlice})
		} else {
			tokens = append(tokens, Token{Text: el, Type: TokenTypeKey})
		}
	}
	return tokens, nil
}
