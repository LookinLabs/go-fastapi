package fastapi

import (
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
)

type ValidInput struct {
	Field string `json:"field"`
}

type ValidOutput struct {
	Field string `json:"field"`
}

func validHandler(ctx *gin.Context, input ValidInput) (ValidOutput, error) {
	return ValidOutput{Field: input.Field}, nil
}

func invalidHandlerNoGinContext(input ValidInput) (ValidOutput, error) {
	return ValidOutput{Field: input.Field}, nil
}

func invalidHandlerWrongReturnType(ctx *gin.Context, input ValidInput) (string, error) {
	return input.Field, nil
}

func invalidHandlerNoErrorReturn(ctx *gin.Context, input ValidInput) (ValidOutput, string) {
	return ValidOutput{Field: input.Field}, ""
}

func invalidHandlerWrongInputType(ctx *gin.Context, input string) (ValidOutput, error) {
	return ValidOutput{Field: input}, nil
}

func TestAddCallConcurrently(t *testing.T) {
	router := &Router{routes: make(map[string]interface{})}

	tests := []struct {
		name        string
		handler     interface{}
		shouldPanic bool
	}{
		{"ValidHandler", validHandler, false},
		{"InvalidHandlerNoGinContext", invalidHandlerNoGinContext, true},
		{"InvalidHandlerWrongReturnType", invalidHandlerWrongReturnType, true},
		{"InvalidHandlerNoErrorReturn", invalidHandlerNoErrorReturn, true},
		{"InvalidHandlerWrongInputType", invalidHandlerWrongInputType, true},
	}

	var wg sync.WaitGroup

	for _, tt := range tests {
		wg.Add(1)
		go func(tt struct {
			name        string
			handler     interface{}
			shouldPanic bool
		}) {
			defer wg.Done()
			t.Run(tt.name, func(t *testing.T) {
				defer func() {
					if r := recover(); (r != nil) != tt.shouldPanic {
						t.Errorf("AddCall() panic = %v, wantPanic = %v", r != nil, tt.shouldPanic)
					}
				}()
				router.AddCall("/test", tt.handler)
			})
		}(tt)
	}

	wg.Wait()
}
