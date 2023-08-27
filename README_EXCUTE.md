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