# 2026-03-21 弹幕小游戏 control_barrage.go TODO 实现深度研究与代码实证

**学习时间**: 2026-03-21 18:29 (Asia/Shanghai)
**领域**: 抖音/快手/TikTok弹幕小游戏开发
**主题**: control_barrage.go 三大TODO项（弹幕/礼物/点赞）的代码级实现方案、规则引擎设计、游戏动作映射

---

## 一、现状分析：control_barrage.go 的三个 TODO

### 1.1 当前代码状态

```go
// ChatMessage 中的 TODO
// todo:弹幕逻辑

// GiftMessage 中的 TODO  
// todo:礼物逻辑

// LikeMessage 中的 TODO
// todo:点赞逻辑
```

三个 TODO 都在 `control_barrage.go` 中，各自处理一种用户互动消息。

### 1.2 上下文：谁在调用这些函数

```
convergeMessage() (zmgr_logic.go)
  ├── GiftMessage()   ← messageList.Type == "GiftMessage"
  ├── ChatMessage()   ← messageList.Type == "ChatMessage"  
  └── LikeMessage()   ← messageList.Type == "LikeMessage" (合并后每1秒触发)
```

每个函数接收 `*pb.MessageList` 和 `*WebsocketAnchorClient`（主播连接）。

### 1.3 已有基础设施

- **礼物配置**: `generateConfig/autoConfig.go` 中的 `GiftConfig`（Id, GiftNum, GetSoldierNum, RescueInjuryNum）
- **Buff配置**: `generateConfig/autoConfig.go` 中的 `BuffConfig`（BuffType, Int_Param, Duration等）
- **玩家数据**: `model/model_playerInfo.go` 中的 `PlayerMgr`，通过 `GetOpenVo(openId)` 获取
- **排行榜**: `model/model_rank_redis.go` 中的 `UpdateGameWeekRank()` / `UpdateGameMonthRank()`
- **Websocket发送**: `ws.WsSend(cmd, proto.Message)` 已实现

---

## 二、TODO #1：弹幕逻辑实现（ChatMessage）

### 2.1 设计目标

弹幕消息的处理需要：
1. **敏感词过滤**（防止违规内容）
2. **弹幕规则匹配**（关键词 → 游戏动作）
3. **分数/积分更新**（活跃度奖励）
4. **游戏动作触发**（如"进攻"→ 发射弹幕攻击）
5. **AI增强**（可选：意图分类 + 情感计算）

### 2.2 弹幕规则引擎设计

