// Package config defines the interface for reading flag values from configuration files.
package config

// SearchResult holds a value found at a config path.
type SearchResult struct {
	// IFace holds the decoded value (type depends on the config format).
	IFace interface{}
}

// Reader searches a parsed config file for values at a given dot-separated path.
// Returns nil, nil if the path is not found in the config.
type Reader interface {
	Search(path string) (*SearchResult, error)
}

// NewReader constructs a [Reader] from a file path. Returns an error if the file
// cannot be read or parsed. A non-existent file is not an error (returns a Reader
// that finds nothing).
type NewReader func(filePath string) (Reader, error)
