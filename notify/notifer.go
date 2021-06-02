package notify

type Notifer interface {
	// 通知
	Notify(msgs []string)
}
