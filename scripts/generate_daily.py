import requests
import json
import os
import base64

API_KEY = os.getenv("ARK_API_KEY")
REFERENCE_IMAGE_PATH = "/root/.openclaw/media/inbound/b919e21b-5e46-4167-9a81-ae5ce5d230f6.jpg"

def image_to_image(prompt):
    # 读取本地参考图转base64
    with open(REFERENCE_IMAGE_PATH, "rb") as f:
        img_data = f.read()
        img_base64 = base64.b64encode(img_data).decode("utf-8")
    
    url = "https://ark.cn-beijing.volces.com/api/v3/images/generations"
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {API_KEY}"
    }
    payload = {
        "model": "doubao-seedream-4-5-251128",
        "prompt": prompt,
        "image": img_base64,
        "strength": 0.2,
        "num_images": 1,
        "size": "1920x1920"
    }
    
    response = requests.post(url, headers=headers, json=payload)
    if response.status_code == 200:
        data = response.json()
        image_url = data["data"][0]["url"]
        img_response = requests.get(image_url)
        output_path = f"/root/.openclaw/workspace/generated/daily_{int(os.times()[4])}.png"
        os.makedirs(os.path.dirname(output_path), exist_ok=True)
        with open(output_path, "wb") as f:
            f.write(img_response.content)
        return output_path
    else:
        print(f"Error: {response.status_code} - {response.text}")
        return None

if __name__ == "__main__":
    output = image_to_image("22岁戴黑框眼镜的女生，日常居家场景，穿休闲卫衣坐在沙发上撸猫，温馨治愈风格，完全保留参考图的所有面部特征，不要改变五官和发型")
    print(output)
