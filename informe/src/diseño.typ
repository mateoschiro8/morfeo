#pagebreak()

#import "@preview/chronos:0.2.1"

#let hl = text.with(fill: orange)


= Funcionamiento

La herramienta *morfeo* cuenta con dos partes: la *CLI* (_command line interface_) y el *server*. 
La *cli* es la encargada de procesar los comandos del usuario, comunicarse con el server y crear los tokens.
El *server* es el encargado de la identificación de los tokens, de detectar las activaciones de los mismos y dar las alertas.

Se incluye a continuación un diagrama que representa el flujo de funcionamiento de la herramienta.

#align(center)[
  #text(font: "Liberation Mono")[
    #chronos.diagram({
      import chronos: * 
      _par("User", show-bottom: false, color: orange)
      _par("CLI", show-bottom: false, color: orange)
      _par("Server", show-bottom: false, color: orange)
      
      _seq("User", "CLI", comment: "./morfeo {tokenType}", comment-align: "center")
      _seq("CLI", "Server", comment: "      POST /tokens      ")
      _seq("Server", "CLI", comment:"{tokenID}", comment-align: "center")
      _seq("CLI", "User", comment: " /{tokenType}/{tokenID} ")
    })
  ]
]

Explicación:
+ El usuario ejecuta #hl[`./morfeo {tokenType}`], donde #hl[`{tokenType}`] es uno de los formatos disponibles (#hl[`qr`], #hl[`bin`], #hl[`css`], etc), con las flags correspondientes del formato escogido.
+ La CLI se comunica con el server, mandando un *POST* a #hl[`tokens`], indicando la creación de un nuevo token.
+ El server se comunica con la base de datos, almacena la información mandada por la CLI, y devuelve el #hl[`tokenID`] generado por la misma.
+ La CLI toma el #hl[`tokenID`] recibido y construye la URL correspondiente, juntándolo con el tipo de token solicitado.

Notar que en el paso 2, al comunicarse con el server, no se indica qué tipo de token se está creando. No existe distinción desde el lado del server en los distintos tipos de tokens. Esto facilita el manejo de las activaciones, pero limita a que todos los tokens generados tengan el mismo método de _call-home_: una solicitud *GET* al servidor.

Cuando un token es activado, la solicitud *GET* es mandada a una URL con la forma #hl[`/{tokenType}/{tokenID}`]. La distinción del tipo de token en la URL permite que el handler ejecute los pasos adicionales (además de la alerta) de los tokens que así lo requieran (por ejemplo, la redirección en el código QR).

#pagebreak()

= Implementación 

== CLI

La *CLI* se encuentra implementada en #link("https://go.dev/")[#hl[Golang]], y fue utilizada una librería de creación de *CLIs* llamada #link("https://pkg.go.dev/github.com/spf13/cobra")[#hl[Cobra]]. La misma permite definir y agregar fácilmente nuevos comandos, además de realizar el _parsing_ de las flags recibidas.

Todos los tokens toman como entrada necesaria dos flags: *msg* y *chat*. La flag *msg* define qué mensaje será mandado cuando se realice la alerta de activación del token, y la flag *chat* es el ID del chat de Telegram donde será mandada la alerta. Cada comando puede tomar también otras flags más específicas de su funcionamiento. 

Independientemente del tipo de token a crear, cada comando recibe los flags correspondientes y utiliza la siguiente función de creación (se omitió el manejo de errores para ahorrar espacio):

```go
func CreateToken(msg string, extra string, chat string) string {
	data := types.UserInput{
		Msg:   msg,
		Extra: extra,
    Chat:  chat,
	}

	body, err := json.Marshal(data)
	resp, err := http.Post(serverURL+"/tokens",
                         "application/json",
                         bytes.NewBuffer(body))
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	tokenID := string(respBytes)
	return tokenID
}
```
Esta función toma tres parámetros y crea un struct de tipo #hl[`UserInput`]. Este struct contiene toda la información que será almacenada del token a crear. El campo #hl[`Msg`] almacena el mensaje a mostrar, el campo #hl[`Extra`] almacena aquella información más específica que necesitan algunos tokens para funcionar (se encuentra vacío para los otros), y el campo #hl[`Chat`], con el ID del chat de Telegram por al cual mandar el aviso.

Luego, transforma esa información a formato JSON, y lo manda en el cuerpo del pedido al server. El server se encarga del registro del token y el almacenamiento de los datos, y devuelve el #hl[`tokenID`] resultante, que esta función devuelve al comando para que continúe con la creación del token.

Esta arquitectura permite que agregar un nuevo formato de honeytoken sea tan simple como agregar un comando nuevo, y que el mismo utilice la función #hl[`createToken`] (además de definir el correspondiente handler en el server). Se detallan a continuación el funcionamiento e implementación de los distintos tipos de tokens disponibles.

=== QR 
#lorem(10)

=== Binario

La idea para crear un honeytoken a partir de un binario compilado es crear un nuevo binario que actúe de _wrapper_ del binario original. Es decir, que el binario resultante realice la alerta al servidor, y luego simule el comportamiento del binario original.

