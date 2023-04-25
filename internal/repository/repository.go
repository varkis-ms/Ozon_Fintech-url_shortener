package repository

import "context"

// Repository Интерфейс, реализующий методы репозиториев
type Repository interface {
	SaveUrl(ctx context.Context, baseUrl, shortUrl string) error
	GetBaseUrlByShort(ctx context.Context, shortUrl string) (string, error)
	GetShortUrlByBase(ctx context.Context, baseUrl string) (string, bool, error)
}
