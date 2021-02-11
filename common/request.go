package common

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Request struct {
	client *http.Client
}

type Resp struct {
	Status     string
	StatusCode int
	Body       []byte
}

func NewRequest() *Request {
	return &Request{
		&http.Client{
			Transport: &http.Transport{
				DisableKeepAlives: false,
				TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: 6 * time.Second,
		},
	}
}

func (r *Request) Get(reqUrl string, data map[string]string) (*Resp, error) {
	var (
		request *http.Request
	)
	u, err := url.ParseRequestURI(reqUrl)
	if err != nil {
		return nil, err
	}

	query := u.Query()
	if data != nil {
		for k, v := range data {
			if len(v) > 0 {
				query.Set(k, v)
			}
		}
	}

	u.RawQuery = query.Encode()
	fmt.Println(u.String())

	request, err = http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	response, err := r.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	resp, err := parseResp(response)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *Request) PostJson(reqUrl string, data map[string]interface{}) (*Resp, error) {
	var (
		request *http.Request
	)

	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", reqUrl, buf)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := r.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	resp, err := parseResp(response)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *Request) PostForm(reqUrl string, data url.Values) (*Resp, error) {
	var (
		request  *http.Request
		postData string
	)
	postData = data.Encode()
	request, err := http.NewRequest("POST", reqUrl, strings.NewReader(postData))

	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := r.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	resp, err := parseResp(response)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func parseResp(response *http.Response) (*Resp, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	resp := &Resp{
		Body:       body,
		Status:     response.Status,
		StatusCode: response.StatusCode,
	}
	return resp, nil
}
