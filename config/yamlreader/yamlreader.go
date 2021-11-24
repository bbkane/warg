package yamlreader

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/bbkane/warg/config"
	"gopkg.in/yaml.v2"
)

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

type configMap = map[interface{}]interface{}

type yamlConfigReader struct {
	data configMap
}

func New(filePath string) (config.Reader, error) {
	cr := &yamlConfigReader{}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		// the file not existing is ok
		return cr, nil
	}

	err = yaml.Unmarshal(content, &cr.data)
	if err != nil {
		return nil, err
	}
	// fmt.Printf("%#v\n", m)
	return cr, nil

}

func (cr *yamlConfigReader) Search(path string) (config.SearchResult, error) {
	data := cr.data
	tokens, err := tokenize(path)
	if err != nil {
		return config.SearchResult{}, err
	}

	lenTokens := len(tokens)
	var current interface{} = data
	for i, token := range tokens {
		if i == lenTokens-2 && token.Type == tokenTypeSlice {
			// we're at the second to last token, so current should also be a slice
			// cast it to a slice, then get all keys from it!
			sliceOfDicts, ok := current.([]interface{})
			if !ok {
				return config.SearchResult{}, fmt.Errorf(
					"expecting []interface{}: \n  actual type %T\n  actual value: %#v\n   path: %v\n  token: %v",
					current, current, path, token,
				)
			}
			finalToken := tokens[lenTokens-1]
			if finalToken.Type != tokenTypeKey {
				return config.SearchResult{}, fmt.Errorf(
					"expected TokenTypeKey for last element: path: %v: token: %v",
					path,
					token,
				)
			}
			var ret []interface{}
			for _, e := range sliceOfDicts {
				cm, ok := e.(configMap)
				if !ok {
					return config.SearchResult{}, fmt.Errorf(
						"expecting ConfigMap: \n  actual type %T\n  actual value: %#v\n  path: %v\n  token: %v",
						current, current, path, token,
					)
				}
				val, exists := cm[finalToken.Text]
				if !exists {
					return config.SearchResult{}, fmt.Errorf(
						"for the slice operator, ALL elements must contain the key: path: %v: key: %v",
						path, finalToken.Text,
					)
				}
				ret = append(ret, val)
			}
			return config.SearchResult{IFace: ret, Exists: true, IsAggregated: true}, nil
		} else {
			// outside the special case, we should be able to just index into this thing, and loop again
			// or, if it's the last one, return
			if token.Type != tokenTypeKey {
				return config.SearchResult{}, fmt.Errorf(
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
				return config.SearchResult{}, fmt.Errorf(
					"expecting ConfigMap: \n  actual type %T\n  actual value: %#v\n  path: %v\n  token: %v",
					current, current, path, token,
				)
			}

			// dumbest thing ever, but it appears I cannot reassign current and assing exists with the same statement
			next, exists := currentMap[token.Text]
			current = next
			if !exists {
				return config.SearchResult{}, nil
			}
		}
	}
	return config.SearchResult{IFace: current, Exists: true, IsAggregated: false}, nil
}
