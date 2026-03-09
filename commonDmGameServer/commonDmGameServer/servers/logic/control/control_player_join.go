package control

import (
	pb "dmGameServer/pb"
	"dmGameServer/zlog"
)

// 玩家加入  (就加入游戏) 需要判断是否新玩家
func PlayerJoin(c *WebsocketAnchorClient, OpenVos []*pb.OpenBaseVo) {
	// BaseBattleJoin 加入的基础 (主要做 记录房间 通知离开房间这些数据)
	pbList, contList := BaseBattleJoin(c, OpenVos)
	// 加入有数据变化的
	var updatePbList []*pb.OpenVo
	for i, pdb := range pbList {
		zlog.Logger.Debug().Msgf("玩家加入请求AccountId[%v]   cont:%v  OpenId:%v", c.AccountId, contList[i], pdb.OpenId)
		updatePbList = append(updatePbList, pdb)
	}
	// 观众加入直播间响应
	roomJoinResp := &pb.PlayerJoinNotify{
		OpenList: pbList,
	}
	zlog.Logger.Debug().Msgf("玩家加入请求AccountId[%v] 玩家pdb[%+v]", c.AccountId, roomJoinResp)
	// 回复玩家加入
	c.WsSend(COMMAND_PLAYER_JOIN_S, roomJoinResp)
}

// 加入的基础
func BaseBattleJoin(c *WebsocketAnchorClient, OpenVos []*pb.OpenBaseVo) (PbList []*pb.OpenVo, contList []string) {
	zlog.Logger.Debug().Msgf("玩家加入请求AccountId[%v] 玩家uid[%+v]", c.AccountId, OpenVos)
	for _, v := range OpenVos {
		pdb, _ := GetOpenVo(v.OpenId)
		if pdb == nil {
			zlog.Logger.Error().Msgf("GetOpenVo err [%v] [%v]", v.OpenId, c.GameId)
			continue
		}
		if c.AnchorV == nil {
			zlog.Logger.Error().Msgf("c.anchorV == nil [%v]", c.AccountId)
			continue
		}
		// 加入玩家的数据
		PbList = append(PbList, pdb)
		contList = append(contList, v.Content)
	}
	return PbList, contList
}
