package url

import (
	"context"
	"errors"
	"fmt"
	"rest_url_shortener/internal/repository"
	"rest_url_shortener/internal/utils"
)

// Service Структура сервиса (содержит интерфейс репозитория)
type Service struct {
	repo repository.Repository
}

// NewService Конструктор для создания сервиса (Service )
func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

/*
ForwardBaseUrl Метод cервиса, который вызывает соответствующий метод выбранного репозитория
для нахождения короткой ссылки по изначальной ссылке (по сути реализация редиректа)
*/
func (r *Service) ForwardBaseUrl(ctx context.Context, shortUrl string) (string, error) {
	baseUrl, err := r.repo.GetBaseUrlByShort(ctx, shortUrl)
	if err != nil {
		return "", err
	}
	return baseUrl, nil
}

/*
AddShortUrl Метод cервиса, который вызывает соответствующий метод выбранного репозитория для
генерирации короткой ссылки, сохраненя данных в хранилище выбранного репозитория и возвращенея короткой ссылкы
Так же проверяет переданную и сгенерированную ссылы на наличие в хранилище выбранного репозитория.
*/
func (r *Service) AddShortUrl(ctx context.Context, baseUrl string) (string, error) {
	var shortenedUrl = utils.Encode63()
	shortUrl, exists, err := r.repo.GetShortUrlByBase(ctx, baseUrl)
	if exists {
		_, checkErr := r.repo.GetBaseUrlByShort(ctx, shortenedUrl)
		if checkErr != errors.New("unknown url") {
			checkErr = r.repo.SaveUrl(ctx, baseUrl, shortenedUrl)
			if checkErr != nil {
				return "", checkErr
			}
			return fmt.Sprintf("http://localhost:8080/%s", shortenedUrl), nil
		}
		if checkErr != nil {
			return "", checkErr
		}
		var ok bool
		for ok != true {
			shortenedUrl = utils.Encode63()
			_, checkErr = r.repo.GetBaseUrlByShort(ctx, shortenedUrl)
			if checkErr != errors.New("unknown url") {
				ok = true
			}
			if checkErr != nil {
				return "", errors.New("server error, please try again later")
			}
		}
	}
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("http://localhost:8080/%s", shortUrl), nil
}
