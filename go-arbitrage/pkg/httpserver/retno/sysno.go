package retno

/**
此文件仅写系统错误码,不允许写入业务相关返回码
0 ~ 9999 为系统错误码，勿用
10000 ~ 为业务码
*/
const (
	//系统
	SYS_ERROR           int64 = 0 //又调皮了
	SYS_OK                    = 1 //1 正确返回 非1 错误返回
	SYS_PATH_INVALID          = 2 //非法请求路径
	SYS_METHOD_NOT_FIND       = 3 //请求接口不存在
)
