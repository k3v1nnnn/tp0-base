# TP0: Documentacion
## Ejercicio 1

ultimo commit -> https://github.com/k3v1nnnn/tp0-base/commit/77a059d37b04d9dda16f0f44cf8c5b753c0faba3

Se agrega el parametro `CLIENTS` al comando `make docker-compose-up`

Ejemplo de uso para generar 2 clientes:
```
make docker-compose-up CLIENTS=2
```
Si no se agrega el parametro o no es un numero se utiliza el valor por defecto (`CLIENTS=1`)

## Ejercicio 2

ultimo commit -> https://github.com/k3v1nnnn/tp0-base/commit/4e8ec715f5064d34dea1a380ebdc92a2ee822033

Reemplazamos el `COPY` por `volumes`

## Ejercicio 3

ultimo commit -> https://github.com/k3v1nnnn/tp0-base/commit/990c09673477268fa9779f3598b3a27c46a240ae

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

ultimo commit -> https://github.com/k3v1nnnn/tp0-base/commit/dc30661bfddf9d9e49f00ba20e508729272f4bd5

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

ultimo commit -> https://github.com/k3v1nnnn/tp0-base/commit/d10a1a471489d7ab393a07b7f2fef1c34ce3510f

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

ultimo commit -> https://github.com/k3v1nnnn/tp0-base/commit/76182635cc6b444c73236ad5b1e3c39334940450

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


## Ejercicio 7

ultimo commit -> https://github.com/k3v1nnnn/tp0-base/commit/5cf520fc30487a4f896d899259c1177a450cef86

Agremos clase Loteria para manejar/contar los ganadores

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
  - **br#** request, nos sirve para saber que el cliente pide los resultados de las apuesta
- El servidor response con **a** (accept) cuando se almacena una apuesta o saber que recibio el mensaje. En otros
  casos responde con el mismo mensaje que le llego
- El servidor responde a las consultas de los resultados de las apuestas con **\*#**
  - **r#** retry, nos sirve para que el cliente espere un rato y vuelva a consultar despues
  - **w#** winner, nos sirve para que el cliente obtenga la cantidad de ganadores
### Ejecucion

- Asegurarse que se tiene descomprimido el archivo dataset.zip en **.data/dataset/archivos.csv**
- Terminal 1
  - `make docker-compose-up CLIENTS=5 `
  - esperamos un momento, revisando la **Terminal 2**
- Terminal 2
  - `make docker-compose-logs`
  - se observan los mensajes que intercambian todos los servicios

## Ejercicio 8

Agregamos hilos, cada hilo maneja la conexion con el cliente donde se ejecuta la subida de apuestas o la peticion de 
ganadores.

### Sincronizacion
- Agregamos un lock para la seccion critica, que en nuestro caso es la funcion que persiste las apuestas.

### Ejecucion

- Asegurarse que se tiene descomprimido el archivo dataset.zip en **.data/dataset/archivos.csv**
- Terminal 1
  - `make docker-compose-up CLIENTS=5 `
  - esperamos un momento, revisando la **Terminal 2**
- Terminal 2
  - `make docker-compose-logs`
  - se observan los mensajes que intercambian todos los servicios


