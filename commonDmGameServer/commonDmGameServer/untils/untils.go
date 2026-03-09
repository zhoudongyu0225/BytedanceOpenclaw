package untils

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/md5"
	rand2 "crypto/rand"
	"dmGameServer/common"
	"dmGameServer/zlog"
	"encoding/base64"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/pierrec/lz4/v4"
	"github.com/rs/xid"
	"io"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
	"unicode/utf8"
)

var globalX int32

// 简短的uid且唯一
func GenerateUID() string {
	globalX++
	globalX = globalX % 10000
	guid := xid.New()
	return fmt.Sprintf("%v%v", globalX, guid.String())
}

// string转int32
func StringToInt32(str string) int32 {
	i, err := strconv.Atoi(str)
	if err != nil {
		zlog.Logger.Error().Msgf("StringToInt32 error: %v", err)
		return 0
	}
	t := int32(i)
	if t < 0 {
		t = 0
	}
	return t
}

var oneUint32 uint32

// 简短的uid且唯一
func GenerateUint32() uint32 {
	oneUint32++
	oneUint32 = oneUint32 % 10000000
	// 百分百大于0
	return oneUint32 + 1
}

var sdkId int64

func GetSdkId() int64 {
	sdkId++
	sdkId = sdkId % 1000000000
	return sdkId
}

// ReadFileAsString 读取指定路径下的文件并将内容以字符串形式返回
func ReadFileAsString(filePath string) string {
	// 读取文件内容
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		zlog.Logger.Error().Msgf("读取文件失败: %v", err)
		return ""
	}
	// 将内容转换为字符串并返回
	return string(content)
}

// 将结构体转换为字节序列并计算 MD5 哈希
func GetMD5(s interface{}) string {
	// 将结构体序列化为 JSON 字节序列
	jsonData, err := json.Marshal(s)
	if err != nil {
		zlog.Logger.Error().Msgf("json.Marshal error: %v", err)
		return ""
	}

	// 计算 MD5 哈希
	hash := md5.Sum(jsonData)

	// 将哈希转换为十六进制字符串
	hashString := hex.EncodeToString(hash[:])

	return hashString
}

// 深拷贝
func DeepCopy(src, newDst interface{}) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		zlog.Logger.Error().Msgf("gob encode error: %v", err)
	}
	err := gob.NewDecoder(&buf).Decode(newDst)
	if err != nil {
		zlog.Logger.Error().Msgf("gob decode error: %v", err)
	}
}

// 更 float的数组arr  权重随机一个数 返回下标
func GetRandomIndexByWeight(arr []float64) int {
	var sum float64
	for _, value := range arr {
		sum += value
	}
	// 生成随机数
	r := GenerateRandomNumberFloat(0, sum)
	var total float64
	for i, value := range arr {
		total += value
		if r <= total {
			return i
		}
	}
	return 0
}

// 函数根据给定的概率返回一个布尔值，表示是否成功。
// 参数应该是一个介于0到1之间的浮点数，表示成功的概率。
// 判断概率是否满足(真随机 效率低)
func ProbabilitySuccess(probability float64) bool {
	if probability < 0 {
		return false
	}
	if probability >= 1 {
		return true
	}
	// 生成一个0到99999之间的随机整数
	randInt, err := rand2.Int(rand2.Reader, big.NewInt(100000))
	if err != nil {
		zlog.Logger.Error().Msgf("生成随机数失败:%v %v", err, probability)
		return false
	}
	// 将随机整数映射到概率范围内，并检查是否成功
	randFloat := float64(randInt.Int64()) / 100000.0
	return randFloat < probability
}

// 真随机 生成范围在 [a, b) 之间的随机整数
func TrueRandomNumber(a, b int) int {
	if a == 0 && b == 0 {
		return 0
	}
	if b == a {
		return a
	}
	if a > b {
		a, b = b, a
	}
	randInt, err := rand2.Int(rand2.Reader, big.NewInt(int64(b-a)))
	if err != nil {
		zlog.Logger.Error().Msgf("生成随机数失败:%v %v %v", err, a, b)
		return 0
	}
	return int(randInt.Int64()) + a
}

