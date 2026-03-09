package untils

import (
	"bytes"
	"dmGameServer/common"
	"dmGameServer/zlog"
	"fmt"
	"github.com/goccy/go-json"
	"net/http"
	"os"
	"runtime"
	"time"
)

// 开播地址
// TapBc 回合制补偿
func KBDz(name, url string) {
	type ContentT struct {
		Text string `json:"text"`
	}
	type PanicMsg struct {
		MsgType string    `json:"msg_type"`
		Content *ContentT `json:"content"`
	}
	// 获取本机主机名
	msg := fmt.Sprintf(`时间%v 主播%v Titok开播地址%v `, time.Now().Format("2006/01/02 - 15:04:05"), name, url)
	zlog.Logger.Error().Msgf("%s", msg)
	webhookURL := "https://open.feishu.cn/open-apis/bot/v2/hook/6386251d-8bee-4fc8-bdc3-f691ccd3772b"
	panicMsg := &PanicMsg{
		MsgType: "text",
		Content: &ContentT{
			Text: msg,
		},
	}
	go func() {
		payload, err := json.Marshal(panicMsg)

		// 发送 POST 请求
		resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}
		defer resp.Body.Close()
	}() // 发送 POST 请求

}

// 直播信息
type LivInfo struct {
	AccountId string
	Rank1     string
	Rank2     string
	Rank3     string
	Time1     string
	Time2     string
	Time3     string
	ShouRu1   string
	ShouRu2   string
	ShouRu3   string
	Dps1      string
	Dps2      string
	Dps3      string
	GameSR1   string // 游戏收入
	GameSR2   string // 游戏收入
	GameSR3   string // 游戏收入
}

// 主播开播
func ZBKPoss(AccountId string, pt, gameId int32, name string, curNum []int32, roomId string, Version string, livInfo *LivInfo) {
	type ContentT struct {
		Text string `json:"text"`
	}
	type PanicMsg struct {
		MsgType string    `json:"msg_type"`
		Content *ContentT `json:"content"`
	}
	ptName := ""
	switch pt {
	case 1:
		ptName = "Tiktok"
	case 2:
		ptName = "视频号"
	case 3:
		ptName = "快手"
	case 4:
		ptName = "哔站"
	case 5:
		ptName = "抖音"
		roomId = fmt.Sprintf("https://webcast.amemv.com/douyin/webcast/reflow/%v", roomId)
	case 6:
		ptName = "QQ"
	case 9:
		ptName = "狼人杀"
	case 10:
		ptName = "迷你派对"
	case 15:
		ptName = "TikTok"
	}
	webhookURL := "https://open.feishu.cn/open-apis/bot/v2/hook/6386251d-8bee-4fc8-bdc3-f691ccd3772b"
	gameName := fmt.Sprintf("%v", gameId)

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	// 获取本机主机名
	hostname, _ := os.Hostname()

	msg := fmt.Sprintf(`☄<%s> %s 
全平台[%d][抖音%d][快手%d][视频号%d][Tk%d]
             日  | 月  | 永久
《%v》排名 : %v  | %v  | %v
直播时长 : %v  | %v  | %v
主播收入 : %v  | %v  | %v
主播DPS 	: %v  | %v  | %v
时间 	: %s
[游戏版本号]:%s [主播id]:%s [主播地址]:%s [hostname]:%s [servicename]:%s [总内存占用]:%v MB [Goroutines数量]:%v
游戏收入 : %v  | %v  | %v`, gameName, name, curNum[4], curNum[0], curNum[1], curNum[2], curNum[3], ptName, livInfo.Rank1, livInfo.Rank2, livInfo.Rank3, livInfo.Time1, livInfo.Time2, livInfo.Time3,
		livInfo.ShouRu1, livInfo.ShouRu2, livInfo.ShouRu3, livInfo.Dps1, livInfo.Dps2, livInfo.Dps3,
		time.Now().Format("2006/01/02 - 15:04:05"),
		Version, AccountId, roomId, hostname, common.ServiceSign, memStats.Alloc/1024/1024, runtime.NumGoroutine(), livInfo.GameSR1, livInfo.GameSR2, livInfo.GameSR3)
	zlog.Logger.Debug().Msgf("%s", msg)
	//
	//os1 := runtime.GOOS
	//switch os1 {
	//case "windows":
	//	return
	//}
	panicMsg := &PanicMsg{
		MsgType: "text",
		Content: &ContentT{
			Text: msg,
		},
	}
	payload, err := json.Marshal(panicMsg)
	// 发送 POST 请求
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

}

