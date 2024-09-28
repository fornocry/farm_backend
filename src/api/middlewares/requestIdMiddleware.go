package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (u MiddlewareServiceImpl) RequestIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := uuid.New()
		c.Writer.Header().Set("X-Request-Id", id.String())
		c.Next()
	}
}
