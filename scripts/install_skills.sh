#!/bin/bash
SKILLS=(
    "ontology"
    "self-improving-agent"
    "summarize"
    "proactive-agent"
    "multi-search-engine"
    "tavily-search"
    "find-skills"
    "youtube-watcher"
    "auto-updater"
    "api-gateway"
    "desktop-control"
    "automation-workflows"
    "browser-use"
)

echo "开始安装技能，已跳过已存在的skill-creator..."
echo "共需要安装 ${#SKILLS[@]} 个技能"
echo "==================================="

for skill in "${SKILLS[@]}"; do
    echo -e "\n🔧 安装技能: $skill"
    npx clawhub@latest install "$skill" --force
    if [ $? -eq 0 ]; then
        echo "✅ $skill 安装成功"
    else
        echo "❌ $skill 安装失败，将尝试从GitHub搜索安装"
        # 尝试从GitHub搜索
        repo=$(gh search repo "$skill openclaw skill" --limit 1 | awk '{print $1}' | head -1)
        if [ -n "$repo" ] && [ "$repo" != "No" ]; then
            echo "找到仓库: $repo，正在克隆..."
            git clone "https://github.com/$repo.git" ~/.openclaw/workspace/skills/$skill
            if [ $? -eq 0 ]; then
                echo "✅ $skill 从GitHub安装成功"
            else
                echo "❌ $skill 所有安装方式失败，后续手动实现核心能力"
            fi
        else
            echo "❌ 未找到 $skill 公共仓库，后续手动实现核心能力"
        fi
    fi
done

echo -e "\n==================================="
echo "技能安装流程完成，已安装的技能列表："
ls -la ~/.openclaw/workspace/skills/
