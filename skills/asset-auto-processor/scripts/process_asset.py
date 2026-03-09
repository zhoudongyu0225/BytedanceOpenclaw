#!/usr/bin/env python3
import os
import sys
from rembg import remove
from PIL import Image
import cv2
import numpy as np

OUTPUT_DIR = "/data/svn/dino-defense/assets/processed/"
os.makedirs(OUTPUT_DIR, exist_ok=True)

def process_image(input_path):
    # 读取图片
    img = Image.open(input_path)
    # 抠除背景
    img_no_bg = remove(img)
    # 生成输出文件名
    base_name = os.path.splitext(os.path.basename(input_path))[0]
    output_path = os.path.join(OUTPUT_DIR, f"{base_name}_transparent.png")
    # 保存带通道的PNG
    img_no_bg.save(output_path, "PNG")
    print(f"素材处理完成，已保存到：{output_path}")
    return output_path

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("用法：python process_asset.py <输入图片路径>")
        sys.exit(1)
    input_path = sys.argv[1]
    if not os.path.exists(input_path):
        print(f"错误：文件 {input_path} 不存在")
        sys.exit(1)
    process_image(input_path)
