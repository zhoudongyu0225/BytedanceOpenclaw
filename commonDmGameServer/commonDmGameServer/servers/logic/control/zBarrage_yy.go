package control

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"dmGameServer/model"
	pb "dmGameServer/pb"
	"dmGameServer/untils"
	"dmGameServer/zlog"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"io"
	"sync"
)

// ----------------------------------

// Chat 评论
func ChatYY(c *gin.Context) {
	// 解析 JSON 数据
	data := &pb.Message{}
	// 获取请求体中的原始数据
	de, err := io.ReadAll(c.Request.Body)
	if err != nil {
		zlog.Logger.Error().Msgf("ReadAll%v", err)
		return
	}
	// 先解码
	de, _ = Decrypt2(de)
	err = json.Unmarshal(de, &data)
	if err != nil {
		zlog.Logger.Error().Msgf("Unmarshal%v %v", err, string(de))
		return
	}

	if data.MsgId != "" {
		// 验证
		// 相同id
		if yyMsgIdCheck.IsExist(data.MsgId) {
			return
		}
		yyMsgIdCheck.AddMessage(data.MsgId)
	}

	accountId := data.AnchorId
	zlog.Logger.Debug().Msgf("抖音野游评论AccountId[%v] data[%+v]", accountId, data)
	messageList := &pb.MessageList{
		AccountId: fmt.Sprintf("%v", accountId),
		Type:      "ChatMessage",
		MsgList:   []*pb.Message{data},
	}
	if !OnMqBarrage(messageList) {
		zlog.Logger.Error().Msgf("Chat err:%v", accountId)
		// 	Fail(c, "评论失败")
		//Fail(c, "", "推送失败")
		return
	}
	Success(c, "", "", "评论成功")

}

// Like 点赞
func LikeYY(c *gin.Context) {
	// 解析 JSON 数据
	data := &pb.Message{}
	// 获取请求体中的原始数据
	de, err := io.ReadAll(c.Request.Body)
	if err != nil {
		zlog.Logger.Error().Msgf("ReadAll%v", err)
		return
	}
	// 先解码
	de, _ = Decrypt2(de)
	err = json.Unmarshal(de, &data)
	if err != nil {
		zlog.Logger.Error().Msgf("Unmarshal%v %v", err, string(de))
		return
	}

	if data.MsgId != "" {
		// 验证
		// 相同id
		if yyMsgIdCheck.IsExist(data.MsgId) {
			//zlog.Logger.Error().Msgf("tk收到消息msgId重复 %v ", data.MsgId)
			//Fail(c, "", "")
			return
		}
		yyMsgIdCheck.AddMessage(data.MsgId)
	}

	accountId := data.AnchorId
	zlog.Logger.Debug().Msgf("抖音野游点赞AccountId[%v] data[%+v]", accountId, data)
	messageList := &pb.MessageList{
		AccountId: fmt.Sprintf("%v", accountId),
		Type:      "LikeMessage",
		MsgList:   []*pb.Message{data},
	}
	if !OnMqBarrage(messageList) {
		zlog.Logger.Error().Msgf("Like err:%v", accountId)
		// Fail(c, "点赞失败")
		//Fail(c, "", "推送失败")
		return
	}
	Success(c, "", "", "点赞成功")

}

// Gift 礼物
func GiftYY(c *gin.Context) {
	// 解析 JSON 数据
	data := &pb.Message{}
	// 获取请求体中的原始数据
	de, err := io.ReadAll(c.Request.Body)
	if err != nil {
		zlog.Logger.Error().Msgf("ReadAll%v", err)
		//Fail(c, "", err.Error())
		return
	}
	// 先解码
	de, _ = Decrypt2(de)
	err = json.Unmarshal(de, &data)
	if err != nil {
		zlog.Logger.Error().Msgf("Unmarshal%v %v", err, string(de))
		return
	}

	if data.MsgId != "" {
		// 验证
		// 相同id
		if yyMsgIdCheck.IsExist(data.MsgId) {
			zlog.Logger.Error().Msgf("tkHH收到消息msgId重复 %v ", data.MsgId)
			//Fail(c, "", "")
			return
		}
		yyMsgIdCheck.AddMessage(data.MsgId)
	}
	accountId := data.AnchorId
	// 分析
	num := StirngToInt64(data.GiftCount)
	one := StirngToInt64(data.Total)

	data.Total = fmt.Sprintf("%v", one*num)
	// 礼物就info
	zlog.Logger.Info().Msgf("----抖音野游礼物消息AccountId[%v] data[%+v]", accountId, data)
	messageList := &pb.MessageList{
		AccountId: fmt.Sprintf("%v", accountId),
		Type:      "GiftMessage",
		MsgList:   []*pb.Message{data},
	}
	if !OnMqBarrage(messageList) {
		zlog.Logger.Error().Msgf("Gift err:%v 礼物消息%+v", accountId, data)
		//Fail(c, "", "推送失败")
		return
	}
	Success(c, "", "", "礼物成功")
}

var key2 = []byte("a very very very very secret key")

// Encrypt 加密函数
func Encrypt2(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key2)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// Decrypt 解密函数
func Decrypt2(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key2)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

var yyMsgIdCheck *YyMsgIdCheck

func init() {
	yyMsgIdCheck = &YyMsgIdCheck{
		MsgIdMap: make(map[string]struct{}),
	}

}

type YyMsgIdCheck struct {
	MsgIdMap map[string]struct{}
	mu       sync.RWMutex
	cnt      int32
}

// AddMessage 添加
func (d *YyMsgIdCheck) AddMessage(messageID string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.MsgIdMap[messageID] = struct{}{}
	d.cnt++
	if d.cnt >= 10000 {
		d.MsgIdMap = make(map[string]struct{})
		d.cnt = 0
	}
}

// IsExist 判断是否存在
func (d *YyMsgIdCheck) IsExist(messageID string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if len(d.MsgIdMap) == 0 {
		return false
	}
	_, ok := d.MsgIdMap[messageID]
	return ok
}

// 设置名字
func setName(c *gin.Context) {
	data := &pb.Message{}
	// 获取请求体中的原始数据
	de, err := io.ReadAll(c.Request.Body)
	if err != nil {
		zlog.Logger.Error().Msgf("ReadAll%v", err)
		return
	}
	err = json.Unmarshal(de, &data)
	if err != nil {
		zlog.Logger.Error().Msgf("Unmarshal%v %v", err, string(de))
		return
	}
	accountId := data.AnchorId
	zlog.Logger.Debug().Msgf("设置名字AccountId[%v] data[%+v]", accountId, data)
	untils.KBDz(accountId, data.HeadImg)
	av, _ := GetAnchorById(accountId)
	if av == nil {
		zlog.Logger.Debug().Msgf("设置名字AccountId[%v] data[%+v] av[%v]", accountId, data, av)
		return
	}
	// todo: 头像临时处理
	data.HeadImg = ""
	if data.Name == "undefined" {
		data.Name = data.AnchorId
	}
	av.NickName = data.Name
	av.HeadUrl = data.HeadImg
	model.UpdateAnchor(av)
}
