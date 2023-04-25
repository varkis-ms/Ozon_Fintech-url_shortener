package url

import (
	"context"
	"log"
	pb "rest_url_shortener/internal/pb"
	"time"
)

/*
SaveUrl Метод контроллера (эндпоинт сервиса), который вызывает соответствующий метод сервиса
для генерации короткой ссылки, сохранения данных в хранилище и возвращения короткой ссылки
*/
func (c *UrlController) SaveUrl(ctx context.Context, req *pb.SaveUrlRequest) (*pb.SaveUrlResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	shortUrl, err := c.Service.AddShortUrl(ctx, req.GetBaseUrl())
	if err != nil {
		log.Printf("SaveUrl  ✘  --> an unexpected error occurred in the method: %s", err)
		return nil, err
	}
	return &pb.SaveUrlResponse{ShortUrl: shortUrl}, nil
}
