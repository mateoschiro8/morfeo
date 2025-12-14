#pagebreak()

#let hl = text.with(fill: orange)

#let codeblock = block.with(
  fill: rgb("#4d360621"),
  width: 100%,
  inset: 8pt,
  radius: 4pt,
  spacing: 1%
)

= Manual de uso

Hay dos maneras de utilizar la herramienta. La primer manera es compilando el binario manualmente a partir del código fuente. La segunda manera es mediante un contenedor de docker que la compila, y utilizar los comandos dentro del contenedor.   

En ambos casos, luego de clonado el repositorio debe ejecutarse el siguiente comando:

#codeblock[`$ cp env-samples .env`]

para tener disponibles las variables de ambiente. Además, para poder recibir las alertas es necesario activar el bot que las manda. Para hacer esto simplemente se debe entrar al #link("https://t.me/morfeoSeguroBot")[#hl[#underline[chat]]] con el mismo, y darle #hl[`start`]. 

== Compilando el código fuente

Teniendo #link("https://go.dev/")[#hl[Golang]] instalado, y estando en la rama #hl[main], se ejecuta en la raiz del repositorio:
#codeblock[`$ go build`]

Esto generará un ejecutable llamado #hl[`morfeo`]. Luego, al ejecutar
#codeblock[`$ ./morfeo`]
desde donde se encuentre, se podrá ver la salida:
#codeblock[#set text(size: 12pt)
```  
Usage:
  morfeo [command]

Available Commands:
  bin         Genera un honeytoken a partir de un binario
  css         Genera el honeytoken de css para paginas clonadas
  help        Help about any command
  image       Genera un honeytoken de imagen
  pdf         Genera el honeytoken de pdf
  qr          Genera el honeytoken de qr      

Flags:
      --chat string   Chat ID al cual enviar la alerta al ser activado
  -h, --help          help for morfeo
      --msg string    Identificador del token

Use "morfeo [command] --help" for more information about a command
```
]



Luego, se ejecuta #hl[`./morfeo [command] [flags]`]. En caso de necesitar más información sobre un comando, ejecutarlo sin flags o con la flag #hl[`--help`] imprime más detalles sobre el mismo. 

Algunos comandos tienen flags que otros no tienen, pero hay dos flags obligatorias que todos los comandos comparten:
+ #hl[`--msg`]: Un mensaje para identificar al token que está siendo creado, que será mandado en el mensaje de alerta cuando el mismo sea activado.
+ #hl[`--chat`]: El *ID* del chat de Telegram al que será mandada la alerta cuando el token que está siendo creado sea activado. Para encontrarlo, se le puede escribir a #link("https://t.me/RawDataBot")[#hl[#underline[este bot]]].

*Ejemplo de uso:*

Para crear un código QR que actúe como honeytoken:
#codeblock[`$ ./morfeo qr --msg "de ejemplo" --chat {chatID}`]

Esto genera un QR que al ser escaneado produce la siguiente alerta mediante Telegram:

#align(center, image("img/ejemplo.jpeg", width: 60%))

== Utilizando un contenedor de Docker

Si se desea en su lugar utilizar un contenedor de #link("https://www.docker.com/")[#hl[Docker]], deben seguirse los siguientes pasos.
Primero, pararse en la rama #hl[DockerTest], y crear los directorios que serán utilizados como entrada y salida con el contenedor:
#codeblock[`$ mkdir tokens input tmp`]

Luego, deben completarse las variables del #hl[`UID`] y #hl[`GID`] del #hl[`.env`] con los valores devueltos por los siguientes comandos, respectivamente: #hl[`id -u`] e #hl[`id -g`]. Una manera simple de realizar este paso es con el siguiente comando:

#codeblock[`$ sed -e "s/^UID=.*/UID=\"$(id -u)\"/" \
      -e "s/^GID=.*/GID=\"$(id -g)\"/" \
      env-sample > .env`]

Luego, se realiza un build del contenedor:
#codeblock[`$ docker compose build`]

Finalmente, cada vez que se desee utilizar la herramienta se ejecuta:
#codeblock[`$ docker compose run --rm morfeo-cli`]

Esto abre una terminal en el directorio #hl[`/app`] del contenedor. En él, se encuentra el binario #hl[`morfeo`] ya compilado, y dos directorios más: #hl[`/app/input`] y #hl[`/app/output`]. Estos directorios se encuentran linkeados a los directorios creados previamente. 

Para pasarle los archivos de entrada, los mismos deben ubicarse en el directorio #hl[`./input`] del host, y los mismos aparecerán en #hl[`/app/input`] del contenedor. De forma inversa, la herramienta generará los archivos en #hl[`/app/output`], y los mismos podrán ser encontrados en #hl[`./tokens`] del host. 

Luego, la utilización de la herramienta es similar a la mencionada previamente.