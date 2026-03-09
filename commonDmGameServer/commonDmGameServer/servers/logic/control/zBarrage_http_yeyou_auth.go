package control

import (
	"dmGameServer/common"
	"dmGameServer/model"
	pb "dmGameServer/pb"
	"dmGameServer/zlog"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"google.golang.org/protobuf/runtime/protoimpl"
	"io"
	"net/http"
)

// Register 注册
func Register(c *gin.Context) {
	// 解析 JSON 数据
	data := &pb.AccountCreateReq{}
	// 获取请求体中的原始数据
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		zlog.Logger.Error().Msgf("ReadAll%v", err)
		Fail(c, "", "Failed to read request body")
		return
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		zlog.Logger.Error().Msgf("Unmarshal%v", err)
		// 处理解析 JSON 失败的错误
		Fail(c, "", "Failed to parse JSON data")
		return
	}
	if data.AccountId == "" {
		zlog.Logger.Error().Msgf("账号不能为空 :%v", data.AccountId)
		Fail(c, "", "账号不能为空")
		return
	}
	// 判断用户是否存在
	v, _ := GetAnchorById(data.AccountId, data.GameId, 0)
	if v != nil {
		zlog.Logger.Error().Msgf("账号已注册 :%v", data.AccountId)
		Fail(c, "", "账号已注册")

		return
	}
	if data.AdminPassword != common.AdminPassWard {
		zlog.Logger.Error().Msgf("管理员密码不对 :%v", data.AdminPassword)
		Fail(c, "", "管理员密码不正确")
		return
	}

	anchor := &pb.AnchorDBInfo{
		AccountId:  data.AccountId,
		NickName:   data.NickName,
		GameId:     data.GameId,
		PlatformId: data.PlatformId,
		ModeId:     data.ModeId,
	}
	if err := model.AddAnchor(anchor); err != nil {
		zlog.Logger.Error().Msgf("Register err:%v", err)
		Response(c, http.StatusUnprocessableEntity, 442, nil, err.Error())
		return
	}

	// [野游]主播新增
	addOverviewOfTransactions := &model.OverviewOfTransactions{
		NewAnchorNum: 1,
	}
	model.UpdateOverviewOfTransactions(addOverviewOfTransactions, fmt.Sprintf(model.CollectionOverviewOfTransactions, data.GameId, data.PlatformId, data.ModeId))

	zlog.Logger.Info().Msgf("野游账号注册成功 :%v", data.AccountId)
	Success(c, "", "", "注册成功")

}

