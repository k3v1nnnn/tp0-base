SEPARATOR = "|"
FIELDS = ["first_name", "last_name", "document", "birthdate", "number", "agency_id"]
FLAGS_FIELDS = ["flag", "agency_id"]


def unserialize_bet(payload):
    values = payload.split(SEPARATOR)
    return dict(zip(FIELDS, values))
    
def unserialize_flag(payload):
    values = payload.split(SEPARATOR)
    return dict(zip(FLAGS_FIELDS, values))

def serialize_winners(winners):
    return SEPARATOR.join(winners)
