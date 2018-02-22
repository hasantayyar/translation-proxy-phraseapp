package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Run translation proxy API
func Run() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Translation Proxy is running")
	})

	api := router.Group("/api/v2")
	{
		api.GET("/projects/:project_id/locales/:id/download", func(c *gin.Context) {
			name := c.Param("name")
			c.String(http.StatusOK, "Hello %s", name)
		})
	}

	router.Run(":8080")
}
