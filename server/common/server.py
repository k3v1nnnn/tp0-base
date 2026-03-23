import socket
import logging
import signal
from common.protocol import recv
from common.lottery import Lottery
from common.handler_factory import HandlerFactory


class Server:
    def __init__(self, port, listen_backlog, agencies):
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._server_running = True
        self._lottery = Lottery(agencies)
        signal.signal(signal.SIGTERM, self.__handle_sigterm)

    def __safe_server_socket_close(self):
        try:
            self._server_socket.close()
        except OSError:
            pass

    def __handle_sigterm(self, *_):
        logging.info("action: receive_sigterm | result: success")
        self._server_running = False
        self.__safe_server_socket_close()

    def run(self):
        while self._server_running:
            try:
                client_sock = self.__accept_new_connection()
                self.__handle_client_connection(client_sock)
            except OSError:
                break
        logging.info("action: server_shutdown | result: success")

    def __handle_client_connection(self, client_sock):
        try:
            message = recv(client_sock)
            handler = HandlerFactory.get_handler(client_sock, self._lottery, message)
            handler.run()
        except OSError as e:
            logging.error(f"action: handle_client_connection | result: fail | error: {e}")
        finally:
            client_sock.close()
            logging.info("action: close_connection | result: success")

    def __accept_new_connection(self):
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c
