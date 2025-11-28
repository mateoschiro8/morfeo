# morfeo
Un generador de Honeytokens. IT IS IN PROCESS

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
    - Rehacer los tokens con la nueva base. Tokens restantes: pdf, img, qrs, css. Proximos: bin 
    - Iniciar Ngrok con el metodo StartServer()
    - Ver de usar una base de datos (La mas simple podría ser MongoDB)
    - Pasar el proyecto a un docker
    - Acomodar las variables globales del cmd (se pueden reutilizar)
    - Hacer una función Avisar() en el package handlers que todos importen y usen
    - Provar alternativas a Ngrok que no manden un text para evitar fishing (posiblemente servero)

# Correr Ngrok
Usamos la flag de inspect-addr para que solamente pueda estar en esa y si el puerto esta ocupado falle 
    ngrok http 8000 --inspect=False

ese comando matara a cualquier ngrok corriendo en el sistema operativo, solo funciona en linux
    killall ngrok

si un paquete tiene ese header cuando se manda a ngrok, el resultado es el archivo directo, sin la pestaño de aviso. 
En el caso de css Firefox arma el paquete http por ello no tiene ese campo, el parche es que si uno abre el inspector, va a la pestaña red y ve el pedido get que fallo puede darle click izquierdo y editar el paquete y reenviar, ahi puede agregar el header:
    : true