package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
)

func main() {
	provider, err := NewTraceProvider("http://localhost:14268/api/traces")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if err := provider.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}()
	otel.SetTracerProvider(provider)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	r := gin.Default()
	SetupHandler(r)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	log.Println("Server is running, press ctrl+C to stop")
	<-ctx.Done()
	stop()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := provider.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server is forced to shutdown: ", err)
	}
	log.Println("Server exiting")
}
