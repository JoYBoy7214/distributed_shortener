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

// Service implements the Shortener gRPC service.
// We export it (capital 'S') or keep it private and return the interface.
// For simplicity here, let's export the struct or use a constructor.
type Service struct {
	pb.UnimplementedShortenerServer // Embed this for forward compatibility
	m                               map[string]string
	mtx                             sync.RWMutex // RWMutex allows multiple readers at once
}

// NewService initializes the map and returns a ready-to-use service.
func NewService() *Service {
	return &Service{
		m: make(map[string]string),
	}
}

func (s *Service) CreateShortUrl(ctx context.Context, req *pb.CreateShortUrlRequest) (*pb.CreateShortUrlResponse, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// Seed the random generator (required for older Go versions, good practice anyway)
	// In production, we would use crypto/rand for secure IDs.
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
	s.mtx.RLock() // Read Lock: multiple threads can read at the same time
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
