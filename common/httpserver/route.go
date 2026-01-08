package httpserver

import "github.com/gin-gonic/gin"

// Route defines a handler and its OpenAPI metadata.
type Route struct {
	RouteSpec
	Handler gin.HandlerFunc
}
