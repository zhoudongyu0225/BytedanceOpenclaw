# 配置&记忆文件备份恢复指南

## 备份规则
### 触发条件
满足任意一个条件自动触发备份：
1. 距离上一次备份时间超过24小时
2. 待备份文件总大小变化超过10KB

### 存储规则
- 备份根目录：`/root/.openclaw/workspace/backups/`
- 按星期循环存储：共7个子目录（monday到sunday），对应周一到周日的备份
- 同名星期目录下的旧备份会被新备份覆盖，最多保留最近7天的备份

### 备份范围
所有核心配置和记忆文件：
- SOUL.md（核心设定）
- IDENTITY.md（身份信息）
- USER.md（用户信息）
- MEMORY.md（长期记忆）
- AGENTS.md（工作规则）
- TOOLS.md（工具配置）
- HEARTBEAT.md（心跳配置）
- memory/ 目录（所有短期记忆日志）

## 自动恢复脚本使用方法
脚本路径：`/root/.openclaw/workspace/scripts/restore.sh`
1. 查看备份状态：`./restore.sh status`
2. 恢复最近一次备份：`./restore.sh latest`
3. 恢复指定星期的备份：`./restore.sh monday`（替换monday为对应星期英文小写）

## 手动恢复方法
1. 进入备份目录：`cd /root/.openclaw/workspace/backups/[对应星期目录]`
2. 直接将需要的文件复制回工作区根目录：`cp SOUL.md /root/.openclaw/workspace/`

## 手动触发备份
运行备份脚本即可立即执行一次备份：
`/root/.openclaw/workspace/scripts/backup.sh`
