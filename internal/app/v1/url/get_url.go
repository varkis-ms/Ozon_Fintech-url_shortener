package url

import (
	"context"
	"log"
	pb "rest_url_shortener/internal/pb"
	"time"
)

/*
GetUrl Метод контроллера (эндпоинт сервиса), который вызывает соответствующий метод сервиса
для нахождения изначальной ссылки по короткой, достав это значение из хранилища, возвращает изначальную ссылку
*/
func (c *UrlController) GetUrl(ctx context.Context, req *pb.GetUrlRequest) (*pb.GetUrlResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	baseUrl, err := c.Service.ForwardBaseUrl(ctx, req.GetShortUrl())
	if err != nil {
		log.Printf("GetUrl  ✘  --> an unexpected error occurred in the method: %s", err)
		return &pb.GetUrlResponse{BaseUrl: baseUrl}, err
	}
	return &pb.GetUrlResponse{BaseUrl: baseUrl}, nil
}
