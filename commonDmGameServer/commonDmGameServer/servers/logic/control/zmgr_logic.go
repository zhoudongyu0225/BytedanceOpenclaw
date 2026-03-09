package control

import (
	"bytes"
	"context"
	"dmGameServer/common"
	"dmGameServer/model"
	pb "dmGameServer/pb"
	"dmGameServer/untils"
	"dmGameServer/zlog"
	"encoding/binary"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"log"
	"math"
	"net"
	"net/http"
	"runtime/debug"
	"strconv"
	"sync"
	"time"
)

// 逻辑管理器
var mgr *ControlLogicMgr

const (
	WsSendChanMaxNum   = 10000 // ws发送消息的管道最大数量
	HttpSendChanMaxNum = 10000 // http发送消息的管道最大数量
)

// GetAnchorClient 获取 主播的结构体
func GetAnchorClient(AccountId string) *AnchorClient {
	if mgr == nil {
		zlog.Logger.Error().Msgf("mgr nil[%v]", AccountId)
		return nil
	}
	if v1, ok := mgr.AnchorMap.Load(AccountId); ok {
		v := v1.(*AnchorClient)
		return v
	} else {
		zlog.Logger.Debug().Msgf("没有找到主播的结构体[%v]", AccountId)
		return nil
	}
}

// GetAnchorClientByRoomId 获取 主播的结构体
func GetAnchorClientByRoomId(RoomId string) *AnchorClient {
	if RoomId == "" {
		zlog.Logger.Error().Msgf("GetAnchorClientByRoomId RoomId 是空的 %v", RoomId)
		return nil
	}
	if mgr == nil {
		zlog.Logger.Error().Msgf("mgr nil[%v]", RoomId)
		return nil
	}
	if acIdI, ok := mgr.RoomIdAccontIdMap.Load(RoomId); ok {
		acId := acIdI.(string)
		return GetAnchorClient(acId)
	} else {
		zlog.Logger.Error().Msgf("没有根据房间-找到对应的主播结构体 roomCode [%v]", RoomId)
		return nil
	}
}

// --------------------------------------------客户端ws---------------------------------

type WebsocketAnchorClient struct {
	Socket             *websocket.Conn
	GameId             int32                // 游戏Id
	Version            string               // 版本号
	AccountId          string               // 账号Id
	IsClientLogOn      bool                 // 客户端是否登录
	WriteCloseChan     chan bool            // 写关闭管道
	ReadWsCloseChan    chan bool            // 关闭主播的读ws的道
	ReadCloseChan      chan bool            // 读弹幕和ws转换后的消息关闭管道
	BarrageMessageChan chan *pb.MessageList // 弹幕消息管道
	Send               chan []byte          // 发给ws的消息管道什
	IsStart            bool                 // 是否开始游戏
	SubNum             int32                // 累计扣除的连胜分数
	RoomId             string               // 房间（登录就会赋值）什
	IsDebug            bool
	StartTime          int64 // 建立连接开始时间戳 （秒）
	LastPingTime       int64 // 最后一次ping的时间戳 （秒）
	AppInfo            *model.AppConfigInfo
	Token              string // 主播token 目前用于抖音
	Access_token       string // 主播access_token 目前用于抖音
	giftV              int64  // 本次游戏收到的礼物价值
	// slg数据
	AnchorV *pb.AnchorDBInfo
	SdkId   int64
	// 上次补包时间
	LastSendTime int64
	// 点赞队列
	LikeMessageList *pb.MessageList
	// 玩家列表
	PlayerMap sync.Map
}

func StirngToInt64(str string) int64 {
	if str == "" {
		return 0
	}
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return int64(i)
}

// 弹幕汇总消息(里面礼物逻辑只能调用一次)
func (ws *WebsocketAnchorClient) convergeMessage(messageList *pb.MessageList) {
	if messageList == nil {
		zlog.Logger.Error().Msgf("convergeMessage err:%v", "message is nil")
		return
	}
	zlog.Logger.Debug().Msgf("弹幕汇总消息 %v", messageList.Type)
	switch messageList.Type {
	case "GiftMessage": // 房间礼物推送响应
		GiftMessage(messageList, ws)
	case "ChatMessage": // 房间评论推送响应
		ChatMessage(messageList, ws)
	case "LikeMessage": // 房间点赞响应
		if ws.LikeMessageList == nil {
			ws.LikeMessageList = &pb.MessageList{}
		}
		// 合并
		ws.LikeMessageList = mergeMessagesLike(messageList, ws.LikeMessageList)
	default:
		zlog.Logger.Error().Msgf("不是期望的类型", messageList.Type)
	}
}

