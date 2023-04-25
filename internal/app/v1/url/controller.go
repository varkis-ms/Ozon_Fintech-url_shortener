package url

import (
	pb "rest_url_shortener/internal/pb"
	"rest_url_shortener/internal/service/url"
)

// UrlController Создание контроллера для вызова методов сервиса и связи с proto
type UrlController struct {
	pb.UnimplementedUrlShortenerServer
	Service *url.Service
}

// NewUrlController Конструктор для создания контроллера (UrlController)
func NewUrlController(service *url.Service) *UrlController {
	return &UrlController{Service: service}
}
