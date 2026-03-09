#!/bin/bash
# 每日进化复盘报告生成脚本
WORKSPACE="/root/.openclaw/workspace"
TODAY=$(date +%Y-%m-%d)
MEMORY_FILE="$WORKSPACE/memory/$TODAY.md"
REPORT_FILE="$WORKSPACE/memory/${TODAY}_evolution_report.md"
CHAT_ID="user:ou_78cc320f0b8cf38abf3f9433655b72da"

# 若当天无记忆文件则退出
if [ ! -f "$MEMORY_FILE" ]; then
  echo "$TODAY 无会话记录，跳过生成报告" >> "$WORKSPACE/backups/.cron_log"
  exit 0
fi

# 生成报告内容
cat > "$REPORT_FILE" << EOF
# $TODAY 每日进化报告
## 今日核心要点回顾
$(grep -E "(决策|需求|规则|问题|解决|学习)" "$MEMORY_FILE" | head -20)

## 今日学习/精进内容
1. 
2. 
3. 

## 今日问题复盘
- 遇到的问题：
- 解决方案：
- 改进点：

## 可固化技能提议（Top3）
1. 
2. 
3. 
EOF

# 调用大模型补全报告内容
REPORT_CONTENT=$(openclaw run prompt "根据以下当天的记忆内容，补全这份进化报告，要求内容真实具体，符合当天实际交互情况，不要编造内容：\n\n记忆内容：\n$(cat $MEMORY_FILE)\n\n报告模板：\n$(cat $REPORT_FILE)")

# 发送报告到飞书
openclaw message send --channel feishu --target "$CHAT_ID" --message "$REPORT_CONTENT"

echo "$TODAY 进化报告已发送" >> "$WORKSPACE/backups/.cron_log"
