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
        batch_size = 0
        try:
            buffer = conn.recv(Message.CONFIG_MAX_LENGTH)
            message = Message(buffer)
            confirmation = message.serialize_confirmation()
            conn.send(confirmation)
            if message.is_config:
                batch_size = message.deserialize_config()
            else:
                conn.close()
                return
            logging.info(f'action: receive_message | result: in_progress | ip: {addr[0]}')
            buffer = conn.recv(batch_size * 8000)
            message = Message(buffer)
            if message.is_bet:
                bets = message.deserialize_bets_batch()
                store_bets(bets)
                confirmation = message.serialize_confirmation()
                conn.send(confirmation)
                while not message.is_last:
                    buffer = conn.recv(batch_size * 8000)
                    message = Message(buffer)
                    bets = message.deserialize_bets_batch()
                    store_bets(bets)
                    confirmation = message.serialize_confirmation()
                    conn.send(confirmation)
                logging.info(f'action: stored_bet | result: success')
            confirmation = message.serialize_confirmation()
            conn.send(confirmation)
        except OSError as e:
            logging.error(f'action: receive_message | result: fail | error: {e.strerror}')
        finally:
            conn.close()
