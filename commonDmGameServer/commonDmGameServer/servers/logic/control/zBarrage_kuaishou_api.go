package control

import (
	pb "dmGameServer/pb"
	"dmGameServer/zlog"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"io"
	"net/http"
)

// KuaishouMessage 发送过来的消息结构体
type KuaishouMessage struct {
	Data       *KuaishouMessageData `json:"data"`
	Message_id string               `json:"message_id"`
	Event      string               `json:"event"`
	App_id     string               `json:"app_id"`
	Timestamp  int64                `json:"timestamp"`
}

// KuaishouMessageData 发送过来的消息结构体
type KuaishouMessageData struct {
	Unique_message_id string                `json:"unique_message_id"` // 唯一的消息id, 可用于幂等消费，第三方需要根据unique_message_id做好幂等控制
	Author_open_id    string                `json:"author_open_id"`    // 主播id
	Room_code         string                `json:"room_code"`         // 房间码
	Push_type         string                `json:"push_type"`         // push的数据标识，不同的pushType 可能对应不同的数据结构
	Payload           []*PayloadMessageData `json:"payload"`
}

type PayloadMessageData struct {
	UserInfo KuiShouPlayerInfo `json:"userInfo"` // 用户的信息
	// ---------送礼物的信息
	UniqueNo       string `json:"uniqueNo"`
	GiftId         string `json:"giftId"`
	GiftName       string `json:"giftName"`
	GiftCount      int32  `json:"giftCount"`
	GiftUnitPrice  int32  `json:"giftUnitPrice"`
	GiftTotalPrice int32  `json:"giftTotalPrice"`
	// ---------评论的信息
	Content string `json:"content"`
	// ---------点赞的信息
	Count int32 `json:"count"`
}

// 快手玩家信息
type KuiShouPlayerInfo struct {
	UserId   string `json:"userId"`   // id
	UserName string `json:"userName"` // 昵称
	HeadUrl  string `json:"headUrl"`  // 头像
}

// 快手弹幕处理
func KuaishouBarrageHandle(c *gin.Context) {
	// 解析 JSON 数据
	kuaishouMessage := &KuaishouMessage{}
	// 获取请求体中的原始数据
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		zlog.Logger.Error().Msgf("快手收到消息ReadAll%v", err)
		Fail(c, "", err.Error())
		return
	}
	err = json.Unmarshal(body, &kuaishouMessage)
	if err != nil {
		zlog.Logger.Error().Msgf("快手收到消息Unmarshal%v ", err)
		Fail(c, "", err.Error())
		return
	}
	if kuaishouMessage.Data == nil {
		zlog.Logger.Error().Msgf("快手收到消息Data err%v ", kuaishouMessage.Data)
		ReturnKuaiShou(c, kuaishouMessage.Message_id)
		return
	}
	messageList := &pb.MessageList{}
	//zlog.Logger.Debug().Msgf("快手收到消息Message[%+v] Data[%+v]", kuaishouMessage, kuaishouMessage.Data)

	switch kuaishouMessage.Data.Push_type {
	case "giftSend": // 礼物
		messageList.Type = "GiftMessage"
		zlog.Logger.Info().Msgf("快手礼物消息Message[%+v] Data[%v]", kuaishouMessage, string(body))
		for _, p := range kuaishouMessage.Data.Payload {
			zlog.Logger.Info().Msgf("快手礼物玩家[%+v] GiftId[%v] GiftTotalPrice[%v]", p.UserInfo.UserName, p.GiftId, p.GiftTotalPrice)
			message := &pb.Message{
				Type:      messageList.Type,
				Name:      p.UserInfo.UserName,
				Uid:       p.UserInfo.UserId,
				HeadImg:   p.UserInfo.HeadUrl,
				Content:   "",
				Count:     "",
				Total:     fmt.Sprintf("%v", p.GiftTotalPrice),
				GiftId:    p.GiftId,
				GiftCount: fmt.Sprintf("%v", p.GiftCount),
				GiftName:  p.GiftName,
			}
			messageList.MsgList = append(messageList.MsgList, message)
		}
	case "liveComment": // 评论
		messageList.Type = "ChatMessage"
		for _, p := range kuaishouMessage.Data.Payload {
			message := &pb.Message{
				Type:    messageList.Type,
				Name:    p.UserInfo.UserName,
				Uid:     p.UserInfo.UserId,
				HeadImg: p.UserInfo.HeadUrl,
				Content: p.Content,
				Count:   "",
			}
			messageList.MsgList = append(messageList.MsgList, message)
		}
		zlog.Logger.Info().Msgf("快手评论消息Message[%+v] Data[%v]", kuaishouMessage, string(body))

	case "liveLike": // 点赞
		messageList.Type = "LikeMessage"
		for _, p := range kuaishouMessage.Data.Payload {
			message := &pb.Message{
				Type:    messageList.Type,
				Name:    p.UserInfo.UserName,
				Uid:     p.UserInfo.UserId,
				HeadImg: p.UserInfo.HeadUrl,
				Content: "",
				Count:   fmt.Sprintf("%v", p.Count),
			}
			messageList.MsgList = append(messageList.MsgList, message)
		}
		zlog.Logger.Info().Msgf("快手点赞消息Message[%+v] Data[%v]", kuaishouMessage, string(body))

	case "follow": // 关注
		zlog.Logger.Info().Msgf("快手关注消息 kuaishouMessage[%+v]", kuaishouMessage)
	default:
		zlog.Logger.Error().Msgf("Push_type err%v ", kuaishouMessage.Data.Push_type)
		ReturnKuaiShou(c, kuaishouMessage.Message_id)
		return
	}
	messageList.RoomId = kuaishouMessage.Data.Room_code

	zlog.Logger.Info().Msgf("[测试]快手弹幕服发过去的消息  Type[%v] RoomId[%v] AccountId[%v]  Author_open_id[%v]", messageList.Type, messageList.RoomId, messageList.AccountId, kuaishouMessage.Data.Author_open_id)

	if !GetMgr().HttpSendMessageToServerLocal(messageList) {
		zlog.Logger.Error().Msgf("HttpSendMessage err:%v", kuaishouMessage.Data.Room_code)
		ReturnKuaiShou(c, kuaishouMessage.Message_id)
		return
	}

	ReturnKuaiShou(c, kuaishouMessage.Message_id)
}

// 返回给快手 这个消息不用再发了
func ReturnKuaiShou(c *gin.Context, msgId string) {
	data := gin.H{
		"result":     1,     // 必填。 1-成功，其他-失败。失败小程序平台会尝试重推此消息
		"message_id": msgId, // 当前消息的message_id
	}
	c.JSON(http.StatusOK, data)
}