// mergeMessages 合并 MessageList 中 uid 相同的消息
func mergeMessagesLike(messageListA, messageListB *pb.MessageList) *pb.MessageList {
	// 创建一个 map 来根据 uid 分组消息
	mergedMessages := make(map[string]*pb.Message)
	// 遍历 messageListA，将消息按 uid 放入 mergedMessages
	for _, msg := range messageListA.MsgList {
		if msg == nil {
			continue
		}
		mergedMessages[msg.Uid] = msg
	}
	// 遍历 messageListB，将消息按 uid 放入 mergedMessages
	for _, msg := range messageListB.MsgList {
		uid := msg.Uid
		if v, ok := mergedMessages[uid]; ok {
			a := untils.StringToInt32(mergedMessages[uid].Count)
			b := untils.StringToInt32(v.Count)
			mergedMessages[uid].Count = fmt.Sprintf("%v", a+b)
		} else {
			mergedMessages[uid] = msg
		}
	}
	newMessageList := &pb.MessageList{}
	newMessageList.Type = messageListA.Type
	newMessageList.AccountId = messageListA.AccountId
	newMessageList.RoomId = messageListA.RoomId
	newMessageList.MsgList = nil
	for _, msg := range mergedMessages {
		newMessageList.MsgList = append(newMessageList.MsgList, msg)
	}
	return newMessageList
}

// 关闭读写管道
func (ws *WebsocketAnchorClient) CloseReadAndWirte() {
	// 关闭主播的读幕管道
	ws.ReadWsCloseChan <- true
	zlog.Logger.Info().Msgf("关闭主播的读ws管道 %v", ws.AccountId)
	// 关闭主播的读弹幕幕管道
	ws.ReadCloseChan <- true
	zlog.Logger.Info().Msgf("关闭主播的读幕管道 %v", ws.AccountId)
	// 写关闭管道
	ws.WriteCloseChan <- true
	zlog.Logger.Info().Msgf("写关闭管道 %v", ws.AccountId)
	// 延迟10毫秒 放置关闭的消息没有发送出去
	time.Sleep(10 * time.Millisecond)
	err := ws.Socket.Close()
	ws.IsClientLogOn = false
	if err != nil {
		zlog.Logger.Error().Msgf("CloseReadAndWirte err:%v", err)
	}
}

// Close 关闭
func (ws *WebsocketAnchorClient) Close(isActive bool) {
	zlog.Logger.Info().Msgf("关闭连接 1111111111%v", ws.AccountId)
	// 关闭读写
	ws.CloseReadAndWirte()
	zlog.Logger.Info().Msgf("关闭弹幕管道 %v", ws.AccountId)

	if ws.AnchorV != nil {
		allTime := time.Now().Unix() - ws.StartTime
		// 添加主播数据
		DelTapMap(ws.AnchorV.PlatformId, ws.AccountId)
		// 更新直播时间
		model.UpdateAnchorLiveTime(ws.AccountId, allTime, ws.AnchorV.PlatformId, ws.AnchorV.GameId)

		untils.ZBGPoss(ws.AccountId, ws.AnchorV.PlatformId, ws.AnchorV.GameId, ws.AnchorV.NickName, allTime, ws.giftV,
			GetTapMapNum(), ws.Version, 0,
			model.GetLiveInfo(ws.AccountId, ws.AnchorV.PlatformId, ws.AnchorV.GameId))
	}
	// 清理
	nickName := "未查命名"
	v := GetAnchorClient(ws.AccountId)
	if v == nil {
		zlog.Logger.Error().Msgf("GetAnchorClient is nil %v", ws.AccountId)
	} else {
		nickName = v.Anchor.NickName
	}

	t := time.Unix(ws.StartTime, 0)
	// 获取年月日
	year, month, day := t.Date()
	// 创建开始直播的凌晨的时间戳
	midnight := time.Date(year, month, day, 0, 0, 0, 0, t.Location()).Unix()
	// 创建开始直播的第二天凌晨的时间戳
	newMidnight := midnight + 86400
	now := time.Now().Unix()
	if now > newMidnight {
		// 跨天分2不
		allTime1 := newMidnight - ws.StartTime
		// ---------------上一天--------
		// 直播时间
		addOverviewOfTransactions := &model.OverviewOfTransactions{
			LiveTime: allTime1,
		}
		// 更新[流水总览] 上一天
		model.UpdateOverviewOfTransactions(addOverviewOfTransactions, GetCollectionOverviewOfTransactionsKey(ws.AccountId), "1")

		// 更新[主播流水]直播时长上一天
		anchorOfTransactions := &model.AnchorOfTransactions{
			AccountId: ws.AccountId,
			Name:      nickName,
		}
		model.UpdateAnchorOfTransactionsLiveTime(anchorOfTransactions, allTime1, GetCollectionAnchorOfTransactionsKey(ws.AccountId), "1")

		// -------------当天--------------------
		allTime2 := now - newMidnight
		// 直播时间
		addOverviewOfTransactions = &model.OverviewOfTransactions{
			LiveTime: allTime2,
		}
		// 更新[流水总览] 当天
		model.UpdateOverviewOfTransactions(addOverviewOfTransactions, GetCollectionOverviewOfTransactionsKey(ws.AccountId))

		// 更新[主播流水]直播时长 当天
		anchorOfTransactions = &model.AnchorOfTransactions{
			AccountId: ws.AccountId,
			Name:      nickName,
		}
		model.UpdateAnchorOfTransactionsLiveTime(anchorOfTransactions, allTime2, GetCollectionAnchorOfTransactionsKey(ws.AccountId))
	} else {
		allTime := time.Now().Unix() - ws.StartTime
		// 直播时间
		addOverviewOfTransactions := &model.OverviewOfTransactions{
			LiveTime: allTime,
		}
		// 更新流水总览
		model.UpdateOverviewOfTransactions(addOverviewOfTransactions, GetCollectionOverviewOfTransactionsKey(ws.AccountId))
		// 更新主播流水直播时长 当天
		anchorOfTransactions := &model.AnchorOfTransactions{
			AccountId: ws.AccountId,
			Name:      nickName,
		}
		model.UpdateAnchorOfTransactionsLiveTime(anchorOfTransactions, allTime, GetCollectionAnchorOfTransactionsKey(ws.AccountId))
	}
	// 删掉
	mgr.AnchorMap.Delete(ws.AccountId)
	if v != nil {
		if v.KuiShouInfo != nil {
			// 快手解绑
			v.KuiShouInfo.RunStop()
		}
	}
	mgr.RoomIdAccontIdMap.Delete(ws.RoomId)
	// 删掉正在直播的主播
	model.DelCurrCurrAnchor(ws.AccountId)
	zlog.Logger.Info().Msgf("关闭连接------2222222 %v", ws.AccountId)
}

