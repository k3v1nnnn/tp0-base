from abc import ABC, abstractmethod


class Handler(ABC):
    def __init__(self, client_sock, lottery, message):
        self.client_sock = client_sock
        self.lottery = lottery
        self.message = message

    @abstractmethod
    def run(self):
        pass
