package gostein

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type errorResponse struct {
	Error string `json:"error"`
}

// GetParams is the parameters for the read operations: Get or search
type GetParams struct {
	Offset int64
	Limit  int64
	Condition map[string]string
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

	if len(gp.Condition) > 0 {
		jsonSearch, err := json.Marshal(gp.Condition)
		if err != nil {
			return ""
		}

		querySearchJSON := url.QueryEscape(string(jsonSearch))

		queryString = queryString + fmt.Sprintf("search=%s", querySearchJSON)
	}

	return removeSuffix(queryString, "&")
}

// UpdateParams is the parameters for the update operation
type UpdateParams struct {
	Condition map[string]string `json:"condition"`
	Set       map[string]string `json:"set"`
	Limit     int64             `json:"limit,omitempty"`
}

// DeleteParams is the parameters for the delete operation
type DeleteParams struct {
	Condition map[string]string `json:"condition"`
	Limit     int64             `json:"limit,omitempty"`
}

type BasicAuth struct {
	Username string
	Password string
}

// Interface is the interface for the stein client
type Interface interface {
	Get(sheet string, params GetParams) ([]map[string]interface{}, error)
	Add(sheet string, rows ...map[string]interface{}) (usedSheetRange string, err error)
	Update(sheet string, params UpdateParams) (countUpdatedRows int64, err error)
	Delete(sheet string, params DeleteParams) (countDeletedRows int64, err error)
}

type stein struct {
	url        string
	httpClient *http.Client
	basicAuth  *BasicAuth
}

// New creates a new stein client
// url is the base url of the stein api. i.e: https://api.steinhq.com/v1/storages/[your-api-id]
// httpClient is the http client to use. If nil, http.DefaultClient will be used
func New(url string, httpClient *http.Client, basicAuth *BasicAuth) Interface {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	url = removeSuffix(url, "/")

	return &stein{
		url:        url,
		httpClient: httpClient,
		basicAuth:  basicAuth,
	}
}

func (s *stein) decodeJSON(r io.Reader, v interface{}) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(v)
}

func (s *stein) doRequest(method string, resourceURL string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, resourceURL, body)
	if err != nil {
		return nil, err
	}

	if s.basicAuth != nil {
		req.SetBasicAuth(s.basicAuth.Username, s.basicAuth.Password)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	return s.httpClient.Do(req)
}

func (s *stein) validateResponse(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrNot2XX{StatusCode: resp.StatusCode, Message: "Unauthorized"}
	}

	var errorResponse errorResponse
	err := s.decodeJSON(resp.Body, &errorResponse)
	if err != nil {
		return err
	}

	return ErrNot2XX{StatusCode: resp.StatusCode, Message: errorResponse.Error}
}

// Get returns the rows in the given sheet
func (s *stein) Get(sheet string, params GetParams) ([]map[string]interface{}, error) {
	resource := fmt.Sprintf("%s/%s", s.url, removePrefix(sheet, "/"))

	queryParams := params.queryString()
	if queryParams != "" {
		resource = resource + "?" + queryParams
	}

	resp, err := s.doRequest(http.MethodGet, resource, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = s.validateResponse(resp)
	if err != nil {
		return nil, err
	}

	data := make([]map[string]interface{}, 0)
	err = s.decodeJSON(resp.Body, &data)
	if err != nil {
		return nil, ErrDecodeJSON{Err: err}
	}

	return data, nil
}

func (s *stein) Add(sheet string, rows ...map[string]interface{}) (usedSheetRange string, err error) {
	var (
		result = struct {
			UsedSheetRange string `json:"updatedRange"`
		}{}

		resource = fmt.Sprintf("%s/%s", s.url, removePrefix(sheet, "/"))
	)

	jsonRow, err := json.Marshal(rows)
	if err != nil {
		return result.UsedSheetRange, err
	}

	resp, err := s.doRequest(http.MethodPost, resource, bytes.NewReader(jsonRow))
	if err != nil {
		return result.UsedSheetRange, err
	}

	defer resp.Body.Close()

	err = s.validateResponse(resp)
	if err != nil {
		return result.UsedSheetRange, err
	}

	err = s.decodeJSON(resp.Body, &result)
	if err != nil {
		return result.UsedSheetRange, ErrDecodeJSON{Err: err}
	}

	return result.UsedSheetRange, nil
}

func (s *stein) Update(sheet string, params UpdateParams) (countUpdatedRows int64, err error) {
	var (
		result = struct {
			TotalUpdatedRows int64 `json:"totalUpdatedRows"`
		}{}

		resource = fmt.Sprintf("%s/%s", s.url, removePrefix(sheet, "/"))
	)

	// prevent error "Expected property `condition` to be of type `object` but received type `null` in object"
	if params.Condition == nil {
		params.Condition = make(map[string]string)
	}

	if params.Set == nil {
		return result.TotalUpdatedRows, ErrParamCantBeNil{ParamName: "set"}
	}

	payload, err := json.Marshal(params)
	if err != nil {
		return result.TotalUpdatedRows, err
	}

	resp, err := s.doRequest(http.MethodPut, resource, bytes.NewReader(payload))
	if err != nil {
		return result.TotalUpdatedRows, err
	}

	defer resp.Body.Close()

	err = s.validateResponse(resp)
	if err != nil {
		return result.TotalUpdatedRows, err
	}

	err = s.decodeJSON(resp.Body, &result)
	if err != nil {
		return result.TotalUpdatedRows, ErrDecodeJSON{Err: err}
	}

	return result.TotalUpdatedRows, nil
}

func (s *stein) Delete(sheet string, params DeleteParams) (countDeletedRows int64, err error) {
	var (
		result = struct {
			CountDeletedRows int64 `json:"clearedRowsCount"`
		}{}

		resource = fmt.Sprintf("%s/%s", s.url, removePrefix(sheet, "/"))
	)

	// prevent error "Expected property `condition` to be of type `object` but received type `null` in object"
	if params.Condition == nil {
		params.Condition = make(map[string]string)
	}

	payload, err := json.Marshal(params)
	if err != nil {
		return countDeletedRows, err
	}

	resp, err := s.doRequest(http.MethodDelete, resource, bytes.NewReader(payload))
	if err != nil {
		return countDeletedRows, err
	}

	defer resp.Body.Close()

	err = s.validateResponse(resp)
	if err != nil {
		return countDeletedRows, err
	}

	err = s.decodeJSON(resp.Body, &result)
	if err != nil {
		return countDeletedRows, ErrDecodeJSON{Err: err}
	}

	return result.CountDeletedRows, nil
}
