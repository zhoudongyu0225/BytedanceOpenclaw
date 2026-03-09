package control

import (
	"dmGameServer/model"
	pb "dmGameServer/pb"
	"dmGameServer/untils"
	"dmGameServer/zlog"
	"fmt"
	"strconv"
	"time"
)

// slg评论
func ChatMessage(messageList *pb.MessageList, ws *WebsocketAnchorClient) {
	now := time.Now().Unix()
	msgJoin(messageList, ws)
	for _, msg := range messageList.MsgList {
		zlog.Logger.Info().Msgf("玩家评论   %v  %v %v %v ", msg.Uid, msg.Name, msg.Content, ws.AccountId)
		if msg.Uid == "" {
			untils.TapErr(fmt.Sprintf("msg.Uid is 空 %v", msg))
			continue
		}
		// 在自己的直播间且主播在
		pdb, _ := GetOpenVo(msg.Uid)
		if pdb == nil {
			untils.TapErr(fmt.Sprintf("pdb is nil %v", msg))
			continue
		}
		if msg.Uid != "" {
			pdb.OpenId = msg.Uid
			pdb.OpenId = msg.Uid
		}
		if msg.Name != "" {
			pdb.NickName = msg.Name
		}
		if msg.HeadImg != "" {
			pdb.AvatarUrl = msg.HeadImg
		}
		if now-pdb.LastUpdateTime > 5 {
			currWeekRank, currWeekScore := model.GetGameWeekRankByRankVo(pdb.OpenId, ws.AnchorV.PlatformId)
			currMonthRank, currMonthScore := model.GetGameMonthRankByRankVo(pdb.OpenId, ws.AnchorV.PlatformId)
			pdb.Score = currWeekScore
			pdb.MonthScore = currMonthScore
			pdb.Rank = int32(currWeekRank)
			pdb.MonthRank = int32(currMonthRank)
			pdb.LastUpdateTime = now
		}
		// todo:弹幕逻辑

	}

}

// slg礼物
func GiftMessage(messageList *pb.MessageList, ws *WebsocketAnchorClient) {
	now := time.Now().Unix()
	msgJoin(messageList, ws)
	//now := time.Now().Unix()
	// 加入的玩家
	// 加入的玩家
	for _, msg := range messageList.MsgList {
		if msg.Uid == "" {
			untils.TapErr(fmt.Sprintf("平台发的 uid is nil %v", msg))
			continue
		}
		// 在自己的直播间且主播在
		pdb, _ := GetOpenVo(msg.Uid)
		if pdb == nil {
			untils.TapErr(fmt.Sprintf("pdb is nil %v", msg))
			continue
		}
		if msg.Uid != "" {
			pdb.OpenId = msg.Uid
		}
		if msg.Name != "" {
			pdb.NickName = msg.Name
		}
		if msg.HeadImg != "" {
			pdb.AvatarUrl = msg.HeadImg
		}

		currWeekRank, currWeekScore := model.GetGameWeekRankByRankVo(pdb.OpenId, ws.AnchorV.PlatformId)
		currMonthRank, currMonthScore := model.GetGameMonthRankByRankVo(pdb.OpenId, ws.AnchorV.PlatformId)
		pdb.Score = currWeekScore
		pdb.MonthScore = currMonthScore
		pdb.Rank = int32(currWeekRank)
		pdb.MonthRank = int32(currMonthRank)
		pdb.LastUpdateTime = now

		// todo:礼物逻辑
		// 单个礼物
	}
}
func StirngToInt32(str string) int32 {
	if str == "" {
		return 0
	}
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return int32(i)
}