// 弹幕登录请求
type ClientLoginReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 数据前缀
	DataFixVo *ClientLoginDataFixVo `protobuf:"bytes,1,opt,name=dataFixVo,proto3" json:"dataFixVo,omitempty"`
	// 账号
	AccountId string `protobuf:"bytes,2,opt,name=accountId,proto3" json:"accountId,omitempty"`
	// 版本
	Version string `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
	// 直播间绑定码 (野游发直播路径，官方发房间id)
	RoomId string `protobuf:"bytes,4,opt,name=roomId,proto3" json:"roomId,omitempty"`
	// token
	Token string `protobuf:"bytes,5,opt,name=token,proto3" json:"token,omitempty"`
}
type ClientLoginDataFixVo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 游戏ID
	GameId int32 `protobuf:"varint,1,opt,name=gameId,proto3" json:"gameId,omitempty"`
	// 平台
	PlatformId int32 `protobuf:"varint,2,opt,name=platformId,proto3" json:"platformId,omitempty"`
	// 地区
	RegionId int32 `protobuf:"varint,3,opt,name=regionId,proto3" json:"regionId,omitempty"`
	// 游戏模式/官游/野游
	ModeId int32 `protobuf:"varint,4,opt,name=modeId,proto3" json:"modeId,omitempty"`
	// 区服
	PartitionId int32 `protobuf:"varint,5,opt,name=partitionId,proto3" json:"partitionId,omitempty"`
}

// Login 登录
func Login(c *gin.Context) {
	// 解析 JSON 数据
	yeYouLoginReq := &ClientLoginReq{}
	// 获取请求体中的原始数据
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		zlog.Logger.Error().Msgf("ReadAll%v", err)
		Fail(c, "", "Failed to read request body")
		return
	}
	err = json.Unmarshal(body, &yeYouLoginReq)
	if err != nil {
		zlog.Logger.Error().Msgf("Unmarshal%v", err)
		// 处理解析 JSON 失败的错误
		Fail(c, "", "Failed to parse JSON data")

		return
	}
	zlog.Logger.Info().Msgf("野游工具账号登录 :%v", yeYouLoginReq.AccountId)
	// 判断用户是否存在
	v, _ := GetAnchorById(yeYouLoginReq.AccountId, yeYouLoginReq.DataFixVo.GameId)
	if v == nil {
		zlog.Logger.Error().Msgf("账号号不存在%v", yeYouLoginReq.AccountId)
		// Response(c, http.StatusUnprocessableEntity, 442, nil, "账号号不存在")
		Fail(c, "", "账号号不存在")
		return
	}
	if yeYouLoginReq.DataFixVo.GameId != v.GameId {
		str := fmt.Sprintf("游戏数据匹配不上[%v][%v][%v]", yeYouLoginReq.AccountId, yeYouLoginReq.DataFixVo.GameId, v.GameId)
		zlog.Logger.Error().Msg(str)
		Fail(c, "", "游戏数据gameId匹配不上")
		return
	}
	if yeYouLoginReq.DataFixVo.ModeId != v.ModeId {
		str := fmt.Sprintf(" 游戏模式数据匹配不上[%v][%v][%v]", yeYouLoginReq.AccountId, yeYouLoginReq.DataFixVo.ModeId, v.ModeId)
		zlog.Logger.Error().Msg(str)
		Fail(c, "", "游戏模式数据匹配不上")
		return
	}
	if yeYouLoginReq.DataFixVo.PlatformId != v.PlatformId {
		str := fmt.Sprintf(" 平台数据匹配不上[%v][%v][%v]", yeYouLoginReq.AccountId, yeYouLoginReq.DataFixVo.PlatformId, v.PlatformId)
		zlog.Logger.Error().Msg(str)
		Fail(c, "", "平台数据匹配不上")
		return
	}

	appInfo := model.GetAppConfigInfo(fmt.Sprintf("%v.%v.%v.%v.%v", yeYouLoginReq.DataFixVo.GameId, yeYouLoginReq.DataFixVo.PlatformId,
		yeYouLoginReq.DataFixVo.RegionId, yeYouLoginReq.DataFixVo.ModeId, yeYouLoginReq.DataFixVo.PartitionId))
	// 开启二次验证
	if appInfo.EnableSecondAuth {
		if !model.IsWhite(yeYouLoginReq.AccountId, yeYouLoginReq.DataFixVo.GameId) {
			zlog.Logger.Error().Msgf("账号不在白名单%v", yeYouLoginReq.AccountId)
			Fail(c, "", "账号没有权限")
			return
		}
	}

	// 是否在黑名单
	if model.IsBlack(yeYouLoginReq.AccountId) {
		zlog.Logger.Error().Msgf("账号在黑名单%v", yeYouLoginReq.AccountId)
		Fail(c, "", "账号已经被封号")
		return
	}
	// 发送token // 工具类  (根据工具登录的 roomId做唯一)
	token, err, _ := GetToken(v.AccountId)
	if err != nil {
		zlog.Logger.Error().Msgf("token err:%v Account%v", err, yeYouLoginReq.AccountId)
		Fail(c, "", "系统token获取异常")
		return
	}

	GetMgr().SetYeyouToken(v.AccountId, token)

	c.Header("Authorization", token)
	yeYouDanMuLoginResp := &pb.YeYouDanMuLoginResp{
		Token:       token,
		ExpiredTime: 0, // 秒
	}
	zlog.Logger.Debug().Msgf("登录成功 :%v", yeYouDanMuLoginResp)
	// // 后面需要在http头信息带token
	SuccessLogin(c, "", yeYouLoginReq.AccountId, "登录成功", yeYouDanMuLoginResp.Token, yeYouDanMuLoginResp.ExpiredTime)
}
