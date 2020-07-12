package source

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	discardURL = "https://localhost:9"
	sizeLimit  = 455_902
)

func generate(n int, seed int64) []string {
	if n > sizeLimit {
		n = sizeLimit
	}

	srand := rand.New(rand.NewSource(seed))

	tokens := make([]string, n)
	for i := 0; i < n; i++ {
		b := make([]byte, 22)
		for i := range b {
			b[i] = charset[srand.Intn(len(charset))]
		}
		tokens[i] = string(b)
	}

	return tokens
}

func newTestSource(t *testing.T, rawurl string) Source {
	cfgURL, err := url.Parse(rawurl)
	if err != nil {
		t.Fatalf("failed to create new source: %v", err)
	}

	cfg := &Config{
		Retry:   1,
		Timeout: 1 * time.Second,
		URL:     cfgURL,
	}

	return New(cfg)
}

func newTestServer(t *testing.T, seed int64) *httptest.Server {
	var handler http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", gin.MIMEHTML)
		res.Header().Set("Server", "Google Frontend")
		res.Header().Set("Date", time.Now().String())

		if req.Method != "POST" || req.URL.Path != "/" {
			http.NotFound(res, req)
			return
		}

		rawsize := defaultQuery(req.URL, "size", "0")
		size, err := strconv.Atoi(rawsize)
		if err != nil {
			// The API should return 400 instead of 200.
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(""))
			return
		}

		var b strings.Builder
		for _, token := range generate(size, seed) {
			fmt.Fprintf(&b, "%s\n", token)
		}

		// NOTE: The API should return 400 instead of 200.
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(b.String()))
	}

	return httptest.NewServer(handler)
}

func defaultQuery(url *url.URL, key, defaultValue string) string {
	if value, ok := getQuery(url, key); ok {
		return value
	}

	return defaultValue
}

func getQuery(url *url.URL, key string) (string, bool) {
	if url != nil {
		if values, ok := url.Query()[key]; ok && len(values) > 0 {
			return values[0], true
		}
	}

	return "", false
}
