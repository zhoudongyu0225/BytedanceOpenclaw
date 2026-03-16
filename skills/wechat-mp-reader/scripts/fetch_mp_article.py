#!/usr/bin/env python3
import sys
import requests
from bs4 import BeautifulSoup
import markdownify
import re

def fetch_wechat_article(url):
    headers = {
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36',
        'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7',
        'Accept-Language': 'zh-CN,zh;q=0.9,en;q=0.8',
        'Cache-Control': 'max-age=0',
        'Connection': 'keep-alive',
    }

    try:
        response = requests.get(url, headers=headers, timeout=10)
        response.encoding = 'utf-8'
        soup = BeautifulSoup(response.text, 'html.parser')
        
        # 提取标题
        title = soup.find('h1', id='activity-name')
        title = title.get_text(strip=True) if title else '无标题'
        
        # 提取作者
        author = soup.find('span', id='js_name')
        author = author.get_text(strip=True) if author else '未知作者'
        
        # 提取发布时间
        publish_time = soup.find('em', id='publish_time')
        publish_time = publish_time.get_text(strip=True) if publish_time else '未知时间'
        
        # 提取正文
        content = soup.find('div', id='js_content')
        if not content:
            return "无法获取文章内容，可能需要手动验证或文章已被删除"
        
        # 转换为Markdown
        md_content = markdownify.markdownify(str(content), heading_style="ATX")
        
        # 清理多余的空行
        md_content = re.sub(r'\n{3,}', r'\n\n', md_content)
        
        # 组装结果
        result = f"# {title}\n\n"
        result += f"作者：{author}\n"
        result += f"发布时间：{publish_time}\n\n"
        result += "---\n\n"
        result += md_content
        
        return result
        
    except Exception as e:
        return f"抓取失败：{str(e)}"

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("用法：python fetch_mp_article.py <公众号文章URL>")
        sys.exit(1)
    
    url = sys.argv[1]
    print(fetch_wechat_article(url))
