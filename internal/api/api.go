package api

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Router struct {
	port   int
	router *gin.Engine
}

func NewRouter(port int) *Router {
	return &Router{
		port:   port,
		router: gin.Default(),
	}
}

func (r *Router) InitializeMiddlewares() {
	r.router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,

		MaxAge: 24 * time.Hour,
	}))
	//apiGroup := NewGroup(r.router, "/api")
	//apiGroup.Group().GET("/hello", func(c *gin.Context) {
	//	c.String(200, "Hello")
	//})
	//
	//r.router.GET("/", func(c *gin.Context) {
	//	if err := setOwnerCookie(c); err != nil {
	//		c.String(500, "Failed to set owner cookie")
	//		return
	//	}
	//
	//	c.String(200, "Hello")
	//})
}
func (r *Router) Router() *gin.Engine {
	return r.router
}
func (r *Router) Serve() {
	if err := r.router.Run(fmt.Sprintf(":%d", r.port)); err != nil {
		panic(err)
	}
}