// WsSend ws消息发送
func (ws *WebsocketAnchorClient) WsSend(cmd int16, m proto.Message) bool {
	if m == nil {
		zlog.Logger.Error().Msgf("WsSend errm is nil :%v", ws.AccountId)
		return false
	}

	//	zlog.Logger.Debug().Msgf("WsSend cmd:%v %v", cmd, ws.AccountId)

	var body []byte
	var err error
	body, err = proto.Marshal(m)
	if err != nil {
		zlog.Logger.Error().Msgf("WsSend err:%v %v", err, ws.AccountId)
		return false
	}

	// 业务数据的长度
	bodyLength := uint32(len(body))
	// 整个数据包长度(创建buffer的)
	datagramLength := bodyLength + 4 + 2 + 1 + 1
	// 准备数据包大小缓冲区
	buffer := bytes.NewBuffer(make([]byte, 0, datagramLength))
	// 1.缓冲区写 业务数据 长度
	err = binary.Write(buffer, binary.LittleEndian, bodyLength)
	// 2.缓冲区写command
	err = binary.Write(buffer, binary.LittleEndian, cmd)
	// 3.加密
	err = binary.Write(buffer, binary.LittleEndian, false)
	// 消息理论上不会大于这么大
	if len(ws.Send) >= 100 {
		zlog.Logger.Error().Msgf("WsSend 通道消息冗余:len%v  AccountId%v", len(ws.Send), ws.AccountId)
	}
	if len(ws.Send) >= WsSendChanMaxNum {
		zlog.Logger.Error().Msgf("WsSend 通道已满:%v AccountId%v", len(ws.Send), ws.AccountId)
		return false
	}
	if ws.Send == nil {
		zlog.Logger.Error().Msgf("ws.Send nil")
		return true
	}
	if datagramLength >= 100 {
		newBody, err1 := untils.GzipCompress(body)
		if err1 != nil {
			zlog.Logger.Error().Msgf("WsSend err:%v %v", err1, ws.AccountId)
			return false
		}
		mbps := float64(datagramLength) / 1024 / 1024 * 8
		gMbps := float64(len(newBody)) / 1024 / 1024 * 8
		if mbps > gMbps {
			// 业务数据的长度
			bodyLength = uint32(len(newBody))
			// 整个数据包长度(创建buffer的)
			datagramLength = bodyLength + 4 + 2 + 1 + 1
			// 准备数据包大小缓冲区
			bufferNew := bytes.NewBuffer(make([]byte, 0, datagramLength))
			// 1.缓冲区写 业务数据 长度
			err = binary.Write(bufferNew, binary.LittleEndian, bodyLength)
			// 2.缓冲区写command
			err = binary.Write(bufferNew, binary.LittleEndian, cmd)
			// 3.加密
			err = binary.Write(bufferNew, binary.LittleEndian, false)
			// 4.压缩
			err = binary.Write(bufferNew, binary.LittleEndian, true)
			// 5.缓冲区写业务数据
			_, err = bufferNew.Write(newBody)
			if err != nil {
				zlog.Logger.Error().Msgf("WsSend err:%v %v", err, ws.AccountId)
				return false
			}
			ws.Send <- bufferNew.Bytes()
		} else {
			// 非压缩
			err = binary.Write(buffer, binary.LittleEndian, false)
			// 5.缓冲区写业务数据
			_, err = buffer.Write(body)
			if err != nil {
				zlog.Logger.Error().Msgf("WsSend err:%v %v", err, ws.AccountId)
				return false
			}
			ws.Send <- buffer.Bytes()
		}
	} else {
		// 非压缩
		err = binary.Write(buffer, binary.LittleEndian, false)
		// 5.缓冲区写业务数据
		_, err = buffer.Write(body)
		if err != nil {
			zlog.Logger.Error().Msgf("WsSend err:%v %v", err, ws.AccountId)
			return false
		}
		ws.Send <- buffer.Bytes()
	}
	return true
}

