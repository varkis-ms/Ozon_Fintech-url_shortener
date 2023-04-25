package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"rest_url_shortener/internal/repository"
)

// postgresRepository Структура Postgres базы данных (содержит подключение к базе данных)
type postgresRepository struct {
	db *pgxpool.Pool
}

// NewRepository Конструктор для создания Postgres базы данных (postgresRepository)
func NewRepository(db *pgxpool.Pool) repository.Repository {
	return &postgresRepository{db: db}
}

// SaveUrl Метод репозитория (postgresRepository), который сохраняет данных в базу данных Postgres
func (r *postgresRepository) SaveUrl(ctx context.Context, baseUrl string, shortUrl string) error {
	_, err := r.db.Exec(ctx,
		"INSERT INTO shortened_urls (url_base, url_short) VALUES ($1, $2);",
		baseUrl, shortUrl)
	if err == nil {
		return err
	}
	return nil
}

/*
GetBaseUrlByShort Метод репозитория (postgresRepository), который по короткой ссылке
находит изначальную ссылку в хранилище
Если такой ссылки нет, то возвращает ошибку unknown url
*/
func (r *postgresRepository) GetBaseUrlByShort(ctx context.Context, shortUrl string) (string, error) {
	var baseUrl string
	err := r.db.QueryRow(ctx,
		"SELECT url_base FROM shortened_urls WHERE url_short = $1;",
		shortUrl).Scan(&baseUrl)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", errors.New("unknown url")
		}
		return "", err
	}
	return baseUrl, nil
}

/*
GetShortUrlByBase Метод репозитория (postgresRepository), который по изначальной ссылке
находит короткую ссылку в хранилище
Так же этот метод служит проверкой на уникальность переданной ссылки в базу данных Postgres
*/
func (r *postgresRepository) GetShortUrlByBase(ctx context.Context, baseUrl string) (string, bool, error) {
	var shortUrl string
	err := r.db.QueryRow(ctx,
		"SELECT url_short FROM shortened_urls WHERE url_base = $1;",
		baseUrl).Scan(&shortUrl)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", true, errors.New("unknown url")
		}
		return "", false, err
	}
	return shortUrl, false, nil
}
