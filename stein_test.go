package stein

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_stein_getFullURL(t *testing.T) {
	s := &stein{}

	t.Run("should return the correct url", func(t *testing.T) {
		s.url = "https://api.steinhq.com/v1/storages/[your-api-id]"
		got := s.getFullURL("Sheet1", map[string]string{"offset": "4", "limit": "10"})
		if !strings.HasPrefix(got, "https://api.steinhq.com/v1/storages/[your-api-id]/Sheet1") ||
			!strings.Contains(got, "offset=4") || 
			!strings.Contains(got, "limit=10") {
			t.Errorf("stein.getFullURL() = %v, want %v", got, "https://api.steinhq.com/v1/storages/[your-api-id]/Sheet1?offset=0&limit=10")
		}
	})

	t.Run("should remove slash on path prefix", func(t *testing.T) {
		s.url = "https://api.steinhq.com/v1/storages/[your-api-id]"
		got := s.getFullURL("/Sheet1", nil)
		assert.Equal(t, "https://api.steinhq.com/v1/storages/[your-api-id]/Sheet1", got)
	})
}

func Test_stein_addParams(t *testing.T) {
	s := &stein{}

	t.Run("should return the full url with zero params", func(t *testing.T) {
		got := s.addParams("http://url.com/sheetname", map[string]string{})
		assert.Equal(t, "http://url.com/sheetname", got)
	})

	t.Run("should return the valid url with mulitple params", func(t *testing.T) {
		got := s.addParams("http://url.com/sheetname", map[string]string{"offset": "4", "limit": "10"})
		segment := strings.Split(got, "?")
		assert.Equal(t, 2, len(segment))
		assert.Equal(t, "http://url.com/sheetname", segment[0])
		assert.Contains(t, segment[1], "offset=4")
		assert.Contains(t, segment[1], "limit=10")
	})
	
}

func Test_stein_Get(t *testing.T) {
	t.Run("should return the correct response", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			jsonBody := `[{
				"title": "Escort Warrior",
				"url": "https://url.com/sheetname"
			}, {
				"title": "Bowblade Spirit",
				"url": "https://url.com/sheetname"
			}]`

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(jsonBody))
		}))
		defer ts.Close()

		sc := New(ts.URL, nil)
		resp, err := sc.Get("/sheetname", BasicQuery{})
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		if len(resp) != 2 {
			t.Errorf("Expected 2 results, got %v", len(resp))
		}

		assert.Equal(t, "Escort Warrior", resp[0]["title"].(string))
		assert.Equal(t, "Bowblade Spirit", resp[1]["title"].(string))
	})

	t.Run("should return error if http code not 2xx", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
		}))
		defer ts.Close()

		sc := New(ts.URL, nil)
		_, err := sc.Get("/sheetname", BasicQuery{})
		assert.NotNil(t, err)
		if !errors.As(err, &ErrNot2XX{}) {
			t.Errorf("Expected ErrNot2XX, got %v", err)
		}
	})

	t.Run("should return error on failed decode the response", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			jsonBody := `bad json`

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(jsonBody))
		}))
		defer ts.Close()

		sc := New(ts.URL, nil)
		_, err := sc.Get("/sheetname", BasicQuery{})
		assert.NotNil(t, err)
		if !errors.As(err, &ErrDecodeJSON{}) {
			t.Errorf("Expected ErrDecode, got %v", err)
		}
	})
}
