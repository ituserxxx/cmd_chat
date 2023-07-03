package comm

var (
	EventSysInfo   = "Sys info"   //系统消息
	EventGuangbo   = "GB"         //广播消息
	EventPublicMsg = "public msg" // 公开消息
	EventInitName  = "init name"  //初始化名称

	EventInputAllUsers = "all user" // 查看所有用户
	EventInputAT       = "@"        //私聊

)

type MsgInfo struct {
	Event string `json:"event"`
	Data  string `json:"data"`
	Code  int    `json:"code"`
}
