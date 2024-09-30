package fastapi

type Router struct {
	routes map[string]interface{}
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]interface{}),
	}
}
