import time
import json
import logging
import requests
from datetime import datetime

# 配置日志
logging.basicConfig(
    level=logging.DEBUG,
    format='%(asctime)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)


class MessageFetcher:
    """消息获取器 - 从API获取消息数据"""
    
    @staticmethod
    def format_response(response_data: dict, offset: int, limit: int) -> dict:                            
        """格式化响应数据为指定格式"""
        return {
            "offset": offset,
            "limit": limit,
            "format": "messages",
            "total": response_data.get("total", 0),
            "messages": response_data.get("messages", [])
        }


def fetch_messages(config: dict) -> dict:
    """从API获取消息数据"""
    # 构建URL
    service_port = config["service_port"]
    domain_id = config["domain_id"]
    channel_id = config["channel_id"]
    offset = config["offset"]
    limit = config["limit"]
    
    url = f"http://localhost:{service_port}/{domain_id}/channels/{channel_id}/messages"
    params = {
        "offset": offset,
        "limit": limit
    }
    
    # 设置请求头
    headers = {
        "Authorization": f"Client {config['client_secret']}",
        "Content-Type": "application/json"
    }
    
    # 打印完整请求信息
    logger.info("===== 开始获取消息请求 =====")
    logger.info(f"请求URL: {url}")
    logger.info(f"请求参数: {params}")
    logger.info(f"请求头: {headers}")
    
    try:
        # 发送GET请求
        response = requests.get(
            url,
            params=params,
            headers=headers,
            timeout=config.get("timeout", 10)
        )
        
        # 打印完整响应信息
        logger.info(f"响应状态码: {response.status_code}")
        if response.text:
            response_data = response.json()
            logger.info(f"响应内容: {json.dumps(response_data, indent=2)}")
        else:
            logger.info("响应内容为空")
        logger.info("===== 请求结束 =====")
        
        # 检查响应状态
        response.raise_for_status()
        
        # 格式化响应数据
        formatted_result = MessageFetcher.format_response(
            response_data, 
            offset, 
            limit
        )
        logger.info("消息获取成功")
        return formatted_result
        
    except requests.exceptions.RequestException as e:
        logger.error(f"获取消息失败: {str(e)}")
        # 补充错误时的响应信息
        if hasattr(e, 'response') and e.response:
            logger.error(f"错误响应状态码: {e.response.status_code}")
            if e.response.text:
                logger.error(f"错误响应内容: {e.response.text}")
        return {"error": str(e)}
    except json.JSONDecodeError as e:
        logger.error(f"解析响应JSON失败: {str(e)}")
        logger.error(f"原始响应: {response.text}")
        return {"error": f"解析响应失败: {e}"}


def main():
    """主函数 - 直接在代码中配置参数"""
    # 配置参数（与API要求完全一致）
    config = {
        "service_port": 9011,                                               # 服务端口 (Timescale: 9011, Postgres: 9009)
        "domain_id": "562d704a-c442-499a-aff3-223f580bf6b3",                # 域ID
        "channel_id": "b0ec13df-9ff0-48b9-9cb6-b3be072e7c99",               # 通道ID
        "offset": 0,                                                        # 偏移量
        "limit": 10,                                                        # 每页数量
        "client_secret": "4fa8890c-5888-48ba-93b8-3e7db1165b65",            # 客户端密钥
        "timeout": 10,                                                      # 请求超时时间(秒)
    }
    
    # 执行模式配置
    fetch_mode = "once"      # 单次获取模式（便于调试）
    # fetch_mode = "interval"  # 定时获取模式
    
    if fetch_mode == "once":
        # 只获取一次
        logger.info("开始单次消息获取（与API示例对比）")
        messages = fetch_messages(config)
        # 打印结果（格式化JSON）
        print(json.dumps(messages, indent=2))
        # 记录结果到日志
        if "error" in messages:
            logger.error(f"获取消息结果: {messages['error']}")
            exit(1)
        else:
            logger.info(f"成功获取 {len(messages['messages'])} 条消息")
            exit(0)
    else:
        # 定时获取
        fetch_interval = 60     # 获取间隔(秒)
        logger.info(f"开始定时获取消息，间隔: {fetch_interval}秒")
        try:
            while True:
                messages = fetch_messages(config)
                # 仅在成功时打印消息内容
                if "error" not in messages:
                    logger.info(f"获取到 {len(messages['messages'])} 条消息")
                time.sleep(fetch_interval)
        except KeyboardInterrupt:
            logger.info("接收到中断信号，程序退出")


if __name__ == "__main__":
    main()