// 读
func (ws *WebsocketAnchorClient) Read() {
	wbMessageModelChan := make(chan *WbMessageModel, HttpSendChanMaxNum)
	// 主播的协程
	go func() {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				untils.PanicPoss1(err, stack, ws.AccountId)
			}
		}()

		// 1秒一次的
		ticker1 := time.NewTicker(time.Second * 1)
		defer ticker1.Stop()

		// 3600秒一次的  1小时
		ticker3600 := time.NewTicker(time.Second * 3600)
		defer ticker3600.Stop()

		// 上一次时间
		lastTime := time.Now().Unix()
		// 主播的协程
		for {
			select {
			case <-ticker3600.C: // 1小时
			case <-ticker1.C: // 1秒
				nowTime := time.Now().Unix()
				// 在线跨天时间
				if untils.IsNewDay(nowTime, lastTime) {
					//todo: 跨天
				}
				if ws != nil {
					// slg的定时器
					ws.SlgTimeTick()
				}
				lastTime = nowTime

				// 点赞
				if ws.LikeMessageList == nil {
					ws.LikeMessageList = &pb.MessageList{}
				}
				// 点赞
				LikeMessage(ws.LikeMessageList, ws)
				ws.LikeMessageList = &pb.MessageList{}

			case wbMessageModel := <-wbMessageModelChan: // 读取到的消息
				if _, ok := WsProcessMap[wbMessageModel.Command]; !ok {
					zlog.Logger.Error().Msgf("错误的消息类型:%v %v", wbMessageModel.Command, ws.AccountId)
					continue
				}
				if wbMessageModel.Command != COMMAND_LOGIN_C2S {
					if !ws.IsClientLogOn {
						zlog.Logger.Error().Msgf("还未登陆  [%v]  %v", ws.AccountId, wbMessageModel.Command)
						// ws通知
						Tips(ws, "还未登陆")
						continue
					}
				}
				// 处理
				WsProcessMap[wbMessageModel.Command](ws, wbMessageModel)
			case msg := <-ws.BarrageMessageChan:
				ws.convergeMessage(msg)
			case <-ws.ReadCloseChan:
				zlog.Logger.Info().Msgf("主播的协程1协程关闭 %v", ws.AccountId)
				return
			}
		}
	}()
	defer func() {
		if err := recover(); err != nil {
			stack := debug.Stack()
			untils.PanicPoss1(err, stack, ws.AccountId)
		}
	}()

	// 读取的 读取的都会给主播的协程
	for {
		select {
		case <-ws.ReadWsCloseChan:
			zlog.Logger.Info().Msgf("主播的协程2协程关闭 %v", ws.AccountId)
			return
		default:
			wsReadTimeout := time.Duration(common.WsRedOut) * time.Second
			if ws.IsClientLogOn {
				wsReadTimeout = time.Duration(20) * time.Second
			}
			if ws.Socket == nil {
				zlog.Logger.Error().Msgf("ws.Socket is nil:%v", ws.AccountId)
				return
			}
			err := ws.Socket.SetReadDeadline(time.Now().Add(wsReadTimeout))
			if err != nil {
				zlog.Logger.Error().Msgf("SetReadDeadline err:%v %v %v", err, ws.AccountId, common.WsRedOut)
			}
			messageType, data, err := ws.Socket.ReadMessage()
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					zlog.Logger.Info().Msgf("读取时间超时读协程关闭 %v", ws.AccountId)
					ws.Close(true)
					return
					// 处理超时错误
				} else {
					// 处理其他错误
					zlog.Logger.Error().Msgf("读取消息错误读协程关闭 err:%v %v", err, ws.AccountId)
					ws.Close(true)
					return
				}
			}
			switch messageType {
			case websocket.BinaryMessage:
				dataBuffer := bytes.NewBuffer(data)
				wbMessageModel := &WbMessageModel{}
				wbMessageModel.Read(dataBuffer)
				wbMessageModelChan <- wbMessageModel
			case websocket.CloseMessage:
				ws.Close(false)
				zlog.Logger.Info().Msgf("被动读协程关闭 err:%v %v", err, ws.AccountId)
				return
			default:
				zlog.Logger.Error().Msgf("错误的消息类型:%v %v", messageType, ws.AccountId)
			}
		}
	}
}

