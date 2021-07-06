package e

var MsgFlags = map[int]string{
	SUCCESS : "ok",
	ERROR : "failed",
	INVALID_PARAMS : "请求参数错误",

	ERROR_EXIST_TAG :"已经存在改标签名称",
	ERROR_NOT_EXIST_TAG : "改标签不存在",
	ERROR_NOT_EXIST_ARTICLE : "改文章不存在",

	ERROR_AUTH_CHECK_TOKEN_FAIL : "Token 鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT : "Token 已超时",
	ERROR_AUTH_TOKNE : "Token 生成失败",
	ERROR_AUTH : "Token 错误",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
