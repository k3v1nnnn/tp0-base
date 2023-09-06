import logging

from .utils import Bet


class Message:
    MAX_LENGTH = 70  # bytes
    CONFIG_MAX_LENGTH = 6

    def __init__(self, buffer):
        self._bet_separator = ","
        self._bet_confirmation = "a"
        self._bet_message_type = ["bc#", "bf#", "bi#"]
        self.is_bet = False
        self.is_last = False
        self.is_config = False
        self._separator = "#"
        self._filler = "@"
        self._message = ""
        self._transform(buffer)

    def _transform(self, buffer):
        message = buffer.rstrip().decode('utf-8')
        message = message.replace(self._filler, "")
        self.is_bet = message[0:3] in self._bet_message_type
        self.is_last = message[0:3] == "bf#"
        self.is_config = message[0:3] == "bi#"
        self._message = message

    def _deserialize_bet(self, _bet) -> Bet:
        info = _bet.split(self._separator)
        return Bet(info[0], info[1], info[2], info[3], info[4], info[5])

    def deserialize_config(self) -> int:
        info = self._message.split(self._separator)
        return int(info[1])


    def deserialize_bet(self) -> Bet:
        info = self._message.split(self._separator)
        return Bet(info[1], info[2], info[3], info[4], info[5], info[6])

    def deserialize_bets_batch(self) -> list[Bet]:
        self._message = self._message[3:]
        bets = self._message.split(self._bet_separator)
        return [*map(self._deserialize_bet, bets)]

    def serialize_confirmation(self):
        if self.is_bet:
            return self._bet_confirmation.encode('utf-8')
        return self._message.encode('utf-8')
