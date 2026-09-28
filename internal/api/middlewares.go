package api

type Middleware struct {
}

func NewMiddleware() *Middleware {
	return &Middleware{}
}
func (m *Middleware) Use(router *Router) {

}

func (r *Router) RegisterMiddlewares(middlewares ...*Endpoint) {

}
