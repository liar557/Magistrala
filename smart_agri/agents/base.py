from abc import ABC, abstractmethod

class AgentBase(ABC):
    @abstractmethod
    def run(self, *args, **kwargs):
        pass