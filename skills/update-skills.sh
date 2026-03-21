#!/bin/bash
# Skill 自动更新脚本 - 每天执行一次
DATE=$(date +%Y-%m-%d)
LOG="=== Skill Update $DATE ==="
SKILLS_DIR="$HOME/.openclaw/workspace/skills"

echo "$LOG" >> $HOME/.openclaw/skills-update.log
cd $SKILLS_DIR

# 已安装的需要更新的skill列表
REPOS=(
  "pskoett/self-improving-agent:self-improving-agent"
  "nesdeq/openclaw-feeds:openclaw-feeds"
  "jimliuxinghai/find-skills:find-skills"
  "spclaudehome/skill-vetter:skill-vetter"
  "savorgbot-exe/focus-mode:focus-mode"
  "thesethrose/agent-browser:agent-browser"
  "jeffjhunter/ai-daily-briefing:ai-daily-briefing"
  "globalcaos/youtube-ultimate:youtube-ultimate"
  "thesethrose/reddit-search:reddit-search"
)

for repo in "${REPOS[@]}"; do
  path=$(echo $repo | cut -d: -f2)
  if [ -d "$path" ]; then
    echo "--- Updating $path ---" >> $HOME/.openclaw/skills-update.log
    cd $path && git pull --force 2>&1 | tail -3 >> $HOME/.openclaw/skills-update.log && cd $SKILLS_DIR
  else
    echo "--- Cloning $path ---" >> $HOME/.openclaw/skills-update.log
    git clone --depth=1 "https://github.com/$(echo $repo | cut -d: -f1).git" $path 2>&1 | tail -3 >> $HOME/.openclaw/skills-update.log
  fi
done

echo "=== Update Done $DATE ===" >> $HOME/.openclaw/skills-update.log
