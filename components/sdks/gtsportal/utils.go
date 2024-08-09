package gtsportal

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/transport/http"
	"io/ioutil"
	stdhttp "net/http"
)

type jsonResult struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// 响应解析
func DecoderResponse() http.DecodeResponseFunc {
	return func(ctx context.Context, res *stdhttp.Response, v interface{}) error {
		if v == nil {
			return nil
		}

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		codec := http.CodecForResponse(res)
		rlt := &jsonResult{}
		err = codec.Unmarshal(data, rlt)
		if err != nil {
			return err
		}

		if rlt.Code != 0 {
			return fmt.Errorf("code:%v,msg:%v", rlt.Code, rlt.Message)
		}

		b, err := codec.Marshal(rlt.Data)
		return codec.Unmarshal(b, v)
	}
}
