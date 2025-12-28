package main

import (
	"log"
	"net"

	pb "github.com/JoYBoy7214/distributed_shortener/api/proto/v1"
	service "github.com/JoYBoy7214/distributed_shortener/internal/shortener"
	"google.golang.org/grpc"
)

func main() {
	ln, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("failed to listen on 50051", err)
	}
	s := grpc.NewServer()
	server := service.NewService()
	pb.RegisterShortenerServer(s, server)
	log.Println("server started on 50051")
	if err := s.Serve(ln); err != nil {
		log.Fatal("failed to start grpc server")
	}
}
