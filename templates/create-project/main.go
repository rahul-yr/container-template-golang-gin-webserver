package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func readynessCheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}

func health(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}

func allStatus(c *gin.Context) {
	hostname, _ := os.Hostname()

	c.JSON(200, gin.H{
		"status":   "ok",
		"hostname": hostname,
		"host":     c.Request.Host,
		"ip":       c.ClientIP(),
		"url":      c.Request.URL,
	})
}

func customRoutes(r *gin.Engine) {
	// Add your custom routes below
	// TODO add your custom routes here

}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(cors.Default())

	// Default routes
	router.GET("/container/readyness", readynessCheck)
	router.GET("/container/health", health)
	router.GET("/container/all", allStatus)

	// Add Custom routes below
	customRoutes(router)

	// Start the server
	port := "8080"
	log.Println("Server started on port " + port)
	router.Run(":" + port)
}
