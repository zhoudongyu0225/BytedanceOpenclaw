package control

import (
	"bytes"
	"dmGameServer/zlog"
	"fmt"
	"github.com/goccy/go-json"
	"io"
	"net/http"
	"time"
)

// DyGetInfoRes
type DyGetInfoRes struct {
	Data struct {
		Info struct {
			RoomId       int64  `json:"room_id"`
			AnchorOpenId string `json:"anchor_open_id"`
			AvatarUrl    string `json:"avatar_url"`
			NickName     string `json:"nick_name"`
		} `json:"info"`
	} `json:"data"`
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

type DyTaskInfo struct {
	ErrNo  int    `json:"err_no"`
	ErrMsg string `json:"err_msg"`
	Logid  string `json:"logid"`
	Data   struct {
		TaskId string `json:"task_id"`
	} `json:"data"`
}

type DyGiftInfo struct {
	ErrNo  int    `json:"err_no"`
	ErrMsg string `json:"err_msg"`
	Logid  string `json:"logid"`
	Data   struct {
		SuccessTopGiftIdList []string `json:"success_top_gift_id_list"`
	} `json:"data"`
}

type DyGiftInfo1 struct {
	ErrNo   int    `json:"err_no"`
	ErrTips string `json:"err_tips"`
	Data    struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	} `json:"data"`
}

