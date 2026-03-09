package control

import (
	pb "dmGameServer/pb"
	"dmGameServer/untils"
	"dmGameServer/zlog"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

// ----------------------------------
// Chat 评论
func Chat1(c *gin.Context) {
	// 解析 JSON 数据
	data := &pb.Message{}
	// 获取请求体中的原始数据
	de, err := io.ReadAll(c.Request.Body)
	if err != nil {
		zlog.Logger.Error().Msgf("ReadAll%v", err)
		//Fail(c, "", err.Error())
		return
	}
	err = json.Unmarshal(de, &data)
	if err != nil {
		zlog.Logger.Error().Msgf("Unmarshal%v %v", err, string(de))
		//Fail(c, "", err.Error())
		return
	}

	if data.MsgId != "" {
		// 验证
		// 相同id
		if yyMsgIdCheck.IsExist(data.MsgId) {
			//zlog.Logger.Error().Msgf("tkHH收到消息msgId重复 %v ", data.MsgId)
			//Fail(c, "", "")
			return
		}
		yyMsgIdCheck.AddMessage(data.MsgId)
	}

	// 是否有映射
	if data.AnchorId != "" {
		//t, ok := mgr.YsLoginMap.Load(data.AnchorId)
		//if ok {
		//	data.AnchorId = t.(string)
		//}
	}

	accountId := data.AnchorId
	zlog.Logger.Debug().Msgf("tk野游评论AccountId[%v] data[%+v]", accountId, data)
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
func Like1(c *gin.Context) {
	// 解析 JSON 数据
	data := &pb.Message{}
	// 获取请求体中的原始数据
	de, err := io.ReadAll(c.Request.Body)
	if err != nil {
		zlog.Logger.Error().Msgf("ReadAll%v", err)
		//Fail(c, "", err.Error())
		return
	}
	err = json.Unmarshal(de, &data)
	if err != nil {
		zlog.Logger.Error().Msgf("Unmarshal%v %v", err, string(de))
		//Fail(c, "", err.Error())

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

	// 是否有映射
	if data.AnchorId != "" {
		//t, ok := mgr.YsLoginMap.Load(data.AnchorId)
		//if ok {
		//	data.AnchorId = t.(string)
		//}
	}
	accountId := data.AnchorId
	zlog.Logger.Debug().Msgf("tk野游点赞AccountId[%v] data[%+v]", accountId, data)
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

type YsLoginReq struct {
	Old string `json:"old"`
	New string `json:"New"`
}

// Gift 礼物
func Gift1(c *gin.Context) {
	// 解析 JSON 数据
	data := &pb.Message{}
	// 获取请求体中的原始数据
	de, err := io.ReadAll(c.Request.Body)
	if err != nil {
		zlog.Logger.Error().Msgf("ReadAll%v", err)
		//Fail(c, "", err.Error())
		return
	}
	err = json.Unmarshal(de, &data)
	if err != nil {
		zlog.Logger.Error().Msgf("Unmarshal%v %v", err, string(de))
		//Fail(c, "", "Failed to parse JSON data")
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

	// 是否有映射
	if data.AnchorId != "" {
		//t, ok := mgr.YsLoginMap.Load(data.AnchorId)
		//if ok {
		//	data.AnchorId = t.(string)
		//}
	}
	accountId := data.AnchorId
	// 分析
	num := StirngToInt64(data.GiftCount)
	zlog.Logger.Info().Msgf("修正前----野游礼物消息AccountId[%v] data[%+v]", accountId, data)

	//  计算礼物总价值
	if !data.RepeatEndhh {
		if data.GroupId == "0" {
			num = 1
			//  连续的
			zlog.Logger.Info().Msgf("连续礼物0 data:%v  num:%v", data, num)
			// 连续的礼物
			data.GiftCount = fmt.Sprintf("%v", num)
		} else {
			// 添加连续返回差量
			num = yyGorupCheck.AddGorup(data.GroupId, data.Uid, int32(num))
			//  连续的
			zlog.Logger.Info().Msgf("连续礼物非0 data:%v  num:%v", data, num)
			// 连续的礼物
			data.GiftCount = fmt.Sprintf("%v", num)
		}
	} else {
		// 不连续的
		zlog.Logger.Info().Msgf("不连续礼物 %v", data)
		// 需要跳过 连续发过的礼物
		if data.GroupId != "0" {
			if yyGorupCheck.IsExist(data.GroupId, data.Uid) {
				zlog.Logger.Error().Msgf("tk收到礼物GroupId重复 %v %v 需要删除 ", data.GroupId, data.Uid)
				yyGorupCheck.Del(data.GroupId, data.Uid)
				return
			}
		}
	}
	one := StirngToInt64(data.Total)
	// 这个地单个礼物的价值
	data.Total = fmt.Sprintf("%v", one*num)
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

var yyGorupCheck *YyGorupIdCheck

func init() {
	yyMsgIdCheck = &YyMsgIdCheck{
		MsgIdMap: make(map[string]struct{}),
	}
	yyGorupCheck = &YyGorupIdCheck{
		GorupIdMap: make(map[string]int32),
	}
}

type YyGorupIdCheck struct {
	GorupIdMap map[string]int32
	mu         sync.RWMutex
}

// AddMessage 连续才添加
func (d *YyGorupIdCheck) AddGorup(gorupID, uid string, cnt int32) int64 {
	d.mu.Lock()
	defer d.mu.Unlock()
	currCnt := int32(0)
	key := fmt.Sprintf("%v_%v", gorupID, uid)
	if v, ok := d.GorupIdMap[key]; ok {
		currCnt = v
	}
	d.GorupIdMap[key] = cnt
	// 返回新的和旧的差值
	t := cnt - currCnt
	zlog.Logger.Debug().Msgf("AddGorup 连续礼物 gorupID:%v cnt:%v currCnt:%v", key, cnt, currCnt)
	if t < 0 {
		t = 1
	}
	return int64(t)
}

// IsExist 判断是否存在
func (d *YyGorupIdCheck) IsExist(gorupID, uid string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	key := fmt.Sprintf("%v_%v", gorupID, uid)
	if len(d.GorupIdMap) == 0 {
		return false
	}
	_, ok := d.GorupIdMap[key]
	return ok
}

// 删除
func (d *YyGorupIdCheck) Del(gorupID, uid string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	key := fmt.Sprintf("%v_%v", gorupID, uid)
	delete(d.GorupIdMap, key)

}

//------------------------------

var IpTkHH = "127.0.0.1:3000"

// 开始TkHH
func StartTkHH(accountId string) {
	url := fmt.Sprintf("http://%v/loginTk?name=%v", IpTkHH, accountId)
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		zlog.Logger.Error().Msgf("请求失败%v", err)
		untils.TapErr(fmt.Sprintf("StartTkHH 请求失败%v", err))
		return
	}
	zlog.Logger.Info().Msgf("开始TkHH请求成功%v", accountId)
	defer resp.Body.Close()
	return

}

// 关闭tkhh
func CloseHHTk(AccountId string) {
	// 关闭tkhh
	go func() {
		zlog.Logger.Info().Msgf("关闭tkhh %v", AccountId)
		url := fmt.Sprintf("http://%v/logoutTk?name=%v", IpTkHH, AccountId)
		// 发送get请求
		resp, err := http.Get(url)
		if err != nil {
			zlog.Logger.Error().Msgf("请求失败%v", err)
		} else {
			// 打印返回值
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				zlog.Logger.Error().Msgf("读取失败%v", err)
			} else {
				zlog.Logger.Info().Msgf("请求成功%v", string(body))
			}
			defer resp.Body.Close()
		}
	}()
}
