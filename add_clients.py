import sys

clients = 1
if len(sys.argv) > 1:
    try:
        clients = int(sys.argv[1])
        clients = clients if clients else 1
    except ValueError:
        clients = 1


services = """version: '3.9'
name: tp0
services:"""

server_service = """
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    environment:
      - PYTHONUNBUFFERED=1
      - LOGGING_LEVEL=DEBUG
    volumes:
      - ./server/config.ini:/config.ini
    networks:
      - testing_net"""

network = """
networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24"""

clients_service = ""

for i in range(1, clients + 1):
    clients_service = clients_service + "\n" + """
  client""" + str(i) + """:
    container_name: client""" + str(i) + """
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID=""" + str(i) + """
      - CLI_LOG_LEVEL=DEBUG
    volumes:
      - ./client/config.yaml:/config.yaml
    networks:
      - testing_net
    depends_on:
      - server"""

file = open("docker-compose-dev.yaml", "w")
file.write(services + server_service + clients_service + "\n" + network + "\n")
file.close()
