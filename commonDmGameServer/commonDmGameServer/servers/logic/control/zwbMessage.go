package control

import (
	"bytes"
	"dmGameServer/zlog"
	"encoding/binary"
)

type WbMessageModel struct {
	Length     uint32
	Command    int16
	IsEncrypt  bool
	IsCompress bool
	Body       []byte
}

// 读取数据缓冲区（后面读取需要优化）
func (messageModel *WbMessageModel) Read(dataBuffer *bytes.Buffer) {
	var length uint32
	var command int16
	var isEncrypt bool
	var isCompress bool
	//
	err := binary.Read(dataBuffer, binary.LittleEndian, &length)
	if err != nil {
		zlog.Logger.Error().Msgf("[ERR] Read 获取length失败: %v", err)
		return
	}
	err = binary.Read(dataBuffer, binary.LittleEndian, &command)
	if err != nil {
		zlog.Logger.Error().Msgf("[ERR] Read 获取command失败: %v", err)
		return
	}
	err = binary.Read(dataBuffer, binary.LittleEndian, &isEncrypt)
	if err != nil {
		zlog.Logger.Error().Msgf("[ERR] Read 获取isEncrypt失败: %v", err)
		return
	}
	err = binary.Read(dataBuffer, binary.LittleEndian, &isCompress)
	if err != nil {
		zlog.Logger.Error().Msgf("[ERR] Read 获取isCompress失败: %v", err)
		return
	}
	// 有两个是指令
	body := make([]byte, length)
	_, err = dataBuffer.Read(body)
	if err != nil {
		zlog.Logger.Error().Msgf("[ERR] Read 获取body失败: %v %v %v", err, length, command)
		return
	}
	messageModel.Length = length
	messageModel.Command = command
	messageModel.IsEncrypt = isEncrypt
	messageModel.IsCompress = isCompress
	messageModel.Body = body
}
