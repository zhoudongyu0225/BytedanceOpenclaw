import requests
import json

APP_ID = "cli_a926e51e7178dcc0"
APP_SECRET = "nWEZIwPq8hayJJgjkgPScfnoBZdC6qDS"
MESSAGE_ID = "om_x100b55f306af9ca0c38ad728b4663d1"

def get_tenant_access_token():
    url = "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal"
    payload = json.dumps({"app_id": APP_ID, "app_secret": APP_SECRET})
    headers = {"Content-Type": "application/json"}
    response = requests.post(url, headers=headers, data=payload)
    return response.json()["tenant_access_token"]

def get_message_content(token):
    url = f"https://open.feishu.cn/open-apis/im/v1/messages/{MESSAGE_ID}"
    headers = {"Authorization": f"Bearer {token}"}
    response = requests.get(url, headers=headers)
    res = response.json()
    content = json.loads(res["data"]["items"][0]["body"]["content"])
    print(f"附件实际内容: {content}")
    return content

if __name__ == "__main__":
    token = get_tenant_access_token()
    get_message_content(token)
