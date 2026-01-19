package config

// SearchResult is the result of trying to search through a config for a value
type SearchResult struct {
	// IFace holds the value if found
	IFace interface{}
}

// Reader searches with a path to try to get a config value.
// if the result is nil, nil, it means the result wasn't found
type Reader interface {
	Search(path string) (*SearchResult, error)
}

// NewReader constructs a ConfigReader
type NewReader func(filePath string) (Reader, error)
