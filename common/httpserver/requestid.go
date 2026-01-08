package httpserver

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIDHeader = "X-Request-ID"

// RequestIDMiddleware ensures every request has a request ID header.
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader(requestIDHeader)
		if reqID == "" {
			reqID = uuid.NewString()
			c.Request.Header.Set(requestIDHeader, reqID)
		}
		c.Header(requestIDHeader, reqID)
		c.Set(requestIDHeader, reqID)
		c.Next()
	}
}
