from .utils import Bet


class Message:
    MAX_LENGTH = 80  # bytes

    def __init__(self, buffer):
        self._bet_confirmation = "a"
        self._bet_message_type = "b#"
        self.is_bet = False
        self._separator = "#"
        self._filler = "@"
        self._message = ""
        self._transform(buffer)

    def _transform(self, buffer):
        self._message = buffer.rstrip().decode('utf-8')
        self._message = self._message.replace(self._filler, "")
        self.is_bet = self._message[0:2] == self._bet_message_type

    def deserialize_bet(self) -> Bet:
        info = self._message.split(self._separator)
        return Bet(info[1], info[2], info[3], info[4], info[5], info[6])

    def serialize_confirmation(self):
        if self.is_bet:
            return self._bet_confirmation.encode('utf-8')
        return self._message.encode('utf-8')
