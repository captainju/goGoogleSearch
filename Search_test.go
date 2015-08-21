package googlesearch

import (
	"testing"
)

func TestSearch(t *testing.T) {

	resultChan := Search("\"golang\"", "com", "en")

	if resultChan == nil {
		t.Errorf("Search(\"golang\") didn't go well")
	}

	res1 := <-resultChan

	if res1.Link != "https://golang.org/" {
		t.Errorf("Search(\"golang\") first result should be \"https://golang.org/\", got %s", res1.Link)
	}


}
