package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

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
	dbUrl = os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		dbUrl = "postgres://ZORO:secretfornow@localhost:5432/shortener"
	}
	fmt.Println(dbUrl)

	s := grpc.NewServer()
	pool, err := pgxpool.New(context.Background(), dbUrl)
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	server := service.NewService(pool, rdb)
	server.CreateUrlTable()
	server.Wtg.Add(1)
	go server.StartWorker()
	pb.RegisterShortenerServer(s, server)
	go func(s *grpc.Server, ln net.Listener) {
		log.Println("server started on 50051")
		if err := s.Serve(ln); err != nil {
			log.Println("failed to start grpc server")
		}
	}(s, ln)
	select {
	case <-ctx.Done():
		s.GracefulStop()
		fmt.Println("Stopping gRPC server...")
		server.GraceFullShutdown()
		fmt.Println("Calling internal shutdown...")

		stop()
	}
}
