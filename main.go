package main

import (
	"embed"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/shiibs/fourCore-project/controller"
)

//go:embed static/*
var content embed.FS

func main() {
	router := gin.Default()

	// Create a new controller with the embedded content
	ctrl := controller.NewController(content)

	// Serve the HTML file
	router.GET("/", ctrl.GetHome)

	// POST API to execute commands
	router.POST("/api/execute", controller.PostExecute)

	// Get execute status
	router.GET("api/status/:id", controller.GetStatus)

	if err := router.Run(":31337"); err != nil {
		log.Fatal(err)
	}
}
