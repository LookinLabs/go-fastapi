package fastapi

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestValidateHandlerInvalidCases(t *testing.T) {
	tests := []struct {
		name    string
		handler interface{}
	}{
		{
			name:    "Wrong number of arguments",
			handler: func(_ struct{}) struct{} { return struct{}{} },
		},
		{
			name:    "Wrong number of return values",
			handler: func(_ *gin.Context, _ struct{}) {},
		},
		{
			name:    "First argument not *gin.Context",
			handler: func(_ struct{}, _ *gin.Context) (struct{}, error) { return struct{}{}, nil },
		},
		{
			name:    "Second argument not a struct",
			handler: func(_ *gin.Context, _ string) (struct{}, error) { return struct{}{}, nil },
		},
		{
			name:    "Second return value not an error",
			handler: func(_ *gin.Context, _ struct{}) (struct{}, string) { return struct{}{}, "" },
		},
		{
			name:    "First return value not a struct",
			handler: func(_ *gin.Context, _ struct{}) (string, error) { return "", nil },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("validateHandler did not panic for case: %s", tt.name)
				}
			}()
			validateHandler(tt.handler)
		})
	}
}

func TestValidateHandlerValidCase(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("validateHandler panicked for valid handler")
		}
	}()
	validateHandler(func(_ *gin.Context, _ struct{}) (struct{}, error) { return struct{}{}, nil })
}
