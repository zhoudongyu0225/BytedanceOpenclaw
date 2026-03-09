---
name: asset-auto-processor
description: 游戏素材自动化处理技能，接收参考图片自动拆解元素，输出带透明通道的PNG素材，适合UE游戏工程使用。触发场景：用户发送参考图要求拆解素材、生成带通道PNG、游戏UI素材处理。
---

# 素材自动化处理技能

## 功能说明
自动处理用户上传的参考图片，执行以下操作：
1. 自动识别图片中的独立元素
2. 抠除背景生成带alpha透明通道的PNG素材
3. 批量导出到SVN素材目录`/data/svn/dino-defense/assets/`
4. 自动命名分类存储，方便直接导入UE工程

## 依赖工具
- opencv-python：图像处理
- rembg：背景抠除
- pillow：图片格式转换

## 使用方法
用户发送参考图片后，自动调用`scripts/process_asset.py`脚本处理，输出结果到SVN目录并告知用户存储路径。
