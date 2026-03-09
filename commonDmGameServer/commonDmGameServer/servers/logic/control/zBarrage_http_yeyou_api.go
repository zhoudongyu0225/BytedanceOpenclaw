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

// http之前做数据解压等操作
func before(c *gin.Context) []byte {
	// 获取请求体中的原始数据
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		zlog.Logger.Error().Msgf("ReadAll%v", err)
		Fail(c, "", err.Error())
		return nil
	}
	return body
}

// Chat 评论
func Chat(c *gin.Context) {
	// 解析 JSON 数据
	data := &pb.Message{}
	de := before(c)
	if de == nil {
		return
	}
	err := json.Unmarshal(de, &data)
	if err != nil {
		zlog.Logger.Error().Msgf("Unmarshal%v %v", err, string(de))
		Fail(c, "", err.Error())
		return
	}
	accountId, ok1 := c.Get("AccountId")
	if !ok1 {
		zlog.Logger.Error().Msgf("Chat err:%v", accountId)
		// Response(c, http.StatusUnprocessableEntity, 442, nil, "参数错误")
		Fail(c, "", "参数错误")
		return
	}
	zlog.Logger.Debug().Msgf("野游评论AccountId[%v] data[%+v]", accountId, data)
	messageList := &pb.MessageList{
		AccountId: fmt.Sprintf("%v", accountId),
		Type:      "ChatMessage",
		MsgList:   []*pb.Message{data},
	}
	if !GetMgr().HttpSendMessageToServerLocal(messageList) {
		zlog.Logger.Error().Msgf("Chat err:%v", accountId)
		// 	Fail(c, "评论失败")
		Fail(c, "", "推送失败")
		return
	}
	Success(c, "", "", "评论成功")

}

// Like 点赞
func Like(c *gin.Context) {
	// 解析 JSON 数据
	data := &pb.Message{}
	de := before(c)
	if de == nil {
		return
	}
	err := json.Unmarshal(de, &data)
	if err != nil {
		zlog.Logger.Error().Msgf("Unmarshal%v %v", err, string(de))
		Fail(c, "", err.Error())

		return
	}
	accountId, ok1 := c.Get("AccountId")
	if !ok1 {
		zlog.Logger.Error().Msgf("Chat err:%v", accountId)
		// Response(c, http.StatusUnprocessableEntity, 442, nil, "参数错误")
		Fail(c, "", "参数错误")
		return
	}
	zlog.Logger.Debug().Msgf("野游点赞AccountId[%v] data[%+v]", accountId, data)
	messageList := &pb.MessageList{
		AccountId: fmt.Sprintf("%v", accountId),
		Type:      "LikeMessage",
		MsgList:   []*pb.Message{data},
	}
	if !GetMgr().HttpSendMessageToServerLocal(messageList) {
		zlog.Logger.Error().Msgf("Like err:%v", accountId)
		// Fail(c, "点赞失败")
		Fail(c, "", "推送失败")
		return
	}
	Success(c, "", "", "点赞成功")

}

// Gift 礼物
func Gift(c *gin.Context) {
	// 解析 JSON 数据
	data := &pb.Message{}
	de := before(c)
	if de == nil {
		return
	}
	err := json.Unmarshal(de, &data)
	if err != nil {
		zlog.Logger.Error().Msgf("Unmarshal%v %v", err, string(de))
		Fail(c, "", "Failed to parse JSON data")
		return
	}
	accountId, ok1 := c.Get("AccountId")
	if !ok1 {
		zlog.Logger.Error().Msgf("Chat err:%v", accountId)
		Response(c, http.StatusUnprocessableEntity, 442, nil, "参数错误")
		return
	}
	// 礼物就info
	zlog.Logger.Info().Msgf("野游礼物消息AccountId[%v] data[%+v]", accountId, data)
	messageList := &pb.MessageList{
		AccountId: fmt.Sprintf("%v", accountId),
		Type:      "GiftMessage",
		MsgList:   []*pb.Message{data},
	}
	if !GetMgr().HttpSendMessageToServerLocal(messageList) {
		zlog.Logger.Error().Msgf("Gift err:%v 礼物消息%+v", accountId, data)
		Fail(c, "", "推送失败")
		return
	}
	Success(c, "", "", "礼物成功")
}
