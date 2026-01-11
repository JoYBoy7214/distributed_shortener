package main

import (
	"context"
	"log"
	"net"
	"os"

	pb "github.com/JoYBoy7214/distributed_shortener/api/proto/v1"
	service "github.com/JoYBoy7214/distributed_shortener/internal/shortener"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

var dbUrl string

func main() {
	ln, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("failed to listen on 50051", err)
	}
	dbUrl = os.Getenv("DB_URL")
	s := grpc.NewServer()
	pool, err := pgxpool.New(context.Background(), dbUrl)
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	server := service.NewService(pool, rdb)
	server.CreateUrlTable()
	pb.RegisterShortenerServer(s, server)
	log.Println("server started on 50051")
	if err := s.Serve(ln); err != nil {
		log.Fatal("failed to start grpc server")
	}
}
