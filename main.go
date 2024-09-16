package main

import (
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/shiibs/fourCore-project/model"
)

func main() {
	router := gin.Default()

	router.StaticFile("/", "./static/index.html")

	router.POST("/api/execute", func(c *gin.Context) {
		var req model.CommandRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cmd := exec.Command("cmd.exe", "/C", req.Command)

		output, err := cmd.CombinedOutput()
		if err != nil {
			c.JSON(http.StatusOK, model.CommandReponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, model.CommandReponse{Output: string(output)})
	})

	router.Run(":31337")
}
