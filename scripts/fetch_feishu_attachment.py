import requests
import json
import os

# 飞书应用配置
APP_ID = "cli_a926e51e7178dcc0"
APP_SECRET = "nWEZIwPq8hayJJgjkgPScfnoBZdC6qDS"
CHAT_ID = "oc_2273a36513ca66dc7cbcebfe65135b21"

def get_tenant_access_token():
    """获取租户访问令牌"""
    url = "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal"
    payload = json.dumps({
        "app_id": APP_ID,
        "app_secret": APP_SECRET
    })
    headers = {
        "Content-Type": "application/json"
    }
    response = requests.request("POST", url, headers=headers, data=payload)
    return response.json()["tenant_access_token"]

def get_recent_messages(token):
    """获取最近的消息列表"""
    url = f"https://open.feishu.cn/open-apis/im/v1/messages?container_id_type=chat&container_id={CHAT_ID}&page_size=20"
    headers = {
        "Authorization": f"Bearer {token}"
    }
    response = requests.request("GET", url, headers=headers)
    res = response.json()
    print("接口返回:", res)
    return res["data"]["items"]

def download_attachment(token, message_id, file_key, save_path):
    """下载附件"""
    url = f"https://open.feishu.cn/open-apis/im/v1/messages/{message_id}/resources/{file_key}?type=file"
    headers = {
        "Authorization": f"Bearer {token}"
    }
    response = requests.request("GET", url, headers=headers)
    with open(save_path, "wb") as f:
        f.write(response.content)
    print(f"附件已下载到: {save_path}")

if __name__ == "__main__":
    token = get_tenant_access_token()
    messages = get_recent_messages(token)
    
    for msg in messages:
        msg_content = json.loads(msg["body"]["content"])
        if "attachments" in msg_content:
            for att in msg_content["attachments"]:
                file_name = att["name"]
                file_key = att["file_key"]
                save_path = f"/root/.openclaw/workspace/{file_name}"
                download_attachment(token, msg["message_id"], file_key, save_path)
                # 下载最近的一个压缩包就退出
                if file_name.endswith((".zip", ".rar", ".7z", ".tar", ".gz", ".bz2", ".xz")):
                    print(f"已找到压缩包: {file_name}")
                    exit(0)
    
    print("未找到最近的压缩包附件")
