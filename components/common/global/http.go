package global

import (
	"components/common/proto"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"net/http"
)

var (
	logger *log.Helper
)

type jsonResult struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Init(loge log.Logger) {
	logger = log.NewHelper(loge)
}

func ResponseEncoder(w http.ResponseWriter, r *http.Request, out interface{}) error {
	if out == nil {
		return nil
	}

	reqbody, err := io.ReadAll(r.Body)
	if err != nil {
		return errors.BadRequest("CODEC", err.Error())
	}

	respbody, err := json.Marshal(out)
	if err != nil {
		return err
	}

	logger.Debugf("[http][remote][ResponseEncoder] %v %v req: %v, params: %v, resp: %v", r.RemoteAddr, r.Method, r.URL.Path, string(reqbody), string(respbody))

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(respbody)
	if err != nil {
		return err
	}
	return nil
}

func ErrorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	reqbody, _ := io.ReadAll(r.Body)

	code, msg := fromError(err)
	se := &jsonResult{
		Code:    code,
		Message: msg,
	}

	body, err := json.Marshal(se)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.Warnf("[http][remote][ErrorEncoder] %v %v req: %v, params: %v, resp: %v", r.RemoteAddr, r.Method, r.URL.Path, string(reqbody), string(body))

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(body)
}

func ToError(code int32, msg string) error {
	errorStatus := status.New(codes.Unknown, fmt.Sprintf("code:%v,msg:%v", code, msg))
	//附加业务错误
	ds, _ := errorStatus.WithDetails(
		&proto.ErrorWithDetails{
			Code:    code,
			Message: msg,
		},
	)
	return ds.Err()
}

func fromError(err error) (code int32, msg string) {
	errStatus, ok := status.FromError(err)
	if !ok {
		return int32(codes.Unknown), err.Error()
	}

	if errStatus.Code() == codes.Unknown {
		//解析业务错误
		convertErr := status.Convert(err)
		for _, d := range convertErr.Details() {
			switch info := d.(type) {
			case *proto.ErrorWithDetails:
				return info.GetCode(), info.GetMessage()
			default:
				return int32(errStatus.Code()), errStatus.Message()
			}
		}
	}
	return int32(errStatus.Code()), errStatus.Message()
}
