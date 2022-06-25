package stein

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type BasicQuery struct {
	Offset int64
	Limit  int64
}

// Interface is the interface for the stein client
type Interface interface {
	Get(sheet string, basicQuery BasicQuery) ([]map[string]interface{}, error)
}

type stein struct {
	url        string
	httpClient *http.Client
}

// New creates a new stein client
// url is the base url of the stein api. i.e: https://api.steinhq.com/v1/storages/[your-api-id]
// httpClient is the http client to use. If nil, http.DefaultClient will be used
func New(url string, httpClient *http.Client) Interface {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	if strings.HasSuffix(url, "/") {
		url = strings.TrimSuffix(url, "/")
	}

	return &stein{
		url:        url,
		httpClient: httpClient,
	}
}

func (s *stein) getFullURL(path string, params map[string]string) string {
	if strings.HasPrefix(path, "/") {
		path = strings.TrimPrefix(path, "/")
	}
	resource := fmt.Sprintf("%s/%s", s.url, path)

	return s.addParams(resource, params)
}

func (s *stein) addParams(resource string, params map[string]string) string {
	if len(params) == 0 {
		return resource
	}

	resource = resource + "?"
	for k, v := range params {
		resource = resource + fmt.Sprintf("%s=%s&", k, v)
	}

	return strings.TrimSuffix(resource, "&")
}

func (s *stein) decodeJSON(r io.Reader, v interface{}) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(v)
}

// Get returns the rows in the given sheet
func (s *stein) Get(sheet string, basicQuery BasicQuery) ([]map[string]interface{}, error) {
	mapQuery := map[string]string{}
	if basicQuery.Offset > 0 {
		mapQuery["offset"] = fmt.Sprintf("%d", basicQuery.Offset)
	}
	if basicQuery.Limit > 0 {
		mapQuery["limit"] = fmt.Sprintf("%d", basicQuery.Limit)
	}

	resource := s.getFullURL(sheet, mapQuery)

	resp, err := s.httpClient.Get(resource)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrNot2XX{StatusCode: resp.StatusCode}
	}

	data := make([]map[string]interface{}, 0)
	err = s.decodeJSON(resp.Body, &data)
	if err != nil {
		return nil, ErrDecodeJSON{Err: err}
	}

	return data, nil
}
