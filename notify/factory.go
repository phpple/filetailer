package notify

// 构建通知器
func BuildNotifer(name string, option map[string]string) Notifer  {
	switch name {
	case "dingtalk":
		notifer := NewDingtalkNotier(option["token"])
		return &notifer
	default:
		notifer := LoggerNotifer{}
		return &notifer
	}
}
