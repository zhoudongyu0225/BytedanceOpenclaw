import requests
import json
import os
import zipfile

APP_ID = "cli_a926e51e7178dcc0"
APP_SECRET = "nWEZIwPq8hayJJgjkgPScfnoBZdC6qDS"
CHAT_ID = "oc_2273a36513ca66dc7cbcebfe65135b21"
FILE_KEY = "file_v3_00vj_6131e214-f900-46ef-b969-966cc870b13g"
FILE_NAME = "commonDmGameServer.zip"

def get_tenant_access_token():
    url = "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal"
    payload = json.dumps({"app_id": APP_ID, "app_secret": APP_SECRET})
    headers = {"Content-Type": "application/json"}
    response = requests.post(url, headers=headers, data=payload)
    return response.json()["tenant_access_token"]

def get_all_messages(token):
    url = f"https://open.feishu.cn/open-apis/im/v1/messages?container_id_type=chat&container_id={CHAT_ID}&page_size=50"
    headers = {"Authorization": f"Bearer {token}"}
    response = requests.get(url, headers=headers)
    res = response.json()
    for msg in res["data"]["items"]:
        print(f"message_id: {msg['message_id']}, type: {msg['msg_type']}, time: {msg['create_time']}")
        if msg["msg_type"] == "file":
            content = json.loads(msg["body"]["content"])
            if content["file_key"] == FILE_KEY:
                print(f"\n找到匹配的附件！message_id: {msg['message_id']}")
                return msg["message_id"]
    return None

def download_attachment(token, message_id):
    url = f"https://open.feishu.cn/open-apis/im/v1/messages/{message_id}/resources/{FILE_KEY}?type=file"
    headers = {"Authorization": f"Bearer {token}"}
    response = requests.get(url, headers=headers)
    save_path = f"/root/.openclaw/workspace/{FILE_NAME}"
    with open(save_path, "wb") as f:
        f.write(response.content)
    print(f"✅ 压缩包下载成功: {save_path}, 大小: {os.path.getsize(save_path)/1024/1024:.2f}MB")
    
    # 解压
    extract_dir = f"/root/.openclaw/workspace/{os.path.splitext(FILE_NAME)[0]}"
    os.makedirs(extract_dir, exist_ok=True)
    try:
        with zipfile.ZipFile(save_path, 'r') as zf:
            zf.extractall(extract_dir)
        print(f"✅ 解压成功: {extract_dir}")
        # 列出目录结构
        print("\n📁 目录结构:")
        for root, dirs, files in os.walk(extract_dir):
            level = root.replace(extract_dir, '').count(os.sep)
            indent = ' ' * 2 * level
            print(f"{indent}{os.path.basename(root)}/")
            subindent = ' ' * 2 * (level + 1)
            for file in files:
                print(f"{subindent}{file}")
    except Exception as e:
        print(f"❌ 解压失败: {e}")
        import subprocess
        subprocess.run(["file", save_path])

if __name__ == "__main__":
    token = get_tenant_access_token()
    message_id = get_all_messages(token)
    if message_id:
        download_attachment(token, message_id)
    else:
        print("❌ 未找到匹配的附件消息")
