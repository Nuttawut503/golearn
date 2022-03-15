package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", func(ctx *gin.Context) {
		delay := rand.Intn(200) + 600
		time.Sleep(time.Millisecond * time.Duration(delay))
		ctx.String(http.StatusOK, "response was sent")
	})
	r.Run()
}
