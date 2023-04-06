package utils

import (
	"bytes"
	"encoding/json"
	"github.com/google/go-querystring/query"
	"io"
	"net/http"
)

func request[R any](uri, method string, headers map[string]string, body interface{}, resp R) error {
	var reader io.Reader
	if body != nil {
		data, _ := json.Marshal(body)
		reader = bytes.NewBuffer(data)
	}
	client := &http.Client{}
	req, err := http.NewRequest(method, uri, reader)
	if err != nil {
		return err
	}
	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)
	dec := json.NewDecoder(res.Body)
	if err := dec.Decode(resp); nil != err {
		return err
	}
	return nil
}

func HttpGet[T any, R any](uri string, headers map[string]string, params T, resp R) error {
	v, _ := query.Values(params)
	return request(uri+"?"+v.Encode(), http.MethodGet, headers, nil, resp)
}

func HttpPost[T any, R any](uri string, headers map[string]string, body T, resp R) error {
	return request(uri, http.MethodPost, headers, body, resp)
}
