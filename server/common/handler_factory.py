import re
from common.bets_handler import BetsHandler
from common.winners_handler import WinnersHandler


class HandlerFactory:
    @staticmethod
    def get_handler(client_sock, lottery, message):
        if re.match(r"WINNERS\|\d+", message):
            return WinnersHandler(client_sock, lottery, message)
        else:
            return BetsHandler(client_sock, lottery, message)
