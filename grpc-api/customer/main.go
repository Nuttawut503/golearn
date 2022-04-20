package main

import (
	"context"
	"grpc-api/customer/server"
	"grpc-api/customer/server/customerpb"
	"log"
	"net"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalln(err)
	}
	srv := grpc.NewServer()
	customerpb.RegisterCustomerServer(srv, &server.CustomerServer{})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		stop()
		log.Println("Server is exiting")
		srv.Stop()
	}()

	log.Println("Server is running...")
	if err := srv.Serve(lis); err != nil && err != grpc.ErrServerStopped {
		log.Fatalf("Error while serving: %v", err)
	}
	log.Println("Server is closed")
}
