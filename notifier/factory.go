package notifier

import (
	"github.com/mitchellh/mapstructure"
	"log"
)

// 构建通知器
func BuildNotifier(name string, option map[string]interface{}) Notifier {
	switch name {
	case "dingtalk":
		var dtOption DingtalkOption
		mapstructure.Decode(option, &dtOption)
		log.Println("dingtalk.option", dtOption)
		notifer := NewDingtalkNotier(dtOption)
		return &notifer
	default:
		notifer := LoggerNotifier{}
		return &notifer
	}
}
