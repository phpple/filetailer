package notify

import "github.com/mitchellh/mapstructure"

// 构建通知器
func BuildNotifer(name string, option map[string]interface{}) Notifer  {
	switch name {
	case "dingtalk":
		var dtOption DingtalkOption
		mapstructure.Decode(option,&dtOption)
		notifer := NewDingtalkNotier(dtOption)
		return &notifer
	default:
		notifer := LoggerNotifer{}
		return &notifer
	}
}
