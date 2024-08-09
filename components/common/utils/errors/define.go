package errors

import "errors"

// 体力不足
var UserEnergyNotEnough = errors.New("user energy not enough")

// 任务条件参数错误
var IndicatorConditionParamsErr = errors.New("indicator condition params error")

// 任务条件参数配置错误
var IndicatorConditionConfigErr = errors.New("indicator condition config error")
