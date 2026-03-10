import os
import sys
import base64
from volcenginesdkarkruntime import Ark

DEFAULT_MODEL = "ep-m-20260309005340-fvp7h"

def image_to_image(prompt: str, image_path: str, strength: float = 0.3):
    """基于参考图生成图片
    
    Args:
        prompt: 生成提示词
        image_path: 参考图本地路径
        strength: 参考图相似度，0-1，越小越接近参考图
    """
    if not os.path.exists(image_path):
        print(f"参考图路径不存在：{image_path}")
        return

    api_key = os.getenv("ARK_API_KEY")
    client = Ark(api_key=api_key)
    
    try:
        # 读取参考图并编码
        with open(image_path, "rb") as f:
            image_base64 = base64.b64encode(f.read()).decode("utf-8")
        
        response = client.images.generate(
            model=os.getenv("MODEL_IMAGE_NAME", DEFAULT_MODEL),
            prompt=prompt,
            image=f"data:image/jpeg;base64,{image_base64}",
            strength=strength
        )

        download_dir = os.getenv("IMAGE_DOWNLOAD_DIR", os.path.expanduser("./"))
        if not os.path.exists(download_dir):
            os.makedirs(download_dir, exist_ok=True)

        for i, image in enumerate(response.data):
            timestamp = int(time.time())
            filename = f"generated_image_ref_{timestamp}_{i}.png"
            filepath = os.path.join(download_dir, filename)
            urllib.request.urlretrieve(image.url, filepath)
            print(f"基于参考图生成成功：{filepath}")
    except Exception as e:
        print(f"图生图失败：{e}")

if __name__ == "__main__":
    if len(sys.argv) < 3:
        print("Usage: python image_to_image.py <prompt> <image_path> [strength=0.3]")
        sys.exit(1)
    prompt = sys.argv[1]
    image_path = sys.argv[2]
    strength = float(sys.argv[3]) if len(sys.argv) > 3 else 0.3
    image_to_image(prompt, image_path, strength)
