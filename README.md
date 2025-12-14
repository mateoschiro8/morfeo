# morfeo
Un generador de Honeytokens. Su deploy se encuentra [acá](https://morfeo-c8s3.onrender.com/). Es necesario entrar y despertar al server antes de empezar a trabajar.

Para más información del sistema, así como un manual de uso de la herramienta, puede consultarse el [informe](https://github.com/mateoschiro8/morfeo/blob/main/informe/Ciberseguros-Tokensnare.pdf).


## Desarrollo local

Correr ```cp env-sample .env```, y luego ```go run main.go server --msg a --chat a``` para levantar el server y ```go run main.go {cmd}``` para la cli.