// 写
func (ws *WebsocketAnchorClient) Write() {
	defer func() {
		if err := recover(); err != nil {
			stack := debug.Stack()
			untils.PanicPoss1(err, stack, ws.AccountId)
		}
	}()
	for {
		select {
		case message, ok := <-ws.Send:
			if !ok {
				zlog.Logger.Error().Msgf("Write err:%v", ws.AccountId)
				continue
			}
			ws.Socket.SetWriteDeadline(time.Now().Add(time.Second * time.Duration(common.WsWriteOut)))
			err := ws.Socket.WriteMessage(websocket.BinaryMessage, message)
			// 太久写不了就关闭
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					zlog.Logger.Info().Msgf("写时间超时写协程关闭  %v", ws.AccountId)
					ws.Close(true)
					return
					// 处理超时错误
				} else {
					// 处理其他错误
					zlog.Logger.Error().Msgf("写消息错误 err:%v %v", err, ws.AccountId)
					ws.Close(true)
					return
				}
			}
		case <-ws.WriteCloseChan:
			zlog.Logger.Info().Msgf("客户端的写协程关闭[%v]", ws.AccountId)
			return
		}
	}
}

// ------------------------------------------------AnchorClient------------------------------------------------

// AnchorClient 与主播客户端交互的
type AnchorClient struct {
	AccountId       string                 // 账号Id
	Anchor          *pb.AnchorDBInfo       // 主播数据
	WebsocketClient *WebsocketAnchorClient // ws管理
	//KcpClient       *KcpAnchorClient       // kcp管理
	KuiShouInfo *KuiShouOfficialInfo // 快手官方信息
}

// -----------------------------------------------ControlLogicMgr-------------------------------------------------

// ControlLogicMgr 逻辑管理器
type ControlLogicMgr struct {
	AnchorMap         sync.Map // 全部主播数据  key:accontId v:*AnchorClient  // 通过id找到对应的AnchorClient
	RoomIdAccontIdMap sync.Map // 主播数据  key:roomdId v:accontId  // 通过roomId找到对应的主播的id

	// 玩家对应的主播
	PlayerRoomIdMap map[string]*AnchorClient // key 玩家uid  玩家对应的主播   // 玩家对应的主播的AnchorClient
	RwLock          sync.RWMutex

	// 平台的主播数量
	tapMap   map[int32]map[string]struct{} // key: gameId.平台   v:主播数量
	tapMapRw sync.RWMutex
}

// GetMgr 获取逻辑管理器
func GetLogicMgr() *ControlLogicMgr {
	return mgr
}

// 如果为 “nil” 说明第一次登录 (不同评定uid会一样吗)
func GetPlayerAc(uid string) *AnchorClient {
	mgr.RwLock.RLock()
	defer mgr.RwLock.RUnlock()
	if v, ok := mgr.PlayerRoomIdMap[uid]; ok {
		return v
	}
	return nil
}

// 设置玩家对应的主播 (主要在发弹幕的时候设置的)
func SetPlayerAc(uid string, anchorClient *AnchorClient) {
	mgr.RwLock.Lock()
	defer mgr.RwLock.Unlock()
	mgr.PlayerRoomIdMap[uid] = anchorClient
}

// 添加主播数
func AddTapMap(gameId int32, id string) {
	mgr.tapMapRw.Lock()
	defer mgr.tapMapRw.Unlock()
	if mgr.tapMap == nil {
		mgr.tapMap = make(map[int32]map[string]struct{})
	}
	if _, ok := mgr.tapMap[gameId]; !ok {
		mgr.tapMap[gameId] = make(map[string]struct{})
	}
	if mgr.tapMap[gameId] == nil {
		mgr.tapMap[gameId] = make(map[string]struct{})
	}
	mgr.tapMap[gameId][id] = struct{}{}
}

