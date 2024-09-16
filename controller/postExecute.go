package controller

import (
	"net/http"
	"os/exec"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shiibs/fourCore-project/model"
)

var (
	commandStatuses = make(map[string]*model.CommandStatus)
	mutex           sync.Mutex
)

func PostExecute(c *gin.Context) {
	var req model.CommandRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.CommandResponse{Error: err.Error()})
		return
	}

	id := time.Now().Format("20060102150405")
	status := &model.CommandStatus{}
	mutex.Lock()
	commandStatuses[id] = status
	mutex.Unlock()

	go executeCommand(req.Command, id)
	c.JSON(http.StatusOK, gin.H{"id": id})
}

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
