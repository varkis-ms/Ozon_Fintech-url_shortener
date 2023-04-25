package cache

import (
	"context"
	"errors"
	"rest_url_shortener/internal/repository"
)

/*
InMemoryRepository Структура In-Memory хранилища
Базируется на hash-map для удобного и быстрого поиска
Используется 2 hash-map, для обработки всех возможных сценариев
*/
type inMemoryRepository struct {
	dbShortToBase map[string]string
	dbBaseToShort map[string]string
}

// NewInMemoryRepository Конструктор для создания In-Memory хранилища (InMemoryRepository)
func NewInMemoryRepository() repository.Repository {
	return &inMemoryRepository{
		dbShortToBase: make(map[string]string),
		dbBaseToShort: make(map[string]string),
	}
}

// SaveUrl Метод репозитория (InMemoryRepository), который сохраняет данных в хранилище
func (m *inMemoryRepository) SaveUrl(ctx context.Context, baseUrl, shortUrl string) error {
	m.dbShortToBase[shortUrl] = baseUrl
	m.dbBaseToShort[baseUrl] = shortUrl
	return nil
}

/*
GetBaseUrlByShort Метод репозитория (InMemoryRepository), который по короткой ссылке
находит изначальную ссылку в хранилище
Если такой ссылки нет, то возвращает ошибку unknown url
*/
func (m *inMemoryRepository) GetBaseUrlByShort(ctx context.Context, shortUrl string) (string, error) {
	if value, ok := m.dbShortToBase[shortUrl]; ok {
		return value, nil
	}
	return "", errors.New("unknown url")
}

/*
GetShortUrlByBase Метод репозитория (InMemoryRepository), который по изначальной ссылке
находит короткую ссылку в хранилище
Так же этот метод служит проверкой на уникальность переданной ссылки в хранилище
*/
func (m inMemoryRepository) GetShortUrlByBase(ctx context.Context, baseUrl string) (string, bool, error) {
	if value, ok := m.dbBaseToShort[baseUrl]; ok {
		return value, false, nil
	}
	return "", true, nil
}
