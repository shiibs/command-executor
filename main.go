package main

import (
	"embed"
	"log"
	"net/http"
	"os/exec"
	"sync"
	"time"

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

type CommandStatus struct {
	Output string
	Error  string
	Done   bool
}

var (
	commandStatuses = make(map[string]*CommandStatus)
	mutex           sync.Mutex
)

func executeCommand(command string, id string) {
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()

	mutex.Lock()
	defer mutex.Unlock()

	status, exists := commandStatuses[id]
	if !exists {
		return
	}

	status.Output = string(output)
	if err != nil {
		status.Error = err.Error()
	}

	status.Done = true
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

		id := time.Now().Format("20060102150405")
		status := &CommandStatus{}
		mutex.Lock()
		commandStatuses[id] = status
		mutex.Unlock()

		go executeCommand(req.Command, id)
		c.JSON(http.StatusOK, gin.H{"id": id})
	})

	router.GET("api/status/:id", func(c *gin.Context) {
		id := c.Param("id")

		mutex.Lock()
		status, exists := commandStatuses[id]
		mutex.Unlock()

		if !exists {
			c.JSON(http.StatusNotFound, CommandResponse{Error: "Command not founc"})
			return
		}

		if status.Done {
			c.JSON(http.StatusOK, CommandResponse{Output: status.Output, Error: status.Error})
		} else {
			c.JSON(http.StatusOK, CommandResponse{Error: "Command is still running"})
		}
	})

	if err := router.Run(":31337"); err != nil {
		log.Fatal(err)
	}
}
