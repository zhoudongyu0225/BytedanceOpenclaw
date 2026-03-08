import requests
import json

APP_ID = "cli_a926e51e7178dcc0"
APP_SECRET = "nWEZIwPq8hayJJgjkgPScfnoBZdC6qDS"
FILE_KEY = "file_v3_00vj_62db74d6-d278-4afe-8358-6011ba2045cg"
FILE_NAME = "commonDmGameServer.zip"
MESSAGE_ID = "om_placeholder" # 可以留空，只要file_key正确即可

def get_tenant_access_token():
    url = "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal"
    payload = json.dumps({"app_id": APP_ID, "app_secret": APP_SECRET})
    headers = {"Content-Type": "application/json"}
    response = requests.post(url, headers=headers, data=payload)
    return response.json()["tenant_access_token"]

def download_attachment(token):
    url = f"https://open.feishu.cn/open-apis/im/v1/messages/{MESSAGE_ID}/resources/{FILE_KEY}?type=file"
    headers = {"Authorization": f"Bearer {token}"}
    response = requests.get(url, headers=headers)
    save_path = f"/root/.openclaw/workspace/{FILE_NAME}"
    with open(save_path, "wb") as f:
        f.write(response.content)
    print(f"压缩包已成功下载到: {save_path}")
    return save_path

if __name__ == "__main__":
    token = get_tenant_access_token()
    save_path = download_attachment(token)
    # 自动解压
    import zipfile
    extract_dir = f"/root/.openclaw/workspace/commonDmGameServer"
    with zipfile.ZipFile(save_path, 'r') as zip_ref:
        zip_ref.extractall(extract_dir)
    print(f"压缩包已解压到: {extract_dir}")
    # 列出目录结构
    import os
    print("\n目录结构：")
    for root, dirs, files in os.walk(extract_dir):
        level = root.replace(extract_dir, '').count(os.sep)
        indent = ' ' * 2 * level
        print(f"{indent}{os.path.basename(root)}/")
        subindent = ' ' * 2 * (level + 1)
        for file in files:
            print(f"{subindent}{file}")
