package control

import (
	"crypto/md5"
	"dmGameServer/untils"
	"dmGameServer/zlog"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"sort"
	"strings"
	"sync"
	"time"
)

// KuiShouOfficialInfo 快手官方信息
type KuiShouOfficialInfo struct {
	AppId       string           // appid
	AppSecret   string           // 秘钥
	AccessToken string           // token
	RoomCode    string           // 房间号
	UserInfo    *KuiShouUserInfo // 主播信息
	RoundId     int32
	IsRunIng    bool     // 循环是否运行
	ginCtxMap   sync.Map // key msgUid value: *gin.Context //
}

func (k *KuiShouOfficialInfo) AddGinCtx(msgUid string, c *gin.Context) {
	k.ginCtxMap.Store(msgUid, c)
}

func (k *KuiShouOfficialInfo) PopGinCtx(msgUid string) *gin.Context {
	if v, ok := k.ginCtxMap.Load(msgUid); ok {
		k.ginCtxMap.Delete(msgUid)
		return v.(*gin.Context)
	}
	return nil
}

// KuiShouLoginRes 快手登录回复
type KuiShouLoginRes struct {
	Result      int32  `json:"result"`       // 如果result 不是 1， 会有error_msg
	AccessToken string `json:"access_token"` // 用于获取隐私资源
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"` // 过期时间（秒）
}

// KuiShouBingReq 快手绑定请求
type KuiShouBingReq struct {
	RoomCode  string `json:"roomCode"`  // 快手直播间绑定码
	Timestamp int64  `json:"timestamp"` // 时间戳
	Sign      string `json:"sign"`      // token
	RoundId   string `json:"roundId"`   // 房间号
}

// KuiShouBingRes 快手绑定返回
type KuiShouBingRes struct {
	Result   int32            `json:"result"`   // 如果result 不是 1， 会有error_msg
	ErrorMsg string           `json:"errorMsg"` // 错误信息
	UserInfo *KuiShouUserInfo `json:"userInfo"` // 主播信息
}

// 快手主播
type KuiShouUserInfo struct {
	AuthorOpenId string `json:"authorOpenId"` // 主播id
	UserName     string `json:"userName"`     // 主播昵称
	HeadUrl      string `json:"headUrl"`      // 主播头像
}

// KuiShouStatus 快手状态返回
type KuiShouStatusRes struct {
	Result     int32  `json:"result"`     // 如果result 不是 1， 会有error_msg
	ErrorMsg   string `json:"errorMsg"`   // 错误信息
	TaskStatus int32  `json:"taskStatus"` // 主播信息 //1-开启状态；2-断开状态
}

// CommonRes
type CommonRes struct {
	Result   int32  `json:"result"`   // 如果result 不是 1， 会有error_msg
	ErrorMsg string `json:"errorMsg"` // 错误信息
}

// KuiShouGiftTopRes
type KuiShouStopRes struct {
	Result     int32  `json:"result"`     // 如果result 不是 1， 会有error_msg
	ErrorMsg   string `json:"errorMsg"`   // 错误信息
	TaskStatus int32  `json:"taskStatus"` // 1-开启状态；2-断开状态
}

func (k *KuiShouOfficialInfo) GetUrl(url string) string {
	return fmt.Sprintf("%v?app_id=%v&access_token=%v", url, k.AppId, k.AccessToken)
}

func (k *KuiShouOfficialInfo) Init() bool {
	// 1.获取token
	signParamsMap := make(map[string]interface{})
	signParamsMap["app_id"] = k.AppId
	signParamsMap["app_secret"] = k.AppSecret
	signParamsMap["grant_type"] = "client_credentials" // 固定值“client_credentials”
	res, err := untils.SendPostRequest2("https://open.kuaishou.com/oauth2/access_token", signParamsMap)
	if err != nil {
		zlog.Logger.Error().Msgf("SendPostRequest err:%v", err)
		return false
	}
	kuiShouLoginRes := &KuiShouLoginRes{}
	err = json.Unmarshal(res, kuiShouLoginRes)
	if err != nil {
		zlog.Logger.Error().Msgf("Unmarshal err:%v", err)
		return false
	}
	if kuiShouLoginRes.Result != 1 {
		zlog.Logger.Error().Msgf("access_token code%v %+v", kuiShouLoginRes.Result, signParamsMap)
		return false
	}
	zlog.Logger.Debug().Msgf("获取token成功:%+v", kuiShouLoginRes)
	k.AccessToken = kuiShouLoginRes.AccessToken
	return true
}