```go
// servers/logic/control/danmaku_rule_engine.go

package control

import (
    "dmGameServer/untils"
    "dmGameServer/zlog"
    "strings"
    "sync"
    "time"
)

// 弹幕动作类型
type DanmakuActionType int

const (
    ActionNone         DanmakuActionType = iota  // 无动作
    ActionAttack                              // 攻击
    ActionDefend                               // 防守
    ActionHeal                                 // 治疗
    ActionSpeedUp                              // 加速
    ActionSpeedDown                            // 减速
    ActionSpecial                              // 特殊技能
    ActionChat                                 // 仅聊天（无游戏效果）
)

// 单条弹幕规则
type DanmakuRule struct {
    Keywords     []string          // 匹配关键词（OR关系）
    Action       DanmakuActionType // 对应动作
    ScoreReward  int64             // 活跃度积分奖励
    CooldownMs   int64             // 冷却时间（毫秒），防止刷屏
    RequireLevel int               // 最低玩家等级要求
    Priority     int               // 规则优先级（越大越高）
}

// 弹幕规则引擎
type DanmakuRuleEngine struct {
    rules      []DanmakuRule      // 规则列表（按Priority降序）
    cooldowns  sync.Map           // openId -> lastTriggerTime
    mutex      sync.RWMutex
}

var danmakuRuleEngine *DanmakuRuleEngine

func InitDanmakuRuleEngine() {
    danmakuRuleEngine = &DanmakuRuleEngine{
        rules: []DanmakuRule{
            // === 攻击类 ===
            {
                Keywords:    []string{"进攻", "攻击", "打", "杀", "冲", "上", "干", "揍", "打他", "打它", "attack", "fire", "hit"},
                Action:      ActionAttack,
                ScoreReward: 2,
                CooldownMs:  3000,  // 3秒冷却
                Priority:    10,
            },
            // === 防守类 ===
            {
                Keywords:    []string{"防守", "防御", "守", "挡", "护", "盾", "defend", "shield", "protect"},
                Action:      ActionDefend,
                ScoreReward: 2,
                CooldownMs:  3000,
                Priority:    10,
            },
            // === 治疗类 ===
            {
                Keywords:    []string{"治疗", "奶", "加血", "回血", "heal", "health", "hp"},
                Action:      ActionHeal,
                ScoreReward: 2,
                CooldownMs:  5000,
                Priority:    10,
            },
            // === 加速类 ===
            {
                Keywords:    []string{"快", "加速", "冲鸭", "冲啊", "speed", "faster", "rush"},
                Action:      ActionSpeedUp,
                ScoreReward: 1,
                CooldownMs:  2000,
                Priority:    8,
            },
            // === 减速类（给敌方） ===
            {
                Keywords:    []string{"慢", "减速", "拖", "delay", "slow"},
                Action:      ActionSpeedDown,
                ScoreReward: 1,
                CooldownMs:  5000,
                Priority:    8,
            },
            // === 特殊技能 ===
            {
                Keywords:    []string{"大招", "绝招", "技能", "放大", "放技能", "skill", "ultimate"},
                Action:      ActionSpecial,
                ScoreReward: 5,
                CooldownMs:  15000,  // 15秒冷却
                Priority:    15,
            },
        },
    }
    zlog.Logger.Info().Msgf("弹幕规则引擎初始化完成，共 %d 条规则", len(danmakuRuleEngine.rules))
}

// 规则匹配
func (e *DanmakuRuleEngine) Match(text string) (DanmakuActionType, int64, bool) {
    text = strings.TrimSpace(text)
    if text == "" {
        return ActionNone, 0, false
    }

    // 预处理：小写化（兼容英文关键词）
    textLower := strings.ToLower(text)

    e.mutex.RLock()
    defer e.mutex.RUnlock()

    var bestRule *DanmakuRule
    for i := range e.rules {
        rule := &e.rules[i]
        for _, kw := range rule.Keywords {
            if strings.Contains(textLower, strings.ToLower(kw)) {
                if bestRule == nil || rule.Priority > bestRule.Priority {
                    bestRule = rule
                }
            }
        }
    }

    if bestRule == nil {
        return ActionChat, 1, false // 无游戏动作，但有聊天活跃积分
    }

    return bestRule.Action, bestRule.ScoreReward, true
}

// 冷却检查
func (e *DanmakuRuleEngine) CheckCooldown(openId string, cooldownMs int64) bool {
    now := time.Now().UnixMilli()
    if v, ok := e.cooldowns.Load(openId); ok {
        lastTime := v.(int64)
        if now-lastTime < cooldownMs {
            return false // 冷却中
        }
    }
    e.cooldowns.Store(openId, now)
    return true
}

// ProcessDanmaku 处理单条弹幕
// 返回：(动作类型, 是否触发游戏效果, 积分奖励)
func (e *DanmakuRuleEngine) ProcessDanmaku(openId string, text string) (DanmakuActionType, bool, int64) {
    // 1. 敏感词过滤
    if !untils.CheckSensitiveWord(text) {
        zlog.Logger.Debug().Msgf("弹幕含敏感词已过滤: %s", text)
        return ActionNone, false, 0
    }

    // 2. 规则匹配
    action, scoreReward, hasGameEffect := e.Match(text)

    // 3. 冷却检查（仅对有游戏效果的action）
    if hasGameEffect && !e.CheckCooldown(openId, e.getCooldown(action)) {
        zlog.Logger.Debug().Msgf("弹幕冷却中，跳过游戏动作: openId=%s action=%v", openId, action)
        return action, false, scoreReward // 仍给积分，但不触发游戏效果
    }

    return action, hasGameEffect, scoreReward
}

func (e *DanmakuRuleEngine) getCooldown(action DanmakuActionType) int64 {
    for i := range e.rules {
        if e.rules[i].Action == action {
            return e.rules[i].CooldownMs
        }
    }
    return 5000 // 默认5秒
}
```

### 2.3 弹幕动作 → 游戏帧协议

需要设计一个将弹幕动作转换为游戏帧广播的机制：

