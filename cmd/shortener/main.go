package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/JoYBoy7214/distributed_shortener/api/proto/v1"
	service "github.com/JoYBoy7214/distributed_shortener/internal/shortener"
	db "github.com/JoYBoy7214/distributed_shortener/internal/storage/postgres"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	s := grpc.NewServer(
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	)
	pool, err := db.NewPostgresStore(dbUrl)
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
	grpc_prometheus.Register(s)
	grpc_prometheus.EnableHandlingTimeHistogram()
	go func(s *grpc.Server, ln net.Listener) {
		log.Println("server started on 50051")
		if err := s.Serve(ln); err != nil {
			log.Println("failed to start grpc server")
		}
	}(s, ln)
	// Start Metrics Server
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		fmt.Println("Metrics server started on :9090")
		http.ListenAndServe(":9090", nil)
	}()
	select {
	case <-ctx.Done():
		s.GracefulStop()
		fmt.Println("Stopping gRPC server...")
		server.GraceFullShutdown()
		fmt.Println("Calling internal shutdown...")

		stop()
	}
}
