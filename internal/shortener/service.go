package shortener

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	pb "github.com/JoYBoy7214/distributed_shortener/api/proto/v1"
	repo "github.com/JoYBoy7214/distributed_shortener/internal/storage"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	pb.UnimplementedShortenerServer
	repo        repo.Repository
	rdb         *redis.Client
	messageChan chan string
	Wtg         sync.WaitGroup
}

const bufferSize int = 500

func NewService(pool repo.Repository, rdb *redis.Client) *Service {
	c := make(chan string, bufferSize)
	return &Service{
		repo:        pool,
		rdb:         rdb,
		messageChan: c,
	}
}
func (s *Service) CreateUrlTable() error {

	err := s.repo.CreateDb(context.Background())
	if err != nil {
		return status.Errorf(codes.Internal, "failed to create the url table")
	}
	fmt.Println("Data base created")
	return nil

}
func (s *Service) CreateShortUrl(ctx context.Context, req *pb.CreateShortUrlRequest) (*pb.CreateShortUrlResponse, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	//fmt.Println("end point is working")
	var shortID string
	for i := 0; i < 5; i++ {
		b := make([]byte, 10)
		for j := range b {
			b[j] = charset[rng.Intn(len(charset))]
		}
		shortID = string(b)
		err := s.repo.SaveURL(ctx, shortID, req.OriginalUrl)
		if err != nil {
			if strings.Contains(err.Error(), "SQLSTATE 23505") && i != 4 {
				continue
			} else {
				return nil, status.Errorf(codes.Internal, "failed to insert the Url %s", req.OriginalUrl)
			}
		} else {
			break
		}
	}
	go s.rdb.Set(ctx, shortID, req.OriginalUrl, 0) //we can ignore error ig

	return &pb.CreateShortUrlResponse{
		ShortUrl: shortID,
	}, nil
}

func (s *Service) GetOriginalUrl(ctx context.Context, req *pb.GetOriginalUrlRequest) (*pb.GetOriginalUrlResponse, error) {
	select {
	case s.messageChan <- req.ShortUrl:
	default:
		//fmt.Println("buffer fulled")
	}
	var originalUrl string
	originalUrl, err := s.rdb.Get(ctx, req.ShortUrl).Result()
	if err == nil {
		return &pb.GetOriginalUrlResponse{
			OriginalUrl: originalUrl,
		}, nil
	}
	if err != redis.Nil {
		log.Print("redis down ", err)
	}
	originalUrl, err = s.repo.GetURL(ctx, req.ShortUrl)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "data not found for the url : %s", req.ShortUrl)
	}
	go s.rdb.Set(ctx, req.ShortUrl, originalUrl, 0)
	return &pb.GetOriginalUrlResponse{
		OriginalUrl: originalUrl,
	}, nil
}

func (s *Service) GetClickCount(ctx context.Context, req *pb.GetOriginalUrlRequest) (*pb.GetClickCountResponse, error) {
	count, err := s.repo.GetClickCount(ctx, req.ShortUrl)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "data not found for the url ;%s", req.ShortUrl)
	}
	return &pb.GetClickCountResponse{
		ClickCount: int64(count),
	}, nil
}

func (s *Service) StartWorker() {
	defer s.Wtg.Done()
	fmt.Println("worker started")

	for id := range s.messageChan {
		d := time.Now().Add(500 * time.Millisecond)
		ctx, cancel := context.WithDeadline(context.Background(), d)
		err := s.repo.IncrementClick(ctx, id)
		if err != nil {
			log.Println("error in updating the click count", err)
		}
		cancel()
	}

}

func (s *Service) GraceFullShutdown() {
	fmt.Println("1. Closing channel...")
	close(s.messageChan)

	fmt.Println("2. Waiting for worker...")
	s.Wtg.Wait()

	fmt.Println("3. Closing DB pool...")
	s.repo.Close()

	fmt.Println("4. Closing Redis...")
	s.rdb.Close()
}
