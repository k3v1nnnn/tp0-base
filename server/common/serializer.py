SEPARATOR = "|"
FIELDS = ["first_name", "last_name", "document", "birthdate", "number", "agency_id"]


def unserialize_bet(payload):
    values = payload.split(SEPARATOR)
    return dict(zip(FIELDS, values))
