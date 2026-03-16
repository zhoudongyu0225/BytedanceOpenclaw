---
name: wechat-mp-reader
description: 微信公众号文章读取工具，支持抓取公众号文章的正文内容、图片，自动解析微信反爬限制。触发场景：用户发送公众号链接、要求读取公众号文章内容、分析公众号文章。
---
# 微信公众号文章读取工具

## 功能
- 自动抓取微信公众号文章的完整正文内容
- 支持提取文章中的图片、标题、作者、发布时间等元信息
- 自动绕过微信公众号反爬机制

## 使用方法
当用户提供微信公众号链接时，运行scripts/fetch_mp_article.py脚本：
```bash
python scripts/fetch_mp_article.py <公众号文章URL>
```
脚本会返回Markdown格式的文章内容。
