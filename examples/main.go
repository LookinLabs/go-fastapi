package examples

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-fastapi"
)

type EchoInput struct {
	Phrase string `json:"phrase"`
}

type EchoOutput struct {
	OriginalInput EchoInput `json:"original_input"`
}

func EchoHandler(ctx *gin.Context, in EchoInput) (out EchoOutput, err error) {
	out.OriginalInput = in
	return
}

func main() {
	r := gin.Default()

	// Use Router for handling API requests and generating OpenAPI definition
	myRouter := fastapi.NewRouter()
	myRouter.AddCall("/echo", EchoHandler)

	r.POST("/api/*path", myRouter.GinHandler) // must have *path parameter

	swagger := myRouter.EmitOpenAPIDefinition()
	swagger.Info.Title = "My awesome API"
	jsonBytes, _ := json.MarshalIndent(swagger, "", "    ")
	fmt.Println(string(jsonBytes))

	// Serve Swagger JSON
	r.GET("/swagger.json", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", jsonBytes)
	})

	// Serve Swagger UI
	r.Static("/swagger-ui", "./swaggerui")

	r.Run()
}

// Try it:
//     $ curl -H "Content-Type: application/json" -X POST --data '{"phrase": "hello"}' localhost:8080/api/echo
//     {"response":{"original_input":{"phrase":"hello"}}}
