package fastapi

import (
	"testing"

	"github.com/gin-gonic/gin"
)

type InnerStruct struct {
	XYZ string `json:"XYZ"`
}

type In struct {
	Input string          `json:"input"`
	X     int             `json:"x"`
	Y     float32         `json:"y"`
	Z     bool            `json:"z"`
	I     []string        `json:"i"`
	J     map[string]int8 `json:"j"`
}

type In2 struct {
	InputTwo string      `json:"input_two"`
	Inner    InnerStruct `json:"-"`
}

type Out struct {
	Output string `json:"output"`
}

func RequestHandler(_ *gin.Context, _ In) (Out, error) {
	return Out{}, nil
}

func RequestHandlerTwo(_ *gin.Context, _ In2) (Out, error) {
	return Out{}, nil
}

func TestOpenAPIDefinition(t *testing.T) {
	openAPIRouter := NewRouter()
	openAPIRouter.AddCall("/ping", RequestHandler)
	openAPIRouter.AddCall("/pong", RequestHandlerTwo)
	sw := openAPIRouter.EmitOpenAPIDefinition()

	if len(sw.Paths.Paths) != 2 {
		t.Fatal("Wrong number of paths")
	}

	if len(sw.Definitions) != 4 {
		t.Fatal("Wrong number of definitions")
	}

	_, innerPresent := sw.Definitions["InnerStruct"]
	if !innerPresent {
		t.Fatal("Nested structure is not present in definitions")
	}
}
