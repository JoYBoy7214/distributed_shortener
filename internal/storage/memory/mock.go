package memory

import (
	"context"
	"errors"
)

type MockDb struct {
	Db  map[string]string
	CDb map[string]int
}

func (p *MockDb) CreateDb(ctx context.Context) error {
	p.Db = make(map[string]string)
	p.CDb = make(map[string]int)
	return nil
}
func (p *MockDb) SaveURL(ctx context.Context, shortCode string, originalUrl string) error {
	_, ok := p.Db[shortCode]
	if ok {
		return errors.New("SQLSTATE 23505")
	} else {
		p.Db[shortCode] = originalUrl
	}
	return nil
}

func (p *MockDb) GetURL(ctx context.Context, shortCode string) (string, error) {
	var originalUrl string
	originalUrl, ok := p.Db[shortCode]
	if ok {
		return originalUrl, nil
	}
	return "", errors.New("Url not found")

}

func (p *MockDb) GetClickCount(ctx context.Context, shortcode string) (int, error) {
	var count int
	count, ok := p.CDb[shortcode]
	if ok {
		return count, nil
	}
	return -1, errors.New("Url not found")

}

func (p *MockDb) IncrementClick(ctx context.Context, shortCode string) error {
	p.CDb[shortCode]++
	return nil
}

func (p *MockDb) Close() {
	clear(p.CDb)
	clear(p.Db)
}