// 减少主播数
func DelTapMap(gameId int32, id string) {
	mgr.tapMapRw.Lock()
	defer mgr.tapMapRw.Unlock()
	if mgr.tapMap == nil {
		mgr.tapMap = make(map[int32]map[string]struct{})
	}
	if _, ok := mgr.tapMap[gameId]; !ok {
		mgr.tapMap[gameId] = make(map[string]struct{})
	}
	// 删除
	delete(mgr.tapMap[gameId], id)

}

// 获取主播数  抖音| 快手| 视频号| 总
func GetTapMapNum() []int32 {
	mgr.tapMapRw.RLock()
	defer mgr.tapMapRw.RUnlock()
	a := make([]int32, 0)
	// 抖音
	if v, ok := mgr.tapMap[int32(pb.PlatformId_DouYin)]; ok {
		a = append(a, int32(len(v)))
	} else {
		a = append(a, 0)
	}
	// 快手
	if v, ok := mgr.tapMap[int32(pb.PlatformId_KuaiShou)]; ok {
		a = append(a, int32(len(v)))
	} else {
		a = append(a, 0)
	}
	// 视频号
	if v, ok := mgr.tapMap[int32(pb.PlatformId_ShiPinHao)]; ok {
		a = append(a, int32(len(v)))
	} else {
		a = append(a, 0)
	}
	// tiktok
	if v, ok := mgr.tapMap[int32(pb.PlatformId_TkHh)]; ok {
		a = append(a, int32(len(v)))
	} else {
		a = append(a, 0)
	}
	// 总数
	b := 0
	for _, v := range mgr.tapMap {
		b += len(v)
	}
	a = append(a, int32(b))
	return a
}

//// 获取slg的主播列表
//func GetSlgAnchorInfoList(c *WebsocketAnchorClient) ([]*pb.ShowSlgAnchorInfo, int32) {
//	var list []*pb.ShowSlgAnchorInfo
//	mgr.AnchorMap.Range(func(key, value interface{}) bool {
//		// 跳过自己
//		if c.AccountId == key.(string) {
//			return true
//		}
//		anchorClient := value.(*AnchorClient)
//		if anchorClient == nil {
//			return true
//		}
//		if anchorClient.Anchor == nil {
//			return true
//		}
//
//		if !model.GetPlatformOpenInfo() {
//			// 关闭跨平台
//			// 跳过不同的平台
//			if c.slgAnchor.AnchorV.PlatformId != anchorClient.Anchor.PlatformId {
//				return true
//			}
//		}
//		ssl := &pb.ShowSlgAnchorInfo{
//			AccountId: key.(string),
//			NickName:  anchorClient.Anchor.NickName,
//			HeadUrl:   anchorClient.Anchor.HeadUrl,
//		}
//		// 主播slg的情况
//		list = append(list, ssl)
//		return true
//	})
//
//	// 排序
//	return list, int32(len(list))
//}

func Distance(x1, y1, x2, y2 int32) float64 {
	dx := float64(x2 - x1)
	dy := float64(y2 - y1)
	return math.Sqrt(dx*dx + dy*dy)
}

// InitLogicMgr 初始化逻辑管理器
func InitLogicMgr() {
	mgr = &ControlLogicMgr{
		PlayerRoomIdMap: make(map[string]*AnchorClient),
		tapMap:          make(map[int32]map[string]struct{}),
	}
}

// 发送弹幕消息
func (m *ControlLogicMgr) BarrageSendMessage(AccountId string, messageList *pb.MessageList) bool {
	if v1, ok := m.AnchorMap.Load(AccountId); ok {
		v := v1.(*AnchorClient)
		if v.WebsocketClient == nil {
			zlog.Logger.Error().Msgf("没有找到主播开启的客户端[%v]", AccountId)
			return false
		}
		zlog.Logger.Debug().Msgf("发送消息HttpSendMessage %v", AccountId)
		if v.WebsocketClient.BarrageMessageChan == nil {
			v.WebsocketClient.BarrageMessageChan = make(chan *pb.MessageList, HttpSendChanMaxNum)
		}
		// 消息理论上不会大于这么大
		if len(v.WebsocketClient.BarrageMessageChan) >= WsSendChanMaxNum/2 {
			zlog.Logger.Error().Msgf("WsSend 通道消息冗余:len%v  AccountId%v", len(v.WebsocketClient.BarrageMessageChan), AccountId)
		}
		if len(v.WebsocketClient.BarrageMessageChan) >= WsSendChanMaxNum {
			zlog.Logger.Error().Msgf("WsSend 通道已满:%v AccountId%v", len(v.WebsocketClient.BarrageMessageChan), AccountId)
			return false
		}
		v.WebsocketClient.BarrageMessageChan <- messageList
	} else {
		zlog.Logger.Error().Msgf("没有找到主播的客户端[%v]", AccountId)
		return false
	}
	return true
}

