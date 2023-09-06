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

netcat = """
  netcat:
    container_name: netcat
    build:
      context: .
      dockerfile: netcat/Dockerfile
    entrypoint: sh netcat_echo.sh
    depends_on:
      - server
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
      - CLI_FIRSTNAME=name""" + str(i) + """
      - CLI_LASTNAME=l_name""" + str(i) + """
      - CLI_DOCUMENT=1234""" + str(i) + """
      - CLI_BIRTHDATE=1999-03-1""" + str(i) + """
      - CLI_NUMBER=757""" + str(i) + """
      - CLI_FILEPATH=./agency-""" + str(i) + """.csv
    volumes:
      - ./client/config.yaml:/config.yaml
      - ./.data/dataset/agency-""" + str(i) + """.csv:/agency-""" + str(i) + """.csv
    networks:
      - testing_net
    depends_on:
      - server"""

file = open("docker-compose-dev.yaml", "w")
file.write(services + server_service + netcat + clients_service + "\n" + network + "\n")
file.close()
