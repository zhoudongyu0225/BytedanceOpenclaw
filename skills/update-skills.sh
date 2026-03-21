#!/bin/bash
# Skill 自动更新脚本 - 每天执行一次
DATE=$(date +%Y-%m-%d)
SKILLS_DIR="$HOME/.openclaw/workspace/skills"
LOG="$HOME/.openclaw/skills-update.log"

echo "=== Skill Update $DATE ===" >> $LOG
cd $SKILLS_DIR

BASE_URL="https://raw.githubusercontent.com/openclaw/skills/main/skills"

# 需要更新的 skills 列表
SKILLS=(
  "jimliuxinghai/find-skills"
  "spclaudehome/skill-vetter"
  "savorgbot-exe/focus-mode"
  "thesethrose/agent-browser"
  "jeffjhunter/ai-daily-briefing"
  "globalcaos/youtube-ultimate"
  "thesethrose/reddit-search"
  "pskoett/self-improving-agent"
  "nesdeq/openclaw-feeds"
)

for item in "${SKILLS[@]}"; do
  name=$(echo $item | cut -d/ -f2)
  
  if [ -d "$name" ] && [ -f "$name/SKILL.md" ]; then
    # 已有skill，只更新SKILL.md（raw文件下载）
    echo "Updating $name..." >> $LOG
    curl -s --max-time 10 "$BASE_URL/$item/SKILL.md" -o "$name/SKILL.md" 2>&1 | head -1 >> $LOG
    echo "  $name SKILL.md updated" >> $LOG
  elif [ -d "$name" ]; then
    # 目录存在但没有SKILL.md，尝试下载
    echo "Fetching $name SKILL.md..." >> $LOG
    curl -s --max-time 10 "$BASE_URL/$item/SKILL.md" -o "$name/SKILL.md" 2>&1 | head -1 >> $LOG
  else
    # 不存在，创建并下载
    mkdir -p "$name"
    echo "Cloning $name..." >> $LOG
    for f in "SKILL.md" "README.md"; do
      curl -s --max-time 10 "$BASE_URL/$item/$f" -o "$name/$f" 2>&1 | head -1 >> $LOG
    done
    echo "  $name installed" >> $LOG
  fi
done

echo "=== Update Done $DATE ===" >> $LOG
