package utils_test

import (
	u "rest_url_shortener/internal/utils"
	"strconv"
	"strings"
	"testing"
)

func TestEncode63(t *testing.T) {
	short := u.Encode63()

	// Correct length
	if len(short) != 10 {
		t.Errorf("length of short url is %d, want 10", len(short))
	}

	// Correct symbols
	for c := range short {
		if !strings.Contains(u.ALPHABET, strconv.Itoa(c)) {
			t.Errorf("invalid character in short url: %c", c)
		}
	}

	// Uniqueness
	var cnt = 100_000
	var shorts = map[string]struct{}{}
	for i := 0; i < cnt; i++ {
		shorts[u.Encode63()] = struct{}{}
	}
	if len(shorts) != cnt {
		t.Errorf("not all values are unique. Expected %d, got %d", cnt, len(shorts))
	}

}
