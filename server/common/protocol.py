import struct


def _send_all(sock, data):
    total_sent = 0
    while total_sent < len(data):
        sent = sock.send(data[total_sent:])
        if sent == 0:
            raise ConnectionError("Socket connection broken")
        total_sent += sent


def _recv_all(sock, n):
    data = b""
    while len(data) < n:
        chunk = sock.recv(n - len(data))
        if not chunk:
            raise ConnectionError("Connection closed by peer")
        data += chunk
    return data


def recv(sock):
    header = _recv_all(sock, 4)
    length = struct.unpack("!I", header)[0]
    return _recv_all(sock, length).decode("utf-8")


def send_response(sock, msg):
    """Send a length-prefixed response to the client."""
    data = msg.encode("utf-8")
    header = struct.pack("!I", len(data))
    _send_all(sock, header + data)