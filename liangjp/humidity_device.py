import time
import json
import random
import logging
import requests
from datetime import datetime

# 配置日志
logging.basicConfig(
    level=logging.DEBUG,
    format='%(asctime)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

class HumiditySensor:
    """湿度设备数据生成器"""
    @staticmethod
    def generate_senml_data(base_name: str) -> list:
        """生成湿度传感器的SenML数据"""
        return [
            {
                "bn": base_name,
                "n": "Humidity",   # 数据名称：湿度
                "u": "%RH",        # 单位：相对湿度百分比
                "v": round(random.uniform(30, 80), 1)  # 湿度范围：30-80%RH
            }
        ]

def upload_senml_data(config: dict) -> bool:
    """上传SenML数据到服务器"""
    # 构建URL
    url = f"{config['base_url']}/http/m/{config['domain_id']}/c/{config['channel_id']}/{config['subtopic']}"
    
    # 生成湿度数据
    senml_data = HumiditySensor.generate_senml_data(config["base_name"])
    json_data = json.dumps(senml_data)
    
    # 请求头
    headers = {
        "Content-Type": "application/senml+json",
        "Authorization": f"Client {config['client_secret']}"
    }
    
    # 打印请求信息
    logger.info(f"===== 湿度设备上传 =====")
    logger.info(f"URL: {url}")
    logger.info(f"请求体: {json_data}")
    
    try:
        # 发送请求
        response = requests.post(
            url,
            data=json_data,
            headers=headers,
            verify=config["ca_cert_path"],
            timeout=config["timeout"]
        )
        response.raise_for_status()
        logger.info(f"上传成功，状态码: {response.status_code}")
        return True
    except requests.exceptions.RequestException as e:
        logger.error(f"上传失败: {str(e)}")
        return False

def main():
    # 湿度设备专属配置
    config = {
        "base_url": "https://localhost",                                    # 基础URL
        "domain_id": "562d704a-c442-499a-aff3-223f580bf6b3",                # 域ID
        "channel_id": "b0ec13df-9ff0-48b9-9cb6-b3be072e7c99",               # 通道ID
        "subtopic": "humidity",                                               # 子主题
        "client_secret": "267d0a2d-4426-4e08-82b7-d17f99ba4bce",            # 客户端密钥
        "ca_cert_path": "CA/ca.crt",                                        # CA证书路径
        "base_name": "test",                                            # SenML基础名称
        "timeout": 10                                                       # 请求超时时间
    }
    
    # 上传模式
    upload_mode = "interval"  # 单次上传：once；定时上传：interval
    upload_interval = 30  # 定时上传间隔（秒）
    
    if upload_mode == "once":
        logger.info("开始单次上传")
        success = upload_senml_data(config)
        exit(0 if success else 1)
    else:
        logger.info(f"开始定时上传（间隔{upload_interval}秒）")
        try:
            while True:
                upload_senml_data(config)
                time.sleep(upload_interval)
        except KeyboardInterrupt:
            logger.info("程序已退出")

if __name__ == "__main__":
    main()