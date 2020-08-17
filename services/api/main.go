package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.MaxMultipartMemory = 8 << 20

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "pong3")
	})
	router.GET("/api/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong3")
	})
	router.Run(":80")
}