// SetClientLogOn 客户端登录成功(会创建新的去覆盖以前老的)
func (m *ControlLogicMgr) SetClientLogOn(AccountId string, Socket *WebsocketAnchorClient, anchor *pb.AnchorDBInfo, KuiShouInfo *KuiShouOfficialInfo, loginReq *pb.LoginReq) bool {
	// 赋值
	Socket.AccountId = AccountId
	Socket.GameId = anchor.GameId
	Socket.Version = loginReq.Version
	Socket.IsClientLogOn = true
	Socket.StartTime = time.Now().Unix()
	Socket.RoomId = loginReq.RoomId
	if v1, ok := m.AnchorMap.Load(AccountId); ok {
		v := v1.(*AnchorClient)
		v.WebsocketClient = Socket
		v.Anchor = anchor
		zlog.Logger.Info().Msgf("----客户端登录以前已经存在登录链接[%v]", AccountId)
	} else {
		ac := &AnchorClient{
			AccountId:       AccountId,
			Anchor:          anchor,
			WebsocketClient: Socket,
		}
		m.AnchorMap.Store(AccountId, ac)

		// ------------------
		switch loginReq.ModeId {
		case int32(pb.ModeId_GuanYou): // 官游
			switch loginReq.PlatformId {
			case int32(pb.PlatformId_Tiktok): // Tiktok
			case int32(pb.PlatformId_ShiPinHao): // 视频号
			case int32(pb.PlatformId_KuaiShou): // 快手
				ac.KuiShouInfo = KuiShouInfo
				zlog.Logger.Info().Msgf("快手官方开始运行 anchor:%+v  KuiShouInfo:%+v  RoomCode:%v", anchor, ac.KuiShouInfo, ac.KuiShouInfo.RoomCode)
				// 运行 (快手可以推送消息)
				if !ac.KuiShouInfo.Run(ac) {
					return false
				}
			case int32(pb.PlatformId_BZhan): // B站
			case int32(pb.PlatformId_DouYin): // 抖音
			case int32(pb.PlatformId_QQ): // QQ
			}
			m.RoomIdAccontIdMap.Store(Socket.RoomId, ac.AccountId)
		}
	}
	// 活跃主播  true存在 false不存在
	if !model.IsHaveDayAnchorInfo(AccountId) {
		// 开始统计活跃主播
		addOverviewOfTransactions := &model.OverviewOfTransactions{
			ActiveAnchorNum: 1,
		}
		// 更新流水总览
		model.UpdateOverviewOfTransactions(addOverviewOfTransactions, GetCollectionOverviewOfTransactionsKey(AccountId))
	}
	// 直播地址
	liveUrl := loginReq.RoomId
	switch loginReq.ModeId {
	case int32(pb.ModeId_GuanYou): // 官游
		switch loginReq.PlatformId {
		case int32(pb.PlatformId_Tiktok): // Tiktok
		case int32(pb.PlatformId_ShiPinHao): // 视频号
		case int32(pb.PlatformId_KuaiShou): // 快手
			liveUrl = fmt.Sprintf("https://live.kuaishou.com/u/%v", liveUrl)
		case int32(pb.PlatformId_BZhan): // B站
		case int32(pb.PlatformId_DouYin): // 抖音
			liveUrl = fmt.Sprintf("https://webcast.amemv.com/douyin/webcast/reflow/%v", liveUrl)
		case int32(pb.PlatformId_QQ): // QQ
		}
	}
	// AddCurrCurrAnchor 正在直播的主播账号集合
	model.UpdateCurrCurrAnchor(&model.CurAnchorOfTransactions{
		AccountId: AccountId,
		LiveUrl:   liveUrl,
	})
	zlog.Logger.Info().Msgf("客户端登录成功[%v]", AccountId)
	return true
}

// 客户端是否登录
func (m *ControlLogicMgr) IsClientLogOn(AccountId string) (bool, *WebsocketAnchorClient) {
	if v1, ok := m.AnchorMap.Load(AccountId); ok {
		v := v1.(*AnchorClient)
		if v.WebsocketClient == nil {
			return false, nil
		}
		return true, v.WebsocketClient
	} else {
		return false, nil
	}
}

// GetCollectionOverviewOfTransactionsKey 流水总览
func GetCollectionOverviewOfTransactionsKey(AccountId string) string {
	if GetAnchorClient(AccountId) == nil {
		zlog.Logger.Error().Msgf("流水总览没有找到对应的主播数据[%v]", AccountId)
		return ""
	}
	Anchor := GetAnchorClient(AccountId).Anchor
	if Anchor == nil {
		zlog.Logger.Error().Msgf("流水总览没有找到对应的主播数据[%v]", AccountId)
		return ""
	}
	newUid := fmt.Sprintf(model.CollectionOverviewOfTransactions, Anchor.GameId, Anchor.PlatformId, Anchor.ModeId)
	return newUid
}

