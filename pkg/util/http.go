package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func HTTPRequest[T interface{}](baseURL string, params map[string]string,
	headers map[string]string, body io.Reader, response *T) error {
	client := &http.Client{}
	queryParams := make([]string, len(params))
	for key, val := range params {
		param := fmt.Sprintf("%s=%s", url.QueryEscape(key), url.QueryEscape(val))
		queryParams = append(queryParams, param)
	}
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s?%s", baseURL, strings.Join(queryParams, "&")),
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
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(response)
	return err
}