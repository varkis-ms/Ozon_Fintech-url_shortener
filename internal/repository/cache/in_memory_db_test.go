package cache_test

import (
	"context"
	"errors"
	c "rest_url_shortener/internal/repository/cache"
	"testing"
)

var baseUrl = "https://test_tring.com/index.html"
var shortenedUrl = "randomString"
var ctx = context.Background()

func TestInMemoryRepository_SaveUrl(t *testing.T) {
	// Init in-memory database
	db := c.NewInMemoryRepository()

	// Check func Save url
	err := db.SaveUrl(ctx, baseUrl, shortenedUrl)
	if err != nil {
		t.Errorf("Save error. Expected: nil, got: %s", err.Error())
	}
	base, err := db.GetBaseUrlByShort(ctx, shortenedUrl)
	if err != nil {
		t.Fatal(err)
	}
	if base != baseUrl {
		t.Errorf("Incorrect response. Expected: %s, got: %s", baseUrl, base)
	}

	// Checking response when baseurl is already in database
	err = db.SaveUrl(ctx, baseUrl, shortenedUrl)
	if err != nil {
		t.Errorf("Error not handled. Expected: nil, got: %s", err)
	}
}

func TestInMemoryRepository_GetBaseUrlByShort(t *testing.T) {
	// Init in-memory database
	db := c.NewInMemoryRepository()

	// Set url that does not exist in the database
	base, err := db.GetBaseUrlByShort(ctx, shortenedUrl)
	ourError := errors.New("unknown url").Error()
	if err.Error() != ourError {
		t.Fatalf("Incorrect response. Expected: %s, got: %s", ourError, err)
	}

	// Save baseUrl
	err = db.SaveUrl(ctx, baseUrl, shortenedUrl)
	if err != nil {
		t.Fatal(err)
	}

	// Check that by baseUrl we find the correct shorUrl
	base, err = db.GetBaseUrlByShort(ctx, shortenedUrl)
	if err != nil {
		t.Fatal(err)
	}
	if base != baseUrl {
		t.Errorf("Incorrect response. Expected: %s, got: %s", baseUrl, base)
	}

}

func TestInMemoryRepository_GetShortUrlByBase(t *testing.T) {
	// Init in-memory database
	db := c.NewInMemoryRepository()

	// Set url that does not exist in the database
	_, exists, _ := db.GetShortUrlByBase(ctx, baseUrl)
	if !exists {
		t.Errorf("Incorrect response. Expected: false, got: %t", exists)
	}

	// Save baseUrl
	err := db.SaveUrl(ctx, baseUrl, shortenedUrl)
	if err != nil {
		t.Fatal(err)
	}

	// Check that by baseUrl we find the correct shorUrl
	short, exists, _ := db.GetShortUrlByBase(ctx, baseUrl)
	if exists {
		t.Errorf("Incorrect response. Expected: false, got: %t", exists)
	}
	if short != shortenedUrl {
		t.Errorf("Incorrect response. Expected: %s, got: %s", shortenedUrl, short)
	}
}