```go
// 弹幕动作结果
type DanmakuResult struct {
    OpenId      string
    NickName    string
    Action      DanmakuActionType
    ScoreReward int64
    Timestamp   int64
}

// BroadcastDanmakuAction 将弹幕动作广播给游戏客户端
func BroadcastDanmakuAction(ws *WebsocketAnchorClient, result *DanmakuResult) {
    // TODO: 这里需要定义一个 GameFrame 协议来发送动作
    // 格式示例：{ type: "barrage_action", openId: "xxx", action: "attack", timestamp: 1234567890 }
    
    // 目前通过 Tips 或自定义协议通知UE4客户端
    // 示例：发送一个Tips通知包含动作信息
    tips := &pb.TipsNotify{
        Tips: fmt.Sprintf("[%s] 发起动作: %v", result.NickName, result.Action),
    }
    ws.WsSend(COMMAND_TIPS_S, tips)
    
    zlog.Logger.Debug().Msgf("弹幕动作广播: openId=%s action=%v", result.OpenId, result.Action)
}

// 动作类型转字符串（用于显示和协议）
func (a DanmakuActionType) String() string {
    switch a {
    case ActionAttack:   return "attack"
    case ActionDefend:   return "defend"
    case ActionHeal:     return "heal"
    case ActionSpeedUp:  return "speed_up"
    case ActionSpeedDown: return "speed_down"
    case ActionSpecial:  return "special"
    default:             return "none"
    }
}
```

### 2.4 ChatMessage 完整实现

```go
// ChatMessage 完整实现（替换原TODO）
func ChatMessage(messageList *pb.MessageList, ws *WebsocketAnchorClient) {
    now := time.Now().Unix()
    msgJoin(messageList, ws)  // 先处理玩家加入

    for _, msg := range messageList.MsgList {
        zlog.Logger.Info().Msgf("玩家评论 %v %v %v %v", msg.Uid, msg.Name, msg.Content, ws.AccountId)
        
        if msg.Uid == "" {
            untils.TapErr(fmt.Sprintf("msg.Uid is 空 %v", msg))
            continue
        }

        // 获取/创建玩家数据
        pdb, _ := GetOpenVo(msg.Uid)
        if pdb == nil {
            untils.TapErr(fmt.Sprintf("pdb is nil %v", msg))
            continue
        }

        // 更新玩家基本信息
        if msg.Uid != "" {
            pdb.OpenId = msg.Uid
        }
        if msg.Name != "" {
            pdb.NickName = msg.Name
        }
        if msg.HeadImg != "" {
            pdb.AvatarUrl = msg.HeadImg
        }

        // 更新排名信息（每5秒一次）
        if now-pdb.LastUpdateTime > 5 {
            currWeekRank, currWeekScore := model.GetGameWeekRankByRankVo(pdb.OpenId, ws.AnchorV.PlatformId)
            currMonthRank, currMonthScore := model.GetGameMonthRankByRankVo(pdb.OpenId, ws.AnchorV.PlatformId)
            pdb.Score = currWeekScore
            pdb.MonthScore = currMonthScore
            pdb.Rank = int32(currWeekRank)
            pdb.MonthRank = int32(currMonthRank)
            pdb.LastUpdateTime = now
        }

        // ========== 【TODO:弹幕逻辑】完整实现 ==========
        
        // 1. 初始化规则引擎（延迟初始化避免循环依赖）
        if danmakuRuleEngine == nil {
            InitDanmakuRuleEngine()
        }

        // 2. 处理弹幕
        content := msg.Content
        if content == "" {
            continue
        }

        actionType, hasGameEffect, scoreReward := danmakuRuleEngine.ProcessDanmaku(msg.Uid, content)

        // 3. 更新活跃积分
        if scoreReward > 0 {
            // 更新周榜分数
            model.UpdateGameWeekRank(msg.Uid, scoreReward, ws.AnchorV.PlatformId)
            // 更新月榜分数
            model.UpdateGameMonthRank(msg.Uid, scoreReward, ws.AnchorV.PlatformId)
            
            // 更新玩家当前分数缓存
            pdb.Score += scoreReward
            
            zlog.Logger.Debug().Msgf("弹幕活跃积分: openId=%s content=%s action=%v score=%d",
                msg.Uid, content, actionType, scoreReward)
        }

        // 4. 触发游戏动作
        if hasGameEffect && ws.IsStart {
            result := &DanmakuResult{
                OpenId:      msg.Uid,
                NickName:    msg.Name,
                Action:      actionType,
                ScoreReward: scoreReward,
                Timestamp:   time.Now().UnixMilli(),
            }
            BroadcastDanmakuAction(ws, result)
        }
    }
}
```

