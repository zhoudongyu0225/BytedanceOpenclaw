package control

import (
	pb "dmGameServer/pb"
	"dmGameServer/zlog"
	"github.com/gin-gonic/gin"
	"net/http"
)

var sing = "!123456amin"

type BBB struct {
	AccountId string `json:"AccountId"`
}

func Notice(c *gin.Context) {
	//MsgInfo
	var data struct {
		AccountId  string `json:"AccountId"`
		PlatformId int32  `json:"PlatformId"`
		Notice     string `json:"Notice"`
		Sin        string `json:"Sin"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, gin.H{"content": "参数错误"})
		return
	}
	if data.Sin != sing {
		c.JSON(http.StatusOK, gin.H{"content": "权限不足"})
		return
	}
	zlog.Logger.Info().Msgf("公告消息 %+v", data)
	if data.AccountId == "" {
		// 给全部人发公告
		GetLogicMgr().AnchorMap.Range(func(key, value interface{}) bool {
			ac := value.(*AnchorClient)
			if ac.Anchor.PlatformId != int32(data.PlatformId) {
				return true
			}
			TipsType(ac.WebsocketClient, data.Notice, pb.TipType_Notice)
			return true
		})
	} else {
		ac := GetAnchorClient(data.AccountId)
		if ac == nil {
			zlog.Logger.Error().Msgf("没有找到主播 :%v", data.AccountId)
			return
		}
		TipsType(ac.WebsocketClient, data.Notice, pb.TipType_Notice)
	}
	return
}
