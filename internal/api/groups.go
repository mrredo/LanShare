package api

import "github.com/gin-gonic/gin"

type SubGroup interface {
	Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup
}
type Group struct {
	relativePath string
	subGroup     SubGroup
	group        *gin.RouterGroup
}

func NewGroup(subGroup SubGroup, relativePath string) *Group {
	return &Group{
		relativePath: relativePath,
		subGroup:     subGroup,
		group:        subGroup.Group(relativePath),
	}
}
func (g *Group) Group() *gin.RouterGroup {
	return g.group
}
