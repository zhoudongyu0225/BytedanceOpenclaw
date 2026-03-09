package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// 根据协议文件生成协议逻辑

func main() {
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}
	// 读取 当下root目录下的servers/logic/control/init_define.go的文本
	f, err := os.Open(currentDir + "/servers/logic/control/init_define.go")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("Error reading file:", err)
	}
	lastDesc := ""
	// 获取包含的 COMMAND 的一行文本 一行一行打印出现
	for _, line := range strings.Split(string(b), "\n") {
		if strings.Contains(line, "COMMAND") {
			line = strings.TrimSpace(line)
			// 去重空格
			line = strings.Replace(line, " ", "", -1)
			// 取掉 int16
			line = strings.Replace(line, "int16", "", -1)
			// 根据=分割
			lineList := strings.Split(line, "=")
			headcmd := lineList[1]
			cmdNum, err := strconv.Atoi(headcmd)
			if err != nil {
				fmt.Println("----", line, headcmd)
				fmt.Println("Error reading file:", err)
				return
			}
			cmd := lineList[0]
			// cmd小写
			cmd = strings.ToLower(cmd)
			//取掉command_
			cmd = strings.Replace(cmd, "command_", "", -1)
			name := ""
			if cmdNum < 0 {
				fmt.Println(fmt.Sprintf("F%v_%v", headcmd, cmd))
				name = fmt.Sprintf("F%v_%v", headcmd, cmd)
			} else {
				fmt.Println(fmt.Sprintf("Z%v_%v", headcmd, cmd))
				name = fmt.Sprintf("Z%v_%v", headcmd, cmd)
			}
			fmt.Println(name, cmd)
			// 当下root目录下的protocol/创建  proto的文件
			// 生成一个文
			createProtoFile(name, currentDir+"/protocol", lastDesc)
		} else {
			lastDesc = line
		}
	}

}

func createProtoFile(input string, dir string, desc string) {
	// 根据输入字符串生成文件名和消息名称
	var fileName, messageName string
	if strings.HasSuffix(input, "_s2c") {
		fileName = replaceLast(input, "_s2c", "Resp.proto")
		messageName = generateMessageName(input, "_s2c")
		messageName = fmt.Sprintf("%vResp", messageName)
	} else if strings.HasSuffix(input, "_c2s") {
		fileName = replaceLast(input, "_c2s", "Req.proto")
		messageName = generateMessageName(input, "_c2s")
		messageName = fmt.Sprintf("%vReq", messageName)
	} else if strings.HasSuffix(input, "_c") {
		fileName = replaceLast(input, "_c", "C.proto")
		messageName = generateMessageName(input, "_c")
		messageName = fmt.Sprintf("%vC", messageName)
	} else {
		fileName = replaceLast(input, "_s", "Notify.proto")
		messageName = generateMessageName(input, "_s")
		messageName = fmt.Sprintf("%vNotify", messageName)
	}
	fileName = replaceOneUnderscore(fileName)

	fs := strings.Split(fileName, "_")
	fileName = fmt.Sprintf("%v_%v.proto", fs[0], messageName)

	// 生成完整文件路径
	filePath := filepath.Join(dir, fileName)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); err == nil {
		fmt.Println("文件已存在，跳过:", filePath)
		return
	}

	// 创建并打开文件
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("无法创建文件:", err)
		return
	}
	defer file.Close()

	// 文件内容
	content := fmt.Sprintf(`syntax = "proto3";
package pb;
option go_package = "./;pb";
%v
message %s{

}`, desc, messageName)

	if strings.HasSuffix(messageName, "Resp") {
		// 文件内容
		content = fmt.Sprintf(`syntax = "proto3";
package pb;
option go_package = "./;pb";
%v
message %s{
 string  errMsg = 1;//消息code 0成功 1失败
}`, desc, messageName)

	}

	// 将内容写入文件
	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println("写入文件失败:", err)
		return
	}

	fmt.Println("文件已生成:", filePath)
}

func generateMessageName(input string, suffix string) string {
	// 移除指定后缀
	base := strings.TrimSuffix(input, suffix)
	// 移除前缀
	base = strings.TrimPrefix(base, "F")
	// 将 '_' 分隔的字符串转换为 CamelCase
	parts := strings.Split(base, "_")
	for i := range parts {
		parts[i] = strings.Title(parts[i])
	}
	return strings.Join(parts[1:], "")
}
func replaceOneUnderscore(s string) string {
	count := 0
	return strings.Map(func(r rune) rune {
		if r == '_' {
			count++
			if count == 1 {
				return '_'
			}
			return -1
		}
		return r
	}, s)
}

// replaceLast 封装函数，传入原始字符串和需要替换的子字符串
func replaceLast(s, old, new string) string {
	if old == "" {
		return s
	}
	index := strings.LastIndex(s, old)
	if index == -1 {
		return s
	}
	return s[:index] + new + s[index+len(old):]
}
