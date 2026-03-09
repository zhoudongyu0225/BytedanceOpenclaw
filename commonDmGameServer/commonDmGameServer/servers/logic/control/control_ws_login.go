package control

import (
	"dmGameServer/common"
	"dmGameServer/model"
	pb "dmGameServer/pb"
	"dmGameServer/untils"
	"dmGameServer/zlog"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/golang/protobuf/proto"
	"strings"
	"time"
)

type KuaiShouLogin struct {
	Result   int    `json:"result"`
	ErrorMsg string `json:"errorMsg"`
	UserInfo struct {
		AuthorOpenId  string `json:"authorOpenId"`
		UserName      string `json:"userName"`
		HeadUrl       string `json:"headUrl"`
		LiveStreamUrl string `json:"liveStreamUrl"`
	} `json:"userInfo"`
}

// 版本是否够ok
func VersionIsOk(reqVersion, version string) bool {
	if reqVersion == "" {
		zlog.Logger.Error().Msgf("VersionIsOk err:%v", "reqVersion is nil")
		return false
	}
	// 服务器没有配置版本号
	if version == "" {
		return true
	}
	if reqVersion == version {
		return true
	}
	reqVersionList := strings.Split(reqVersion, ".")
	versionList := strings.Split(version, ".")
	if len(reqVersionList) != len(versionList) {
		return false
	}
	// 版本号都是大于
	for i, v := range reqVersionList {
		if v > versionList[i] {
			return true
		}
		if v < versionList[i] {
			return false
		}
	}
	return true
}

// 主播信息
type CommUserInfo struct {
	AppId       string
	AppSecret   string
	UserName    string `json:"userName"`
	HeadUrl     string `json:"headUrl"`
	Id          string `json:"id"`
	RoomId      string
	TopGiftList []string `json:"topGiftList"` // 礼物置顶
}

// 提示
func Tips(c *WebsocketAnchorClient, msg string) {
	c.WsSend(COMMAND_TIPS_S, &pb.TipsNotify{
		TipMsg: msg,
	})
}

// 更加类型提示
func TipsType(c *WebsocketAnchorClient, msg string, TipsType pb.TipType) {
	c.WsSend(COMMAND_TIPS_S, &pb.TipsNotify{
		TipMsg: msg,
		Type:   TipsType,
	})
}

