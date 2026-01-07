package shortener

import (
	"context"
	"math/rand"
	"strings"
	"time"

	pb "github.com/JoYBoy7214/distributed_shortener/api/proto/v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	pb.UnimplementedShortenerServer
	pool *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{
		pool: pool,
	}
}
func (s *Service) CreateUrlTable() error {
	_, err := s.pool.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS urls(id SERIAL UNIQUE,short_code VARCHAR(10) PRIMARY KEY UNIQUE NOT NULL,original_url TEXT NOT NULL)")
	if err != nil {
		return status.Errorf(codes.Internal, "failed to create the url table")
	}
	return nil

}
func (s *Service) CreateShortUrl(ctx context.Context, req *pb.CreateShortUrlRequest) (*pb.CreateShortUrlResponse, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var shortID string
	for i := 0; i < 5; i++ {
		b := make([]byte, 10)
		for j := range b {
			b[j] = charset[rng.Intn(len(charset))]
		}
		shortID = string(b)
		_, err := s.pool.Exec(ctx, "INSERT INTO urls (short_code, original_url) VALUES ($1, $2)", shortID, req.OriginalUrl)
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
	return &pb.CreateShortUrlResponse{
		ShortUrl: shortID,
	}, nil
}

func (s *Service) GetOriginalUrl(ctx context.Context, req *pb.GetOriginalUrlRequest) (*pb.GetOriginalUrlResponse, error) {
	Row := s.pool.QueryRow(ctx, "SELECT original_url FROM urls WHERE short_code = $1", req.ShortUrl)
	var originalUrl string

	err := Row.Scan(&originalUrl)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "data not found for the url : %s", req.ShortUrl)
	}

	return &pb.GetOriginalUrlResponse{
		OriginalUrl: originalUrl,
	}, nil
}
