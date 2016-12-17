package qim


const (
	//用户就医到位（用户扫描医生的二维码,将会把用户消息发送到医生端）
	CONTENT_TYPE_USER_SIGN  = 20001
	//抢处方
	CONTENT_TYPE_PRESCRIPTION_GRAB  = 20002
	//检查到位
	CONTENT_TYPE_CHECK_SIGN = 20003
	//排队通知当前排队用户
	CONTENT_TYPE_QUEUEUP_NOTIFY_CURRENT = 20004
)
