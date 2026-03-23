import logging
from common.lottery_status import LotteryStatus
from common.utils import store_bets, load_bets, has_won


class Lottery:
    def __init__(self, total_participants):
        self.total_participants = total_participants
        self.agencies = {}
        self.status = LotteryStatus.WAITING

    def add_agency(self, agency):
        self.agencies[int(agency)] = []
        if len(self.agencies) == self.total_participants:
            self._start()
            logging.info("action: sorteo | result: success")

    def _start(self):
        for bet in load_bets():
            if has_won(bet):
                self.agencies[bet.agency].append(bet.document)
        self.status = LotteryStatus.FINISHED

    def winners(self, agency):
        return self.agencies.get(int(agency), [])

    def add_bets(self, bets):
        store_bets(bets)

    def is_finish(self):
        return self.status == LotteryStatus.FINISHED
