package gostein

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		resp, err := sc.Get("/sheetname", SearchParams{})
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
		_, err := sc.Get("/sheetname", SearchParams{})
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
		_, err := sc.Get("/sheetname", SearchParams{})
		assert.NotNil(t, err)
		if !errors.As(err, &ErrDecodeJSON{}) {
			t.Errorf("Expected ErrDecode, got %v", err)
		}
	})
}

func TestSearchParams_queryString(t *testing.T) {
	t.Run("should return the correct query string", func(t *testing.T) {
		params := SearchParams{
			Offset: 20,
			Limit:  10,
			Conditions: map[string]string{
				"column_1": "value_column_1",
				"column_2": "value_column_2",
			},
		}

		queryString := params.queryString()
		assert.Equal(t, `offset=20&limit=10&search=%7B%22column_1%22%3A%22value_column_1%22%2C%22column_2%22%3A%22value_column_2%22%7D`, queryString)
	})
}
