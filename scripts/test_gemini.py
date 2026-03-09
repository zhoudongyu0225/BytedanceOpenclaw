import requests
import json

API_KEY = "AIzaSyDcrlnSG29fg2GJHLhZ1gDEfwSDctuaLiM"

def list_models():
    url = f"https://generativelanguage.googleapis.com/v1/models?key={API_KEY}"
    response = requests.get(url)
    if response.status_code == 200:
        models = response.json()["models"]
        print("可用的Gemini模型列表：")
        print("-" * 80)
        for model in models:
            if "gemini" in model["name"].lower():
                print(f"模型ID: {model['name'].split('/')[-1]}")
                print(f"显示名称: {model['displayName']}")
                print(f"描述: {model['description']}")
                print(f"支持的方法: {', '.join(model['supportedGenerationMethods'])}")
                print(f"最大输入token: {model.get('inputTokenLimit', 'N/A')}")
                print(f"最大输出token: {model.get('outputTokenLimit', 'N/A')}")
                print("-" * 80)
    else:
        print(f"查询失败，状态码: {response.status_code}")
        print(f"错误信息: {response.text}")

def test_quota():
    # 测试简单调用看配额
    url = f"https://generativelanguage.googleapis.com/v1/models/gemini-1.5-pro-latest:generateContent?key={API_KEY}"
    payload = {
        "contents": [{"parts": [{"text": "你好，简单回复一句即可"}]}]
    }
    headers = {"Content-Type": "application/json"}
    response = requests.post(url, headers=headers, json=payload)
    if response.status_code == 200:
        print("\n✅ Gemini API调用正常，配额可用")
        usage = response.json().get("usageMetadata", {})
        if usage:
            print(f"本次调用消耗：输入token {usage.get('promptTokenCount',0)}，输出token {usage.get('candidatesTokenCount',0)}，总token {usage.get('totalTokenCount',0)}")
    else:
        print(f"\n❌ 调用失败: {response.status_code} - {response.text}")

if __name__ == "__main__":
    list_models()
    test_quota()