// slg点赞
func LikeMessage(messageList *pb.MessageList, ws *WebsocketAnchorClient) {
	now := time.Now().Unix()
	msgJoin(messageList, ws)
	// now := time.Now().Unix()
	// 加入的玩家
	// 加入的玩家
	for _, msg := range messageList.MsgList {
		if msg.Uid == "" {
			untils.TapErr(fmt.Sprintf("平台发的 uid is nil %v", msg))
			continue
		}
		// 在自己的直播间且主播在
		pdb, _ := GetOpenVo(msg.Uid)
		if pdb == nil {
			untils.TapErr(fmt.Sprintf("pdb is nil %v", msg))
			continue
		}
		//--
		if now-pdb.LastUpdateTime > 60 {
			currWeekRank, currWeekScore := model.GetGameWeekRankByRankVo(pdb.OpenId, ws.AnchorV.PlatformId)
			currMonthRank, currMonthScore := model.GetGameMonthRankByRankVo(pdb.OpenId, ws.AnchorV.PlatformId)
			pdb.Score = currWeekScore
			pdb.MonthScore = currMonthScore
			pdb.Rank = int32(currWeekRank)
			pdb.MonthRank = int32(currMonthRank)
			pdb.LastUpdateTime = now
		}

		// todo:点赞逻辑
	}
}

// 检查玩家的信息
func CheckOpenVoName(pdb *pb.OpenVo, msg *pb.Message) {
	if pdb == nil {
		return
	}
	if msg.Uid != "" {
		pdb.OpenId = msg.Uid
	}
	if msg.Name != "" {
		pdb.NickName = msg.Name
	}
	if msg.HeadImg != "" {
		pdb.AvatarUrl = msg.HeadImg
	}
}

// 消息加入 用于礼物和点赞 和任意评论
func msgJoin(messageList *pb.MessageList, ws *WebsocketAnchorClient) {
	// 加入的玩家
	var JisonOpenVos []*pb.OpenBaseVo
	for _, msg := range messageList.MsgList {
		if msg.Uid == "" {
			untils.TapErr(fmt.Sprintf("平台发的 uid is nil %v", msg))
			continue
		}
		// 在自己的直播间且主播在
		pdb, _ := GetOpenVo(msg.Uid)
		if pdb == nil {
			untils.TapErr(fmt.Sprintf("pdb is nil %v", msg))
			continue
		}
		// 检查玩家的信息
		CheckOpenVoName(pdb, msg)
		// 玩家获取上一个的主播
		lastAc := GetPlayerAc(msg.Uid)
		if lastAc == nil {
			// 1.玩家没有在主播
			// 加入
			JisonOpenVos = append(JisonOpenVos, &pb.OpenBaseVo{
				OpenId:    msg.Uid,
				AvatarUrl: msg.HeadImg,
				NickName:  msg.Name,
				Content:   msg.Content,
			})
		} else {
			if lastAc.AccountId != ws.AccountId || lastAc.WebsocketClient == nil || !lastAc.WebsocketClient.IsClientLogOn {
				// 加入
				JisonOpenVos = append(JisonOpenVos, &pb.OpenBaseVo{
					OpenId:    msg.Uid,
					AvatarUrl: msg.HeadImg,
					NickName:  msg.Name,
					Content:   msg.Content,
				})
				// 通知其他直播间这个推送离开
				if lastAc.WebsocketClient != nil && lastAc.WebsocketClient.IsClientLogOn {
					// 推送离开
					ts2c := &pb.PlayerLeaveNotify{
						OpenId: msg.Uid,
					}
					zlog.Logger.Info().Msgf("通知玩家这个推送离开 %v %v to  %v", msg.Uid, lastAc.AccountId, ws.AccountId)
					// 删掉玩家的
					lastAc.WebsocketClient.DelPlayer(msg.Uid)
					lastAc.WebsocketClient.WsSend(COMMAND_PLAYER_LEAVE_S, ts2c)
				}
			}
		}
	}
	// 玩家加入
	if len(JisonOpenVos) > 0 {
		// 玩家加入 (带回复 有存档)
		PlayerJoin(ws, JisonOpenVos)
	}
}
