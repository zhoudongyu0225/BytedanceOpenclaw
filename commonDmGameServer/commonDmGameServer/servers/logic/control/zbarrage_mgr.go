package control

import (
	pb "dmGameServer/pb"
	"dmGameServer/zlog"
	"sync"
)

// 弹幕逻辑管理器
var barrageMgr *ControlMgr

// ControlMgr 逻辑管理器
type ControlMgr struct {
	yeyouTokenMap sync.Map // 管理野游token的问题 key:id value token
}

// GetMgr 获取逻辑管理器
func GetMgr() *ControlMgr {
	return barrageMgr
}

func InitBarrageMgr() {
	barrageMgr = &ControlMgr{}
}

// SetYeyouToken 设置
func (m *ControlMgr) SetYeyouToken(accountId, token string) {
	barrageMgr.yeyouTokenMap.Store(accountId, token)
}

func (m *ControlMgr) IsYeyouToken(accountId, token string) bool {
	t, ok := barrageMgr.yeyouTokenMap.Load(accountId)
	if !ok {
		zlog.Logger.Error().Msgf("野游token不存在 %v", accountId)
		return false
	}
	_token := t.(string)
	if _token != token {
		return false
	}
	return true
}

// 本服直接发送消息
func (m *ControlMgr) HttpSendMessageToServerLocal(messageList *pb.MessageList) bool {
	//  接受弹幕服来的消息
	OnMqBarrage(messageList)
	return true
}
