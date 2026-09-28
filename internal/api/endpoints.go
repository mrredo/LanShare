package api

type Endpoint struct {
}

func NewEndpoint() *Endpoint {
	return &Endpoint{}
}
func (e *Endpoint) Use(router *Router) {

}
func (r *Router) RegisterEndpoints(endpoints ...*Endpoint) {

}
