package cache

import (
	"net/http"
	"testing"
)

func TestSize(t *testing.T) {
	type testType struct {
		status  int
		header  http.Header
		str     string
		data    []byte
		numbers []int
	}

	testObj := testType{
		header: http.Header{
			"Use-Agent": []string{"abc", "de", "d"},
		},
		str:     "a",
		data:    []byte{'a'},
		numbers: []int{2, 3},
	}
	s := size(testObj)
	t.Log("size:", s)
}