La función de creación de estos tokens tiene la siguiente forma:

```go
func generateBinaryWrapper(cmd *cobra.Command, args []string) {

	tokenID := CreateToken(msg, "", chat)

	data, err := os.ReadFile(in)
	b64 := base64.StdEncoding.EncodeToString(data)

	code := strings.ReplaceAll(wrapperTemplate, "{{B64}}", b64)
	code = strings.ReplaceAll(code, "{{Endpoint}}", 
														serverURL+"/bins/"+tokenID)

	os.WriteFile("tmp.go", []byte(code), 0644)

	outCmd := exec.Command("go", "build", "-o", out, "tmp.go")
	outCmd.Stdout = os.Stdout
	outCmd.Stderr = os.Stderr
	outCmd.Run()

	os.Remove("tmp.go")
}
```

Primero, llama a la función de creación de tokens mencionada anteriormente, y consigue el tokenID del nuevo token. 

Luego, lee el binario compilado que fue pasado como flag, y lo codifica en base64. Luego, toma el contenido de #hl[`wrapperTemplate`] y le "inyecta" el binario codificado y la URL a la que debe dar aviso.

Finalmente, crea un archivo llamado #hl[`tmp.go`] con los contenidos de dicho template (_wrapper_ + binario original), crea un comando que lo compila, lo ejecuta, y remueve el archivo fuente. Este binario compilado (llamado #hl[`tmp`] a menos que se utilice la flag _out_) es el honeytoken final.

El template que se compila para crear el binario final tiene la siguiente forma:

```go
const encoded = "{{B64}}"
const endpoint = "{{Endpoint}}"

func sendAlert() {
	client := http.Client{
    	Timeout: 2 * time.Second,
	}
	client.Get(endpoint)
}

func main() {
    
	sendAlert()

	data, _ := base64.StdEncoding.DecodeString(encoded)
	tmpDir, _ := os.MkdirTemp("", "honey-*")
	real := filepath.Join(tmpDir, "realbin")
	os.WriteFile(real, data, 0755)

	cmd := exec.Command(real, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin  = os.Stdin

	err := cmd.Run()
	if err != nil {
			if e, ok := err.(*exec.ExitError); ok {
					os.Exit(e.ExitCode())
			}
			panic(err)
	}
}
```
Lo primero que hace al ser ejecutado es mandar la alerta al servidor, a la URL inyectada por el generador. Con el objetivo de pasar más desapercibido, tiene un timeout de 2 segundos (por si el servidor no contesta). 

Luego, decodifica el binario original compilado, crea un directorio temporal en #hl[`/tmp`] con un sufijo aleatorio, y crea un archivo ahí dentro con los datos del binario. 

Finalmente, crea un comando de ejecución del archivo recién creado, le pasa los mismos _file descriptors_ de entrada, salida y error, y lo ejecuta. Además, en caso de que el binario original finalice con algún código de error, el binario wrapper finaliza de la misma forma.

=== PDF 
#lorem(10)

=== IMG 

En primer lugar, el comando acepta los siguientes flags:

-- `msg` (obligatorio): Identificador único del honeytoken

-- `chat` (obligatorio): ID del chat de Telegram donde se recibirá la alerta al activarse

-- `in` (opcional): Path a una imagen existente que se desee utilizar para generar el honeytoken

-- `out` (opcional): Ruta del archivo HTML de salida (por defecto: honeytoken_image.html)

Cuando se pasa el parametro *`--in`*, el programa:

 + Se fija si existe la imagen

+ Extrae las dimensiones de la imagen usando `image.DecodeConfig()`

+ Genera un archivo HTML que contiene:

  - La imagen original visible

  - Un pixel invisible de 1x1 que apunta a la URL indicada

Ademas, se genera un archivo SVG con estructura similar:

  - Elemento `<image>` principal con la imagen original
  
  - Elemento `<image>` de 1x1
  
  - ViewBox configurado según las dimensiones originales

Si no se proporciona imagen de entrada, el sistema crea:

  - HTML con únicamente un pixel invisible sobre un fondo blanco
  
  - SVG de 1x1 píxel conteniendo solo el honeytoken

Ambos formatos utilizan la técnica de _tracking pixel_: cuando se abre con algun editor de imagenes (no todos) el HTML o SVG, automáticamente realiza una petición HTTP GET a la URL embebida en la imagen

=== CSS
El canary token de CSS es particular por que no te avisa cuando el archivo es usado sino cuando tu pagina web fue clonada, es por ello que en este caso a las flags ya mencionadas se le agregan #hl[`in`], #hl[`out`] y #hl[`dominio`]. Las primeras dos flags indican cual es el arhivo CSS que se desea tokenizar y el nombre del token final (de forma predeterminada crea un archivo con el mismo nombre pero que arranca con new\_). Por su parte la flag #hl[`dominio`] indica cual es el dominio de la pagina del usuario. 
El funcionamiento del mismo consta de insertar al final del archivo CSS lo siguiente:
```css
body {
		background: url(https://morfeo-c8s3.onrender.com/fondo/{token_ID}) !important; 
}

```
Esta instruccion de CSS tendra el efecto de hacer un pedido GET a nuestro servidor en busca de una imagen, en este pedido los buscadores agregan un header llamado #hl[Referer] el cual indica desde que dominio se hizo el pedido. Por otra parte la flag !important le indica al buscador que no se debe saltear este background, asegurando que siempre se pida.

