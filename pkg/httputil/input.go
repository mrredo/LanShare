package httputil

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func GetInput(c *gin.Context, key string) (string, bool) {
	if c == nil {
		return "", false
	}

	if value := c.Query(key); strings.TrimSpace(value) != "" {
		return value, true
	}

	if value := c.PostForm(key); strings.TrimSpace(value) != "" {
		return value, true
	}

	if c.Request != nil && c.Request.Body != nil {
		var body map[string]any
		if err := c.ShouldBindBodyWith(&body, binding.JSON); err == nil && body != nil {
			if rawVal, exists := body[key]; exists && rawVal != nil {
				switch v := rawVal.(type) {
				case string:
					return v, true
				default:
					return fmt.Sprint(v), true
				}
			}
		}
	}

	return "", false
}

func GetInputDefault(c *gin.Context, key string, defaultVal string) string {
	if val, ok := GetInput(c, key); ok {
		return val
	}
	return defaultVal
}
