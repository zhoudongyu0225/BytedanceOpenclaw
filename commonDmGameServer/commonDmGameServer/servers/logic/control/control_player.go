package control

import (
	"dmGameServer/model"
	pb "dmGameServer/pb"
)

// 获取玩家
func GetOpenVo(openId string) (*pb.OpenVo, bool) {
	openVo, isNew := model.GetPlayerMgr().ModelGetOpenVo(openId)
	return openVo, isNew
}
