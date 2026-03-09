import requests
import json
import os
import base64

API_KEY = os.getenv("ARK_API_KEY")
REFERENCE_IMAGE_PATH = "/root/.openclaw/media/inbound/30f020d0-8f16-4677-a446-cb931016788a.jpg"  # 用户给的参考形象图

def image_to_image(prompt):
    # 读取参考图转base64
    with open(REFERENCE_IMAGE_PATH, "rb") as f:
        img_base64 = base64.b64encode(f.read()).decode("utf-8")
    
    url = "https://ark.cn-beijing.volces.com/api/v3/images/generations"
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {API_KEY}"
    }
    payload = {
        "model": "doubao-seedream-4-5-251128",
        "prompt": prompt,
        "image": img_base64,
        "image_weight": 0.9,  # 参考图权重拉满，最大程度保留面部特征
        "num_images": 1,
        "size": "1024x1024"
    }
    
    response = requests.post(url, headers=headers, json=payload)
    if response.status_code == 200:
        data = response.json()
        image_url = data["data"][0]["url"]
        # 下载图片
        img_response = requests.get(image_url)
        output_path = f"/root/.openclaw/workspace/generated/role_{int(os.times()[4])}.png"
        os.makedirs(os.path.dirname(output_path), exist_ok=True)
        with open(output_path, "wb") as f:
            f.write(img_response.content)
        print(f"✅ 生成成功: {output_path}")
        return output_path
    else:
        print(f"❌ 生成失败: {response.status_code} - {response.text}")
        return None

if __name__ == "__main__":
    # 测试生成工作场景图
    output = image_to_image("22岁戴黑框眼镜的女极客，坐在电脑前敲代码，背景是科技感工作室，写实二次元风格，保留参考图的所有面部特征")
    print(output)
