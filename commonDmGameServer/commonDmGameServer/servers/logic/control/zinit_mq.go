package control

import (
	"dmGameServer/model"
	pb "dmGameServer/pb"
	"dmGameServer/untils"
	"dmGameServer/zlog"
	"fmt"
	"runtime/debug"
)

// 接受弹幕服来的消息 (注意这个不在主播的协程上)
func OnMqBarrage(messageList *pb.MessageList) bool {
	defer func() {
		if err := recover(); err != nil {
			stack := debug.Stack()
			untils.PanicPoss1(err, stack, "接受弹幕服来的消息")
		}
	}()
	accountId := ""
	if messageList.AccountId != "" {
		accountId = messageList.AccountId
		if GetAnchorClient(accountId) == nil {
			return false
		}
	}
	if messageList.RoomId != "" {
		ac := GetAnchorClientByRoomId(messageList.RoomId)

		if ac == nil {
			return false
		}
		accountId = ac.AccountId
		messageList.AccountId = accountId
	}

	if _, ok := mgr.AnchorMap.Load(accountId); ok {
		// 活跃玩家
		activePlayerNum := int32(0)
		// 新玩家数量
		newPlayerNum := int32(0)
		// 主播数据下活跃
		zbActivePlayerNum := int32(0)
		// 主播的新玩家
		zbNewPlayerNum := int32(0)
		// 主播数据
		anchor := GetAnchorClient(messageList.AccountId).Anchor
		for _, message := range messageList.MsgList {
			if message.Uid == "" {
				untils.TapErr(fmt.Sprintf("平台发的 uid is nil %v", message))
				continue
			}
			// 第一次也会获取
			pdb, isNew := GetOpenVo(message.Uid)
			if pdb == nil {
				untils.TapErr(fmt.Sprintf("pdb is nil %v", message))
				continue
			}

			if message.Uid != "" {
				pdb.OpenId = message.Uid
			}
			if message.Name != "" {
				pdb.NickName = message.Name
			}
			if message.HeadImg != "" {
				pdb.AvatarUrl = message.HeadImg
			}
			if isNew {
				newPlayerNum++
			}
			// 是否有当天的活跃
			if !model.IsHaveDayPlayerInfo(message.Uid) {
				activePlayerNum++
			}
			// 当日主播下活跃的玩家
			if !model.IsAnchorActivePlayer(message.Uid, messageList.AccountId) {
				zbActivePlayerNum++
			}
			// 是否主播第一次玩
			if !model.IsAnchorPlayer(message.Uid, messageList.AccountId) {
				zbNewPlayerNum++
			}
		}
		// 新增玩家数
		if newPlayerNum > 0 {
			// 礼物值添加
			addOverviewOfTransactions := &model.OverviewOfTransactions{
				NewPlayerNum: newPlayerNum,
			}
			// 更新流水总览
			model.UpdateOverviewOfTransactions(addOverviewOfTransactions, GetCollectionOverviewOfTransactionsKey(messageList.AccountId))
		}
		// 开始统计活跃玩家
		if activePlayerNum > 0 {
			// 开始统计活跃玩家
			addOverviewOfTransactions := &model.OverviewOfTransactions{
				ActivePlayerNum: activePlayerNum,
			}
			// 更新流水总览
			model.UpdateOverviewOfTransactions(addOverviewOfTransactions, GetCollectionOverviewOfTransactionsKey(messageList.AccountId))
		}
		// 当日主播下的活跃玩家数
		if zbActivePlayerNum > 0 {
			if GetAnchorClient(messageList.AccountId) != nil && GetAnchorClient(messageList.AccountId).Anchor != nil {
				// 当日活跃玩家数
				anchorOfTransactions := &model.AnchorOfTransactions{
					AccountId: messageList.AccountId,
					Name:      GetAnchorClient(messageList.AccountId).Anchor.NickName,
				}
				model.UpdateAnchorOfTransactionsActivePlayerNum(anchorOfTransactions, zbActivePlayerNum, GetCollectionAnchorOfTransactionsKey(messageList.AccountId))
			}
		}
		//  当日新增玩家数
		if zbNewPlayerNum > 0 {
			if GetAnchorClient(messageList.AccountId) != nil && GetAnchorClient(messageList.AccountId).Anchor != nil {
				// 当日活跃玩家数
				anchorOfTransactions := &model.AnchorOfTransactions{
					AccountId: messageList.AccountId,
					Name:      GetAnchorClient(messageList.AccountId).Anchor.NickName,
				}
				model.UpdateAnchorOfTransactionsNewPlayerNum(anchorOfTransactions, zbNewPlayerNum, GetCollectionAnchorOfTransactionsKey(messageList.AccountId))
			}
		}
		// 新玩家---等等逻辑---
		// 统计礼物值
		if messageList.Type == "GiftMessage" {
			// ------------------统计
			updateOpenVoList := make([]*pb.OpenVo, 0)
			// 礼物总价值
			allGiftValue := int64(0)
			// 会有很多条
			for _, message := range messageList.MsgList {
				if message.Uid == "" {
					untils.TapErr(fmt.Sprintf("平台发的 uid is nil %v", message))
					continue
				}
				// 第一次也会获取
				pdb, _ := GetOpenVo(message.Uid)
				total := StirngToInt64(message.Total)
				// 礼物值和礼物价值
				pdb.GiftValue += total

				allGiftValue += total
				if anchor != nil {
					anchor.AllGiftValue += float64(allGiftValue)
				}
				// 更新玩家的流水
				playerDetailOfTransactions := &model.PlayerDetailOfTransactions{
					Id:   pdb.OpenId,
					Name: pdb.NickName,
				}
				updateOpenVoList = append(updateOpenVoList, pdb)
				// 1.更新玩家的流水
				model.UpdatePlayerDetailOfTransactions(playerDetailOfTransactions, pdb.GiftValue, total, GetCollectionPlayerOfTransactionsKey(messageList.AccountId))
				// 2.玩家下的礼物值的排行榜
				AddAnchorPlayerRank(messageList.AccountId, message.Uid, float64(pdb.GiftValue))
			}
			// 礼物值添加
			addOverviewOfTransactions := &model.OverviewOfTransactions{
				GiftValue: float64(allGiftValue),
			}
			// 3.更新流水总览
			model.UpdateOverviewOfTransactions(addOverviewOfTransactions, GetCollectionOverviewOfTransactionsKey(messageList.AccountId))

			if anchor != nil {
				// 主播添加收入排行榜
				model.UpdateAnchorRankList(accountId, allGiftValue, anchor.PlatformId, anchor.GameId)
				// 更新更新主播流水礼物值
				anchorOfTransactions := &model.AnchorOfTransactions{
					AccountId: messageList.AccountId,
					Name:      anchor.NickName,
				}
				// 4.更新更新主播流水礼物值
				model.UpdateAnchorOfTransactionsGiftValue(anchorOfTransactions, anchor.AllGiftValue, float64(allGiftValue), GetCollectionAnchorOfTransactionsKey(messageList.AccountId))
				// 5.更新主播信息
				model.UpdateAnchor(anchor)
			}
			// 3.更新礼物值
			model.UpdateOpenVo(updateOpenVoList)
			ac := GetAnchorClient(messageList.AccountId)

			if ac != nil && ac.WebsocketClient != nil {
				GetAnchorClient(messageList.AccountId).WebsocketClient.giftV += allGiftValue
			}
		}
		GetLogicMgr().BarrageSendMessage(accountId, messageList)
	} else {
		zlog.Logger.Info().Msgf("[测试]接受弹幕服来的消息  Type[%v] RoomId[%v] AccountId[%v] -----不存在----", messageList.Type, messageList.RoomId, messageList.AccountId)
	}
	return true
}
