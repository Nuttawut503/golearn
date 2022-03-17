package main

import (
	"context"
	"gogrpc/server"
	"gogrpc/server/temppb"
	"log"
	"net"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalln(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := grpc.NewServer()
	temppb.RegisterTempServer(srv, &server.Server{})
	go func() {
		log.Println("Server running")
		if err := srv.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			log.Fatalf("Error while serving: %v", err)
		}
	}()

	<-ctx.Done()
	stop()
	srv.Stop()
	log.Println("Server exiting")
}
