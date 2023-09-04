import logging
import socket

from .message import Message
from .utils import store_bets


class Connection:
    def __init__(self, port, listen_backlog):
        self._port = port
        self._max_connections = listen_backlog
        self._connection = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

    def start(self):
        self._connection.bind(('', self._port))
        self._connection.listen(self._max_connections)

    def end(self):
        self._connection.close()

    def accept(self):
        logging.info("action: accept_connections | result: in_progress")
        conn, addr = self._connection.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        logging.info(f'action: handler_connection | result: in_progress | ip: {addr[0]}')
        self._handle_connection(conn, addr)
        logging.info(f'action: handler_connection | result: success | ip: {addr[0]}')

    def _handle_connection(self, conn, addr):
        try:
            logging.info(f'action: receive_message | result: in_progress | ip: {addr[0]}')
            buffer = conn.recv(Message.MAX_LENGTH)
            logging.info(f'action: receive_message | result: success | ip: {addr[0]}')
            message = Message(buffer)
            if message.is_bet:
                bet = message.deserialize_bet()
                store_bets([bet])
                logging.info(f'action: stored_bet | result: success | dni: {bet.document} | number: {bet.number}')
            confirmation = message.serialize_confirmation()
            conn.send(confirmation)
        except OSError as e:
            logging.error(f'action: receive_message | result: fail | error: {e.strerror}')
        finally:
            conn.close()
