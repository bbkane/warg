package jsonreader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"go.bbkane.com/warg/config"
	"go.bbkane.com/warg/config/tokenize"
)

// This must be a type declaration or the tests fail
type configMap = map[string]interface{}

type jsonConfigReader struct {
	data configMap
}

func New(filePath string) (config.Reader, error) {
	cr := &jsonConfigReader{
		data: nil,
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		// file not existing is ok
		// TODO: explicitly check for file not found instead of any error
		return cr, nil
	}

	// https://stackoverflow.com/a/16946478/2958070
	d := json.NewDecoder(bytes.NewBuffer(content))
	d.UseNumber()
	err = d.Decode(&cr.data)

	if err != nil {
		return nil, err
	}
	return cr, nil
}

func (cr *jsonConfigReader) Search(path string) (*config.SearchResult, error) {
	data := cr.data
	tokens, err := tokenize.Tokenize(path)
	if err != nil {
		return nil, err
	}

	lenTokens := len(tokens)
	var current interface{} = data
	for i, token := range tokens {
		if i == lenTokens-2 && token.Type == tokenize.TokenTypeSlice {
			// we're at the second to last token, so current should also be a slice
			// cast it to a slice, then get all keys from it!
			sliceOfDicts, ok := current.([]interface{})
			if !ok {
				return nil, fmt.Errorf(
					"expecting []interface{}: \n  actual type %T\n  actual value: %#v\n   path: %v\n  token: %v",
					current, current, path, token,
				)
			}
			finalToken := tokens[lenTokens-1]
			if finalToken.Type != tokenize.TokenTypeKey {
				return nil, fmt.Errorf(
					"expected TokenTypeKey for last element: path: %v: token: %v",
					path,
					token,
				)
			}
			var ret []interface{}
			for _, e := range sliceOfDicts {
				cm, ok := e.(configMap)
				if !ok {
					return nil, fmt.Errorf(
						"expecting ConfigMap: \n  actual type %T\n  actual value: %#v\n  path: %v\n  token: %v",
						current, current, path, token,
					)
				}
				val, exists := cm[finalToken.Text]
				if !exists {
					return nil, fmt.Errorf(
						"for the slice operator, ALL elements must contain the key: path: %v: key: %v",
						path, finalToken.Text,
					)
				}
				ret = append(ret, val)
			}
			return &config.SearchResult{IFace: ret, IsAggregated: true}, nil
		} else {
			// outside the special case, we should be able to just index into this thing, and loop again
			// or, if it's the last one, return
			if token.Type != tokenize.TokenTypeKey {
				return nil, fmt.Errorf(
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
				return nil, fmt.Errorf(
					"expecting ConfigMap: \n  actual type %T\n  actual value: %#v\n  path: %v\n  token: %v",
					current, current, path, token,
				)
			}

			// dumbest thing ever, but it appears I cannot reassign current and assing exists with the same statement
			next, exists := currentMap[token.Text]
			current = next
			if !exists {
				return nil, nil
			}
		}
	}
	return &config.SearchResult{IFace: current, IsAggregated: false}, nil
}
