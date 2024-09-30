package fastapi

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

func (router *Router) AddCall(path string, handler interface{}) {
	validateHandler(handler)
	router.routes[path] = handler
}

func validateHandler(handler interface{}) {
	handlerType := reflect.TypeOf(handler)
	ginCtxType := reflect.TypeOf(&gin.Context{})
	errorInterface := reflect.TypeOf((*error)(nil)).Elem()

	switch {
	case handlerType.NumIn() != 2:
		panic("Wrong number of arguments")
	case handlerType.NumOut() != 2:
		panic("Wrong number of return values")
	case !handlerType.In(0).ConvertibleTo(ginCtxType):
		panic("First argument should be *gin.Context!")
	case handlerType.In(1).Kind() != reflect.Struct:
		panic("Second argument must be a struct")
	case !handlerType.Out(1).Implements(errorInterface):
		panic("Second return value should be an error")
	case handlerType.Out(0).Kind() != reflect.Struct:
		panic("First return value must be a struct")
	}
}
