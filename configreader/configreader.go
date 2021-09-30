package configreader

type ConfigSearchResult struct {
	IFace        interface{}
	Exists       bool
	IsAggregated bool
}

type ConfigReader interface {
	Search(path string) (ConfigSearchResult, error)
}

type NewConfigReader = func(filePath string) (ConfigReader, error)

type ConfigReaderFunc func(path string) (ConfigSearchResult, error)

func (f ConfigReaderFunc) Search(path string) (ConfigSearchResult, error) {
	return f(path)
}

// NOTE: see https://karthikkaranth.me/blog/functions-implementing-interfaces-in-go/
// for how to use ConfigReaderFunc in tests
