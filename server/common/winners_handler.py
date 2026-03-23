import logging
from common.handler import Handler
from common.protocol import send_response
from common.serializer import unserialize_flag, serialize_winners


class WinnersHandler(Handler):
    def run(self):
        flag_data = unserialize_flag(self.message)
        agency_id = flag_data['agency_id']
        if not self.lottery.is_finish():
            send_response(self.client_sock, "WAIT")
            logging.info(f"action: send_wait | result: success | client_id: {agency_id}")
            return
        winners = self.lottery.winners(agency_id)
        send_response(self.client_sock, serialize_winners(winners))
        logging.info(f"action: consulta_ganadores | result: success | client_id: {agency_id} | cant_ganadores: {len(winners)}")
