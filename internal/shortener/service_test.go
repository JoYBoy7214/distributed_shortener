package shortener

import (
	"context"
	"fmt"
	"log"
	"testing"

	pb "github.com/JoYBoy7214/distributed_shortener/api/proto/v1"
	"github.com/JoYBoy7214/distributed_shortener/internal/storage/memory"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func TestCreateShortUrl(t *testing.T) {

	mockStore := &memory.MockDb{
		Db:  make(map[string]string),
		CDb: make(map[string]int),
	}

	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	svc := NewService(mockStore, rdb)

	req := &pb.CreateShortUrlRequest{
		OriginalUrl: "https://google.com",
		UserId:      "user_1",
	}
	resp, err := svc.CreateShortUrl(context.Background(), req)

	if err != nil {
		log.Fatal("test failed")
	}
	if len(resp.ShortUrl) > 0 {
		fmt.Println(resp.ShortUrl)
	} else {
		log.Fatal("test failed")
	}
	_, ok := mockStore.Db[resp.ShortUrl]
	if !ok {
		log.Fatal("test failed")
	}
	fmt.Println("Test Passed")

}
