package config

// ConfigSearchResult is the result of trying to search through a config for a value
type ConfigSearchResult struct {
	// IFace holds the value if found
	IFace interface{}
	// Exists indicates whether the value was found or not
	Exists bool
	// IsAggregated indicates whether the value was stitched together from child config elements
	// For example, consider a config with the following content: {"subreddits": [{"name": "earthporn"}, {"name": "wallpapers"}]}
	// We can get all the names with a path like "subreddits[].name". These results are aggregated into a list, and IsAggregated will be set to true
	IsAggregated bool
}

// ConfigReader searches with a path to try to get a config value.
type ConfigReader interface {
	Search(path string) (ConfigSearchResult, error)
}

// NewConfigReader constructs a ConfigReader
type NewConfigReader = func(filePath string) (ConfigReader, error)
