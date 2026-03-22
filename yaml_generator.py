import sys
def client(id):
    return f"""\
  client{id}:
    container_name: client{id}
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID={id}
      - CLI_FILE=/data/agency-{id}.csv
    volumes:
      - ./client/config.yaml:/config.yaml
      - ./.data/agency-{id}.csv:/data/agency-{id}.csv
    networks:
      - testing_net
    depends_on:
      - server
"""

def yaml(n_clients):
    server = """\
name: tp0
services:
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    environment:
      - PYTHONUNBUFFERED=1
    volumes:
      - ./server/config.ini:/config.ini
    networks:
      - testing_net
"""
    clients = "\n".join(client(i) for i in range(1, n_clients + 1))
    networks = """\
networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24
"""
    return server + "\n" + clients + "\n" + networks

def main():
    if len(sys.argv) == 3:
        file = sys.argv[1]
        n_clients = int(sys.argv[2])
        if n_clients >= 0:
            yaml_file = yaml(n_clients)
            with open(file, "w") as f:
                f.write(yaml_file)

if __name__ == "__main__":
    main()