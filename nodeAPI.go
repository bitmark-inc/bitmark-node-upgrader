package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/google/logger"
)

// NodeAPI run bitmarkd api
func NodeAPI(endpoint string, action NodeAction) ([]byte, error) {
	postBody := ""
	api := ""
	url := ""
	switch action {
	case bitmarkdStart:
		postBody = "{\"option\": \"start\"}\n"
		api = "/api/bitmarkd"
	case bitmarkdStop:
		postBody = "{\"option\": \"stop\"}\n"
		api = "/api/bitmarkd"
	case recorderdStart:
		postBody = "{\"option\": \"start\"}\n"
		api = "/api/recorderd"
	case recorderdStop:
		postBody = "{\"option\": \"stop\"}\n"
		api = "/api/recorderd"
	default:
		return nil, errors.New("no such bitmarkd action")
	}

	if len(endpoint) != 0 {
		url = endpoint + api
	} else {
		url = "http://127.0.0.1:9980" + api
	}

	log.Info(url)

	payload := strings.NewReader(postBody)
	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