// RunDyRoom 启动抖音的
func RunDyRoom(c *WebsocketAnchorClient, commUserInfo *CommUserInfo) bool {
	// 获取  access_token
	// 获取  access_token
	token, ok := GetDyToken(commUserInfo.AppId, commUserInfo.AppSecret)
	if !ok {
		return false
	}
	c.Access_token = token

	// 获取抖音的roomId
	signParamsMap := make(map[string]interface{})
	signParamsMap["token"] = c.Token
	data, err := json.Marshal(signParamsMap)
	if err != nil {
		zlog.Logger.Error().Msgf("Marshal err:%v", err)
		return false
	}
	res, err := SendPostRequesDy2("https://webcast.bytedance.com/api/webcastmate/info", data, c.Access_token)
	if err != nil {
		zlog.Logger.Error().Msgf("SendPostRequesDy err:%v", err)
		return false
	}

	zlog.Logger.Info().Msgf("获取抖音的roomId res:%v data:%v", string(res), res)

	dyGetInfoRes := &DyGetInfoRes{}
	err = json.Unmarshal(res, dyGetInfoRes)
	if err != nil {
		zlog.Logger.Error().Msgf("Marshal err:%v %v", err, res)
		return false
	}
	zlog.Logger.Info().Msgf("获取抖音的信息SendPostRequesDy res:%v data:%v", string(res), res)
	if dyGetInfoRes.Errcode > 0 {
		zlog.Logger.Error().Msgf("  dyGetInfoRes:%v", dyGetInfoRes)
		return false
	}
	commUserInfo.UserName = dyGetInfoRes.Data.Info.NickName
	commUserInfo.Id = dyGetInfoRes.Data.Info.AnchorOpenId
	commUserInfo.RoomId = fmt.Sprintf("%v", dyGetInfoRes.Data.Info.RoomId)
	commUserInfo.UserName = dyGetInfoRes.Data.Info.NickName
	commUserInfo.HeadUrl = dyGetInfoRes.Data.Info.AvatarUrl

	//2.启动推送 评论
	signParamsMap = make(map[string]interface{})
	signParamsMap["roomid"] = commUserInfo.RoomId
	signParamsMap["appid"] = commUserInfo.AppId
	signParamsMap["msg_type"] = "live_comment"
	data, err = json.Marshal(signParamsMap)
	if err != nil {
		zlog.Logger.Error().Msgf("Marshal err:%v", err)
		return false
	}
	res, err = SendPostRequesDy("https://webcast.bytedance.com/api/live_data/task/start", data, c.Access_token)
	if err != nil {
		zlog.Logger.Error().Msgf("SendPostRequesDy err:%v", err)
		return false
	}
	zlog.Logger.Info().Msgf("启动 评论推送SendPostRequesDy res:%v", string(res))

	dyTaskInfo := &DyTaskInfo{}
	err = json.Unmarshal(res, dyTaskInfo)
	if err != nil {
		zlog.Logger.Error().Msgf("Marshal err:%v", err)
		return false
	}
	if dyTaskInfo.ErrNo > 0 {
		zlog.Logger.Error().Msgf("  dyTaskInfo:%v", dyTaskInfo)
		return false
	}
	//3.启动推送 礼物
	signParamsMap = make(map[string]interface{})
	signParamsMap["roomid"] = commUserInfo.RoomId
	signParamsMap["appid"] = commUserInfo.AppId
	signParamsMap["msg_type"] = "live_gift"
	data, err = json.Marshal(signParamsMap)
	if err != nil {
		zlog.Logger.Error().Msgf("Marshal err:%v", err)
		return false
	}
	res, err = SendPostRequesDy("https://webcast.bytedance.com/api/live_data/task/start", data, c.Access_token)
	if err != nil {
		zlog.Logger.Error().Msgf("SendPostRequesDy err:%v", err)
		return false
	}
	zlog.Logger.Info().Msgf("启动 礼物推送SendPostRequesDy res:%v", string(res))

	dyTaskInfo = &DyTaskInfo{}
	err = json.Unmarshal(res, dyTaskInfo)
	if err != nil {
		zlog.Logger.Error().Msgf("Marshal err:%v", err)
		return false
	}
	if dyTaskInfo.ErrNo > 0 {
		zlog.Logger.Error().Msgf("  dyTaskInfo:%v", dyTaskInfo)
		return false
	}
	//4.启动推送 点赞
	signParamsMap = make(map[string]interface{})
	signParamsMap["roomid"] = commUserInfo.RoomId
	signParamsMap["appid"] = commUserInfo.AppId
	signParamsMap["msg_type"] = "live_like"
	data, err = json.Marshal(signParamsMap)
	if err != nil {
		zlog.Logger.Error().Msgf("Marshal err:%v", err)
		return false
	}
	res, err = SendPostRequesDy("https://webcast.bytedance.com/api/live_data/task/start", data, c.Access_token)
	if err != nil {
		zlog.Logger.Error().Msgf("SendPostRequesDy err:%v", err)
		return false
	}
	zlog.Logger.Info().Msgf("启动点赞推送SendPostRequesDy res:%v", string(res))

	dyTaskInfo = &DyTaskInfo{}
	err = json.Unmarshal(res, dyTaskInfo)
	if err != nil {
		zlog.Logger.Error().Msgf("Marshal err:%v", err)
		return false
	}
	if dyTaskInfo.ErrNo > 0 {
		zlog.Logger.Error().Msgf("  dyTaskInfo:%v", dyTaskInfo)
		return false
	}

	// 5.启动置顶礼物
	signParamsMap = make(map[string]interface{})
	signParamsMap["room_id"] = commUserInfo.RoomId
	signParamsMap["app_id"] = commUserInfo.AppId
	signParamsMap["sec_gift_id_list"] = commUserInfo.TopGiftList
	data, err = json.Marshal(signParamsMap)
	if err != nil {
		zlog.Logger.Error().Msgf("Marshal err:%v", err)
		return false
	}
	res, err = SendPostRequesDy3("https://webcast.bytedance.com/api/gift/top_gift", data, c.Access_token)
	if err != nil {
		zlog.Logger.Error().Msgf("SendPostRequesDy err:%v", err)
		return false
	}
	dyGiftInfo := &DyGiftInfo{}
	err = json.Unmarshal(res, dyGiftInfo)
	if err != nil {
		zlog.Logger.Error().Msgf("Marshal err:%v", err)
		return false
	}
	if dyTaskInfo.ErrNo > 0 {
		zlog.Logger.Error().Msgf("  dyGiftInfo:%v", dyGiftInfo)
		return false
	}
	zlog.Logger.Info().Msgf("启动置顶礼物 signParamsMap%v TopGiftList:%v dyGiftInfo%v res:%v", signParamsMap, commUserInfo.TopGiftList, dyGiftInfo, string(res))
	return true
}

func SendPostRequesDy1(url string, data []byte) ([]byte, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		zlog.Logger.Error().Msgf("NewRequest%v", err)
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	response, err := client.Do(request)
	if err != nil {
		zlog.Logger.Error().Msgf("http Do%v", err)
		return nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		zlog.Logger.Error().Msgf("ReadAll%v", err)
		return nil, err
	}
	return body, nil
}

