import logging
from common.handler import Handler
from common.protocol import recv, send_response
from common.serializer import unserialize_bet
from common.batcher import split_batch


class BetsHandler(Handler):
    def run(self):
        message = self.message
        agency_id = None
        while message != 'END':
            success = True
            bets = list(map(unserialize_bet, split_batch(message)))
            if agency_id is None and bets:
                agency_id = bets[0].agency
            try:
                self.lottery.add_bets(bets)
                logging.info(f"action: apuesta_recibida | result: success | cantidad: {len(bets)}")
            except Exception as e:
                logging.error(f"action: apuesta_recibida | result: fail | cantidad: {len(bets)} | error: {e}")
                success = False
            send_response(self.client_sock, "OK" if success else "ERROR")
            logging.info(f"action: send_batch_response | result: success | client_id: {agency_id}")
            message = recv(self.client_sock)
        send_response(self.client_sock, "OK")
        logging.info(f"action: send_end_response | result: success | client_id: {agency_id}")
        self.lottery.add_agency(agency_id)