---

## 三、TODO #2：礼物逻辑实现（GiftMessage）

### 3.1 设计目标

礼物消息的处理需要：
1. **礼物配置映射**（GiftId → 游戏效果：得兵数/治疗量）
2. **分数/积分更新**（礼物价值 → 排行榜积分）
3. **主播收益统计**（累计礼物价值）
4. **游戏动作触发**（高价值礼物 → 特殊效果）

### 3.2 礼物配置解析

从 `generateJsons/Gift.json` 可以看到：
- `Id`: 礼物ID
- `GiftNum`: 礼物数量（数组，第0个是主值）
- `GetSoldierNum`: 获得士兵数量
- `RescueInjuryNum`: 治疗/伤害数值 [治疗, 伤害]

### 3.3 GiftMessage 完整实现

```go
// GiftMessage 完整实现（替换原TODO）
func GiftMessage(messageList *pb.MessageList, ws *WebsocketAnchorClient) {
    now := time.Now().Unix()
    msgJoin(messageList, ws)  // 处理玩家加入

    for _, msg := range messageList.MsgList {
        if msg.Uid == "" {
            untils.TapErr(fmt.Sprintf("平台发的 uid is nil %v", msg))
            continue
        }

        // 获取/创建玩家数据
        pdb, _ := GetOpenVo(msg.Uid)
        if pdb == nil {
            untils.TapErr(fmt.Sprintf("pdb is nil %v", msg))
            continue
        }

        // 更新玩家基本信息
        if msg.Uid != "" {
            pdb.OpenId = msg.Uid
        }
        if msg.Name != "" {
            pdb.NickName = msg.Name
        }
        if msg.HeadImg != "" {
            pdb.AvatarUrl = msg.HeadImg
        }

        // 更新排名信息
        currWeekRank, currWeekScore := model.GetGameWeekRankByRankVo(pdb.OpenId, ws.AnchorV.PlatformId)
        currMonthRank, currMonthScore := model.GetGameMonthRankByRankVo(pdb.OpenId, ws.AnchorV.PlatformId)
        pdb.Score = currWeekScore
        pdb.MonthScore = currMonthScore
        pdb.Rank = int32(currWeekRank)
        pdb.MonthRank = int32(currMonthRank)
        pdb.LastUpdateTime = now

        // ========== 【TODO:礼物逻辑】完整实现 ==========
        
        // 解析礼物ID（msg.GiftId 或从Content解析）
        giftId := extractGiftId(msg)
        giftCfg := config.GetGiftConfig(float64(giftId))
        if giftCfg == nil {
            zlog.Logger.Warn().Msgf("礼物配置不存在: giftId=%d msg=%+v", giftId, msg)
            // 使用默认配置继续处理
            giftCfg = &config.GiftConfig{
                Id:             int32(giftId),
                GetSoldierNum:  []int32{1},
                GiftNum:        []int32{1},
                RescueInjuryNum: []int32{0, 0},
            }
        }

        // 计算礼物数量（支持批量礼物）
        giftNum := extractGiftCount(msg, giftCfg.GiftNum)
        // 计算总士兵数
        getSoldierNum := int64(0)
        if len(giftCfg.GetSoldierNum) > 0 {
            getSoldierNum = int64(giftCfg.GetSoldierNum[0]) * int64(giftNum)
        }
        // 治疗/伤害值
        rescueNum, injuryNum := int64(0), int64(0)
        if len(giftCfg.RescueInjuryNum) >= 2 {
            rescueNum = int64(giftCfg.RescueInjuryNum[0]) * int64(giftNum)
            injuryNum = int64(giftCfg.RescueInjuryNum[1]) * int64(giftNum)
        }

        // 更新排行榜（礼物积分 = 士兵数）
        if getSoldierNum > 0 {
            model.UpdateGameWeekRank(msg.Uid, getSoldierNum, ws.AnchorV.PlatformId)
            model.UpdateGameMonthRank(msg.Uid, getSoldierNum, ws.AnchorV.PlatformId)
            pdb.Score += getSoldierNum
        }

        // 更新主播礼物统计
        giftValue := calculateGiftValue(giftId, giftNum)
        ws.giftV += giftValue  // 累计礼物价值

        // 触发游戏动作（广播给UE4）
        if ws.IsStart && getSoldierNum > 0 {
            giftResult := &GiftResult{
                OpenId:        msg.Uid,
                NickName:      msg.Name,
                GiftId:        giftId,
                GiftNum:       giftNum,
                SoldierNum:    getSoldierNum,
                RescueNum:     rescueNum,
                InjuryNum:     injuryNum,
                TotalGiftValue: giftValue,
                Timestamp:     time.Now().UnixMilli(),
            }
            BroadcastGiftAction(ws, giftResult)
        }

        // 加入礼物排行榜（3人实时榜）
        if giftValue > 0 {
            model.AddAnchorPlayerRank(ws.AccountId, msg.Uid, giftValue)
        }

        zlog.Logger.Info().Msgf("礼物消息: openId=%s giftId=%d num=%d soldier=%d rescue=%d injury=%d value=%.2f",
            msg.Uid, giftId, giftNum, getSoldierNum, rescueNum, injuryNum, giftValue)
    }
}

// 礼物结果
type GiftResult struct {
    OpenId        string
    NickName      string
    GiftId        int32
    GiftNum       int32
    SoldierNum    int64
    RescueNum     int64
    InjuryNum     int64
    TotalGiftValue float64
    Timestamp     int64
}

// 从消息中提取礼物ID（proto定义：GiftId 是 string 类型）
func extractGiftId(msg *pb.Message) int32 {
    // 优先使用专用字段（string类型，需转换）
    if msg.GiftId != "" {
        if id, err := strconv.Atoi(msg.GiftId); err == nil {
            return int32(id)
        }
    }
    // 从Content解析（格式如 "gift:2:10" 表示 id:num）
    parts := strings.Split(msg.Content, ":")
    if len(parts) >= 1 {
        if id, err := strconv.Atoi(parts[0]); err == nil {
            return int32(id)
        }
    }
    return 0
}

// 从消息中提取礼物数量（proto定义：GiftCount 是 string 类型）
func extractGiftCount(msg *pb.Message, defaultVal []int32) int32 {
    if msg.GiftCount != "" {
        if num, err := strconv.Atoi(msg.GiftCount); err == nil && num > 0 {
            return int32(num)
        }
    }
    parts := strings.Split(msg.Content, ":")
    if len(parts) >= 2 {
        if num, err := strconv.Atoi(parts[1]); err == nil && num > 0 {
            return int32(num)
        }
    }
    if len(defaultVal) > 0 {
        return defaultVal[0]
    }
    return 1
}

// 计算礼物价值（用于主播收益统计）
func calculateGiftValue(giftId int32, giftNum int32) float64 {
    // 简化估算：礼物ID越大价值越高
    // 实际应该查表
    baseValue := map[int32]float64{
        1: 0.1,   // 1元
        2: 1.0,   // 10元
        3: 2.0,   // 20元
        4: 5.0,   // 50元
        5: 10.0,  // 100元
        6: 20.0,  // 200元
        7: 30.0,  // 300元
        8: 52.0,  // 520元
    }
    if v, ok := baseValue[giftId]; ok {
        return v * float64(giftNum)
    }
    return 0.1 * float64(giftNum)
}

// BroadcastGiftAction 将礼物动作广播给游戏客户端
func BroadcastGiftAction(ws *WebsocketAnchorClient, result *GiftResult) {
    // TODO: 发送礼物动作帧给UE4
    // 示例：发送Tips包含礼物特效信息
    tips := &pb.TipsNotify{
        Tips: fmt.Sprintf("[%s] 送出礼物 #%d ×%d，获得 %d 士兵！",
            result.NickName, result.GiftId, result.GiftNum, result.SoldierNum),
    }
    ws.WsSend(COMMAND_TIPS_S, tips)
    zlog.Logger.Debug().Msgf("礼物动作广播: openId=%s giftId=%d num=%d", 
        result.OpenId, result.GiftId, result.GiftNum)
}
```

