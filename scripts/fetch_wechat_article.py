import requests
from bs4 import BeautifulSoup
import sys

def fetch_wechat(url):
    headers = {
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36'
    }
    response = requests.get(url, headers=headers)
    if response.status_code != 200:
        print(f"请求失败，状态码：{response.status_code}")
        return
    
    soup = BeautifulSoup(response.text, 'html.parser')
    # 提取微信文章正文
    content = soup.find(id='js_content')
    if not content:
        print("未找到文章内容")
        return
    
    # 提取文本
    text = content.get_text(strip=True, separator='\n')
    print(text)

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("使用方法：python fetch_wechat_article.py <微信文章链接>")
        sys.exit(1)
    fetch_wechat(sys.argv[1])
