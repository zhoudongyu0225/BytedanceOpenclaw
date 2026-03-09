import requests
import json
import os

API_KEY = os.getenv("ARK_API_KEY")
# 先把参考图上传到临时图床获取公网链接，或者直接用你发的图片的飞书链接
REFERENCE_IMAGE_PATH = "/root/.openclaw/media/inbound/b919e21b-5e46-4167-9a81-ae5ce5d230f6.jpg"

def image_to_image(prompt):
    url = "https://ark.cn-beijing.volces.com/api/v3/images/generations"
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {API_KEY}"
    }
    payload = {
        "model": "doubao-seedream-4-5-251128",
        "prompt": prompt,
        "image_url": REFERENCE_IMAGE_URL,
        "strength": 0.2,  # 低强度，最大程度保留原图特征
        "num_images": 1,
        "size": "1920x1920"
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
    output = image_to_image("22岁戴黑框眼镜的女极客，坐在电脑前敲代码，背景是科技感工作室，写实二次元风格，完全保留参考图的所有面部特征，不要改变五官和发型")
    print(output)