---

## 四、TODO #3：点赞逻辑实现（LikeMessage）

### 4.1 设计目标

点赞消息的处理需要：
1. **批量合并**（已实现：每1秒合并一次同uid的点赞）
2. **活跃积分**（点赞数 → 积分）
3. **游戏buff触发**（点赞达到阈值 → 全局buff）
4. **广播特效**（向游戏客户端发送点赞特效）

### 4.2 LikeMessage 完整实现

```go
// LikeMessage 完整实现（替换原TODO）
func LikeMessage(messageList *pb.MessageList, ws *WebsocketAnchorClient) {
    now := time.Now().Unix()
    msgJoin(messageList, ws)

    for _, msg := range messageList.MsgList {
        if msg.Uid == "" {
            untils.TapErr(fmt.Sprintf("平台发的 uid is nil %v", msg))
            continue
        }

        // 获取/创建玩家数据
        pdb, _ := GetOpenVo(msg.Uid)
        if pdb == nil {
            untils.TapErr(fmt.Sprintf("pdb is nil %v", msg))
            continue
        }

        // 更新排名信息（每60秒一次，节省Redis调用）
        if now-pdb.LastUpdateTime > 60 {
            currWeekRank, currWeekScore := model.GetGameWeekRankByRankVo(pdb.OpenId, ws.AnchorV.PlatformId)
            currMonthRank, currMonthScore := model.GetGameMonthRankByRankVo(pdb.OpenId, ws.AnchorV.PlatformId)
            pdb.Score = currWeekScore
            pdb.MonthScore = currMonthScore
            pdb.Rank = int32(currWeekRank)
            pdb.MonthRank = int32(currMonthRank)
            pdb.LastUpdateTime = now
        }

        // ========== 【TODO:点赞逻辑】完整实现 ==========
        
        // 解析点赞数量（合并后的Count字段）
        likeCount := extractLikeCount(msg)
        if likeCount <= 0 {
            likeCount = 1
        }

        // 计算积分（每点赞1次 = 1积分，批量时乘以数量）
        scoreReward := int64(likeCount) * 1

        // 更新排行榜
        if scoreReward > 0 {
            model.UpdateGameWeekRank(msg.Uid, scoreReward, ws.AnchorV.PlatformId)
            model.UpdateGameMonthRank(msg.Uid, scoreReward, ws.AnchorV.PlatformId)
            pdb.Score += scoreReward
        }

        // 触发游戏Buff（点赞达到阈值时）
        // 从Buff配置中查找点赞相关的buff
        triggerLikeBuff(ws, msg.Uid, msg.Name, likeCount)

        zlog.Logger.Debug().Msgf("点赞消息: openId=%s count=%d totalScore=%d",
            msg.Uid, likeCount, scoreReward)
    }
}

// 点赞Buff触发器（基于Buff.json配置）
func triggerLikeBuff(ws *WebsocketAnchorClient, openId, nickName string, likeCount int32) {
    // BuffType=1: 被动增益 (如增加攻击力)
    // BuffType=2: 持续恢复 (如持续回血)
    // BuffType=3: 移动速度
    // BuffType=4: 攻击速度
    // 根据点赞数量触发不同级别的Buff

    var buffType int32
    var buffParam int32
    var buffDesc string

    switch {
    case likeCount >= 100:
        buffType = 4  // 攻击速度+5
        buffParam = 5
        buffDesc = "疯狂点赞！攻速+5"
    case likeCount >= 50:
        buffType = 3  // 移动速度+5
        buffParam = 5
        buffDesc = "疯狂点赞！移速+5"
    case likeCount >= 20:
        buffType = 1  // 攻击力+3
        buffParam = 3
        buffDesc = "点赞攻势！攻击+3"
    case likeCount >= 10:
        buffType = 2  // 持续回血
        buffParam = 5
        buffDesc = "点赞鼓励！持续回血中"
    default:
        // 低于10次不触发特殊buff，但仍然有积分
        return
    }

    // 广播点赞特效
    if ws.IsStart {
        tips := &pb.TipsNotify{
            Tips: fmt.Sprintf("❤️ [%s] 点赞×%d %s", nickName, likeCount, buffDesc),
        }
        ws.WsSend(COMMAND_TIPS_S, tips)
    }

    zlog.Logger.Info().Msgf("点赞触发Buff: openId=%s count=%d buffType=%d param=%d",
        openId, likeCount, buffType, buffParam)
}

// 从Message中提取点赞数量
// 注意：pb.Message.Count 是 string 类型
func extractLikeCount(msg *pb.Message) int32 {
    if msg.Count != "" {
        if count, err := strconv.Atoi(msg.Count); err == nil && count > 0 {
            return int32(count)
        }
    }
    if msg.Content != "" {
        if count, err := strconv.Atoi(msg.Content); err == nil && count > 0 {
            return int32(count)
        }
    }
    return 1
}
```

