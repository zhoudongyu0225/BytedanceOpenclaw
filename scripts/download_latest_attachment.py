import requests
import json
import os
import zipfile

APP_ID = "cli_a926e51e7178dcc0"
APP_SECRET = "nWEZIwPq8hayJJgjkgPScfnoBZdC6qDS"
CHAT_ID = "oc_2273a36513ca66dc7cbcebfe65135b21"

def get_tenant_access_token():
    url = "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal"
    payload = json.dumps({"app_id": APP_ID, "app_secret": APP_SECRET})
    headers = {"Content-Type": "application/json"}
    response = requests.post(url, headers=headers, data=payload)
    return response.json()["tenant_access_token"]

def get_latest_message(token):
    # 拉取最新10条消息
    url = f"https://open.feishu.cn/open-apis/im/v1/messages?container_id_type=chat&container_id={CHAT_ID}&page_size=10"
    headers = {"Authorization": f"Bearer {token}"}
    response = requests.get(url, headers=headers)
    res = response.json()
    return res["data"]["items"]

def download_attachment(token, message_id, file_key, file_name):
    url = f"https://open.feishu.cn/open-apis/im/v1/messages/{message_id}/resources/{file_key}?type=file"
    headers = {"Authorization": f"Bearer {token}"}
    response = requests.get(url, headers=headers)
    save_path = f"/root/.openclaw/workspace/{file_name}"
    with open(save_path, "wb") as f:
        f.write(response.content)
    print(f"压缩包已成功下载到: {save_path}")
    return save_path

if __name__ == "__main__":
    token = get_tenant_access_token()
    messages = get_latest_message(token)
    
    for msg in messages:
        if msg["msg_type"] == "file":
            msg_content = json.loads(msg["body"]["content"])
            file_key = msg_content["file_key"]
            file_name = msg_content["file_name"]
            message_id = msg["message_id"]
            print(f"找到附件: {file_name}, message_id: {message_id}")
            
            # 下载
            save_path = download_attachment(token, message_id, file_key, file_name)
            
            # 解压
            extract_dir = f"/root/.openclaw/workspace/{os.path.splitext(file_name)[0]}"
            os.makedirs(extract_dir, exist_ok=True)
            try:
                with zipfile.ZipFile(save_path, 'r') as zip_ref:
                    zip_ref.extractall(extract_dir)
                print(f"压缩包已解压到: {extract_dir}")
                # 列出目录结构
                print("\n目录结构：")
                for root, dirs, files in os.walk(extract_dir):
                    level = root.replace(extract_dir, '').count(os.sep)
                    indent = ' ' * 2 * level
                    print(f"{indent}{os.path.basename(root)}/")
                    subindent = ' ' * 2 * (level + 1)
                    for file in files:
                        print(f"{subindent}{file}")
            except Exception as e:
                print(f"解压失败: {e}")
                # 查看文件类型
                import subprocess
                subprocess.run(["file", save_path])
            
            exit(0)
    
    print("未找到附件消息")
