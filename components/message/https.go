package message

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

func HttpGet(url string, timeout int64) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	httpClient := http.Client{Timeout: time.Second * time.Duration(timeout)}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func HttpPostRequest(url string, body []byte) ([]byte, error) {

	client := http.Client{Timeout: 10 * time.Second}

	ioBody := bytes.NewReader(body)

	requestData, err := http.NewRequest("POST", url, ioBody)
	if err != nil {
		return nil, err
	}

	requestData.Header.Set("Content-Type", "application/json")

	responseData, err := client.Do(requestData)
	if err != nil {
		//l4g.Error("responseData err is %v", err)
		return nil, err
	}

	defer func() {
		if responseData != nil && responseData.Body != nil {
			responseData.Body.Close()
		}
	}()

	return ioutil.ReadAll(responseData.Body)
}
