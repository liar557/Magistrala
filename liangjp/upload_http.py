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

class SenMLGenerator:
    """SenML消息生成器"""
    
    @staticmethod
    def generate_exact_data(base_name: str) -> list:
        """生成与示例完全一致的SenML消息"""
        return [
            {
                "bn": base_name,
                "n": "lumen",
                "u": "CD",
                "v": round(random.uniform(10, 100), 1)
            },
            {
                "bn": base_name,  # 保留bn以保持结构一致
                "n": "speed",
                "u": "mach",
                "v": round(random.uniform(10, 100), 1)
            }
        ]

def upload_senml_data(config: dict) -> bool:
    """上传SenML数据到HTTP服务器"""
    # 构建URL
    url = f"{config['base_url']}/http/m/{config['domain_id']}/c/{config['channel_id']}/{config['subtopic']}"
    
    # 从配置中获取base_name
    base_name = config.get("base_name", "test")
    logger.info(f"使用base_name: {base_name}")
    
    # 生成精确格式的SenML数据
    senml_data = SenMLGenerator.generate_exact_data(base_name)
    json_data = json.dumps(senml_data)
    
    # 设置请求头
    headers = {
        "Content-Type": "application/senml+json",
        "Authorization": f"Client {config['client_secret']}"
    }
    
    # 打印完整请求信息
    logger.info(f"===== 开始上传请求 =====")
    logger.info(f"请求URL: {url}")
    logger.info(f"请求头: {headers}")
    logger.info(f"请求体: {json_data}")
    
    try:
        # 发送请求
        response = requests.post(
            url,
            data=json_data,
            headers=headers,
            verify=config.get("ca_cert_path", True),
            timeout=config.get("timeout", 10)
        )
        
        # 打印完整响应信息
        logger.info(f"响应状态码: {response.status_code}")
        if response.text:
            logger.info(f"响应内容: {response.text}")
        else:
            logger.info("响应内容为空")
        logger.info("===== 请求结束 =====")
        
        # 检查响应状态
        response.raise_for_status()
        logger.info("上传请求成功")
        return True
    except requests.exceptions.RequestException as e:
        logger.error(f"上传请求失败: {str(e)}")
        # 补充错误时的响应信息
        if hasattr(e, 'response') and e.response:
            logger.error(f"错误响应状态码: {e.response.status_code}")
            if e.response.text:
                logger.error(f"错误响应内容: {e.response.text}")
        return False

def main():
    """主函数 - 直接在代码中配置参数"""
    # 配置参数（与curl命令完全一致）
    config = {
        "base_url": "https://localhost",                                    # 基础URL
        "domain_id": "562d704a-c442-499a-aff3-223f580bf6b3",                # 域ID
        "channel_id": "b0ec13df-9ff0-48b9-9cb6-b3be072e7c99",               # 通道ID
        "subtopic": "volume",                                               # 子主题
        "client_secret": "4fa8890c-5888-48ba-93b8-3e7db1165b65",            # 客户端密钥
        "ca_cert_path": "CA/ca.crt",                                        # CA证书路径
        "base_name": "ljp",                                                 # SenML基础名称
        "timeout": 10                                                       # 请求超时时间
    }
    
    # 上传模式配置
    upload_mode = "interval"       # 单次上传模式（便于调试）
    # upload_mode = "once"  # 定时上传模式
    
    if upload_mode == "once":
        # 只上传一次
        logger.info("开始单次上传（与curl命令对比）")
        success = upload_senml_data(config)
        exit_code = 0 if success else 1
        logger.info(f"上传结果: {'成功' if success else '失败'}")
        exit(exit_code)
    else:
        # 定时上传
        upload_interval = 30      # 上传间隔(秒)
        logger.info(f"开始定时上传，间隔: {upload_interval}秒")
        try:
            while True:
                upload_senml_data(config)
                time.sleep(upload_interval)
        except KeyboardInterrupt:
            logger.info("接收到中断信号，程序退出")

if __name__ == "__main__":
    main()