---

## 五、关键缺失：GameFrame 广播协议设计

### 5.1 当前问题

现有的 `WsSend()` 只能发送 `COMMAND_TIPS_S`（TipsNotify），无法满足游戏动作的需求。

需要定义新的协议来让UE4客户端收到弹幕/礼物/点赞的动作消息。

### 5.2 协议设计建议

```protobuf
// serverMsg.proto 新增

// 弹幕动作帧
message BarrageActionNotify {
    string open_id = 1;        // 玩家openid
    string nick_name = 2;      // 玩家昵称
    int32 action_type = 3;     // 动作类型: 1=攻击 2=防守 3=治疗 4=加速 5=减速 6=特殊
    int64 score_reward = 4;    // 积分奖励
    int64 timestamp = 5;       // 时间戳
}

// 礼物动作帧
message GiftActionNotify {
    string open_id = 1;
    string nick_name = 2;
    int32 gift_id = 3;
    int32 gift_num = 4;
    int64 soldier_num = 5;     // 获得士兵数
    int64 rescue_num = 6;      // 治疗量
    int64 injury_num = 7;      // 伤害值
    double total_value = 8;     // 礼物总价值
    int64 timestamp = 9;
}

// 点赞Buff帧
message LikeBuffNotify {
    string open_id = 1;
    string nick_name = 2;
    int32 like_count = 3;
    int32 buff_type = 4;       // Buff类型
    int32 buff_param = 5;      // Buff参数
    string buff_desc = 6;       // Buff描述
    int64 timestamp = 7;
}
```

