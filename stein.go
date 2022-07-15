package gostein

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type GetParams struct {
	Offset int64
	Limit  int64
	Search map[string]string
}

// builds the query string from the given params with query string escaping
func (gp GetParams) queryString() string {
	queryString := ""
	if gp.Offset > 0 {
		queryString = queryString + fmt.Sprintf("offset=%d&", gp.Offset)
	}

	if gp.Limit > 0 {
		queryString = queryString + fmt.Sprintf("limit=%d&", gp.Limit)
	}

	if len(gp.Search) > 0 {
		jsonSearch, err := json.Marshal(gp.Search)
		if err != nil {
			return ""
		}

		querySearchJSON := url.QueryEscape(string(jsonSearch))

		queryString = queryString + fmt.Sprintf("search=%s", querySearchJSON)
	}

	return removeSuffix(queryString, "&")
}

type WriteResponse struct {
	UpdatedRange string `json:"updatedRange"`
}

type UpdateParams struct {
	Condition map[string]string `json:"condition"`
	Set       map[string]string `json:"set"`
	Limit     int64             `json:"limit,omitempty"`
}

// Interface is the interface for the stein client
type Interface interface {
	Get(sheet string, params GetParams) ([]map[string]interface{}, error)
	Add(sheet string, rows ...map[string]interface{}) (WriteResponse, error)
	Update(sheet string, params UpdateParams) (WriteResponse, error)
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
func (s *stein) Get(sheet string, params GetParams) ([]map[string]interface{}, error) {
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

func (s *stein) Add(sheet string, rows ...map[string]interface{}) (WriteResponse, error) {
	var (
		result WriteResponse
		resource = fmt.Sprintf("%s/%s", s.url, removePrefix(sheet, "/"))
	)

	jsonRow, err := json.Marshal(rows)
	if err != nil {
		return result, err
	}

	resp, err := s.httpClient.Post(resource, "application/json", strings.NewReader(string(jsonRow)))
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, ErrNot2XX{StatusCode: resp.StatusCode}
	}

	err = s.decodeJSON(resp.Body, &result)
	if err != nil {
		return result, ErrDecodeJSON{Err: err}
	}

	return result, nil
}

func (s *stein) Update(sheet string, params UpdateParams) (WriteResponse, error) {
	var (
		result WriteResponse
		resource = fmt.Sprintf("%s/%s", s.url, removePrefix(sheet, "/"))
	)

	payload, err := json.Marshal(params)
	if err != nil {
		return result, err
	}

	req, err := http.NewRequest(http.MethodPut, resource, bytes.NewBuffer(payload))
	if err != nil {
		return result, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, ErrNot2XX{StatusCode: resp.StatusCode}
	}

	err = s.decodeJSON(resp.Body, &result)
	if err != nil {
		return WriteResponse{}, ErrDecodeJSON{Err: err}
	}

	return result, nil
}
