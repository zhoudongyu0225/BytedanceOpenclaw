import requests
import json

APP_ID = "cli_a926e51e7178dcc0"
APP_SECRET = "nWEZIwPq8hayJJgjkgPScfnoBZdC6qDS"

def get_tenant_access_token():
    url = "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal"
    payload = json.dumps({"app_id": APP_ID, "app_secret": APP_SECRET})
    headers = {"Content-Type": "application/json"}
    response = requests.post(url, headers=headers, data=payload)
    return response.json()["tenant_access_token"]

def get_chats(token):
    url = "https://open.feishu.cn/open-apis/im/v1/chats?page_size=50"
    headers = {"Authorization": f"Bearer {token}"}
    response = requests.get(url, headers=headers)
    res = response.json()
    print("会话列表：")
    for chat in res["data"]["items"]:
        print(f"chat_id: {chat['chat_id']}, 名称: {chat['name']}, 类型: {chat['chat_type']}")
    return res["data"]["items"]

if __name__ == "__main__":
    token = get_tenant_access_token()
    get_chats(token)