// 登录
func WsLogin(c *WebsocketAnchorClient, weMessage *WbMessageModel) {
	// 登录
	s2c := &pb.LoginResp{}
	// 1.解析数据
	loginReq := &pb.LoginReq{}
	err := proto.Unmarshal(weMessage.Body, loginReq)
	if err != nil {
		zlog.Logger.Error().Msgf("Unmarshal err:%v [%v]", err, loginReq.AccountId)
		s2c.ErrMsg = err.Error()
		c.WsSend(COMMAND_LOGIN_S2C, s2c)
		// 关闭读写f
		c.CloseReadAndWirte()
		return
	}
	// 2.从客户端获取app信息
	appInfo := model.GetAppConfigInfo(fmt.Sprintf("%v.%v.%v", loginReq.GameId, loginReq.PlatformId,
		loginReq.ModeId))
	zlog.Logger.Info().Msgf("主播客户端登录AccountId[%v] data[%+v] appInfo[%+v]", loginReq.AccountId, loginReq, appInfo)
	// 版本号判断
	if !VersionIsOk(loginReq.Version, appInfo.Version) {
		zlog.Logger.Error().Msgf("版本号不对[%v][%v]", loginReq.Version, appInfo.Version)
		s2c.ErrMsg = "版本号过低，请更新版本"
		c.WsSend(COMMAND_LOGIN_S2C, s2c)
		// 关闭读写
		c.CloseReadAndWirte()
		return
	}
	// 绑定appInfo的信息
	c.AppInfo = appInfo
	// 创建快手的信息
	k := &KuiShouOfficialInfo{
		AppId:       appInfo.AppId,
		AppSecret:   appInfo.Secret,
		AccessToken: "",
		RoomCode:    loginReq.RoomId,
		UserInfo:    &KuiShouUserInfo{},
	}
	// 通用主播信息
	commUserInfo := &CommUserInfo{
		AppId:       appInfo.AppId,
		AppSecret:   appInfo.Secret,
		UserName:    "",
		HeadUrl:     "",
		Id:          loginReq.AccountId,
		RoomId:      loginReq.RoomId,
		TopGiftList: appInfo.TopGiftList,
	}
	// 3.游戏模式判断
	switch loginReq.ModeId {
	case int32(pb.ModeId_GuanYou): // 官游逻辑
		if !GuanGame(c, loginReq, k, commUserInfo) {
			s2c.ErrMsg = "登录游戏失败"
			c.WsSend(COMMAND_LOGIN_S2C, s2c)
			// 关闭读写
			c.CloseReadAndWirte()
			return
		}
	case int32(pb.ModeId_YeYou): // 野游逻辑
		if !YeYouGameLogin(c, loginReq, k, commUserInfo) {
			s2c.ErrMsg = "登录游戏失败"
			c.WsSend(COMMAND_LOGIN_S2C, s2c)
			// 关闭读写
			c.CloseReadAndWirte()
			return
		}
	// 不处理
	default:
		zlog.Logger.Error().Msgf("登录游戏类型错误:%v [%v]", loginReq.AccountId, loginReq.ModeId)
		// ws通知
		s2c.ErrMsg = "登录游戏类型错误"
		c.WsSend(COMMAND_LOGIN_S2C, s2c)
		// 关闭读写
		c.CloseReadAndWirte()
		return
	}

	// IsClientLogOn 客户端是否登录
	ok, oldWs := GetLogicMgr().IsClientLogOn(loginReq.AccountId)
	if ok {
		untils.TapErr(fmt.Sprintf("重复登陆:%v", loginReq.AccountId))
		// ws通知
		s2c.ErrMsg = "登录被挤掉"
		c.WsSend(COMMAND_LOGIN_S2C, s2c)
		// 主动关闭链接
		oldWs.Close(true)
		s2c.ErrMsg = "存在重复登陆,请重新登录"
		zlog.Logger.Error().Msgf("重复登陆两个链接都关闭[%+v]", loginReq.AccountId)
		// 关闭读写
		c.CloseReadAndWirte()
		return
	}

	// 开启二次验证
	if appInfo.EnableSecondAuth {
		if !model.IsWhite(loginReq.AccountId, loginReq.GameId) {
			zlog.Logger.Error().Msgf("账号不在白名单%v", loginReq.AccountId)
			// ws通知
			s2c.ErrMsg = "账号没有权限"
			c.WsSend(COMMAND_LOGIN_S2C, s2c)
			return
		}
	}

	// 账号封号
	// 账号封号
	if model.IsBlack(loginReq.AccountId) {
		zlog.Logger.Error().Msgf("账号在黑名单%v", loginReq.AccountId)
		s2c.ErrMsg = "账号已经被封号"
		c.WsSend(COMMAND_LOGIN_S2C, s2c)
		return
	}
	zlog.Logger.Info().Msgf(" WsLogin 登录[%+v]", loginReq)
	// 设置ws登录 Socket *WebsocketAnchorClient, anchor *model.Anchor, KuiShouInfo *KuiShouOfficialInfo
	anchorV, _ := GetAnchorById(loginReq.AccountId)
	// 填充
	c.AnchorV = anchorV
	if !GetLogicMgr().SetClientLogOn(loginReq.AccountId, c, anchorV, k, loginReq) {
		zlog.Logger.Error().Msgf(" WsLogin 登录失败[%+v]", k)
		// ws通知
		s2c.ErrMsg = "登录失败"
		c.WsSend(COMMAND_LOGIN_S2C, s2c)
		return
	}
	if commUserInfo.Id == "" {
		zlog.Logger.Error().Msgf(" WsLogin 登录失败[%+v]", commUserInfo.Id)
		s2c.ErrMsg = "登录失败"
		c.WsSend(COMMAND_LOGIN_S2C, s2c)
		return
	}
	if loginReq.GameId != int32(common.GameId) {
		zlog.Logger.Error().Msgf(" WsLogin 登录失败,游戏类型错误[%+v] %v", commUserInfo.Id, loginReq.GameId)
		// ws通知
		s2c.ErrMsg = "登录失败,游戏类型错误"
		c.WsSend(COMMAND_LOGIN_S2C, s2c)
		return
	}
	c.SdkId = untils.GetSdkId()
	s2c = &pb.LoginResp{
		RoomId:     loginReq.RoomId,
		ServerTime: int32(time.Now().Unix()),
		AnchorDB:   anchorV,
	}

	//---------------登录成功后-------

	// 排序 活跃度
	// ws通知
	c.WsSend(COMMAND_LOGIN_S2C, s2c)

	// 获取平台开放的数据
	model.UpdateSlgCrossPlatformOpen()
	//----------------飞书群用于--统计数据-----------------
	if anchorV != nil {
		// 添加主播数据
		AddTapMap(anchorV.PlatformId, anchorV.AccountId)
		//  主播开播
		untils.ZBKPoss(c.AccountId, anchorV.PlatformId, anchorV.GameId, anchorV.NickName, GetTapMapNum(), loginReq.RoomId, c.Version,
			model.GetLiveInfo(c.AccountId, c.AnchorV.PlatformId, c.GameId))
	}
}

