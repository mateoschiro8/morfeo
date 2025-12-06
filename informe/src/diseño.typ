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
#lorem(10)

=== PDF 
#lorem(10)

=== IMG 
#lorem(10)

=== CSS
#lorem(10)

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

Se construye luego la URL de conexión a la base de datos usando variables de ambiente, y se realiza la conexión. Se obtiene la colección de los tokens, y para facilitar el uso de la misma se crea un #hl[`tokenController`], que es guardado en el contexto para que esté disponible en todos los handlers.

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