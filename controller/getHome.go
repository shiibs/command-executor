package controller

import (
	"embed"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Controller struct to hold the embedded file system
type Controller struct {
	content embed.FS
}

// Create a new controller with the embedded file system
func NewController(content embed.FS) *Controller {
	return &Controller{content: content}
}

// Handler to serve the HTML file
func (ctrl *Controller) GetHome(c *gin.Context) {
	file, err := ctrl.content.ReadFile("static/index.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error reading file: %v", err)
		return
	}
	c.Data(http.StatusOK, "text/html", file)
}
