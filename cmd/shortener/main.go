package main

import (
	"context"
	"log"
	"net"

	pb "github.com/JoYBoy7214/distributed_shortener/api/proto/v1"
	service "github.com/JoYBoy7214/distributed_shortener/internal/shortener"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

const dbUrl = "postgres://ZORO:secretfornow@localhost:5432/shortener"

func main() {
	ln, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("failed to listen on 50051", err)
	}
	s := grpc.NewServer()
	pool, err := pgxpool.New(context.Background(), dbUrl)
	server := service.NewService(pool)
	server.CreateUrlTable()
	pb.RegisterShortenerServer(s, server)
	log.Println("server started on 50051")
	if err := s.Serve(ln); err != nil {
		log.Fatal("failed to start grpc server")
	}
}
