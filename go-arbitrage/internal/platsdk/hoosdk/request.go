/**
 * request.go
 * ============================================================================
 *
 * ============================================================================
 * author: soyitech
 * createtime: 2018/11/12 14:26
 */

package hoosdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func (self *Client) executeRequest(method string, bodyParameters url.Values) ([]byte, error) {
	var (
		err      error
		response *http.Response
		reqUrl   string
		data     []byte
	)

	//defer func() {
	//	if err != nil {
	//		l4g.Error("request url:%v params:%+v response:%+v err:%v", method, bodyParameters, "", err)
	//	} else {
	//		l4g.Debug("request url:%v params:%+v response:%+v", method, bodyParameters, "")
	//	}
	//}()

	if bodyParameters == nil {
		reqUrl = method
	} else {
		reqUrl = fmt.Sprintf("%s?%s", method, bodyParameters.Encode())
	}

	request, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{}

	response, err = client.Do(request)
	if err != nil {
		return nil, err
	}

	data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (self *Client) executePostRequest(method string, body interface{}) ([]byte, error) {
	var (
		err      error
		response *http.Response
		data     []byte
	)

	//defer func() {
	//	if err != nil {
	//		_ = l4g.Error("request url:%v params:%+v response:%+v err:%v", method, body, "", err)
	//	} else {
	//		l4g.Debug("request url:%v params:%+v response:%+v", method, body, "")
	//	}
	//}()

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", method, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	client := http.Client{}

	response, err = client.Do(request)
	if err != nil {
		return nil, err
	}

	data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
