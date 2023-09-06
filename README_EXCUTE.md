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

## Ejercicio 5

Separamos un poco las responsabilidades, Ahora Tenemos clases como Conexion , Mensaje y Apuesta

### Procolo
- Para la comunicacion usamos TCP
- El tamano del mensaje es fijo de unos 80 bytes
- Se completa con el caracter **@** en el caso de faltar caracteres
- La serializacion de la apuesta es la union de todos sus datos mediante el caracter
**\#** como ejemplo tenemos **{agencia}#{nombre}#{apellido}#{fecha de nacimeinto}#{documento}#{numero}**
- El mensaje que finalmente envia el cliente es la union de **b#** (bet) y la serializacion
- El servidor utiliza **b#** para reconocer a la apuesta
- El servidor response con **a** cuando se almacena una apuesta en otros
casos responde con el mismo mensaje que le llego
### Ejecucion

- Terminal 1
  - `make docker-compose-up CLIENTS=5 `
  - esperamos un momento, revisando la **Terminal 2**
- Terminal 2
  - `make docker-compose-logs`
  - se observan los mensajes que intercambian todos los servicios

## Ejercicio 6

Agremos clase para Archivo para poder leer y parsear el archivo

### Procolo
- Para la comunicacion usamos TCP
- El tamano de cada apuesta es fijo de unos 70 bytes
- El tamano del mensaje es fijo pero depende del tamano del batch, lo calculamos como
  (tamano del batch * 70)
- Se completa con el caracter **@** en el caso de faltar caracteres
- La serializacion de la apuesta es la union de todos sus datos mediante el caracter
  **\#** como ejemplo tenemos **{agencia}#{nombre}#{apellido}#{fecha de nacimeinto}#{documento}#{numero}**
- El mensaje que finalmente envia el cliente es la union de **b\*#** (bet + tipo de evento) y la serializacion
- El servidor utiliza **b\*#** para reconocer a las apuestas y saber cuando finalizo
  - **bi#** initial , nos sirve para saber el tamano de batch
  - **bc#** continue , nos sirve para saber que seguimos esperamo apuestas
  - **bf#** final , nos sirve para saber que son las ultimas apuestas
- El servidor response con **a** (accept) cuando se almacena una apuesta o saber que recibio el mensaje en otros
  casos responde con el mismo mensaje que le llego
### Ejecucion

- Terminal 1
  - `make docker-compose-up CLIENTS=5 `
  - esperamos un momento, revisando la **Terminal 2**
- Terminal 2
  - `make docker-compose-logs`
  - se observan los mensajes que intercambian todos los servicios



