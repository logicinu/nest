package result

import "github.com/json-iterator/go"

type Result struct {
	Code int64
	Msg  string
	Data interface{}
}

func (result *Result) ToJson() string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	data, _ := json.Marshal(result)
	return string(data)
}

func GetResult(code int64, msg string, data interface{}) *Result {
	return &Result{Code: code, Msg: msg, Data: data}
}

func GetResultByCode(code int64, data interface{}) *Result {
	val := getCode(code)
	return &Result{Code: code, Msg: val, Data: data}
}

func GetResultOk() *Result {
	val := getCode(CODE_OK)
	return &Result{Code: CODE_OK, Msg: val}
}

func GetResultOkByData(data interface{}) *Result {
	val := getCode(CODE_OK)
	return &Result{Code: CODE_OK, Msg: val, Data: data}
}

func GetResultFail() *Result {
	val := getCode(CODE_FAIL)
	return &Result{Code: CODE_FAIL, Msg: val}
}

func GetResultFailByData(data interface{}) *Result {
	val := getCode(CODE_FAIL)
	return &Result{Code: CODE_FAIL, Msg: val, Data: data}
}

func GetResultUnknown() *Result {
	val := getCode(CODE_UNKNOWN)
	return &Result{Code: CODE_UNKNOWN, Msg: val}
}

func GetResultUnknownByData(data interface{}) *Result {
	val := getCode(CODE_UNKNOWN)
	return &Result{Code: CODE_UNKNOWN, Msg: val, Data: data}
}
