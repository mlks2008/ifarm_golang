package utils

import (
	"bytes"
	"github.com/json-iterator/go"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func HttpGetRequest(url string) ([]byte, error) {

	client := http.Client{Timeout: 30 * time.Second}

	requestData, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	responseData, err := client.Do(requestData)
	if err != nil {
		return nil, err
	}

	defer func() {
		if responseData != nil && responseData.Body != nil {
			responseData.Body.Close()
		}
	}()

	return ioutil.ReadAll(responseData.Body)
}

func HttpPostRequest(url string, body []byte) ([]byte, error) {

	client := http.Client{Timeout: 30 * time.Second}

	ioBody := bytes.NewReader(body)

	requestData, err := http.NewRequest("POST", url, ioBody)
	if err != nil {
		return nil, err
	}

	responseData, err := client.Do(requestData)
	if err != nil {
		//l4g.Error("responseData err is", err)

		return nil, err
	}

	defer func() {
		if responseData != nil && responseData.Body != nil {
			responseData.Body.Close()
		}
	}()

	return ioutil.ReadAll(responseData.Body)
}

func UnDeserialize(r io.Reader, model interface{}) error {

	body, err := ioutil.ReadAll(r)
	if err != nil {
		//l4g.Error("UnDeserialize:%v", err)
		return err
	}

	err = jsoniter.Unmarshal(body, &model)
	if err != nil {
		//l4g.Error("UnDeserialize:%s", body)
		return err
	}

	return nil
}