// 官游
func GuanGame(c *WebsocketAnchorClient, loginReq *pb.LoginReq, k *KuiShouOfficialInfo, commUserInfo *CommUserInfo) bool {
	switch loginReq.PlatformId {
	case int32(pb.PlatformId_Tiktok): // Tiktok
	case int32(pb.PlatformId_ShiPinHao): // 视频号
	case int32(pb.PlatformId_KuaiShou): // 快手
		zlog.Logger.Info().Msgf("快手登录RoomId[%v] data[%+v]", loginReq.RoomId, loginReq)
		if !k.Init() {
			zlog.Logger.Error().Msgf("快手登录RoomId[%v] data[%+v] Init err", loginReq.RoomId, loginReq)
			// ws通知
			Tips(c, "快手登录失败")
			return false
		}
		// 获取直播间信息
		signParamsMap := make(map[string]interface{})
		signParamsMap["roomCode"] = k.RoomCode
		signParamsMap["timestamp"] = time.Now().UnixMilli()
		signParamsMap["sign"] = CalcSign(signParamsMap, k)
		data, err := json.Marshal(signParamsMap)
		if err != nil {
			zlog.Logger.Error().Msgf("Marshal err:%v", err)
			Tips(c, "快手登录失败")
			return false
		}
		res, err := untils.SendPostRequest(k.GetUrl("https://open.kuaishou.com/openapi/developer/live/data/interactive/live/stream/info"), data)
		if err != nil {
			zlog.Logger.Error().Msgf("SendPostRequest err:%v", err)
			Tips(c, "快手登录失败")
			return false
		}
		kuaiShouLogin := &KuaiShouLogin{}
		err = json.Unmarshal(res, kuaiShouLogin)
		if err != nil {
			zlog.Logger.Error().Msgf("Marshal err:%v", err)
			Tips(c, "快手登录失败")
			return false
		}
		if kuaiShouLogin.Result != 1 {
			zlog.Logger.Error().Msgf("commonRes code%v err:%v", kuaiShouLogin.Result, kuaiShouLogin.ErrorMsg)
			Tips(c, "快手登录失败")
			return false
		}
		zlog.Logger.Debug().Msgf("快手获取直播间信息成功:%+v", kuaiShouLogin)
		// 绑定主播AccountId
		loginReq.AccountId = kuaiShouLogin.UserInfo.AuthorOpenId
		if k.UserInfo == nil {
			k.UserInfo = &KuiShouUserInfo{}
		}
		k.UserInfo.UserName = kuaiShouLogin.UserInfo.UserName
		k.UserInfo.HeadUrl = kuaiShouLogin.UserInfo.HeadUrl
		k.UserInfo.AuthorOpenId = kuaiShouLogin.UserInfo.AuthorOpenId

		commUserInfo.UserName = kuaiShouLogin.UserInfo.UserName
		commUserInfo.Id = kuaiShouLogin.UserInfo.AuthorOpenId
		commUserInfo.HeadUrl = k.UserInfo.HeadUrl
	// --------------
	case int32(pb.PlatformId_BZhan): // B站
	case int32(pb.PlatformId_DouYin): // 抖音
		if loginReq.Token != "" {
			c.Token = loginReq.Token
			isOK := false
			isOK = RunDyRoom(c, commUserInfo)
			if !isOK {
				zlog.Logger.Error().Msgf("抖音登录失败[%v]", loginReq.Token)
				Tips(c, "抖音登录失败")
				return false
			}
			loginReq.RoomId = commUserInfo.RoomId
		}
	case int32(pb.PlatformId_TkHh): //
	case int32(pb.PlatformId_DyHh): //
		// 是否有注册
		aa, _ := model.GetYeYouInfoById(loginReq.AccountId)
		if aa == nil {
			zlog.Logger.Error().Msgf("抖音野游失败[%v]", loginReq.AccountId)
			Tips(c, "没有账号")
			return false
		}
		aa.Name = untils.GetSensitiveFilter().ReplaceSensitiveWords(aa.Name, 0)
		//// 敏感词校验
		commUserInfo.UserName = aa.Name
	default:
		zlog.Logger.Error().Msgf("登录游戏平台错误:%v [%v]", loginReq.AccountId, loginReq.PlatformId)
		Tips(c, "登录游戏平台错误")
		return false
	}
	// --------------!!!!!!!!!必须要有loginReq.AccountId的值!!!!!!!!!!!!!!----------------
	loginReq.AccountId = commUserInfo.Id
	if loginReq.AccountId == "" {
		zlog.Logger.Error().Msgf("登录游戏错误:%v  AccountId == nil", loginReq.AccountId)
		// ws通知
		Tips(c, "登录错误")
		return false
	}
	// 官游统一处理
	// 判断用户是否存在数据库里面不
	anchorV, _ := GetAnchorById(loginReq.AccountId)
	if anchorV == nil || anchorV.AccountId == "" {
		anchorV = &pb.AnchorDBInfo{
			AccountId:  commUserInfo.Id,
			NickName:   commUserInfo.UserName,
			GameId:     loginReq.GameId,
			PlatformId: loginReq.PlatformId,
			ModeId:     loginReq.ModeId,
			HeadUrl:    commUserInfo.HeadUrl,
		}
		// 添加主播到数据库
		if err := model.AddAnchor(anchorV); err != nil {
			zlog.Logger.Error().Msgf("添加主播到数据库失败 err:%v", err)
			Tips(c, "玩家登录失败")
			return false
		}
		zlog.Logger.Info().Msgf("官游账号注册成功 账号:%v 平台:%v", loginReq.AccountId, loginReq.PlatformId)
		// [官游]主播新增
		addOverviewOfTransactions := &model.OverviewOfTransactions{
			NewAnchorNum: 1,
		}
		// 更新流水总览
		model.UpdateOverviewOfTransactions(addOverviewOfTransactions, fmt.Sprintf(model.CollectionOverviewOfTransactions, loginReq.GameId, loginReq.PlatformId, loginReq.ModeId))
	} else {
		// 官游更新主播
		anchorV.AccountId = commUserInfo.Id
		if commUserInfo.UserName != "" {
			anchorV.NickName = commUserInfo.UserName
		}
		if commUserInfo.HeadUrl != "" {
			anchorV.HeadUrl = commUserInfo.HeadUrl
		}
		// 更新主播到数据库
		model.UpdateAnchor(anchorV)
	}
	return true
}

