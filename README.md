# morfeo
Un generador de Honeytokens. IT IS IN PROCESS

## Desarrollo local

Correr ```cp env-sample .env```, completar las variables y luego ```docker compose up``` para levantar el server y ```go run main.go {cmd}``` para la cli.


LINKS:

https://github.com/spf13/cobra/blob/v1.10.1/site/content/user_guide.md
https://pkg.go.dev/github.com/spf13/cobra#section-readme

# IDEA

Cuando se quiere hacer un nuevo honeytoken, se le hace un POST a /tokens, y en el cuerpo del pedido se ponen los datos necesarios para el funcionamiento del token.
Obligatorio el string msg de identificación. Luego, los más específicos de cada uno. Los campos que no apliquen al caso del token IGUAL DEBEN METERSE CON COSAS RANDOM.
Del lado del server se crea un struct general que guarda toda esta info. TODOS LOS DATOS DE UserInput DEBEN SER MANDADOS DESDE EL FRONT. Sino no anda el desjsonlisador.
Luego de crear el token, el server responde con un identificador tipo o2wrui3ehn (o lo que sea).
El front luego se encarga de agarrar ese identificador y meterlo en una url específica de ese tipo de token. Por ej, si es un qr, el front devuelve la url "/qrs/o2wrui3ehn". 
El server entonces, al recibir un GET en "/qrs/{id}", usa el id para traer los datos del token, y como sabe que está en qr, sabé qué datos tiene disponibles para mostrar, 
a donde redirigir, etc etc


# TODOs
    - Que los tokens guarden el CHATID del usuario a quién avisar en en telegram
    - Rehacer los tokens con la nueva base, y pasar a que usen el GetToken como el de qrs/bins. Tokens restantes: pdf, img, css. Proximos: 
    - Acomodar las variables globales del cmd (se pueden reutilizar)


# CAMBIOS

Instrucciones de Uso

Crear Directorios, para tener los permisos correctos:

    mkdir -p tokens/ input/ tmp/

Ejecutar:
    cp env-sample .env

Cambiar variables en .env UID y GID por las tuyas, puedes obtenerlas con:
    Para UID:
        id -u

    Para GID:
        id -g
El siguiente comando hace los anteriores pasos en un solo comando:
    sed -e "s/UID=\"\"/UID=\"$(id -u)\"/" -e "s/GID=\"\"/GID=\"$(id -g)\"/" env-sample > .env

Para levantar el server:
    docker compose up morfeo-server

Para levantar la cli:
    docker compose run -rm morfeo-cli


Para correrlo de forma Local hay que:

    1- en .env cambiar la url a "http://localhost:8000"

    2- en compose.yaml en el servicio de morfeo-cli hay que hardocdear "http://morfeo-server:8000"

    3- en el honey token que te devuelve la cli hay que volver a poner localhost en donde dice morfeo-server

Si se hacen cambios en el codigo se debe correr:
    docker compose build