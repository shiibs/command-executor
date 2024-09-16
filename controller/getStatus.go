package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shiibs/fourCore-project/model"
)

func GetStatus(c *gin.Context) {
	id := c.Param("id")

	mutex.Lock()
	status, exists := commandStatuses[id]
	mutex.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, model.CommandResponse{Error: "Command not founc"})
		return
	}

	if status.Done {
		c.JSON(http.StatusOK, model.CommandResponse{Output: status.Output, Error: status.Error})
	} else {
		c.JSON(http.StatusOK, model.CommandResponse{Error: "Command is still running"})
	}
}
