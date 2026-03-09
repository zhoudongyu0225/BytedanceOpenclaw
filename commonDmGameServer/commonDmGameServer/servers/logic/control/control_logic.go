package control

// 添加玩家
func (ws *WebsocketAnchorClient) AddPlayer(openId string) {
	ws.PlayerMap.Store(openId, openId)
}

// 删掉玩家
func (ws *WebsocketAnchorClient) DelPlayer(openId string) {
	ws.PlayerMap.Delete(openId)
}

// 1秒的定时器
func (ws *WebsocketAnchorClient) SlgTimeTick() {
	// todo:可以写定时的逻辑

}

//----------------------玩家相关-------------------------
