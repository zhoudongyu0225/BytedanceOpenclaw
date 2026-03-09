#!/bin/bash
# 配置&记忆文件自动备份脚本
WORKSPACE="/root/.openclaw/workspace"
BACKUP_ROOT="$WORKSPACE/backups"
STATE_FILE="$BACKUP_ROOT/.backup_state"
# 要备份的文件/目录列表
BACKUP_TARGETS=(
  "SOUL.md"
  "IDENTITY.md" 
  "USER.md"
  "MEMORY.md"
  "AGENTS.md"
  "TOOLS.md"
  "HEARTBEAT.md"
  "memory/"
)

# 计算待备份文件总大小（KB）
calculate_total_size() {
  total=0
  for target in "${BACKUP_TARGETS[@]}"; do
    if [ -e "$WORKSPACE/$target" ]; then
      size=$(du -sk "$WORKSPACE/$target" | awk '{print $1}')
      total=$((total + size))
    fi
  done
  echo $total
}

# 获取当前星期（小写英文）
get_week_day() {
  date +%A | tr '[:upper:]' '[:lower:]'
}

# 读取上次备份状态
if [ -f "$STATE_FILE" ]; then
  source "$STATE_FILE"
else
  LAST_BACKUP_TIME=0
  LAST_BACKUP_SIZE=0
fi

CURRENT_SIZE=$(calculate_total_size)
CURRENT_TIME=$(date +%s)
TIME_DIFF=$((CURRENT_TIME - LAST_BACKUP_TIME))
SIZE_DIFF=$((CURRENT_SIZE - LAST_BACKUP_SIZE))
# 大小差取绝对值
if [ $SIZE_DIFF -lt 0 ]; then
  SIZE_DIFF=$((-SIZE_DIFF))
fi

# 触发条件：超过24小时（86400秒） 或 大小变化超过10KB
if [ $TIME_DIFF -gt 86400 ] || [ $SIZE_DIFF -gt 10 ]; then
  WEEK_DAY=$(get_week_day)
  BACKUP_DIR="$BACKUP_ROOT/$WEEK_DAY"
  # 清空对应星期的旧备份
  rm -rf "$BACKUP_DIR"/*
  # 执行备份
  for target in "${BACKUP_TARGETS[@]}"; do
    if [ -e "$WORKSPACE/$target" ]; then
      cp -r "$WORKSPACE/$target" "$BACKUP_DIR/"
    fi
  done
  # 更新状态文件
  echo "LAST_BACKUP_TIME=$CURRENT_TIME" > "$STATE_FILE"
  echo "LAST_BACKUP_SIZE=$CURRENT_SIZE" >> "$STATE_FILE"
  echo "LAST_BACKUP_DATE=$(date '+%Y-%m-%d %H:%M:%S')" >> "$STATE_FILE"
  echo "[$(date '+%Y-%m-%d %H:%M:%S')] 备份完成，存储至 $BACKUP_DIR" >> "$BACKUP_ROOT/.backup_log"
fi
