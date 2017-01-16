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

	//服务类消息
	CONTENT_TYPE_SERVICE = 20010
)

//会话类型
const (
	//单聊
	SESSIONTYPE_SINGLE = 1
	//群聊
	SESSIONTYPE_GROUP = 2
)
const (
	//默认（好友，群）
	OPERATION_TYPE_COMM = 0
	//经费群
	OPERATION_TYPE_FUNDSGROUP = 3
	//公众号
	OPERATION_TYPE_PUBLICNO = 4
)