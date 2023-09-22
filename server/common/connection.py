import logging
import socket
import threading

from .message import Message
from .utils import store_bets


class Connection:
    LOCK = threading.Lock()

    def __init__(self, port, listen_backlog, lottery):
        self._port = port
        self._max_connections = listen_backlog
        self._lottery = lottery
        self._connection = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._threads = []

    def start(self):
        self._connection.bind(('', self._port))
        self._connection.listen(self._max_connections)

    def end(self):
        for thread in self._threads:
            thread.join()
        self._connection.close()

    def accept(self):
        logging.info("action: accept_connections | result: in_progress")
        conn, addr = self._connection.accept()
        thread = threading.Thread(target=self._handle_connection, args=(conn, addr))
        thread.start()
        self._threads.append(thread)
        logging.info(f'action: handler_connection | result: success | ip: {addr[0]}')

    def _send_bets(self, conn, batch_size):
        buffer = conn.recv(batch_size * Message.MAX_LENGTH)
        message = Message(buffer)
        bets = message.deserialize_bets_batch()
        Connection.LOCK.acquire()
        store_bets(bets)
        Connection.LOCK.release()
        confirmation = message.serialize_confirmation()
        conn.send(confirmation)
        return message

    def _handle_config(self, message, conn, addr):
        confirmation = message.serialize_confirmation()
        conn.send(confirmation)
        batch_size, client_id = message.deserialize_config()
        logging.info(f'action: receive_message | result: in_progress | ip: {addr[0]}')
        message = self._send_bets(conn, batch_size)
        while not message.is_last:
            message = self._send_bets(conn, batch_size)
        logging.info(f'action: stored_bet | result: success | client_id: {client_id}')
        Connection.LOCK.acquire()
        self._lottery.add_finalized_agency(client_id)
        Connection.LOCK.release()

    def _handle_request(self, message, conn):
        response = message.serialize_empty_response()
        if self._lottery.has_winners():
            agency = message.deserialize_request()
            document_winners = self._lottery.get_winners(agency)
            response = message.serialize_winners_response(document_winners)
        conn.send(response)
        conn.close()

    def _handle_confirm(self, message, conn):
        confirmation = message.serialize_confirmation()
        conn.send(confirmation)
        conn.close()

    def _handle_connection(self, conn, addr):
        try:
            buffer = conn.recv(Message.CONFIG_MAX_LENGTH)
            message = Message(buffer)
            if message.is_config:
                self._handle_config(message, conn, addr)
            elif message.is_request:
                self._handle_request(message, conn)
            else:
                self._handle_confirm(message, conn)
        except OSError as e:
            logging.error(f'action: receive_message | result: fail | error: {e.strerror}')
        finally:
            conn.close()