// 主播关播 btime 当场直播  CurrGiftValue当场补偿的值  giftV当场礼物
func ZBGPoss(AccountId string, pt, gameId int32, name string, btime, giftV int64, curNum []int32, Version string, CurrGiftValue int64, livInfo *LivInfo) {

	type ContentT struct {
		Text string `json:"text"`
	}
	type PanicMsg struct {
		MsgType string    `json:"msg_type"`
		Content *ContentT `json:"content"`
	}
	webhookURL := "https://open.feishu.cn/open-apis/bot/v2/hook/6386251d-8bee-4fc8-bdc3-f691ccd3772b"
	ptName := fmt.Sprintf("平台:%v", pt)
	switch pt {
	case 1:
		ptName = "Tiktok"
	case 2:
		ptName = "视频号"
	case 3:
		ptName = "快手"
	case 4:
		ptName = "哔站"
	case 5:
		ptName = "抖音"
	case 6:
		ptName = "QQ"
	case 9:
		ptName = "狼人杀"
	case 10:
		ptName = "迷你派对"
	case 15:
		ptName = "TikTok"
	case 16:
		ptName = "抖音野游"
	}

	gameName := fmt.Sprintf("%v", gameId)

	bt := fmt.Sprintf("%v秒", btime)
	if btime >= 86400 {
		bt = fmt.Sprintf("%.v天", btime/86400)
	}
	if btime >= 3600 {
		bt = fmt.Sprintf("%.v小时", btime/3600)
	}
	if btime >= 60 {
		bt = fmt.Sprintf("%.v分钟", btime/60)
	}

	//dailyTime = dailyTime / 60
	//if dailyTime < 1 {
	//	dailyTime = 1
	//}

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	//fmt.Printf("分配内存: %v MB\n", memStats.Alloc/1024/1024)
	//fmt.Printf("总内存: %v MB\n", memStats.TotalAlloc/1024/1024)
	//fmt.Printf("当前Goroutines数量: %v\n", runtime.NumGoroutine())

	// 获取本机主机名
	hostname, _ := os.Hostname()
	msg := fmt.Sprintf(`♨关播 :<%s>%v
全平台[%d][抖音%d][快手%d][视频号%d][Tk%d]
            日  | 月  | 永久
《%v》排名 : %v  | %v  | %v
直播时长 : %v  | %v  | %v
主播收入 : %v  | %v  | %v
主播DPS 	: %v  | %v  | %v
时间 	: %s
[游戏版本号]:%s [主播id]:%s [本场时长]:%v分钟 [本场礼物值]:%v [本场补偿礼物值]:%v [hostname]:%s [servicename]:%s [总内存占用]:%v MB [Goroutines数量]:%v
游戏收入  %v  | %v  | %v`,
		gameName, name, curNum[4], curNum[0], curNum[1], curNum[2], curNum[3], ptName, livInfo.Rank1, livInfo.Rank2, livInfo.Rank3, livInfo.Time1, livInfo.Time2, livInfo.Time3,
		livInfo.ShouRu1, livInfo.ShouRu2, livInfo.ShouRu3, livInfo.Dps1, livInfo.Dps2, livInfo.Dps3,
		time.Now().Format("2006/01/02 - 15:04:05"),
		Version, AccountId, bt, giftV, CurrGiftValue, hostname, common.ServiceSign, memStats.Alloc/1024/1024, runtime.NumGoroutine(), livInfo.GameSR1, livInfo.GameSR2, livInfo.GameSR3)
	zlog.Logger.Info().Msgf("%s", msg)
	//os1 := runtime.GOOS
	//switch os1 {
	//case "windows":
	//	return
	//}
	panicMsg := &PanicMsg{
		MsgType: "text",
		Content: &ContentT{
			Text: msg,
		},
	}

	payload, err := json.Marshal(panicMsg)
	// 发送 POST 请求
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

}

