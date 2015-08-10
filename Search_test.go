package googlesearch

import "testing"

func TestSearch(t *testing.T) {

	results := Search("golang")
	if results == nil {
		t.Errorf("Search(\"golang\") didn't go well")
	}
	if len(results) == 0 {
		t.Errorf("Search(\"golang\") should return at least one result")
	}
	if len(results) < 10 {
		t.Errorf("Search(\"golang\") should return 10+ results, got %d", len(results))
	}
}
