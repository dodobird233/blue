package main

import (
	"blue/initialize"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	err := initialize.InitDB()
	if err != nil {
		os.Exit(-1)
	}
	r := gin.Default()

	initRouter(r)

	r.Run("192.168.235.1:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
