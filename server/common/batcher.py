BATCH_SEPARATOR = "#"


def split_batch(payload):
    return payload.split(BATCH_SEPARATOR)
