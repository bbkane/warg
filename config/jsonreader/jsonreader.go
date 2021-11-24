package jsonreader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/bbkane/warg/config"
)

type jsonConfigReader struct {
	data configMap
}

func NewJSONConfigReader(filePath string) (config.ConfigReader, error) {
	cr := &jsonConfigReader{}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		// file not existing is ok
		// TODO: explicitly check for file not found instead of any error
		return cr, nil
	}

	err = json.Unmarshal(content, &cr.data)
	if err != nil {
		return nil, err
	}
	return cr, nil
}

type configMap = map[string]interface{}

const tokenTypeKey = "tokenTypeKey"
const tokenTypeSlice = "tokenTypeSlice"

type token struct {
	Text string
	Type string
}

func tokenize(path string) ([]token, error) {
	// TODO: make this better :) - I'm only checking for [] at the end of strings and not looking for escaped dots or anything
	pathElements := strings.Split(path, ".")
	lenPathElemennts := len(pathElements)

	var tokens []token
	for i, el := range pathElements {
		if strings.HasSuffix(el, "[]") {
			if i != lenPathElemennts-2 {
				return nil, fmt.Errorf("[] is only allowed as an element before the last element: path: %#v", path)
			}
			tokens = append(tokens, token{Text: el[:len(el)-2], Type: tokenTypeKey})
			tokens = append(tokens, token{Text: "[]", Type: tokenTypeSlice})
		} else {
			tokens = append(tokens, token{Text: el, Type: tokenTypeKey})
		}
	}
	return tokens, nil
}

func (cr *jsonConfigReader) Search(path string) (config.ConfigSearchResult, error) {
	data := cr.data
	tokens, err := tokenize(path)
	if err != nil {
		return config.ConfigSearchResult{}, err
	}

	lenTokens := len(tokens)
	var current interface{} = data
	for i, token := range tokens {
		if i == lenTokens-2 && token.Type == tokenTypeSlice {
			// we're at the second to last token, so current should also be a slice
			// cast it to a slice, then get all keys from it!
			sliceOfDicts, ok := current.([]interface{})
			if !ok {
				return config.ConfigSearchResult{}, fmt.Errorf(
					"expecting []interface{}: \n  actual type %T\n  actual value: %#v\n   path: %v\n  token: %v",
					current, current, path, token,
				)
			}
			finalToken := tokens[lenTokens-1]
			if finalToken.Type != tokenTypeKey {
				return config.ConfigSearchResult{}, fmt.Errorf(
					"expected TokenTypeKey for last element: path: %v: token: %v",
					path,
					token,
				)
			}
			var ret []interface{}
			for _, e := range sliceOfDicts {
				cm, ok := e.(configMap)
				if !ok {
					return config.ConfigSearchResult{}, fmt.Errorf(
						"expecting ConfigMap: \n  actual type %T\n  actual value: %#v\n  path: %v\n  token: %v",
						current, current, path, token,
					)
				}
				val, exists := cm[finalToken.Text]
				if !exists {
					return config.ConfigSearchResult{}, fmt.Errorf(
						"for the slice operator, ALL elements must contain the key: path: %v: key: %v",
						path, finalToken.Text,
					)
				}
				ret = append(ret, val)
			}
			return config.ConfigSearchResult{IFace: ret, Exists: true, IsAggregated: true}, nil
		} else {
			// outside the special case, we should be able to just index into this thing, and loop again
			// or, if it's the last one, return
			if token.Type != tokenTypeKey {
				return config.ConfigSearchResult{}, fmt.Errorf(
					"expected TokenTypeKey for last element: path: %v: token: %v",
					path,
					token,
				)
			}

			// JSON needs this
			currentMap, ok := current.(configMap)

			// YAML needs this
			// currentMap, ok := current.(map[interface{}]interface{})

			// Sticking with JSON now because it works and is sprinkled over the code :)
			// but see ~/warg_configreader.md - I'm going to create a new package to do that

			if !ok {
				return config.ConfigSearchResult{}, fmt.Errorf(
					"expecting ConfigMap: \n  actual type %T\n  actual value: %#v\n  path: %v\n  token: %v",
					current, current, path, token,
				)
			}

			// dumbest thing ever, but it appears I cannot reassign current and assing exists with the same statement
			next, exists := currentMap[token.Text]
			current = next
			if !exists {
				return config.ConfigSearchResult{}, nil
			}
		}
	}
	return config.ConfigSearchResult{IFace: current, Exists: true, IsAggregated: false}, nil
}
