from .sensor import SensorAgent
from .analysis import AnalysisAgent
from .irrigation import IrrigationAgent
from .permission import permission_manager

class OrchestratorAgent:
    def __init__(self, agents):
        self.agents = agents

    def run(self, user_input):
        sensor_data = self.agents["SensorAgent"].run()
        advice = self.agents["AnalysisAgent"].run(sensor_data)
        result = self.agents["IrrigationAgent"].run(advice)
        return result

def register_agents():
    permission_manager.register("SensorAgent", True)
    permission_manager.register("AnalysisAgent", True)
    permission_manager.register("IrrigationAgent", True)
    agents = {
        "SensorAgent": SensorAgent(),
        "AnalysisAgent": AnalysisAgent(),
        "IrrigationAgent": IrrigationAgent()
    }
    return agents