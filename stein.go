package gostein

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type SearchParams struct {
	Offset int64
	Limit  int64
	Conditions map[string]string
}

// builds the query string from the given params with query string escaping
func (sp SearchParams) queryString() string {
	queryString := ""
	if sp.Offset > 0 {
		queryString = queryString + fmt.Sprintf("offset=%d&", sp.Offset)
	}

	if sp.Limit > 0 {
		queryString = queryString + fmt.Sprintf("limit=%d&", sp.Limit)
	}

	if len(sp.Conditions) > 0 {
		jsonSearch, err := json.Marshal(sp.Conditions)
		if err != nil {
			return ""
		}

		querySearchJSON := url.QueryEscape(string(jsonSearch))

		queryString = queryString + fmt.Sprintf("search=%s", querySearchJSON)
	}

	return strings.TrimSuffix(queryString, "&")
}

// Interface is the interface for the stein client
type Interface interface {
	Get(sheet string, params SearchParams) ([]map[string]interface{}, error)
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

	url = removeSuffix(url, "/")

	return &stein{
		url:        url,
		httpClient: httpClient,
	}
}

func (s *stein) decodeJSON(r io.Reader, v interface{}) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(v)
}

// Get returns the rows in the given sheet
func (s *stein) Get(sheet string, params SearchParams) ([]map[string]interface{}, error) {
	resource := fmt.Sprintf("%s/%s", s.url, removePrefix(sheet, "/"))
	
	queryParams := params.queryString()
	if queryParams != "" {
		resource = resource + "?" + queryParams 
	}

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
