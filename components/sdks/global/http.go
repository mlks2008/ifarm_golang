package global

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"io"
	stdhttp "net/http"
	"reflect"
)

type jsonResult struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

func EncoderRequest(log *log.Helper) http.EncodeRequestFunc {
	return func(ctx context.Context, contentType string, in interface{}) (body []byte, err error) {
		codec := encoding.GetCodec("json")
		body, err = codec.Marshal(in)
		if err != nil {
			log.Warnf("[http][local][EncoderRequest] %v:%v", reflect.TypeOf(in).Elem().Name(), in)
			return nil, errors.BadRequest("CODEC", err.Error())
		} else {
			log.Debugf("[http][local][EncoderRequest] %v:%s", reflect.TypeOf(in).Elem().Name(), body)
		}
		return
	}
}

func ResponseDecoder(log *log.Helper) http.DecodeResponseFunc {
	return func(ctx context.Context, res *stdhttp.Response, out interface{}) error {
		if out == nil {
			return nil
		}

		var data []byte
		var err error

		defer func() {
			var reqparam []byte
			if res.Request.Body != nil {
				var reqbody, err = res.Request.GetBody()
				if err != nil {
					reqparam, _ = io.ReadAll(reqbody)
				}
			}
			if err == nil {
				log.Debugf("[http][local][ResponseDecoder] %v req: %v, params: %s, resp: %s", res.Request.Method, res.Request.URL.Path, reqparam, data)
			} else {
				log.Warnf("[http][local][ResponseDecoder] %v req: %v, params: %s, resp: %s err:%v", res.Request.Method, res.Request.URL.Path, reqparam, data, err)
			}
			//vals := make(map[interface{}]interface{})
			//c.printContextInternals(ctx, true, vals)
			//if len(vals) > 0 {
			//	var reqs string
			//	for _, req := range vals {
			//		reqs += fmt.Sprintf("%+v;", req)
			//	}
			//	c.logger.Debugf("[http][local][response] req: %v params: %v; resp: %s", reqs, "", data)
			//} else {
			//	c.logger.Debugf("[http][local][response] req:%v params: %v; resp: %s", "", "", data)
			//}
		}()

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		rlt := &jsonResult{}
		err = json.Unmarshal(data, rlt)
		if err != nil {
			return err
		}
		if rlt.Code != 0 {
			return fmt.Errorf("sdk code:%v,msg:%v", rlt.Code, rlt.Message)
		}

		err = json.Unmarshal(data, out)
		return err
	}
}

//func (c *Client) printContextInternals(ctx context.Context, inner bool, vals map[interface{}]interface{}) {
//	if reflect.TypeOf(ctx).String() == "context.backgroundCtx" {
//		return
//	}
//
//	contextValues := reflect.ValueOf(ctx).Elem()
//	contextKeys := reflect.TypeOf(ctx).Elem()
//
//	if !inner {
//		//fmt.Printf("\nFields for %s.%s\n", contextKeys.PkgPath(), contextKeys.Name())
//	}
//
//	//var keys []interface{}
//	if contextKeys.Kind() == reflect.Struct {
//		for i := 0; i < contextValues.NumField(); i++ {
//			reflectValue := contextValues.Field(i)
//			reflectValue = reflect.NewAt(reflectValue.Type(), unsafe.Pointer(reflectValue.UnsafeAddr())).Elem()
//
//			reflectField := contextKeys.Field(i)
//
//			if reflectField.Name == "Context" {
//				c.printContextInternals(reflectValue.Interface().(context.Context), inner, vals)
//			} else {
//				if reflectField.Name == "key" {
//					//vals = append(vals, ctx.Value(reflectValue.Interface()))
//					vals[reflectValue.Interface()] = ctx.Value(reflectValue.Interface())
//					//keys = append(keys, reflectValue.Interface())
//				}
//			}
//		}
//	}
//
//	//for _, key := range keys {
//	//	fmt.Printf("ctx.%s = %+v\n", key, ctx.Value(key))
//	//}
//}
