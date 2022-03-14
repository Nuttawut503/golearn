package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var _livenessfile = "/tmp/myservice-healthy"

func main() {
	r := gin.Default()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	r.GET("/health", func(ctx *gin.Context) {
		// check if all clients' connection is still able to work
		if err := rdb.Ping(ctx).Err(); err != nil {
			ctx.String(http.StatusInternalServerError, "Redis client isn't working")
			return
		}
		ctx.String(http.StatusOK, "Clients are available")
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// before running create a file to check if service is running
	os.Create(_livenessfile)
	// remove after the service has stopped
	defer os.Remove(_livenessfile)

	// to do liveness check, it requires to implement graceful-shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
	log.Println("Server exiting")
}
