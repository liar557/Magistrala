import logging
from .base import AgentBase

class WeatherAgent(AgentBase):
    def run(self, location="default"):
        try:
            # 这里可以调用第三方天气API，例如和风天气、OpenWeatherMap等
            # 示例伪代码
            weather_data = {
                "temperature": 28,
                "humidity": 70,
                "rainfall": 0,
                "condition": "晴"
            }
            logging.info(f"Weather data: {weather_data}")
            return weather_data
        except Exception as e:
            logging.error(f"WeatherAgent error: {e}")
            return {"error": str(e)}