### 5.3 在 init_define.go 中添加新指令

```go
// servers/logic/control/init_define.go 新增

const (
    // ... 现有指令 ...
    
    // 弹幕动作帧
    COMMAND_BARRAGE_ACTION_S int16 = 21
    // 礼物动作帧
    COMMAND_GIFT_ACTION_S int16 = 22
    // 点赞Buff帧
    COMMAND_LIKE_BUFF_S int16 = 23
)
```

---

## 六、集成计划

### 6.1 文件变更清单

| 文件 | 变更类型 | 说明 |
|------|----------|------|
| `servers/logic/control/control_barrage.go` | 修改 | 替换3个TODO为完整实现 |
| `servers/logic/control/danmaku_rule_engine.go` | 新增 | 弹幕规则引擎 |
| `servers/logic/control/init_define.go` | 修改 | 添加新指令常量 |
| `pb/serverMsg.proto` | 修改 | 添加BarrageActionNotify/GiftActionNotify/LikeBuffNotify |
| `pb/serverMsg.pb.go` | 重新生成 | 从proto生成 |

### 6.2 实施优先级

**P0（立即实施）：**
1. 弹幕规则引擎 + ChatMessage 实现
2. GiftMessage 实现
3. LikeMessage 实现
4. 基础GameFrame广播（通过TipsNotify临时方案）

**P1（尽快完成）：**
1. protobuf新协议定义与生成
2. 完整的GameFrame广播实现
3. 冷却系统接入Redis（分布式支持）

**P2（后续迭代）：**
1. AI意图分类增强
2. 情感计算积分加成
3. 弹幕弹幕历史追踪

---

## 七、代码完整性验证

### 7.1 需要确认的proto字段

以下字段需要对照实际的 `pb.Message` 定义确认：
- `msg.GetGiftId()` - 礼物ID字段
- `msg.Count` - 点赞数量字段
- `msg.Content` - 内容字段的解析格式

### 7.2 需要确认的函数

以下函数需要确认已存在：
- `untils.CheckSensitiveWord(text)` - 敏感词检查
- `model.AddAnchorPlayerRank(anchorId, playerId, value)` - 玩家礼物排行榜

### 7.3 需要确认的包导入

- `dmGameServer/config` - 礼物配置包（确认autoConfig是否export为config包名）

---

## 八、总结

通过本次深度代码实证研究，三个TODO的实现路径已经清晰：

1. **弹幕逻辑**：规则引擎（关键词匹配 → 动作映射 → 冷却管理 → 积分更新 → 广播）
2. **礼物逻辑**：配置驱动（GiftId → 士兵/治疗/伤害计算 → 积分 → 广播）
3. **点赞逻辑**：阈值触发（合并计数 → Buff等级判定 → 特效广播）

最关键的缺失是 **GameFrame 广播协议**（protobuf新消息类型），这是让UE4客户端感知弹幕动作的核心通道。下一步应优先完成proto定义和代码生成。

---

*本笔记为代码实证研究，聚焦 control_barrage.go TODO 实现方案*
