import logging
from .base import AgentBase
from langchain_ollama import ChatOllama
import re
import json

class AnalysisAgent(AgentBase):
    def __init__(self, model_name="qwen3:8b"):
        # 初始化 Ollama LLM
        self.llm = ChatOllama(model=model_name)

    def run(self, data):
        try:
            # 构造提示词
            prompt = (
                f"当前土壤湿度为 {data.get('soil_moisture', '未知')}%。"
                "请判断是否需要灌溉，并以如下JSON格式输出建议："
                '{"action": "irrigate", "duration": 分钟数, "area": "区域"}。'
                "如果不需要灌溉，action为'none'。"
            )
            # 调用 Ollama LLM 进行分析
            advice = self.llm.invoke(prompt)
            logging.info(f"Analysis result: {advice}")
            # 尝试解析为字典
            try:
                result = json.loads(advice)
            except Exception:
                # 可选：用正则或其他方式提取
                result = {"action": "none"}
            return result
        except Exception as e:
            logging.error(f"AnalysisAgent error: {e}")
            return {"action": "none", "error": str(e)}