// TapErr 服务器找错判断
func TapErr(msgs string) {
	type ContentT struct {
		Text string `json:"text"`
	}
	type PanicMsg struct {
		MsgType string    `json:"msg_type"`
		Content *ContentT `json:"content"`
	}
	// 获取调用者的信息
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		zlog.Logger.Error().Msgf("runtime.Caller() failed")
	}
	// 获取本机主机名
	hostname, _ := os.Hostname()
	msg := fmt.Sprintf(`
[hostname]:%s
[ServiceSign]:%s
[msg]:%s
[time]:%s
[callerName]:%s
[file]:%s
[line]:%v`, hostname, common.ServiceSign, msgs, time.Now().Format("2006/01/02 - 15:04:05"), runtime.FuncForPC(pc).Name(), file, line)
	zlog.Logger.Error().Msgf("%s", msg)
	webhookURL := "https://open.feishu.cn/open-apis/bot/v2/hook/6386251d-8bee-4fc8-bdc3-f691ccd3772b"
	os1 := runtime.GOOS
	switch os1 {
	case "windows":
		webhookURL = "https://open.feishu.cn/open-apis/bot/v2/hook/6386251d-8bee-4fc8-bdc3-f691ccd3772b"
	}
	panicMsg := &PanicMsg{
		MsgType: "text",
		Content: &ContentT{
			Text: msg,
		},
	}
	go func() {
		payload, err := json.Marshal(panicMsg)
		// 发送 POST 请求
		resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			fmt.Println("Error sending request:", err)
			return

		}
		defer resp.Body.Close()
	}()

}

// TapBc 回合制补偿
func TapBc(msgs string) {
	type ContentT struct {
		Text string `json:"text"`
	}
	type PanicMsg struct {
		MsgType string    `json:"msg_type"`
		Content *ContentT `json:"content"`
	}
	// 获取本机主机名
	hostname, _ := os.Hostname()
	msg := fmt.Sprintf(`补偿信息
[hostname]:%s
[ServiceSign]:%s
[msg]:%s
[time]:%s`, hostname, common.ServiceSign, msgs, time.Now().Format("2006/01/02 - 15:04:05"))
	zlog.Logger.Error().Msgf("%s", msg)
	webhookURL := "https://open.feishu.cn/open-apis/bot/v2/hook/6386251d-8bee-4fc8-bdc3-f691ccd3772b"
	panicMsg := &PanicMsg{
		MsgType: "text",
		Content: &ContentT{
			Text: msg,
		},
	}
	go func() {
		payload, err := json.Marshal(panicMsg)
		// 发送 POST 请求
		resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}
		defer resp.Body.Close()
	}() // 发送 POST 请求

}

func PanicPoss(err2 any, stack []byte) {
	type ContentT struct {
		Text string `json:"text"`
	}
	type PanicMsg struct {
		MsgType string    `json:"msg_type"`
		Content *ContentT `json:"content"`
	}
	// 获取本机主机名
	hostname, _ := os.Hostname()
	msg := fmt.Sprintf(`
[hostname]:%s
[ServiceSign]:%s
[Panic Recovery]:%s 
panic recovered:
%s
stack:
%s
      `, hostname, common.ServiceSign, time.Now().Format("2006/01/02 - 15:04:05"), err2, stack)
	zlog.Logger.Error().Msgf("%s", msg)
	webhookURL := "https://open.feishu.cn/open-apis/bot/v2/hook/6386251d-8bee-4fc8-bdc3-f691ccd3772b"
	panicMsg := &PanicMsg{
		MsgType: "text",
		Content: &ContentT{
			Text: msg,
		},
	}
	payload, err := json.Marshal(panicMsg)
	// 发送 POST 请求
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

}
