package control

// CMD指令集合
const (
	// 提示(一般错误)
	COMMAND_TIPS_S int16 = 4
	// 心跳
	COMMAND_PING_C2S int16 = 5
	// 心跳
	COMMAND_PONG_S2C int16 = 6
	// 客户端登录请求
	COMMAND_LOGIN_C2S int16 = 11
	// 客户端登录回复
	COMMAND_LOGIN_S2C int16 = 12
	// 观众加入
	COMMAND_PLAYER_JOIN_S int16 = 14
	// 观众离开
	COMMAND_PLAYER_LEAVE_S int16 = 15
	// 查看排行榜请求
	COMMAND_RANK_C2S int16 = 19
	// 查看排行榜回复
	COMMAND_RANK_S2C int16 = 20
)
