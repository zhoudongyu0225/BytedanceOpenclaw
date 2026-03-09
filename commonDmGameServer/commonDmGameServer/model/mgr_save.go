package model

import (
	pb "dmGameServer/pb"
	"dmGameServer/untils"
	"dmGameServer/zlog"
	"runtime/debug"
	"sync"
	"time"
)

var saveMgr *SaveMgr

const (
	MinT     = 1000 // 最小时间
	MaxT     = 3000 // 最大时间
	MaxDBNum = 500  // 批量存储最大数量
	SaveTime = 5    // 存储去找到的毫秒时间
)

var IsStop = false

// 保存的主播信息的结构体
type SaveAnchor struct {
	Anchor     *pb.AnchorDBInfo
	UpdateTime int64 // 需要存储的时间
}

// 保存的主播信息的结构体
type SavePlayer struct {
	OpenVo     *pb.OpenVo
	UpdateTime int64 // 需要存储的时间
}

type SaveMgr struct {
	// 是否停止接受
	SaveStop chan bool
	// 保存的主播信息的map
	SaveAnchorMap sync.Map // key id,value SaveAnchor
	// 保存的玩家信息的map
	SaveSlgPlayerMap sync.Map // key id ,value SavePlayer
}

// 添加主播存储
func UpdateAnchor(anchor *pb.AnchorDBInfo) {
	CachedTimeMilli := time.Now().UnixMilli()
	// 存在就不存了
	if save, ok := saveMgr.SaveAnchorMap.Load(anchor.AccountId); ok {
		if save.(*SaveAnchor) == nil {
			untils.TapErr("存储中心存储主播信息为空")
			return
		}
		// 获取时间
		updateTime := save.(*SaveAnchor).UpdateTime
		saveMgr.SaveAnchorMap.Store(anchor.AccountId, &SaveAnchor{
			Anchor:     anchor,
			UpdateTime: updateTime,
		})
	} else {
		// 保存主播信息
		saveMgr.SaveAnchorMap.Store(anchor.AccountId, &SaveAnchor{
			Anchor:     anchor,
			UpdateTime: CachedTimeMilli + int64(untils.GenerateRandomNumber(MinT, MaxT)),
		})
	}
}

// 存储主播玩家局外的
func UpdateOpenVo(tOpenVoList []*pb.OpenVo) bool {
	for _, Openvo := range tOpenVoList {
		if Openvo == nil {
			continue
		}
		if Openvo.OpenId == "" {
			zlog.Logger.Info().Msgf("存储中心存储玩家id信息为空")
			continue
		}
		CachedTimeMilli := time.Now().UnixMilli()
		// 存在就不存了
		key := Openvo.OpenId
		// 存在就不存了
		if save, ok := saveMgr.SaveSlgPlayerMap.Load(key); ok {
			if save.(*SavePlayer) == nil {
				untils.TapErr("存储中心存储玩家信息为空")
				continue
			}
			// 获取时间
			updateTime := save.(*SavePlayer).UpdateTime
			saveMgr.SaveSlgPlayerMap.Store(key, &SavePlayer{
				OpenVo:     Openvo,
				UpdateTime: updateTime,
			})
		} else {
			// 保存玩家信息
			saveMgr.SaveSlgPlayerMap.Store(key, &SavePlayer{
				OpenVo:     Openvo,
				UpdateTime: CachedTimeMilli + int64(untils.GenerateRandomNumber(MinT, MaxT)),
			})
		}
	}
	return true

}

func DBCore() {
	defer func() {
		if err := recover(); err != nil {
			stack := debug.Stack()
			untils.PanicPoss(err, stack)
			untils.Go2(DBCore)
		}
	}()
	//  等待多少毫秒存储
	ticker1 := time.NewTicker(time.Millisecond * SaveTime)
	defer ticker1.Stop()
	t := 0
	for {
		select {
		case <-ticker1.C: // 存储----
			t++
			t = t % 2
			switch t {
			case 0:
				k := 0
				CachedTimeMilli := time.Now().UnixMilli()
				saveList := make([]*pb.AnchorDBInfo, 0)
				// 查看主播的数据
				saveMgr.SaveAnchorMap.Range(func(key, value interface{}) bool {
					save := value.(*SaveAnchor)
					if save == nil {
						return true
					}
					// 服务器停服了 直接存
					if IsStop {
						save.UpdateTime = 0
					}
					if CachedTimeMilli > save.UpdateTime {
						k++
						saveList = append(saveList, save.Anchor)
						// 删掉
						saveMgr.SaveAnchorMap.Delete(key)
					}
					return true
				})

				// 保存
				for _, v := range saveList {
					updateAnchorDB(v)
				}
			case 1:
				// --玩家--
				CachedTimeMilli := time.Now().UnixMilli()
				updateList := make([]*pb.OpenVo, 0)
				saveMgr.SaveSlgPlayerMap.Range(func(key, value interface{}) bool {
					save := value.(*SavePlayer)
					if save == nil {
						saveMgr.SaveSlgPlayerMap.Delete(key)
						return true
					}
					// 服务器停服了 直接存
					if IsStop {
						save.UpdateTime = 0
					}
					// 跳过存储太多了
					if len(updateList) > MaxDBNum {
						zlog.Logger.Info().Msgf("这个时间要存的数据太多跳过")
						return false
					}
					if CachedTimeMilli > save.UpdateTime {
						// 保存玩家信息
						updateList = append(updateList, save.OpenVo)
						saveMgr.SaveSlgPlayerMap.Delete(key)
					}
					return true
				})
				if len(updateList) > 0 {
					GetPlayerMgr().updateOpenVoDB(updateList)
				}
			}

		case <-saveMgr.SaveStop:
			return
		}
	}
}

// 注意
func InitSaveMgr() {
	saveMgr = &SaveMgr{
		SaveStop: make(chan bool, 2),
	}
	untils.Go2(DBCore)
}
