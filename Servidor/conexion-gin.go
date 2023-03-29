package main

import (
	"Handlers"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

func main() {
	// Set the router as the default one shipped with Gin
	router := gin.Default()

	//Habrá que crear un melody por cada paquete de ws
	chat_lobby := melody.New()
	prueba := melody.New()
	chat := melody.New()

	router.LoadHTMLFiles("chan.html")
	router.Use(static.Serve("/", static.LocalFile(".", true)))

	// Setup route group for the API
	api := router.Group("/api")
	{
		//Procesa una petición de login
		api.POST("/auth/login", Handlers.PostLogin)

		//Procesa una petición de registro
		api.POST("/auth/register", Handlers.PostRegister)

		//Devuelve lista de amigos confirmados
		api.GET("/amistad/get/:code", Handlers.GetAmistadList)

		//Elimina una relación de amistad
		api.POST("/amistad/remove", Handlers.PostAmistadRm)

		//Modifica nombre, foto y descripción de un jugador
		api.POST("/jugador/mod", Handlers.PostModJug)

		//Modifica la contraseña del jugador
		api.POST("/auth/mod-login", Handlers.PostModLogin)

		//Manda una solicitud de amistad
		api.POST("/amistad/add", Handlers.PostAmistadAdd)

		//Acepta una solicitud de amistad
		api.POST("/amistad/accept", Handlers.PostAmistadAccept)

		//Devuelve la información del usuario
		api.GET("api/jugador/get/:email", Handlers.GetInfoUsuario)

		//Rechaza una solicitud de amistad
		api.POST("/amistad/deny", Handlers.PostAmistadDeny)

		//Devulve los mensajes de un usuario
		api.GET("/msg/get/:code", Handlers.GetMsgList)

		//Pone a leidos los mensajes recibidos por el receptor del emidor
		api.POST("/msg/leer", Handlers.PostLeer)

		//WebSocket del chat entre amigos
		api.GET("/ws/chat/:code",func(c * gin.Context){
			chat.HandleRequest(c.Writer, c.Request)
		})

		//----------------Ejemplos-----------------------------//

		//Ejemplo de paso de parametros por url
		api.GET("/prueba/:param", getParam)

		//Ejemplo para devolver más de un dato
		api.GET("/prueba/names", getNames)

		//Ejemplo para devolver structs de datos
		//api.GET("/prueba/users", getUsers)

		//Carga la página del chat/lobby
		api.GET("/channel/:lobby", func(c *gin.Context) {
			c.HTML(http.StatusOK, "chan.html", nil)
		})

		api.GET("/ws/chat/lobby/:lobby", func(c *gin.Context) {
			//Pasa la petición al ws
			chat_lobby.HandleRequest(c.Writer, c.Request)
		})

		api.GET("/ws/prueba/patricia", func(c *gin.Context) {
			//Pasa la petición al ws
			prueba.HandleRequest(c.Writer, c.Request)
		})

	}

	//Retransmite lo enviado a todos cuya URL sea la misma (lobby)
	chat_lobby.HandleMessage(func(s *melody.Session, msg []byte) {
		chat_lobby.BroadcastFilter(msg, func(q *melody.Session) bool { //Envia la información a todos con la misma url
			return q.Request.URL.Path == s.Request.URL.Path
		})
	})

	//Retransmite el mensaje al ws del receptor del mensaje
	chat.HandleMessage(func(s *melody.Session, msg []byte) {
		//Hay que transformar el msg a JSON para obtener el receptor
		
		//Añadir el mensaje a la base de datos como no leido
		
		//Retransmitir el mensaje al receptor
		chat.BroadcastFilter(msg, func(q *melody.Session) bool { 
			return q.Request.URL.Path == ("api/ws/chat" + receptor)
		})
	})


	//Retransmite como JSON
	prueba.HandleMessage(func(s *melody.Session, msg []byte) {

		fmt.Println(s.Request.URL.Path)
		fmt.Println(prueba.Sessions())

		type M_rcp struct { //He movido el struct aquí para que no sea global
			Tipo      string `json:"tipo"`
			Contenido string `json:"contenido"`
		}

		var M M_rcp

		//Vemos que ha enviado el usuario //Prueba de Patricia
		err := json.Unmarshal(msg, &M)
		if err != nil {
			fmt.Println("Error al decodificar mensaje:", err)
			return
		} else {

			fmt.Println("Mensaje recibido: ", M.Tipo, M.Contenido)

			if M.Contenido == "Hola" {
				//Si necesitamos estructurar los datos a enviar
				type M_env struct {
					Msg string `json:"msg"` //parametros del struct empiezan con mayuscula
					//en json: ponemos el nombre del atributo json
				}
				m := M_env{Msg: "Hola esto es una prueba!"}

				b, _ := json.MarshalIndent(&m, "", "\t")

				prueba.Broadcast(b)

			} else {
				//Si no necesitamos manda información estructurada entre [] si son varios mensajes
				prueba.Broadcast([]byte(`{"msg": "Adios"}`))
			}
		}
	})

	// Start and run the server
	router.Run(":3001")
}

func getParam(c *gin.Context) {
	param := c.Param("param")
	c.JSON(http.StatusOK, gin.H{
		"name": param,
	})

}

func getNames(c *gin.Context) {
	var param = []string{"Adolfo", "Adrián", "Agustín", "Aitor", "Aitor-tilla"}
	c.JSON(http.StatusOK, gin.H{
		"name": param,
	})

}

/*
func getUsers(c *gin.Context) {
	var users = []Login{{"a@gmail.com", "1234"}, {"b@gmail.com", "5678"}, {"c@gmail.com", "91011"}, {"d@gmail.com", "1234"}}
	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}
*/
