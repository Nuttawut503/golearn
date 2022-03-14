package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	buildcommit = ""
	buildtime   = ""
)

func main() {
	greeting, name := "Hello", "Albert"
	if v, ok := os.LookupEnv("greeting"); ok {
		greeting = v
	}
	if v, ok := os.LookupEnv("name"); ok {
		name = v
	}
	r := gin.Default()
	r.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, fmt.Sprintf("%s, %s!", greeting, name))
	})
	r.GET("/x", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"buildcommit": buildcommit,
			"buildtime":   buildtime,
		})
	})
	r.Run()
}
