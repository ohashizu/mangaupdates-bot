package utils

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
)

var (
	client = &fasthttp.Client{}
)

func getRequest(url string) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod("GET")
	req.SetRequestURI(url)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := client.Do(req, resp)
	if err != nil {
		return nil, err
	}

	bodyCopy := make([]byte, len(resp.Body()))
	copy(bodyCopy, resp.Body())

	return bodyCopy, nil
}

func muGetRequest(url, token string) ([]byte, int, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod("GET")
	req.SetRequestURI(url)
	req.Header.Set("Authorization", "Bearer "+token)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := client.Do(req, resp)
	if err != nil {
		return nil, 0, err
	}

	bodyCopy := make([]byte, len(resp.Body()))
	copy(bodyCopy, resp.Body())

	statusCode := resp.StatusCode()

	return bodyCopy, statusCode, nil
}

func muPostRequest(url, token string, body interface{}) ([]byte, int, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod("POST")
	req.SetRequestURI(url)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, 0, err
	}
	req.SetBody(jsonBody)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err = client.Do(req, resp)
	if err != nil {
		return nil, 0, err
	}

	bodyCopy := make([]byte, len(resp.Body()))
	copy(bodyCopy, resp.Body())

	statusCode := resp.StatusCode()

	return bodyCopy, statusCode, nil
}
