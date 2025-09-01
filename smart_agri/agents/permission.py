class PermissionManager:
    def __init__(self):
        self.permissions = {}

    def register(self, agent_name, enabled=True):
        self.permissions[agent_name] = enabled

    def check(self, agent_name):
        return self.permissions.get(agent_name, False)

permission_manager = PermissionManager()