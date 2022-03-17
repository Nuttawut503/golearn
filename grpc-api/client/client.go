package main

import (
	"context"
	"gogrpc/server/temppb"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	c := temppb.NewTempClient(conn)
	p, _ := c.TempStream(context.Background(), &temppb.TempRequest{})
	for {
		v, err := p.Recv()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("receive: ", v.GetTemperature())
	}
}
