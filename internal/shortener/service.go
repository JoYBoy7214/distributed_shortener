package shortener

import (
	"context"
	"math/rand"
	"sync"
	"time"

	pb "github.com/JoYBoy7214/distributed_shortener/api/proto/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	pb.UnimplementedShortenerServer
	m   map[string]string
	mtx sync.RWMutex
}

func NewService() *Service {
	return &Service{
		m: make(map[string]string),
	}
}

func (s *Service) CreateShortUrl(ctx context.Context, req *pb.CreateShortUrlRequest) (*pb.CreateShortUrlResponse, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, 10)
	for i := range b {
		b[i] = charset[rng.Intn(len(charset))]
	}
	shortID := string(b)

	s.mtx.Lock()
	s.m[shortID] = req.OriginalUrl
	s.mtx.Unlock()

	return &pb.CreateShortUrlResponse{
		ShortUrl: shortID,
	}, nil
}

func (s *Service) GetOriginalUrl(ctx context.Context, req *pb.GetOriginalUrlRequest) (*pb.GetOriginalUrlResponse, error) {
	s.mtx.RLock()
	originalURL, exists := s.m[req.ShortUrl]
	s.mtx.RUnlock()

	if !exists {
		// Return a proper gRPC error code
		return nil, status.Errorf(codes.NotFound, "short url %s not found", req.ShortUrl)
	}

	return &pb.GetOriginalUrlResponse{
		OriginalUrl: originalURL,
	}, nil
}
