import logging
from .base import AgentBase
from .permission import permission_manager

class IrrigationAgent(AgentBase):
    def run(self, command):
        try:
            if not permission_manager.check("IrrigationAgent"):
                return "无权限执行灌溉操作"
            if command.get("action") == "irrigate":
                duration = command.get("duration", 0)
                area = command.get("area", "全部区域")
                # 这里调用实际设备控制接口
                result = f"已启动{area}灌溉，时长{duration}分钟"
            else:
                result = "无需灌溉"
            logging.info(f"Irrigation result: {result}")
            return result
        except Exception as e:
            logging.error(f"IrrigationAgent error: {e}")
            return {"error": str(e)}