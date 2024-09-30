package fastapi

import (
	"log"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

func (router *Router) GinHandler(ctx *gin.Context) {
	path := ctx.Param("path")
	log.Print(path)

	handlerFuncPtr, present := router.routes[path]
	if !present {
		respondWithError(ctx, http.StatusNotFound, "handler not found")
		return
	}

	inputVal, err := bindInput(ctx, handlerFuncPtr)
	if err != nil {
		respondWithError(ctx, http.StatusBadRequest, "invalid request")
		return
	}

	outputVal, err := callHandler(ctx, handlerFuncPtr, inputVal)
	if err != nil {
		respondWithError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"response": outputVal.Interface(),
	})
}

func (router *Router) GetRoutes() map[string]interface{} {
	return router.routes
}

func respondWithError(ctx *gin.Context, code int, message interface{}) {
	ctx.JSON(code, gin.H{
		"error": message,
		"code":  code,
	})
}

func bindInput(ctx *gin.Context, handlerFuncPtr interface{}) (interface{}, error) {
	handlerType := reflect.TypeOf(handlerFuncPtr)
	inputType := handlerType.In(1)
	inputVal := reflect.New(inputType).Interface()
	if err := ctx.BindJSON(inputVal); err != nil {
		return nil, err
	}
	return inputVal, nil
}

func callHandler(c *gin.Context, handlerFuncPtr interface{}, inputVal interface{}) (reflect.Value, error) {
	toCall := reflect.ValueOf(handlerFuncPtr)
	outputVal := toCall.Call(
		[]reflect.Value{
			reflect.ValueOf(c),
			reflect.ValueOf(inputVal).Elem(),
		},
	)

	returnedErr := outputVal[1].Interface()
	if returnedErr != nil || !outputVal[1].IsNil() {
		return reflect.Value{}, returnedErr.(error)
	}

	return outputVal[0], nil
}
