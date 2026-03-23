# MEMORY.md - 长期记忆

## 模型切换规则（2026-03-08 白老师制定）
1. 简单日常问答、确认类、信息查询类轻量需求：使用 `ark/doubao-seed-2.0-lite` 模型
2. 需求分析、方案设计、复杂逻辑推理、规则制定类需求：使用 `ark/doubao-seed-2.0-pro` 模型
3. 代码编写、调试、技术问题排查、性能优化类开发需求：使用 `ark/doubao-seed-2.0-code` 模型

## 核心项目信息（恐龙攻守塔防弹幕游戏）
### 技术栈
- 客户端：UE4（C++），帧同步方案，所有游戏逻辑在前端实现
- 后端：Go + MongoDB + Redis，当前机器为弹幕服务器
- 覆盖平台：抖音、快手国内先行，后续上线TikTok（需做多语言适配）

### 核心玩法
恐龙题材双人主播PK塔防：双方直播间观众送礼物生成恐龙兵种进攻对方基地，主播可切换攻防视角，基地被攻破即失败。

### 项目当前状态（2026-03-20更新）
- ✅ 已完成：基础WebSocket框架、Protobuf协议（Z11/Z12/Z14等）、敏感词过滤、雪花ID、排行榜系统、玩家数据模型
- ⚠️ 待完成：弹幕逻辑（`// todo:弹幕逻辑`）、礼物触发游戏指令（`// todo:礼物逻辑`）
- ❌ 未涉及：VeADK-Go AI集成、TikTok海外平台适配

### 后端职责范围
1. 抖音/快手/TikTok平台SDK对接，弹幕/礼物/直播间事件监听
2. 主播PK匹配、房间管理、并发连接管理
3. 帧数据回传、合法性校验、反外挂检测
4. 玩家数据、礼物数据、对战数据持久化存储
5. Web管理后台：主播管理、玩家管理、GM后台、模拟礼物工具
6. 素材自动化处理：参考图拆解、带通道PNG生成自动化流程
7. **【新增待开发】弹幕/礼物 → 游戏指令转换逻辑**

### 资源存储
SVN目录统一存放所有工程、代码、文档、素材

### 关联仓库
- OpenClaw 相关仓库：https://github.com/zhoudongyu0225/BytedanceOpenclaw

## 项目深度研究更新（2026-03-23）
### control_barrage.go 三个TODO状态更新
- ✅ ChatMessage弹幕逻辑（line 48）：完整代码方案已定稿，含敏感词过滤+规则引擎+冷却+积分+广播
- ✅ GiftMessage礼物逻辑（line 90）：完整代码方案已定稿，含computeSoldierNum+computeRescueNum+价值校验+GAME_GIFT广播
- ✅ LikeMessage点赞逻辑（line 134）：完整代码方案已定稿，含累计点赞+阈值Buff触发+GAME_BUFF广播

### Protobuf关键发现
- Message.count = 点赞数量（用于LikeMessage）
- Message.total/giftId/giftCount/giftName = 礼物相关字段（用于GiftMessage）
- TipsNotify.Type: Normal=0, Notice=1, RoomRoll=2, RoomMarquee=3, ServerMarquee=4
- 广播前缀约定：GAME_ACTION|uid|name|action、GAME_GIFT|uid|name|giftName|soldiers|rescue、GAME_BUFF|uid|name|buffType|buffValue

### 待开发模块
- danmakuRuleEngine.go：弹幕规则引擎（DanmakuRule.json驱动）
- VeADK-Go AI审核Agent集成
- GM管理后台
