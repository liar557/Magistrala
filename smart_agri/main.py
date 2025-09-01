import logging
from agents.orchestrator import OrchestratorAgent, register_agents

logging.basicConfig(level=logging.INFO)

if __name__ == "__main__":
    agents = register_agents()
    orchestrator = OrchestratorAgent(agents)
    user_input = "请根据当前土壤湿度自动灌溉"
    result = orchestrator.run(user_input)
    print(result)