// GetCollectionAnchorOfTransactionsKey 主播流水
func GetCollectionAnchorOfTransactionsKey(AccountId string) string {
	if GetAnchorClient(AccountId) == nil {
		zlog.Logger.Error().Msgf("没有找到对应的主播数据[%v]", AccountId)
		return ""
	}
	Anchor := GetAnchorClient(AccountId).Anchor
	if Anchor == nil {
		zlog.Logger.Error().Msgf("主播流水没有找到对应的主播数据[%v]", AccountId)
		return ""
	}
	newUid := fmt.Sprintf(model.CollectionAnchorOfTransactions, Anchor.GameId, Anchor.PlatformId, Anchor.ModeId)
	return newUid
}

// GetCollectionPlayerOfTransactionsKey 玩家流水
func GetCollectionPlayerOfTransactionsKey(AccountId string) string {
	if GetAnchorClient(AccountId) == nil {
		zlog.Logger.Error().Msgf("没有找到对应的主播数据[%v]", AccountId)
		return ""
	}
	Anchor := GetAnchorClient(AccountId).Anchor
	if Anchor == nil {
		zlog.Logger.Error().Msgf("玩家流水没有找到对应的主播数据[%v]", AccountId)
		return ""
	}
	newUid := fmt.Sprintf(model.CollectionPlayerOfTransactions, Anchor.GameId, Anchor.PlatformId, Anchor.ModeId)
	return newUid
}

// AddAnchorPlayerRank 添加到排行榜 3个玩家的排行榜
func AddAnchorPlayerRank(anchorId, playerId string, giftValue float64) bool {
	redisKey := fmt.Sprintf(model.AnchorPlayerRank, anchorId)
	client := model.GetRedisClient()
	ctx := context.Background()
	// 先将玩家添加到有序集5
	_, err := client.ZAdd(ctx, redisKey, &redis.Z{
		Score:  float64(giftValue),
		Member: playerId,
	}).Result()

	if err != nil {
		zlog.Logger.Error().Msgf("添加到排行榜err %v", err)
		return false
	}
	// 3
	maxNum := int64(3)
	// 如果有超过3个玩家，删除排名最后的玩家
	count := model.GetRedisClient().ZCard(ctx, redisKey).Val()
	if count > maxNum {
		model.GetRedisClient().ZRemRangeByRank(ctx, redisKey, 0, count-maxNum-1)
	}
	return true
}

// SaveOut
func SaveOut() {
	mgr.AnchorMap.Range(func(key, value interface{}) bool {
		anchorClient := value.(*AnchorClient)
		if anchorClient == nil {
			return true
		}
		if anchorClient.WebsocketClient != nil {
			anchorClient.WebsocketClient.Close(true)
		}
		return true
	})
	untils.TapErr(fmt.Sprintf("----------slg服务器关闭存储完毕-----------"))
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  8192,
	WriteBufferSize: 8192,
	CheckOrigin: func(r *http.Request) bool {
		// 在这里添加自定义的来源验证逻辑s
		// 检查请求的来源是否符合您的要求，如果符合返回true，否则返回false
		// 例如，您可以检查请求中的Origin头部或其他标识来验证来源
		// 这里是一个简单的示例，拒绝所有请求
		return true
	},
}

func HandleConnection(c *gin.Context) {
	w := c.Writer
	r := c.Request
	// 升级HTTP连接为WebSocket连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error", err)
		return
	}
	zlog.Logger.Info().Msgf("新的连接 %v", conn.RemoteAddr().String())
	client := &WebsocketAnchorClient{
		Socket:             conn,
		Send:               make(chan []byte, WsSendChanMaxNum),
		WriteCloseChan:     make(chan bool, 2),                             // 写关闭管道
		ReadCloseChan:      make(chan bool, 2),                             // 读弹幕和ws转换后的消息关闭管道
		ReadWsCloseChan:    make(chan bool, 2),                             // 关闭主播的读ws的道
		BarrageMessageChan: make(chan *pb.MessageList, HttpSendChanMaxNum), // 弹幕消息管道
	}
	// 一个主播只有两个协程  一个读一个写
	go client.Read()
	go client.Write()
}

// PingPong 心跳
func PingPong(c *WebsocketAnchorClient, weMessage *WbMessageModel) {
	if c.RoomId != "" {
		// 更新
		mgr.RoomIdAccontIdMap.Store(c.RoomId, c.AccountId)
	}
	now := time.Now().Unix()
	c.LastPingTime = now
	// ws通知
	c.WsSend(COMMAND_PONG_S2C, &pb.PongResp{
		ServerTime: now,
	})
}
