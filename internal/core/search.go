package core

type (
	SearchEngine interface {
		Search(query string) ([]string, error)
	}

	searchEngine struct {
	}
)

func NewSearchEngine(store Store) SearchEngine {
	return &searchEngine{}
}

func (se *searchEngine) Search(query string) ([]string, error) {
	// Implement search logic here
	return []string{"result1", "result2"}, nil
}
