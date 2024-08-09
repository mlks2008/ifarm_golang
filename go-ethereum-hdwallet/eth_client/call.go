package eth

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func (this *Client) callNode(apirequest *http.Request, method string, bodyid uint, bodyParameters []interface{}, model interface{}) error {
	var err error

	var resBody []byte

	body, err := NewBody(method, int64(bodyid), bodyParameters)
	if err != nil {
		err = fmt.Errorf("err1 %v", err.Error())
		return err
	}

	request, err := http.NewRequest(http.MethodPost, this.getNodeUrl(method), bytes.NewReader(body))
	if err != nil {
		err = fmt.Errorf("err2 %v", err.Error())
		return err
	}

	request.Header.Set("Content-Type", "application/json;charset=utf-8")

	//if apirequest != nil {
	//	if val := apirequest.Header.Get("Authorization"); val != "" {
	//		request.Header.Set("Authorization", val)
	//	}
	//}

	client := http.Client{}
	if find := strings.Contains(this.getNodeUrl(method), "blockchain.ramestta.com"); find {
		client = http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	}

	response, err := client.Do(request)
	if err != nil {
		err = fmt.Errorf("err3 %v", err.Error())
		return err
	}
	defer response.Body.Close()

	resBody, err = ioutil.ReadAll(response.Body)
	if err != nil {
		err = fmt.Errorf("err4 %v", err.Error())
		return err
	}

	if err = json.Unmarshal(resBody, &model); err != nil {
		err = fmt.Errorf("err5 %v", err.Error())
		return errors.New(fmt.Sprintf("err:%v response:%s", err.Error(), resBody))
	}

	return nil
}

func (this *Client) getNodeUrl(method string) string {
	if this.mainName == "OPTIMISM" && method == "eth_sendRawTransaction" {
		return this.otherUrl
	}

	if this.nodeIndex < len(this.nodeUrls) {
		var url = this.nodeUrls[this.nodeIndex]
		this.nodeIndex++
		return url
	} else {
		var url = this.nodeUrls[0]
		this.nodeIndex = 0
		return url
	}
}

func (this *Client) stat() {

}