// 野游
func YeYouGameLogin(c *WebsocketAnchorClient, loginReq *pb.LoginReq, k *KuiShouOfficialInfo, commUserInfo *CommUserInfo) bool {

	anchorV, _ := GetAnchorById(loginReq.AccountId)
	if anchorV == nil || anchorV.AccountId == "" {
		anchorV = &pb.AnchorDBInfo{
			AccountId:  commUserInfo.Id,
			NickName:   commUserInfo.UserName,
			GameId:     loginReq.GameId,
			PlatformId: loginReq.PlatformId,
			ModeId:     loginReq.ModeId,
			HeadUrl:    commUserInfo.HeadUrl,
		}
		if anchorV.NickName == "" {
			anchorV.NickName = commUserInfo.Id
		}
		// 添加主播到数据库
		if err := model.AddAnchor(anchorV); err != nil {
			zlog.Logger.Error().Msgf("添加主播到数据库失败 err:%v", err)
			Tips(c, "玩家登录失败")
			return false
		}
		zlog.Logger.Info().Msgf("官游账号注册成功 账号:%v 平台:%v", loginReq.AccountId, loginReq.PlatformId)
		// [官游]主播新增
		addOverviewOfTransactions := &model.OverviewOfTransactions{
			NewAnchorNum: 1,
		}
		// 更新流水总览
		model.UpdateOverviewOfTransactions(addOverviewOfTransactions, fmt.Sprintf(model.CollectionOverviewOfTransactions, loginReq.GameId, loginReq.PlatformId, loginReq.ModeId))
	} else {
		// 官游更新主播
		anchorV.AccountId = commUserInfo.Id
		if commUserInfo.UserName != "" {
			anchorV.NickName = commUserInfo.UserName
		}
		if commUserInfo.HeadUrl != "" {
			anchorV.HeadUrl = commUserInfo.HeadUrl
		}
		// 更新主播到数据库
		model.UpdateAnchor(anchorV)
	}

	return true
}
