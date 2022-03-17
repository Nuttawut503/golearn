package server

import (
	"gogrpc/server/temppb"
	"log"
	"math/rand"
	"time"
)

type Server struct {
	temppb.UnimplementedTempServer
}

func (s *Server) TempStream(in *temppb.TempRequest, stream temppb.Temp_TempStreamServer) error {
	log.Println("One joined")
	for {
		if err := stream.Send(&temppb.TempResponse{
			Temperature: int32(rand.Intn(10)) + 25,
		}); err != nil {
			log.Printf("Can't send a message: %v\n", err)
			log.Println("One closed")
			return nil
		}
		time.Sleep(time.Second)
	}
}
