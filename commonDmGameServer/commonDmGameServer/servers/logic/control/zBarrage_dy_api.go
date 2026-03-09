package control

import (
	"crypto/md5"
	pb "dmGameServer/pb"
	"dmGameServer/zlog"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
)

var dyMsgIdCheck *DyMsgIdCheck

func init() {
	dyMsgIdCheck = &DyMsgIdCheck{
		MsgIdMap: make(map[string]struct{}),
	}
}

// 抖音重复id校验
type DyMsgIdCheck struct {
	MsgIdMap map[string]struct{}
	mu       sync.RWMutex
	cnt      int32
}

// AddMessage 添加
func (d *DyMsgIdCheck) AddMessage(messageID string) {
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
func (d *DyMsgIdCheck) IsExist(messageID string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if len(d.MsgIdMap) == 0 {
		return false
	}
	_, ok := d.MsgIdMap[messageID]
	return ok
}

type DyMessageList []struct {
	MsgId     string `json:"msg_id"`
	SecOpenid string `json:"sec_openid"`
	Content   string `json:"content"`
	AvatarUrl string `json:"avatar_url"`
	Nickname  string `json:"nickname"`
	Timestamp int    `json:"timestamp"`

	SecGiftId string `json:"sec_gift_id"` // 加密的礼物id
	GiftNum   int    `json:"gift_num"`
	GiftValue int    `json:"gift_value"`

	LikeNum int32 `json:"like_num"`
}

func DyBarrageHandle(c *gin.Context) {
	// 获取特定的请求头字段
	RoomId := c.Request.Header.Get("x-roomid")
	msgType := c.Request.Header.Get("x-msg-type")
	dyMessageList := &DyMessageList{}
	// 获取请求体中的原始数据
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		zlog.Logger.Error().Msgf("抖音收到消息ReadAll%v", err)
		// 返回消息
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "success",
		})
		return
	}

	// 将原始数据反序列化成Message结构体
	err = json.Unmarshal(body, dyMessageList)
	if err != nil {
		zlog.Logger.Error().Msgf("抖音收到消息Unmarshal%v", err, string(body))
		// 返回消息
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "success",
		})
		return
	}
	signature1 := c.Request.Header.Get("x-signature")
	nonce := c.Request.Header.Get("x-nonce-str")
	t := c.Request.Header.Get("x-timestamp")
	header := map[string]string{
		"x-nonce-str": nonce,
		"x-timestamp": t,
		"x-roomid":    RoomId,
		"x-msg-type":  msgType,
	}
	signature2 := signature(header, string(body))
	if signature1 != signature2 {
		zlog.Logger.Error().Msgf("伪造的抖音消息 RoomId%v  x-signature:%v signature :%v body:%v", RoomId, signature1, signature2, string(body))
		return
	}
	// 打印接收到的消息
	messageList := &pb.MessageList{}
	switch msgType {
	case "live_gift": // 礼物
		messageList.Type = "GiftMessage"
		for _, v := range *dyMessageList {
			zlog.Logger.Info().Msgf("抖音收到礼物消息 RoomId:%v  msgType:%v signature1:%v signature2:%v body:%v dyMessageList %v", RoomId, msgType, signature1, signature2, string(body), dyMessageList)

			// 相同id
			if dyMsgIdCheck.IsExist(v.MsgId) {
				zlog.Logger.Error().Msgf("抖音收到消息msgId重复 RoomId:%v  msgType:%v MsgId %v ", RoomId, msgType, v.MsgId)
				continue
			}
			dyMsgIdCheck.AddMessage(v.MsgId)
			// if v.MsgId
			message := &pb.Message{
				Type:      messageList.Type,
				Name:      v.Nickname,
				Uid:       v.SecOpenid,
				HeadImg:   v.AvatarUrl,
				GiftCount: fmt.Sprintf("%v", v.GiftNum),
				Total:     fmt.Sprintf("%v", v.GiftValue/10),
				GiftId:    v.SecGiftId,
			}
			messageList.MsgList = append(messageList.MsgList, message)
		}
	case "live_comment": // 评论
		zlog.Logger.Info().Msgf("抖音收到评论消息 RoomId:%v  msgType:%v signature1:%v signature2:%v body:%v dyMessageList %v", RoomId, msgType, signature1, signature2, string(body), dyMessageList)

		messageList.Type = "ChatMessage"
		for _, v := range *dyMessageList {
			if dyMsgIdCheck.IsExist(v.MsgId) {
				zlog.Logger.Error().Msgf("抖音收到消息msgId重复 RoomId:%v  msgType:%v MsgId %v ", RoomId, msgType, v.MsgId)
				continue
			}
			dyMsgIdCheck.AddMessage(v.MsgId)
			message := &pb.Message{
				Type:    messageList.Type,
				Name:    v.Nickname,
				HeadImg: v.AvatarUrl,
				Uid:     v.SecOpenid,
				Content: v.Content,
			}
			messageList.MsgList = append(messageList.MsgList, message)
		}
	case "live_like": // 点赞
		messageList.Type = "LikeMessage"
		for _, v := range *dyMessageList {
			if dyMsgIdCheck.IsExist(v.MsgId) {
				zlog.Logger.Error().Msgf("抖音收到消息msgId重复 RoomId:%v  msgType:%v MsgId %v ", RoomId, msgType, v.MsgId)
				continue
			}
			dyMsgIdCheck.AddMessage(v.MsgId)

			message := &pb.Message{
				Type:    messageList.Type,
				Name:    v.Nickname,
				HeadImg: v.AvatarUrl,
				Uid:     v.SecOpenid,
				Count:   fmt.Sprintf("%v", v.LikeNum),
			}
			messageList.MsgList = append(messageList.MsgList, message)
		}
	default:
		zlog.Logger.Error().Msgf("抖音收到消息msgType错误 RoomId:%v  msgType:%v ", RoomId, msgType)
		// 返回消息
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "success",
		})
	}
	messageList.RoomId = RoomId

	if !GetMgr().HttpSendMessageToServerLocal(messageList) {
		zlog.Logger.Error().Msgf("HttpSendMessage err:%v", RoomId)
		// 返回消息
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "success",
		})
		return
	}

	// 返回消息
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
	})
}

func signature(header map[string]string, bodyStr string) string {
	// 自定义的密钥
	secret := "abcfhdg"
	keyList := make([]string, 0, 4)
	for key, _ := range header {
		keyList = append(keyList, key)
	}
	sort.Slice(keyList, func(i, j int) bool {
		return keyList[i] < keyList[j]
	})
	kvList := make([]string, 0, 4)
	for _, key := range keyList {
		kvList = append(kvList, key+"="+header[key])
	}
	urlParams := strings.Join(kvList, "&")
	rawData := urlParams + bodyStr + secret
	md5Result := md5.Sum([]byte(rawData))
	return base64.StdEncoding.EncodeToString(md5Result[:])
}
