# 弹幕技术实战指南

## 一、WebSocket弹幕系统

### 1.1 服务端实现（Node.js）

```javascript
// 弹幕服务器核心代码示例
const WebSocket = require('ws');
const wss = new WebSocket.Server({ port: 8080 });

// 存储所有连接
const clients = new Map();

wss.on('connection', (ws, req) => {
  const clientId = generateId();
  clients.set(clientId, ws);
  
  ws.on('message', (message) => {
    // 解析弹幕消息
    const danmaku = JSON.parse(message);
    
    // 广播给所有客户端
    broadcast({
      type: 'danmaku',
      data: {
        text: danmaku.text,
        color: danmaku.color,
        position: danmaku.position,
        timestamp: Date.now()
      }
    });
  });
  
  ws.on('close', () => {
    clients.delete(clientId);
  });
});

function broadcast(message) {
  const data = JSON.stringify(message);
  clients.forEach(client => {
    if (client.readyState === WebSocket.OPEN) {
      client.send(data);
    }
  });
}
```

### 1.2 客户端实现

```javascript
// 弹幕渲染示例（Canvas）
class DanmakuRenderer {
  constructor(canvas) {
    this.canvas = canvas;
    this.ctx = canvas.getContext('2d');
    this.bullets = [];
    this.running = true;
    this.startRenderLoop();
  }
  
  addBullet(text, options = {}) {
    this.bullets.push({
      text,
      x: this.canvas.width,
      y: options.y || Math.random() * (this.canvas.height - 30),
      speed: options.speed || 3,
      color: options.color || '#ffffff',
      font: options.font || '16px Arial'
    });
  }
  
  render() {
    this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
    
    this.bullets = this.bullets.filter(bullet => {
      bullet.x -= bullet.speed;
      
      if (bullet.x < -200) return false;
      
      this.ctx.font = bullet.font;
      this.ctx.fillStyle = bullet.color;
      this.ctx.fillText(bullet.text, bullet.x, bullet.y);
      
      return true;
    });
    
    if (this.running) {
      requestAnimationFrame(() => this.render());
    }
  }
  
  startRenderLoop() {
    this.running = true;
    this.render();
  }
  
  stop() {
    this.running = false;
  }
}
```

## 二、Cocos Creator小游戏基础

### 2.1 项目结构
```
project/
├── assets/
│   ├── scenes/          # 场景文件
│   ├── scripts/         # TypeScript脚本
│   ├── prefabs/         # 预制体
│   └── textures/        # 纹理资源
├── build/               # 构建输出
└── project.json         # 项目配置
```

### 2.2 基础脚本示例

```typescript
// PlayerController.ts
import { _decorator, Component, Vec3, Input, input, EventKeyboard, KeyCode } from 'cc';

const { ccclass, property } = _decorator;

@ccclass('PlayerController')
export class PlayerController extends Component {
    @property
    speed: number = 300;
    
    private _direction: Vec3 = new Vec3();
    
    onLoad() {
        input.on(Input.EventType.KEY_DOWN, this.onKeyDown, this);
        input.on(Input.EventType.KEY_UP, this.onKeyUp, this);
    }
    
    update(deltaTime: number) {
        const pos = this.node.position;
        this.node.setPosition(
            pos.x + this._direction.x * this.speed * deltaTime,
            pos.y + this._direction.y * this.speed * deltaTime,
            pos.z
        );
    }
    
    onKeyDown(event: EventKeyboard) {
        switch(event.keyCode) {
            case KeyCode.ARROW_LEFT:
                this._direction.x = -1;
                break;
            case KeyCode.ARROW_RIGHT:
                this._direction.x = 1;
                break;
        }
    }
    
    onKeyUp(event: EventKeyboard) {
        this._direction.x = 0;
    }
}
```

## 三、抖音小游戏发布流程

### 3.1 账号要求
- 抖音创作者账号
- 完成实名认证
- 企业号或个人号均可

### 3.2 提审材料
1. 游戏包（.apk 或 .rpk）
2. 软件著作权证书（可选但建议）
3. 隐私政策文档
4. 游戏截图/Icon

### 3.3 审核要点
- 包体大小：主包≤4MB
- 加载时间：首屏≤3秒
- 支付合规：仅支持抖音支付
- 内容合规：无敏感话题

---
*技术指南持续更新中*
