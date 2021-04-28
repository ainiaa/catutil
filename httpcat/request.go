package httpcat

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

//通用http请求
func Request(ctx context.Context, uri string, method string,
	data string, header map[string]string, timeout int64) (re []byte) {
	if method == "" {
		method = http.MethodPost
	}
	req, err := http.NewRequest(method, uri, strings.NewReader(data))
	if err != nil {
		return re
	}
	defer req.Body.Close()
	for k, v := range header {
		req.Header.Add(k, v)
	}
	if req.Header.Get("Content-Type") == "" {
		req.Header.Add("Content-Type", "application/json")
	}

	if timeout == 0 {
		//默认超时2000ms
		timeout = 2000
	}

	clt := http.Client{
		Timeout: time.Duration(timeout) * time.Millisecond,
	}
	resp, err := clt.Do(req)
	if err != nil {
		return re
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp != nil && resp.Body != nil {
		_ = resp.Body.Close()
	}
	return body
}

func RequestWithResp(ctx context.Context, uri string, method string, data string, header map[string]string, timeout int64) (re []byte, response *http.Response) {
	if method == "" {
		method = "POST"
	}
	req, err := http.NewRequest(method, uri, strings.NewReader(data))
	if err != nil {
		return re, nil
	}
	defer req.Body.Close()
	if header != nil {
		for k, v := range header {
			req.Header.Add(k, v)
		}
	}
	if req.Header.Get("Content-Type") == "" {
		req.Header.Add("Content-Type", "application/json")
	}

	if timeout == 0 {
		//默认超时2000ms
		timeout = 2000
	}

	clt := http.Client{
		Timeout: time.Duration(timeout) * time.Millisecond,
	}
	resp, err := clt.Do(req)
	if err != nil {
		return re, nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return body, resp
}

func RequestUrlencoded(ctx context.Context, uri string, method string, data map[string]string, header map[string]string, timeout int64) (re []byte) {
	if method == "" {
		method = "POST"
	}
	dataUrlVal := url.Values{}
	for key, val := range data {
		dataUrlVal.Add(key, val)
	}
	req, err := http.NewRequest(method, uri, strings.NewReader(dataUrlVal.Encode()))
	if err != nil {
		return re
	}
	defer req.Body.Close()
	for k, v := range header {
		req.Header.Add(k, v)
	}
	if req.Header.Get("Content-Type") == "" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	if timeout == 0 {
		//默认超时2000ms
		timeout = 2000
	}

	clt := http.Client{
		Timeout: time.Duration(timeout) * time.Millisecond,
	}
	resp, err := clt.Do(req)
	if err != nil {
		return re
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp != nil && resp.Body != nil {
		_ = resp.Body.Close()
	}
	return body
}
