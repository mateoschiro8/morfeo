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

Para poder usar la herramienta, es necesario clonar el repositorio y tener #link("https://go.dev/")[#hl[Golang]] instalado. Luego, se deben ejecutar los siguientes comandos en la raiz del repositorio:

#codeblock[`$ cp env-samples .env`]
#codeblock[`$ go build`]

Esto generará un ejecutable llamado #hl[`morfeo`]. Luego, al ejecutar
#codeblock[`$ ./morfeo`]
desde donde se encuentre, se podrá ver la salida:
#codeblock[#set text(size: 10pt)
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

Para poder recibir las alertas, es necesario activar el bot que las manda. Para hacer esto, se entra al #link("https://t.me/morfeoSeguroBot ")[#hl[#underline[chat]]] con el mismo, y se le da #hl[`start`]. 

Luego, simplemente se ejecuta #hl[`./morfeo [command] [flags]`]. En caso de necesitar más información sobre un comando, ejecutarlo sin flags o con la flag #hl[`--help`] imprime más detalles sobre el mismo. 

Algunos comandos tienen flags que otros no tienen, pero hay dos flags que todos los comandos comparten:
+ #hl[`--msg`]: Un mensaje para identificar al token que está siendo creado, que será mandado en el mensaje de alerta cuando el mismo sea activado.
+ #hl[`--chat`]: El *ID* del chat de Telegram al que será mandada la alerta cuando el token que está siendo creado sea activado. Para encontrarlo, se le puede escribir a #link("https://t.me/RawDataBot")[#hl[#underline[este bot]]].

*Ejemplo de uso:*

Para crear un código QR que actúe como honeytoken:
#codeblock[`$ ./morfeo qr --msg "de ejemplo" --chat {chatID}`]

Esto genera un QR que al ser escaneado produce la siguiente alerta mediante Telegram:

#align(center, image("img/ejemplo.jpeg", width: 60%))
