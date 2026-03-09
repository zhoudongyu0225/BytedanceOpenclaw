#!/bin/bash
# 配置&记忆文件自动恢复脚本
WORKSPACE="/root/.openclaw/workspace"
BACKUP_ROOT="$WORKSPACE/backups"

if [ $# -eq 0 ]; then
  echo "用法："
  echo "  恢复最近一次备份：./restore.sh latest"
  echo "  恢复指定星期备份：./restore.sh [monday|tuesday|wednesday|thursday|friday|saturday|sunday]"
  echo "  查看所有备份状态：./restore.sh status"
  exit 0
fi

if [ "$1" = "status" ]; then
  if [ -f "$BACKUP_ROOT/.backup_state" ]; then
    source "$BACKUP_ROOT/.backup_state"
    echo "=== 备份状态 ==="
    echo "最后备份时间：$LAST_BACKUP_DATE"
    echo "备份文件总大小：$LAST_BACKUP_SIZE KB"
    echo -e "\n=== 现有备份 ==="
    for day in monday tuesday wednesday thursday friday saturday sunday; do
      count=$(ls -A "$BACKUP_ROOT/$day" 2>/dev/null | wc -l)
      if [ $count -gt 0 ]; then
        echo "$day: 存在备份"
      else
        echo "$day: 无备份"
      fi
    done
  else
    echo "暂无备份记录"
  fi
  exit 0
fi

# 确定要恢复的目录
if [ "$1" = "latest" ]; then
  # 找最新修改的备份目录
  RESTORE_DIR=$(find "$BACKUP_ROOT" -maxdepth 1 -type d ! -path "$BACKUP_ROOT" -printf '%T@ %p\n' | sort -n | tail -1 | awk '{print $2}')
  if [ -z "$RESTORE_DIR" ] || [ $(ls -A "$RESTORE_DIR" | wc -l) -eq 0 ]; then
    echo "错误：没有可用的备份"
    exit 1
  fi
else
  RESTORE_DIR="$BACKUP_ROOT/$1"
  if [ ! -d "$RESTORE_DIR" ] || [ $(ls -A "$RESTORE_DIR" | wc -l) -eq 0 ]; then
    echo "错误：$1 没有可用备份"
    exit 1
  fi
fi

echo "即将从 $RESTORE_DIR 恢复备份，会覆盖当前工作区的配置和记忆文件，确认继续？(y/N)"
read confirm
if [ "$confirm" != "y" ] && [ "$confirm" != "Y" ]; then
  echo "已取消恢复"
  exit 0
fi

# 执行恢复
cp -r "$RESTORE_DIR"/* "$WORKSPACE/"
echo "恢复完成，已将 $RESTORE_DIR 下的文件覆盖到工作区"
