import signal
import logging
from .connection import Connection
from .lottery import Lottery


class Server:
    def __init__(self, port, listen_backlog):
        self._lottery = Lottery()
        self._connection = Connection(port, listen_backlog, self._lottery)
        self._is_running = True
        signal.signal(signal.SIGTERM, self.__exit)

    def run(self):
        self._connection.start()
        while self._is_running:
            try:
                self._connection.accept()
                if self._lottery.can_start() and not self._lottery.has_winners():
                    self._lottery.start()
            except OSError as e:
                logging.error("action: new_client_connection | result: fail | error: " + e.strerror)
        self._connection.end()

    def __exit(self, signal_num, frame):
        logging.info('action: graceful_exit | result: in_progress')
        self._is_running = False
        self._connection.end()
        logging.info('action: graceful_exit | result: success')