Por su parte nuestro servidor al recibir el pedido GET revisara el campo #hl[Referer] y comparara su contenido con el dominio del token, luego si el campo esta vacio se alerta al usuario indicando que es posible que se clonase la pagina y si el campo no coincide con el dominio de la pagina original se alerta que la pagina fue clonada indicando el dominio de la pagina clonada.

El token funciona pero se le podrian realizar algunas mejoras:
	- Compatibilidad con otros formatos: En desarollo web se suelen usar otros lenguajes que luego son traducidos a CSS, como es el caso de SASS. Una posible mejora seria hacer que el token sea compatible con dichos formatos para que no se tenga que volver a crear el token cada ves que se recompile el formato de alto nivel.
	- Dificultar la busqueda: Podriamos tomar algunas medidas para dificultar encontrar el token, por ejemplo se podria colocar en un lugar aleatorio del codigo CSS (no cambia mucho pero es algo), usar clases ya armadas es decir no crear una nueva definicio de body{} si ya existe una, etc.
	- Reducir Alertas: Podriamos hacer que si se detecta varias veces un clon desde el mismo dominio solo se alerte una ves, ademas se podria usar un timer para que pasado un tiempo se permita volver a alertar. 
#pagebreak()

== Server

El *server* también se encuentra implementado en #link("https://go.dev/")[#hl[Golang]], y fue utilizado un framework de manejo de peticiones *HTTP* y creación de aplicaciones web llamado #link("https://gin-gonic.com/")[#hl[Gin]]. Además, se utiliza una base de datos online llamada #link("https://www.mongodb.com/products/platform/atlas-database")[#hl[MongoDB Atlas]].

La función que inicia el server es la siguiente (omitiendo manejo de errores):

```go
func StartServer() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.Data(200, "text/html; charset=utf-8", []byte(morfeoString))
	})

	mongoURL := [...] 
	client, err := mongo.Connect(context.Background(), 
                               options.Client().ApplyURI(mongoURL))

	collection := client.Database("fcen").Collection("tokens")
	tokenController := handlers.NewTokenController(collection)

	// Para que el controller esté disponible en los handlers
	r.Use(func(c *gin.Context) {
		c.Set("tokenController", tokenController)
		c.Next()
	})

	r.POST("/tokens", handleNewToken)

	handlers.HandleQRs(r)
	handlers.HandleIMGs(r)
	handlers.HandleCSS(r)
	handlers.HandlePDFs(r)
	handlers.HandleBINs(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	r.Run(":" + port)
}
```

En ella, primero se define un *router*, donde se van a definir los endpoints del server y sus respectivos handlers. 

El *GET* a #hl[`/`] devuelve *`morfeoString`*, que no es más que un simple HTML que muestra el nombre del sistema. Debido a que la plataforma que fue utilizada para hacer el deploy "duerme" a las aplicaciones que no son utilizadas por un tiempo, este endpoint es utilizado para "despertar" a la aplicación.

Se construye luego la URL de conexión a la base de datos usando variables de ambiente (omitida por espacio), y se realiza la conexión. Se obtiene la colección de los tokens, y para facilitar el uso de la misma se crea un #hl[`tokenController`], que es guardado en el contexto para que esté disponible en todos los handlers.

Luego, se agrega en el *POST* a #hl[`/tokens`] la función de registro de los mismos, #hl[`handleNewToken`]. Esta función simplemente recupera los datos recibidos en el cuerpo del pedido, introduce un nuevo documento con los mismos en la base de datos y envía en la respuesta el #hl[`tokenID`] generado. 

Finalmente, se definen los handlers para las alertas de los tokens, y se levanta el server.

=== Handlers

Todos los handlers siguen una estructura como la siguiente, definiendo cada uno en su caso respuestas distintas o chequeos adicionales:

```go
func HandleTokenType(r *gin.Engine) {
	r.GET("/tokenType/:tokenID", func(c *gin.Context) {
		tokenID := c.Param("tokenID")

		controller := c.MustGet("tokenController").(*TokenController)
		token, err := controller.GetToken(tokenID)
		
    chat := token.Chat
		alertText := "Fue activado el token " + 
                 strings.ToLower(token.Msg) + 
                 " desde la IP: " + c.ClientIP()
		Alert(alertText, chat)
	})
}
```
En ellos, se obtiene el #hl[`tokenID`] de la URL, se recupera el #hl[`tokenController`] del contexto, y se lo utiliza para obtener la información del token correspondiente. Luego, se llama al método #hl[`Alert`], que recibe el mensaje y se encarga de hacer la alerta (en este caso, mediante un mensaje de Telegram al ID guardado).