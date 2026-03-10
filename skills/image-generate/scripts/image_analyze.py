import os
import sys
import base64
from volcenginesdkarkruntime import Ark

DEFAULT_VISION_MODEL = "ep-20260310152101-7r8tp"

def analyze_image(image_path: str):
    """分析图片内容，提取详细特征
    
    Args:
        image_path: 本地图片路径
    """
    if not os.path.exists(image_path):
        print(f"图片路径不存在：{image_path}")
        return None
    
    api_key = os.getenv("ARK_API_KEY")
    client = Ark(api_key=api_key)
    
    try:
        # 读取图片并编码为base64
        with open(image_path, "rb") as f:
            image_base64 = base64.b64encode(f.read()).decode("utf-8")
        
        response = client.chat.completions.create(
            model=os.getenv("MODEL_VISION_NAME", DEFAULT_VISION_MODEL),
            messages=[
                {
                    "role": "user",
                    "content": [
                        {"type": "text", "text": "请详细描述这个动漫角色形象的所有特征：包括发型、发色、发饰、五官（眼睛颜色、形状）、服饰（款式、颜色、细节）、画风风格、角色整体特征，描述越详细越好，方便后续用这个描述生成完全一致的形象"},
                        {"type": "image_url", "image_url": {"url": f"data:image/jpeg;base64,{image_base64}"}}
                    ]
                }
            ]
        )
        
        result = response.choices[0].message.content
        print("=== 图片特征分析结果 ===")
        print(result)
        return result
    except Exception as e:
        print(f"分析图片失败：{e}")
        return None

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python image_analyze.py <image_path>")
        sys.exit(1)
    image_path = sys.argv[1]
    analyze_image(image_path)
