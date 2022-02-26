package config

// SearchResult is the result of trying to search through a config for a value
type SearchResult struct {
	// IFace holds the value if found
	IFace interface{}
	// Exists indicates whether the value was found or not
	Exists bool
	// IsAggregated indicates whether the value was stitched together from child config elements
	// For example, consider a config with the following content: {"subreddits": [{"name": "earthporn"}, {"name": "wallpapers"}]}
	// We can get all the names with a path like "subreddits[].name". These results are aggregated into a list, and IsAggregated will be set to true
	IsAggregated bool
}

// Reader searches with a path to try to get a config value.
type Reader interface {
	Search(path string) (SearchResult, error)
}

// NewReader constructs a ConfigReader
type NewReader func(filePath string) (Reader, error)