// 真随机 生成范围在 [a, b) 之间的随机浮点数
func TrueRandomNumberFloat(a, b float64) float64 {
	if a == 0 && b == 0 {
		return 0
	}
	if b == a {
		return a
	}
	if a > b {
		a, b = b, a
	}
	randInt, err := rand2.Int(rand2.Reader, big.NewInt(100000))
	if err != nil {
		zlog.Logger.Error().Msgf("生成随机数失败:%v %v %v", err, a, b)
		return 0
	}
	return float64(randInt.Int64())/100000.0*(b-a) + a

}

// 判断 float64的概率是否成功
// 判断概率是否满足(伪随机 效率高)
func IsSuccessByWeight(probability float64) bool {
	probability = probability * 100
	r := GenerateRandomNumberFloat(0, 100)
	AddProbabilitySeedMap(int64(r))
	if r < probability {
		// 成功
		return true
	}
	return false
}

// GenerateRandomNumber 生成范围在 [a, b) 之间的随机整数
func GenerateRandomNumber(a, b int) int {
	if a == 0 && b == 0 {
		return 0
	}
	if b == a {
		return a
	}
	if a > b {
		a, b = b, a
	}
	return rand.Intn(b-a) + a
}

// GenerateRandomNumberFloat 生成范围在 [a, b) 之间的随机浮点数
func GenerateRandomNumberFloat(a, b float64) float64 {
	if a == 0 && b == 0 {
		return 0
	}
	if b == a {
		return a
	}
	if a > b {
		a, b = b, a
	}
	return rand.Float64()*(b-a) + a
}

// GetRandomIndex 返回本次随机出来的下标
// -1：元数据空 -2：数据不合法，超出总和
func GetRandomIndex(original []int) int {
	if original == nil {
		return -1
	}

	var sum int
	for _, value := range original {
		sum += value
	}
	r := GenerateRandomNumber(1, sum+1)
	var total int
	for i, value := range original {
		total += value
		if r <= total {
			return i
		}
	}

	return -2
}

// -------------------

// PKCS7 padding
func padPKCS7(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	pad := make([]byte, padding)
	for i := range pad {
		pad[i] = byte(padding)
	}
	return append(data, pad...)
}

// PKCS7 unpadding
func unpadPKCS7(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}

