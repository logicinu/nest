package result

var codeMap map[int64]string

func init() {
	codeMap = make(map[int64]string, 3)

	codeMap[CODE_UNKNOWN] = "未知错误"
	codeMap[CODE_OK] = "操作成功"
	codeMap[CODE_FAIL] = "操作失败"
}

func getCode(key int64) string {
	val, ok := codeMap[key]
	if ok {
		return val
	}

	return ""
}
