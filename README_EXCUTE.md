# TP0: Documentacion
## Ejercicio 1

Se agrega el parametro `CLIENTS` al comando `make docker-compose-up`

Ejemplo de uso para generar 2 clientes:
```
make docker-compose-up CLIENTS=2
```
Si no se agrega el parametro o no es un numero se utiliza el valor por defecto (`CLIENTS=1`)

## Ejercicio 2

Reemplazamos el `COPY` por `volumes`

## Ejercicio 3

Agregamos un nuevo servicio `netcat` que ejecuta el comando `nc -v` para
mostrar el mensaje **Hello 7574**.

Forma de ejecutarlo
- Terminal 1
    - `make docker-compose-up`
    - esperamos un momento, revisando la **Terminal 2**
- Terminal 2
    - `make docker-compose-logs`
    - se observa `netcat   | Hello 7574`

## Ejercicio 4

Agregamos el modulo `signal` en el Cliente y Servidor, donde se define los manejadores
para la senal `SIGTERM`, que en resumen cierran la conexion.

Forma de ejecutarlo 
- Terminal 1
  - `make docker-compose-up`
  - esperamos un momento, revisando la **Terminal 2**
  - `make docker-compose-down`
- Terminal 2
  - `make docker-compose-logs`
  - se observan `exited with code 0` en el cliente y servidor

