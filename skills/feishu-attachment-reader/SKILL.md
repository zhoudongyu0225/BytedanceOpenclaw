---
name: feishu-attachment-reader
description: 自动读取飞书会话中的所有类型附件，支持下载、解压、内容提取。触发场景：用户发送附件（图片/视频/文档/压缩包等）、要求读取/解析附件内容。
---
# Feishu Attachment Reader Skill
## 功能说明
自动识别并处理飞书会话中的所有类型附件，无需手动操作：
1. **支持的附件类型**：
   - 压缩包：zip/rar/7z/tar/gz/bz2/xz → 自动解压到指定目录
   - 文档：docx/xlsx/ppt/pdf/md/txt → 自动提取文本内容
   - 媒体：图片(jpg/png/gif/webp)/视频(mp4/avi/mov)/音频 → 自动识别元信息，支持生成缩略图
   - 其他：所有类型文件自动下载保存

2. **核心能力**：
   - 手动触发：拉取历史消息中的所有附件
   - 自动监听：新消息收到附件时自动下载处理
   - 内容提取：结构化输出文档/压缩包的内容清单
   - 持久化：所有附件默认保存到 `/root/.openclaw/workspace/attachments/` 目录

## 使用方法
- 直接发送附件到会话，自动处理
- 手动触发：`读取最近的附件` / `解析[文件名]`
- 输出内容清单：`列出[文件名]的内容`

## 依赖安装
```bash
apt install -y unrar p7zip-full poppler-utils
pip install python-docx openpyxl PyPDF2 pillow
```
