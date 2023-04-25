package postgres_test

import (
	"context"
	"errors"
	pgdb "rest_url_shortener/internal/repository/postgres"
	"rest_url_shortener/internal/utils"
	"testing"
)

var (
	cfg = utils.StorageConfig{
		PgUsername: "test",
		PgPassword: "123454321",
		PgHost:     "localhost",
		PgPort:     "5433",
		PgDatabase: "url_shortener_db_test",
	}
	baseUrl      = "http.ozon.ru/internship"
	shortenedUrl = "randomString123_"
	ctx          = context.Background()
)

// Get connect to db for test
var pool, err = utils.GetConnectToPg(context.Background(), &cfg)
var db = pgdb.NewRepository(pool)

func TestPostgresRepository_SaveUrl(t *testing.T) {
	// database connection check
	if err != nil {
		t.Fatal(err)
	}
	// Check func Save url
	err = db.SaveUrl(ctx, baseUrl, shortenedUrl)
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

func TestPostgresRepository_GetBaseUrlByShort(t *testing.T) {
	// database connection check
	if err != nil {
		t.Fatal(err)
	}
	// Set url that does not exist in the database
	base, err := db.GetBaseUrlByShort(ctx, "randomFakeString")
	ourError := errors.New("unknown url").Error()
	if err.Error() != ourError {
		t.Fatalf("Incorrect response. Expected: %s, got: %s", ourError, err)
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

func TestPostgresRepository_GetShortUrlByBase(t *testing.T) {
	// database connection check
	if err != nil {
		t.Fatal(err)
	}

	// Set url that does not exist in the database
	_, exists, _ := db.GetShortUrlByBase(ctx, "https://randomFakeUrl.com/err_")
	if !exists {
		t.Errorf("Incorrect response. Expected: false, got: %t", exists)
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
