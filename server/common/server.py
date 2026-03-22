import socket
import logging
import signal
from common.protocol import recv, send_response
from common.serializer import unserialize_bet
from common.batcher import split_batch
from common.utils import Bet, store_bets


class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._server_running = True
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
            while message != 'END':
                agency_id = None
                success = True
                bets_data = list(map(unserialize_bet, split_batch(message)))
                bets = []
                for bet_data in bets_data:
                    if agency_id is None:
                        agency_id = bet_data["agency_id"]
                    bets.append(
                        Bet(
                            agency=bet_data["agency_id"],
                            first_name=bet_data["first_name"],
                            last_name=bet_data["last_name"],
                            document=bet_data["document"],
                            birthdate=bet_data["birthdate"],
                            number=bet_data["number"],
                        )
                    )
                try:
                    store_bets(bets)
                    logging.info(f"action: apuesta_recibida | result: success | cantidad: {len(bets)}")
                except Exception as e:
                    logging.error(f"action: apuesta_recibida | result: fail | cantidad: {len(bets)}")
                    success = False
                send_response(client_sock, "OK" if success else "ERROR")
                logging.info(f"action: send_response | result: success | client_id: {agency_id}")
                message = recv(client_sock)
            send_response(client_sock, "OK")
            logging.info(f"action: send_response | result: success | client_id: {agency_id}")
        except OSError as e:
            logging.error(f"action: receive_bet | result: fail | error: {e}")
        finally:
            client_sock.close()

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c
