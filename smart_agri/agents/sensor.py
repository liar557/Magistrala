import logging
from .base import AgentBase

class SensorAgent(AgentBase):
    def run(self):
        try:
            data = {"temperature": 25, "humidity": 60, "soil_moisture": 30}
            logging.info(f"Sensor data: {data}")
            return data
        except Exception as e:
            logging.error(f"SensorAgent error: {e}")
            return {"error": str(e)}