func (k *KuiShouOfficialInfo) Run(a *AnchorClient) bool {
	// 绑定
	if !k.RunStart(a) {
		return false
	}
	// 快捷礼物和礼物置顶
	if !k.InteractiveGift(a) {
		return false
	}
	// 快捷加入
	if !k.InteractiveJion(a) {
		return false
	}

	return true
}

// RunStart 绑定
func (k *KuiShouOfficialInfo) RunStart(a *AnchorClient) bool {
	k.RoundId++
	signParamsMap := make(map[string]interface{})
	signParamsMap["roomCode"] = k.RoomCode
	signParamsMap["timestamp"] = time.Now().UnixMilli()
	signParamsMap["roundId"] = fmt.Sprintf("%v", k.RoundId)
	signParamsMap["sign"] = CalcSign(signParamsMap, k)
	data, err := json.Marshal(signParamsMap)
	if err != nil {
		zlog.Logger.Error().Msgf("Marshal err:%v", err)
		return false
	}
	res, err := untils.SendPostRequest(k.GetUrl("https://open.kuaishou.com/openapi/developer/live/data/task/start"), data)
	if err != nil {
		zlog.Logger.Error().Msgf("SendPostRequest err:%v", err)
		return false
	}
	kuiShouBingRes := &KuiShouBingRes{}
	err = json.Unmarshal(res, kuiShouBingRes)
	if err != nil {
		zlog.Logger.Error().Msgf("Unmarshal err:%v", err)
		return false
	}
	if kuiShouBingRes.Result != 1 {
		zlog.Logger.Error().Msgf("快手绑定kuiShouBingRes code%v err:%v", kuiShouBingRes.Result, kuiShouBingRes.ErrorMsg)
		return false
	}
	a.KuiShouInfo.UserInfo = kuiShouBingRes.UserInfo
	zlog.Logger.Info().Msgf("快手绑定成功  AccountId:%v UserInfo:%+v  RoomCode:%v", a.AccountId, kuiShouBingRes.UserInfo, a.KuiShouInfo.RoomCode)
	k.IsRunIng = true
	return true
}

// 快捷加入
func (k *KuiShouOfficialInfo) InteractiveJion(a *AnchorClient) bool {
	appConfigInfo := a.WebsocketClient.AppInfo
	if appConfigInfo.InteractiveJoin == "" {
		zlog.Logger.Info().Msgf("快手互动加入为空 %v", a.AccountId)
		return true
	}
	signParamsMap := make(map[string]interface{})
	signParamsMap["roomCode"] = k.RoomCode
	signParamsMap["timestamp"] = time.Now().UnixMilli()
	signParamsMap["roundId"] = fmt.Sprintf("%v", k.RoundId)
	signParamsMap["type"] = "1"
	signParamsMap["data"] = appConfigInfo.InteractiveJoin
	signParamsMap["sign"] = CalcSign(signParamsMap, k)
	data, err := json.Marshal(signParamsMap)
	if err != nil {
		zlog.Logger.Error().Msgf("Marshal err:%v", err)
		return false
	}
	res, err := untils.SendPostRequest(k.GetUrl("https://open.kuaishou.com/openapi/developer/live/data/interactive/start"), data)
	if err != nil {
		zlog.Logger.Error().Msgf("SendPostRequest err:%v", err)
		return false
	}
	commonRes := &CommonRes{}
	err = json.Unmarshal(res, commonRes)
	if err != nil {
		zlog.Logger.Error().Msgf("Marshal err:%v", err)
		return false
	}
	if commonRes.Result != 1 {
		zlog.Logger.Error().Msgf("commonRes code:%v err:%v  signParamsMap:%v", commonRes.Result, commonRes.ErrorMsg, signParamsMap)
		return true
	}
	zlog.Logger.Info().Msgf("快手快捷加入设置成功commonRes:%+v  RoomCode:%v  signParamsMap:%v", commonRes, k.RoomCode, signParamsMap)
	return true
}