// 发送post请求 application/json
func SendPostRequesDy2(url string, data []byte, token string) ([]byte, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		zlog.Logger.Error().Msgf("NewRequest%v", err)
		return nil, err
	}
	// 设置请求头
	request.Header.Set("X-Token", token)
	request.Header.Set("content-Type", "application/json;charset=UTF-8")
	response, err := client.Do(request)
	if err != nil {
		zlog.Logger.Error().Msgf("http Do%v", err)
		return nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		zlog.Logger.Error().Msgf("ReadAll%v", err)
		return nil, err
	}
	return body, nil
}

// 发送post请求 application/json
func SendPostRequesDy3(url string, data []byte, token string) ([]byte, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		zlog.Logger.Error().Msgf("NewRequest%v", err)
		return nil, err
	}
	// 设置请求头
	request.Header.Set("x-token", token)
	request.Header.Set("content-Type", "application/json;charset=UTF-8")
	response, err := client.Do(request)
	if err != nil {
		zlog.Logger.Error().Msgf("http Do%v", err)
		return nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		zlog.Logger.Error().Msgf("ReadAll%v", err)
		return nil, err
	}
	return body, nil
}

// 发送post请求 application/json
func SendPostRequesDy(url string, data []byte, token string) ([]byte, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		zlog.Logger.Error().Msgf("NewRequest%v", err)
		return nil, err
	}
	// 设置请求头
	request.Header.Set("access-token", token)
	request.Header.Set("content-Type", "application/json;charset=UTF-8")
	response, err := client.Do(request)
	if err != nil {
		zlog.Logger.Error().Msgf("http Do%v", err)
		return nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		zlog.Logger.Error().Msgf("ReadAll%v", err)
		return nil, err
	}
	return body, nil
}

var Token = ""
var LastGetTokenTime int64 = 0

func GetDyToken(AppId string, AppSecret string) (string, bool) {
	type DyBSGiftInfo1 struct {
		ErrNo   int    `json:"err_no"`
		ErrTips string `json:"err_tips"`
		Data    struct {
			AccessToken string `json:"access_token"`
			ExpiresIn   int    `json:"expires_in"`
		} `json:"data"`
	}
	now := time.Now().Unix()
	if Token != "" && now-LastGetTokenTime < 60 {
		return Token, true
	}
	LastGetTokenTime = now
	// 获取  access_token
	signParamsMap := make(map[string]interface{})
	signParamsMap["appid"] = AppId
	signParamsMap["secret"] = AppSecret
	signParamsMap["grant_type"] = "client_credential"
	data, err := json.Marshal(signParamsMap)
	if err != nil {
		zlog.Logger.Error().Msgf("Marshal err:%v", err)
		return "", false
	}
	res, err := XBSSendPostRequesDy1("https://developer.toutiao.com/api/apps/v2/token", data)
	if err != nil {
		zlog.Logger.Error().Msgf("SendPostRequesDy err:%v", err)
		return "", false
	}
	dyBSGiftInfo1 := &DyBSGiftInfo1{}
	err = json.Unmarshal(res, dyBSGiftInfo1)
	if err != nil {
		zlog.Logger.Error().Msgf("Marshal err:%v", err)
		return "", false
	}
	if dyBSGiftInfo1.ErrNo > 0 {
		zlog.Logger.Error().Msgf("  dyBSGiftInfo1:%v  %v", dyBSGiftInfo1, string(res))
		return "", false
	}
	zlog.Logger.Info().Msgf("获取抖音的access_token SendPostRequesDy res:%v %v", string(res), dyBSGiftInfo1)
	Token = dyBSGiftInfo1.Data.AccessToken
	return dyBSGiftInfo1.Data.AccessToken, true
}
func XBSSendPostRequesDy1(url string, data []byte) ([]byte, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		zlog.Logger.Error().Msgf("NewRequest%v", err)
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	response, err := client.Do(request)
	if err != nil {
		zlog.Logger.Error().Msgf("http Do%v", err)
		return nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		zlog.Logger.Error().Msgf("ReadAll%v", err)
		return nil, err
	}
	return body, nil
}
