import logging
import socket

from .message import Message
from .utils import store_bets


class Connection:
    def __init__(self, port, listen_backlog, lottery):
        self._port = port
        self._max_connections = listen_backlog
        self._lottery = lottery
        self._connection = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

    def start(self):
        self._connection.bind(('', self._port))
        self._connection.listen(self._max_connections)

    def end(self):
        self._connection.close()

    def accept(self):
        logging.info("action: accept_connections | result: in_progress")
        conn, addr = self._connection.accept()
        self._handle_connection(conn, addr)
        logging.info(f'action: handler_connection | result: success | ip: {addr[0]}')

    def _handle_connection(self, conn, addr):
        try:
            buffer = conn.recv(Message.CONFIG_MAX_LENGTH)
            message = Message(buffer)
            if message.is_config:
                confirmation = message.serialize_confirmation()
                conn.send(confirmation)
                batch_size, client_id = message.deserialize_config()
            elif message.is_request:
                response = message.serialize_empty_response()
                if self._lottery.has_winners():
                    agency = message.deserialize_request()
                    amount_winners = self._lottery.get_winners(agency)
                    response = message.serialize_winners_response(amount_winners)
                conn.send(response)
                conn.close()
                return
            else:
                confirmation = message.serialize_confirmation()
                conn.send(confirmation)
                conn.close()
                return
            logging.info(f'action: receive_message | result: in_progress | ip: {addr[0]}')
            buffer = conn.recv(batch_size * 8000)
            message = Message(buffer)
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
            self._lottery.add_finalized_agency(client_id)
        except OSError as e:
            logging.error(f'action: receive_message | result: fail | error: {e.strerror}')
        finally:
            conn.close()