// Encrypt encrypts plaintext using AES in ECB mode with PKCS7 padding
func Encrypt(plaintext []byte) (string, error) {

	key := []byte("ThisIsASecretKeyz123456Z")
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plaintext = padPKCS7(plaintext, aes.BlockSize)
	ciphertext := make([]byte, len(plaintext))

	// ECB mode does not use an IV
	for i := 0; i < len(plaintext); i += aes.BlockSize {
		block.Encrypt(ciphertext[i:i+aes.BlockSize], plaintext[i:i+aes.BlockSize])
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts ciphertext using AES in ECB mode with PKCS7 padding
func Decrypt(ciphertext string) (string, error) {
	key := []byte("ThisIsASecretKeyz123456Z")

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	decodedCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	plaintext := make([]byte, len(decodedCiphertext))

	// ECB mode does not use an IV
	for i := 0; i < len(decodedCiphertext); i += aes.BlockSize {
		block.Decrypt(plaintext[i:i+aes.BlockSize], decodedCiphertext[i:i+aes.BlockSize])
	}

	plaintext = unpadPKCS7(plaintext)

	return string(plaintext), nil
}

// 压缩
func CompressData(data []byte) ([]byte, error) {
	var compressedBuffer bytes.Buffer
	writer := gzip.NewWriter(&compressedBuffer)
	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	return compressedBuffer.Bytes(), nil
}

// 解压
func DecompressData(compressedData []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	decompressedData, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return decompressedData, nil
}

// Gzip compresses a string using gzip and then encodes it in Base64
func GzipCompressBase64(data string) (string, error) {
	var compressedBuffer bytes.Buffer
	writer := gzip.NewWriter(&compressedBuffer)

	_, err := writer.Write([]byte(data))
	if err != nil {
		return "", err
	}

	err = writer.Close()
	if err != nil {
		return "", err
	}

	// Encode the compressed data in Base64
	return base64.StdEncoding.EncodeToString(compressedBuffer.Bytes()), nil
}

// Gzip decompresses a Base64-encoded and gzip-compressed string
func GzipDecompressBase64(compressedBase64Data string) (string, error) {
	// Decode the Base64 string
	decodedData, err := base64.StdEncoding.DecodeString(compressedBase64Data)
	if err != nil {
		return "", err
	}
	// Create a reader for the decoded data
	reader, err := gzip.NewReader(bytes.NewReader(decodedData))
	if err != nil {
		return "", err
	}
	defer reader.Close()

	// Read the decompressed data
	decompressedData, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(decompressedData), nil
}

// 发送post请求 application/json
func SendPostRequest(url string, data []byte) ([]byte, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		zlog.Logger.Error().Msgf("NewRequest%v", err)
		return nil, err
	}
	// 设置请求头
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

// 发送post请求 x-www-form-urlencoded
func SendPostRequest2(url2 string, signParamsMap map[string]interface{}) ([]byte, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 将 map 转换为 url.Values
	data := make(url.Values)
	for key, value := range signParamsMap {
		data.Add(key, fmt.Sprint(value))
	}

	// 将数据编码为 x-www-form-urlencoded 格式
	body2 := bytes.NewBufferString(data.Encode())

	request, err := http.NewRequest("POST", url2, body2)
	if err != nil {
		zlog.Logger.Error().Msgf("NewRequest%v", err)
		return nil, err
	}
	// 设置请求头
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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

func ListContainsfloat64(l []float64, value float64) bool {
	for _, v := range l {
		if v == value {
			return true
		}
	}
	//switch value.(type) {
	//case int:
	//	list1 := l.([]int)
	//	for _, v := range list1 {
	//		if v == value.(int) {
	//			return true
	//		}
	//	}
	//case int32:
	//	list1 := l.([]int32)
	//	for _, v := range list1 {
	//		if v == value.(int32) {
	//			return true
	//		}
	//	}
	//
	//case int64:
	//	list1 := l.([]int64)
	//	for _, v := range list1 {
	//		if v == value.(int64) {
	//			return true
	//		}
	//	}
	//case string:
	//	list1 := l.([]string)
	//	for _, v := range list1 {
	//		if v == value.(string) {
	//			return true
	//		}
	//	}
	//case float32:
	//	list1 := l.([]float32)
	//	for _, v := range list1 {
	//		if v == value.(float32) {
	//			return true
	//		}
	//
	//	}
	//case float64:
	//	list1 := l.([]float64)
	//	for _, v := range list1 {
	//		if v == value.(float64) {
	//			return true
	//		}
	//	}
	//}
	return false
}

// 概率统计
var ProbabilityMap sync.Map // key:概率  value: 次数
// 概率种子记录
var ProbabilitySeedMap sync.Map // key:种子 value: 次数

func AddProbabilityMap(probability float64) {
	if v, ok := ProbabilityMap.Load(probability); ok {
		ProbabilityMap.Store(probability, v.(int)+1)
	} else {
		ProbabilityMap.Store(probability, 1)
	}
}

func AddProbabilitySeedMap(seed int64) {
	if v, ok := ProbabilitySeedMap.Load(seed); ok {
		ProbabilitySeedMap.Store(seed, v.(int)+1)
	} else {
		ProbabilitySeedMap.Store(seed, 1)
	}
}

// 获取概率统计
func GetProbabilityMap() {
	m1 := make(map[float64]int)
	ProbabilityMap.Range(func(key, value interface{}) bool {
		m1[key.(float64)] = value.(int)
		return true
	})
	m2 := make(map[int64]int)
	m3 := make(map[int64]int)
	ProbabilitySeedMap.Range(func(key, value interface{}) bool {
		t := key.(int64)
		m2[t/10] += value.(int)
		m3[t/5] += value.(int)
		return true
	})
	zlog.Logger.Info().Msgf("概率统计:%v   种子记录(10个):%v   种子记录(5个):%v ", m1, m2, m3)
}

func PanicPoss1(err2 any, stack []byte, AccountId string) {
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
[AccountId]:%s
[Panic Recovery]:%s 
panic recovered:
%s
stack:
%s
      `, hostname, common.ServiceSign, AccountId, time.Now().Format("2006/01/02 - 15:04:05"), err2, stack)
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
	}()

}

type Player30Gift struct {
	GiftV int64
	lt    int64 //上次时间
}

var playerDailyGiftMap = make(map[string]*Player30Gift)
var playerDailyGiftMapRw sync.RWMutex

// PlayerGiftV 玩家礼物值
func PlayerGiftV(AccountId string, pt int32, name string, roomId string, playerId string, playerName string, giftV int64) {

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
		giftV = giftV / 10
		roomId = fmt.Sprintf("https://webcast.amemv.com/douyin/webcast/reflow/%v", roomId)
	case 6:
		ptName = "QQ"
	}

	now := time.Now().Unix()

	playerDailyGiftMapRw.Lock()
	_, ok := playerDailyGiftMap[playerId]
	if !ok {
		playerDailyGiftMap[playerId] = &Player30Gift{
			GiftV: 0,
			lt:    now,
		}
	}
	// 30分钟
	if now-playerDailyGiftMap[playerId].lt > 1800 {
		playerDailyGiftMap[playerId] = &Player30Gift{
			GiftV: 0,
			lt:    now,
		}
	}
	playerDailyGiftMap[playerId].GiftV += giftV
	dv := playerDailyGiftMap[playerId].GiftV
	defer playerDailyGiftMapRw.Unlock()

	// 礼物值小于1w跳过
	if dv < 100000 {
		return
	}

	webhookURL := "https://open.feishu.cn/open-apis/bot/v2/hook/6386251d-8bee-4fc8-bdc3-f691ccd3772b"
	gameName := ""
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	msg := fmt.Sprintf(`%s
!!!!土豪来咯!!!!! :%s
[玩家id]:%s
[游戏平台]:%s
[游戏名字]:%s 
[主播id]:%s
[主播名字]:%s
[30分钟内礼物值(单位角)]:%v
[单次礼物值(单位角)]:%v
[主播地址]:%s
`, playerName, time.Now().Format("2006/01/02 - 15:04:05"), playerId, ptName,
		gameName, AccountId, name, dv, giftV, roomId)
	if AccountId == "_000Vw1JfmFdV85iKH2jPVvO-i71TGSBoB8y" {
		msg = fmt.Sprintf("---------------------审核他来咯--------------------  \n%v", msg)
	}

	zlog.Logger.Debug().Msgf("%s", msg)

	os1 := runtime.GOOS
	switch os1 {
	case "windows":
		return
	}

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

func CompressBuffer(input *bytes.Buffer) (*bytes.Buffer, error) {
	var compressedBuffer bytes.Buffer
	// 创建一个gzip写入器，将数据写入到compressedBuffer中
	gzipWriter := gzip.NewWriter(&compressedBuffer)
	defer gzipWriter.Close()

	// 将输入数据写入gzip写入器
	_, err := gzipWriter.Write(input.Bytes())
	if err != nil {
		return nil, err
	}

	return &compressedBuffer, nil
}

// Compress 使用 LZ4 压缩数据
func LZ4Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	writer := lz4.NewWriter(&b)
	_, err := writer.Write(data)
	if err != nil {
		return nil, fmt.Errorf("压缩错误: %w", err)
	}
	writer.Close()
	return b.Bytes(), nil
}

// Decompress 使用 LZ4 解压缩数据
func LZ4Decompress(compressedData []byte) ([]byte, error) {
	reader := lz4.NewReader(bytes.NewReader(compressedData))
	var decompressedData bytes.Buffer
	_, err := io.Copy(&decompressedData, reader)
	if err != nil {
		return nil, fmt.Errorf("解压缩错误: %w", err)
	}
	return decompressedData.Bytes(), nil
}

// Compress 使用 gzip 压缩数据
func GzipCompress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	writer := gzip.NewWriter(&b)
	_, err := writer.Write(data)
	if err != nil {
		return nil, fmt.Errorf("压缩错误: %w", err)
	}
	writer.Close()
	return b.Bytes(), nil
}

// Decompress 使用 gzip 解压缩数据
func GzipsDecompress(compressedData []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return nil, fmt.Errorf("解压缩错误: %w", err)
	}
	defer reader.Close()

	var decompressedData bytes.Buffer
	_, err = io.Copy(&decompressedData, reader)
	if err != nil {
		return nil, fmt.Errorf("解压缩错误: %w", err)
	}
	return decompressedData.Bytes(), nil
}

// 服务器服关闭找错判断
func ServerClose(msgs string) {
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
服务器关闭完毕
[hostname]:%s
[ServiceSign]:%s
[msg]:%s
[time]:%s
`, hostname, common.ServiceSign, msgs, time.Now().Format("2006/01/02 - 15:04:05"))
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

func RemoveBOM(s string) string {
	bom := '\ufeff'
	if len(s) > 0 && utf8.RuneCountInString(s) > 0 {
		firstRune, size := utf8.DecodeRuneInString(s)
		if firstRune == bom {
			return s[size:]
		}
	}
	return s
}
