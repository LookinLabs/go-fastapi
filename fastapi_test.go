package fastapi

import (
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
)

func runTest(t *testing.T, wg *sync.WaitGroup, testFunc func()) {
	defer wg.Done()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Did not panic")
		}
	}()
	testFunc()
}

func TestHandlersConcurrently(t *testing.T) {
	var wg sync.WaitGroup

	tests := []func(){
		func() { NewRouter().AddCall("x", func(_ struct{}) struct{} { return struct{}{} }) },
		func() { NewRouter().AddCall("x", func(_ *gin.Context, _ struct{}) struct{} { return struct{}{} }) },
		func() { NewRouter().AddCall("x", func(_ struct{}) (struct{}, error) { return struct{}{}, nil }) },
		func() {
			NewRouter().AddCall("x", func(_ struct{}, _ struct{}) (struct{}, error) { return struct{}{}, nil })
		},
		func() {
			NewRouter().AddCall("x", func(_ *gin.Context, _ string) (struct{}, error) { return struct{}{}, nil })
		},
		func() { NewRouter().AddCall("x", func(_ *gin.Context, _ struct{}) (string, error) { return "", nil }) },
	}

	for _, test := range tests {
		wg.Add(1)
		go runTest(t, &wg, test)
	}

	wg.Wait()
}

func TestCorrectHandler(_ *testing.T) {
	NewRouter().AddCall("x", func(_ *gin.Context, _ struct{}) (struct{}, error) { return struct{}{}, nil })
}
