//go:build integration
// +build integration

package integration_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"os"
	pb "rest_url_shortener/internal/pb"
	"testing"
)

type UrlShortenerSuite struct {
	suite.Suite
	ctx             context.Context
	shortenerConn   *grpc.ClientConn
	shortenerClient pb.UrlShortenerClient
	shortUrl        string
	baseUrl         string
}

func (s *UrlShortenerSuite) SetupSuite() {
	shortenerHost := os.Getenv("SHORTENER_SERVER_HOST")
	if shortenerHost == "" {
		shortenerHost = "127.0.0.1:8080"
	}
	shortenerConn, err := grpc.Dial(shortenerHost, grpc.WithInsecure())
	s.Require().NoError(err)

	s.ctx = context.Background()
	s.shortenerClient = pb.NewUrlShortenerClient(shortenerConn)
	s.baseUrl = "https://ozon.ru"
	s.shortUrl = "RaNd_url23"
}

func (s *UrlShortenerSuite) TestSaveUrl() {
	req := &pb.SaveUrlRequest{
		BaseUrl: s.baseUrl,
	}

	// Check that there is a response, and it does not match the BaseUrl
	resp1, err := s.shortenerClient.SaveUrl(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotEqual(resp1.GetShortUrl(), s.baseUrl)

	// Checking that the response to the same baseUrl will be the same
	resp2, err := s.shortenerClient.SaveUrl(s.ctx, req)
	s.Require().NoError(err)
	s.Require().Equal(resp1.GetShortUrl(), resp2.GetShortUrl())
}

func (s *UrlShortenerSuite) TestGetUrl() {
	req := &pb.GetUrlRequest{
		ShortUrl: s.shortUrl,
	}
	// Check the handling of the situation when an unknown shortUrl was passed to
	_, err := s.shortenerClient.GetUrl(s.ctx, req)
	s.Require().ErrorContains(err, "unknown url")

	// Check that there is a response, and it matches the BaseUrl
	reqSave := &pb.SaveUrlRequest{
		BaseUrl: s.baseUrl,
	}

	// Check that there is a response, and it does not match the BaseUrl
	resp1, err := s.shortenerClient.SaveUrl(s.ctx, reqSave)
	s.Require().NoError(err)
	fmt.Println(resp1, resp1.GetShortUrl())
	short := resp1.GetShortUrl()

	checkReq := &pb.GetUrlRequest{
		ShortUrl: short[len(short)-10:],
	}

	resp2, err := s.shortenerClient.GetUrl(s.ctx, checkReq)
	s.Require().NoError(err)
	s.Require().Equal(resp2.GetBaseUrl(), s.baseUrl)

}

func TestUrlShortenerSuite(t *testing.T) {
	suite.Run(t, new(UrlShortenerSuite))
}
