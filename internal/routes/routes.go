package routes

import (
	"net/http"

	"{{MODULE_PATH}}/common/httpserver"
	"{{MODULE_PATH}}/internal/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RouteDependencies contains all dependencies needed by route handlers.
type RouteDependencies struct {
	ItemService *services.ItemService
}

// BuildRoutes constructs all route definitions for the API.
func BuildRoutes(deps RouteDependencies, logger *zap.Logger) []httpserver.Route {
	routes := []httpserver.Route{
		{
			RouteSpec: httpserver.RouteSpec{
				Method:  http.MethodGet,
				GinPath: "/health",
				Summary: "Health check",
				Tags:    []string{"health"},
				Secure:  false,
			},
			Handler: func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "healthy"})
			},
		},
		{
			RouteSpec: httpserver.RouteSpec{
				Method:  http.MethodGet,
				GinPath: "/items",
				Summary: "List all items",
				Tags:    []string{"items"},
				Secure:  false,
			},
			Handler: func(c *gin.Context) {
				if deps.ItemService == nil {
					c.JSON(http.StatusOK, gin.H{"items": []interface{}{}})
					return
				}
				items, err := deps.ItemService.List(c.Request.Context())
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"items": items})
			},
		},
		{
			RouteSpec: httpserver.RouteSpec{
				Method:  http.MethodGet,
				GinPath: "/items/:id",
				Summary: "Get an item by ID",
				Tags:    []string{"items"},
				Secure:  false,
			},
			Handler: func(c *gin.Context) {
				id := c.Param("id")
				if deps.ItemService == nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
					return
				}
				item, err := deps.ItemService.Get(c.Request.Context(), id)
				if err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, item)
			},
		},
	}

	return routes
}
