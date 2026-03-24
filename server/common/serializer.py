from common.utils import Bet

SEPARATOR = "|"
FIELDS = ["first_name", "last_name", "document", "birthdate", "number", "agency_id"]
FLAGS_FIELDS = ["flag", "agency_id"]


def unserialize_bet(payload):
    values = payload.split(SEPARATOR)
    data = dict(zip(FIELDS, values))
    return Bet(
        agency=data["agency_id"],
        first_name=data["first_name"],
        last_name=data["last_name"],
        document=data["document"],
        birthdate=data["birthdate"],
        number=data["number"],
    )
    
def unserialize_flag(payload):
    values = payload.split(SEPARATOR)
    return dict(zip(FLAGS_FIELDS, values))

def serialize_winners(winners):
    return SEPARATOR.join(winners)
