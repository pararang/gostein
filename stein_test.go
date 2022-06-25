package stein

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_stein_getFullURL(t *testing.T) {
	type fields struct {
		url        string
		httpClient *http.Client
	}
	type args struct {
		path   string
		params map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "should return the full url with expected input",
			fields: fields{
				url: "https://api.steinhq.com/v1/storages/[your-api-id]",
			},
			args: args{
				path:   "Sheet1",
				params: map[string]string{"offset": "0", "limit": "10"},
			},
			want: "https://api.steinhq.com/v1/storages/[your-api-id]/Sheet1?offset=0&limit=10",
		},
		{
			name: "should remove slash on path prefix",
			fields: fields{
				url: "http://example.com",
			},
			args: args{
				path:   "/sheetName",
				params: map[string]string{},
			},
			want: "http://example.com/sheetName",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &stein{}
			s.url = tt.fields.url
			if got := s.getFullURL(tt.args.path, tt.args.params); got != tt.want {
				t.Errorf("stein.getFullURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stein_addParams(t *testing.T) {
	type args struct {
		resource string
		params   map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should return the full url with zero params",
			args: args{
				resource: "http://url.com/sheetname",
				params:   map[string]string{},
			},
			want: "http://url.com/sheetname",
		},
		{
			name: "should return the valid url with mulitple params",
			args: args{
				resource: "http://url.com/sheetname",
				params: map[string]string{
					"limit": "10",
					"sort":  "asc",
				},
			},
			want: "http://url.com/sheetname?limit=10&sort=asc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &stein{}
			if got := s.addParams(tt.args.resource, tt.args.params); got != tt.want {
				t.Errorf("stein.addParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

type clientMock struct{}

func (c *clientMock) Do(req *http.Request) (*http.Response, error) {
	jsonBody := `[{
		"title": "Escort Warrior",
		"url": "https://url.com/sheetname"
	}, {
		"title": "Bowblade Spirit",
		"url": "https://url.com/sheetname"
	}]`

	resp := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBufferString(jsonBody)),
	}

	return resp, nil
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
