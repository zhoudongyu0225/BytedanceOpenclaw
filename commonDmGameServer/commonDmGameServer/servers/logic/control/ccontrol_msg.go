package control

import (
	"dmGameServer/model"
	pb "dmGameServer/pb"
	"dmGameServer/zlog"
	"github.com/golang/protobuf/proto"
)

// 获取排行榜
func GetRank(c *WebsocketAnchorClient, weMessage *WbMessageModel) {
	c2s := &pb.RankReq{}
	s2c := &pb.RankResp{}
	err := proto.Unmarshal(weMessage.Body, c2s)
	if err != nil {
		zlog.Logger.Error().Msgf("Unmarshal err:%v [%v]", err, c2s)
		// ws通知
		s2c.ErrMsg = "排行榜协议错误"
		c.WsSend(COMMAND_RANK_S2C, s2c)
		return
	}
	s2c.RankType = c2s.RankType
	rankList := make([]*model.RankVo, 0)

	switch c2s.RankType {
	case pb.RankType_PlayerRankWeek: // 周榜
		rankList = model.GetGameWeekRank(c.AnchorV.PlatformId)
	case pb.RankType_PlayerRankMonth: // 月榜
		rankList = model.GetGameMonthRank(c.AnchorV.PlatformId)
	}
	for _, v := range rankList {
		if v == nil {
			continue
		}
		if v.Value <= 0 {
			continue
		}
		pdb, _ := GetOpenVo(v.OpenId)
		if pdb == nil {
			zlog.Logger.Error().Msgf("玩家信息错误 %v ", v.OpenId)
			continue
		}
		s2c.RanInfoList = append(s2c.RanInfoList, &pb.RanInfo{
			Name:     pdb.NickName,
			HeadUrl:  pdb.AvatarUrl,
			Score:    int64(v.Value),
			WinPoint: pdb.WinPoint,
		})
	}
	zlog.Logger.Info().Msgf("获取排行榜 %v  %v", c2s, s2c)
	c.WsSend(COMMAND_RANK_S2C, s2c)
}
