package yamlreader

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"go.bbkane.com/warg/config"
	"go.bbkane.com/warg/config/internal/tokenize"
)

// go-yaml uses string keys by default
type configMap = map[string]interface{}

type yamlConfigReader struct {
	data configMap
}

func New(filePath string) (config.Reader, error) {
	cr := &yamlConfigReader{
		data: nil,
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		// the file not existing is ok
		return cr, nil
	}

	// from what I can tell, go-yaml automatically uses map[string]interface{} for maps
	// I tried adding an integer key to the test YAML file, and it did not blow up
	err = yaml.UnmarshalWithOptions(content, &cr.data, yaml.Strict())
	if err != nil {
		return nil, err
	}
	// fmt.Printf("%#v\n", m)
	return cr, nil

}

func (cr *yamlConfigReader) Search(path string) (*config.SearchResult, error) {
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
						"expecting map[string]interface{}: \n  actual type %T\n  actual value: %#v\n  path: %v\n  token: %v",
						e, e, path, token,
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

			currentMap, ok := current.(configMap)

			if !ok {
				return nil, fmt.Errorf(
					"expecting map[string]interface{}: \n  actual type %T\n  actual value: %#v\n  path: %v\n  token: %v",
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

	if currentConfigMap, ok := current.(configMap); ok {
		current = currentConfigMap
	}
	return &config.SearchResult{IFace: current, IsAggregated: false}, nil
}