// 快捷礼物
func (k *KuiShouOfficialInfo) InteractiveGift(a *AnchorClient) bool {
	signParamsMap := make(map[string]interface{})
	signParamsMap["roomCode"] = k.RoomCode
	appInfo := a.WebsocketClient.AppInfo
	giftList := ""
	for i, id := range appInfo.TopGiftList {
		if i == len(appInfo.TopGiftList)-1 {
			giftList += id
			break
		}
		giftList += id + ","
	}
	zlog.Logger.Info().Msgf("快手互动请求giftList:%v %v %v", giftList, appInfo.TopGiftList, k.RoomCode)
	signParamsMap["giftList"] = giftList
	if appInfo.InteractiveGift != "" {
		signParamsMap["giftExtendInfo"] = appInfo.InteractiveGift
	}
	signParamsMap["timestamp"] = time.Now().UnixMilli()
	signParamsMap["sign"] = CalcSign(signParamsMap, k)

	zlog.Logger.Debug().Msgf("快手互动请求signParamsMap:%+v", signParamsMap)

	data, err := json.Marshal(signParamsMap)
	if err != nil {
		zlog.Logger.Error().Msgf("Marshal err:%v", err)
		return false
	}
	res, err := untils.SendPostRequest(k.GetUrl("https://open.kuaishou.com/openapi/developer/live/interactive/gift/top"), data)
	if err != nil {
		zlog.Logger.Error().Msgf("SendPostRequest err:%v", err)
		return false
	}
	commonRes := &CommonRes{}
	err = json.Unmarshal(res, commonRes)
	if err != nil {
		zlog.Logger.Error().Msgf("Marshal err:%v", err)
		return false
	}
	if commonRes.Result != 1 {
		zlog.Logger.Error().Msgf("commonRes code%v err:%v", commonRes.Result, commonRes.ErrorMsg)
		return false
	}
	zlog.Logger.Info().Msgf("快手互动快捷礼物设置成功commonRes:%+v", commonRes)
	return true
}

// RunStop 解绑
func (k *KuiShouOfficialInfo) RunStop() bool {
	tt := time.Now().UnixMilli()
	signParamsMap := make(map[string]interface{})
	signParamsMap["roomCode"] = k.RoomCode
	signParamsMap["timestamp"] = tt
	signParamsMap["roundId"] = fmt.Sprintf("%v", k.RoundId)
	signParamsMap["sign"] = CalcSign(signParamsMap, k)
	k.IsRunIng = false
	zlog.Logger.Info().Msgf("解绑signParamsMap:%+v  RoomCode:%v", signParamsMap, k.RoomCode)
	data, err := json.Marshal(signParamsMap)
	if err != nil {
		zlog.Logger.Error().Msgf("Marshal err:%v", err)
		return false
	}

	res, err := untils.SendPostRequest(k.GetUrl("https://open.kuaishou.com/openapi/developer/live/data/task/stop"), data)
	if err != nil {
		zlog.Logger.Error().Msgf("SendPostRequest err:%v", err)
		return false
	}
	kuiShouStopRes := &KuiShouStopRes{}
	err = json.Unmarshal(res, kuiShouStopRes)
	if err != nil {
		zlog.Logger.Error().Msgf("Unmarshal err:%v", err)
		return false
	}
	if kuiShouStopRes.Result != 1 {
		zlog.Logger.Debug().Msgf("kuiShouBingRes code%v err:%v", kuiShouStopRes.Result, kuiShouStopRes.ErrorMsg)
		return false
	}

	zlog.Logger.Info().Msgf("kuiShouStopRes:%+v ", kuiShouStopRes)

	return true
}

// CalcSign 获取快手的的签名结果
func CalcSign(signParamsMap map[string]interface{}, k *KuiShouOfficialInfo) string {
	// 添加用于签名
	signParamsMap["app_id"] = k.AppId
	// 按照字母排序
	var keys []string
	for key := range signParamsMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	// 组装成待签名字符串
	var paramStrBuilder strings.Builder
	for _, key := range keys {
		paramStrBuilder.WriteString(key)
		paramStrBuilder.WriteString("=")
		paramStrBuilder.WriteString(fmt.Sprintf("%v", signParamsMap[key]))
		paramStrBuilder.WriteString("&")
	}
	paramStr := paramStrBuilder.String()
	paramStr = paramStr[:len(paramStr)-1] // 去掉最后一个"&"
	// 组装签名字符串
	signStr := paramStr + k.AppSecret
	// 生成签名返回
	hash := md5.New()
	zlog.Logger.Debug().Msgf("signStr:%v", signStr)
	// 删掉
	delete(signParamsMap, "app_id")
	hash.Write([]byte(signStr))
	return hex.EncodeToString(hash.Sum(nil))
}
