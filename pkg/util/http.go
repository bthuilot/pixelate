package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// HTTPRequest will make an HTTP `method` request to `uri`, and include `query` as the query params,
// `headers` as the headers and `body` as the body. It will parse the result into the response
func HTTPRequest[T interface{}](method, uri string, query map[string]string,
	headers map[string]string, body io.Reader, response *T) error {
	client := &http.Client{}
	req, err := http.NewRequest(
		method,
		addQueryParams(uri, query),
		body)
	for header, value := range headers {
		req.Header.Set(header, value)
	}
	if err != nil {
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("request to %s did not return 200", uri)
	}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(response)
	return err
}

// addQueryParams will add the key and values of `query` as query params to the `uri`
func addQueryParams(uri string, query map[string]string) string {
	queryParams := make([]string, len(query))
	for key, val := range query {
		param := fmt.Sprintf("%s=%s", url.QueryEscape(key), url.QueryEscape(val))
		queryParams = append(queryParams, param)
	}
	return fmt.Sprintf("%s?%s", uri, strings.Join(queryParams, "&"))
}
