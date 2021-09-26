package configpath

import (
	"fmt"
	"strings"
)

type ConfigMap = map[string]interface{}

const TokenTypeKey = "TokenTypeKey"
const TokenTypeSlice = "TokenTypeSlice"

type Token struct {
	Text string
	Type string
}

func tokenize(path string) ([]Token, error) {
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

type FollowPathResult struct {
	IFace      interface{}
	Exists     bool
	Aggregated bool
}

// FollowPath takes a map and a path with elements separated by dots
// and retrieves the interface at the end of it. If the interface
// doesn't exist, then the bool value is false
func FollowPath(data ConfigMap, path string) (FollowPathResult, error) {
	tokens, err := tokenize(path)
	if err != nil {
		return FollowPathResult{}, err
	}

	lenTokens := len(tokens)
	var current interface{} = data
	for i, token := range tokens {
		if i == lenTokens-2 && token.Type == TokenTypeSlice {
			// we're at the second to last token, so current should also be a slice
			// cast it to a slice, then get all keys from it!
			sliceOfDicts, ok := current.([]interface{})
			if !ok {
				return FollowPathResult{}, fmt.Errorf(
					"expecting []interface{}: \n  actual type %T\n  actual value: %#v\n   path: %v\n  token: %v",
					current, current, path, token,
				)
			}
			finalToken := tokens[lenTokens-1]
			if finalToken.Type != TokenTypeKey {
				return FollowPathResult{}, fmt.Errorf(
					"expected TokenTypeKey for last element: path: %v: token: %v",
					path,
					token,
				)
			}
			var ret []interface{}
			for _, e := range sliceOfDicts {
				cm, ok := e.(ConfigMap)
				if !ok {
					return FollowPathResult{}, fmt.Errorf(
						"expecting ConfigMap: \n  actual type %T\n  actual value: %#v\n  path: %v\n  token: %v",
						current, current, path, token,
					)
				}
				val, exists := cm[finalToken.Text]
				if !exists {
					return FollowPathResult{}, fmt.Errorf(
						"for the slice operator, ALL elements must contain the key: path: %v: key: %v",
						path, finalToken.Text,
					)
				}
				ret = append(ret, val)
			}
			return FollowPathResult{IFace: ret, Exists: true, Aggregated: true}, nil
		} else {
			// outside the special case, we should be able to just index into this thing, and loop again
			// or, if it's the last one, return
			if token.Type != TokenTypeKey {
				return FollowPathResult{}, fmt.Errorf(
					"expected TokenTypeKey for last element: path: %v: token: %v",
					path,
					token,
				)
			}
			currentMap, ok := current.(ConfigMap)
			if !ok {
				return FollowPathResult{}, fmt.Errorf(
					"expecting ConfigMap: \n  actual type %T\n  actual value: %#v\n  path: %v\n  token: %v",
					current, current, path, token,
				)
			}
			// dumbest thing ever, but it appears I reassign current and assing exists with the same statement
			next, exists := currentMap[token.Text]
			current = next
			if !exists {
				return FollowPathResult{}, nil
			}
		}
	}
	return FollowPathResult{IFace: current, Exists: true, Aggregated: false}, nil
}
