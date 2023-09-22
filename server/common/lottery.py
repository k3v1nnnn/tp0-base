from .utils import load_bets
from .utils import has_won


class Lottery:
    def __init__(self):
        self._agencies = [1, 2, 3, 4, 5]
        self._finalized_agencies = []
        self._winners = {}

    def add_finalized_agency(self, agency):
        self._finalized_agencies.append(agency)

    def can_start(self):
        for agency in self._agencies:
            if agency not in self._finalized_agencies:
                return False
        return True

    def get_winners(self, agency):
        return self._winners.get(agency, [])

    def has_winners(self):
        return len(self._winners) != 0

    def start(self):
        bets = load_bets()
        for bet in bets:
            if has_won(bet):
                documents = self._winners.get(bet.agency, [])
                if documents:
                    self._winners[bet.agency].append(bet.document)
                else:
                    self._winners[bet.agency] = [bet.document]
