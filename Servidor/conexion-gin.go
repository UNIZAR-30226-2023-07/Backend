package main

import (
	"DB/DAO"
	"DB/VO"
	"Handlers"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"Servidor/Juego/partida"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

func main() {
	// canales de todas las partidas -> clave: codigo_partida, valor: canal de la partida
	partidas := make(map[string]chan string)

	// Set the router as the default one shipped with Gin
	router := gin.Default()
	router.Use(cors.Default())

	//Habrá que crear un melody por cada paquete de ws
	chat_lobby := melody.New()
	prueba := melody.New()
	chat := melody.New()
	partidaNueva := melody.New()

	router.LoadHTMLFiles("chan.html")
	router.Use(static.Serve("/", static.LocalFile(".", true)))

	// Setup route group for the API
	api := router.Group("/api")
	{
		//Procesa una petición de login
		api.POST("/auth/login", Handlers.PostLogin)

		//Procesa una petición de registro
		api.POST("/auth/register", Handlers.PostRegister)

		//Modifica la contraseña del jugador
		api.POST("/auth/mod-login", Handlers.PostModLogin)

		//Devuelve lista de amigos confirmados
		api.GET("/amistad/get/:code", Handlers.GetAmistadList)

		//Elimina una relación de amistad
		api.POST("/amistad/remove", Handlers.PostAmistadRm)

		//Manda una solicitud de amistad
		api.POST("/amistad/add", Handlers.PostAmistadAdd)

		//Acepta una solicitud de amistad
		api.POST("/amistad/accept", Handlers.PostAmistadAccept)

		//Rechaza una solicitud de amistad
		api.POST("/amistad/deny", Handlers.PostAmistadDeny)

		//Devuelve la lista de solicitudes pendientes
		api.GET("/amistad/get/pendientes/:code", Handlers.GetPendientesList)

		//Devuelve la información del usuario
		api.GET("/jugador/get/:email", Handlers.GetInfoUsuario)

		//Devuelve la información del usuario
		api.GET("/jugador/get2/:code", Handlers.GetInfoUsuario2)

		//Modifica nombre, foto y descripción de un jugador
		api.POST("/jugador/mod", Handlers.PostModJug)

		//Devulve los mensajes de un usuario
		api.GET("/msg/get/:code", Handlers.GetMsgList)

		//Pone a leidos los mensajes recibidos por el receptor del emidor
		api.POST("/msg/leer", Handlers.PostLeer)

		//WebSocket del chat entre amigos
		api.GET("/ws/chat/:code", func(c *gin.Context) {
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

		api.GET("/ws/partida", func(c *gin.Context) {
			//Pasa la petición al ws
			partidaNueva.HandleRequest(c.Writer, c.Request)
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

		//Estructuramos el mesaje para sacar el receptor del mismo
		type M_rcp struct {
			Emisor    string `json:"emisor"`
			Receptor  string `json:"receptor"`
			Contenido string `json:"contenido"`
		}

		var M M_rcp

		json.Unmarshal(msg, &M)
		fmt.Println(M)

		//Guardamos el mensaje como no leido en la BD
		mVO := VO.NewMensajesVO(M.Emisor, M.Receptor, M.Contenido, 0)
		mDAO := DAO.MensajesDAO{}
		mDAO.AddMensaje(*mVO)

		//Retransmitir el mensaje al receptor
		chat.BroadcastFilter(msg, func(q *melody.Session) bool {
			return q.Request.URL.Path == ("/api/ws/chat/" + M.Receptor)
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

	//crea una nueva partida y envia el código de la misma por un canal
	partidaNueva.HandleMessage(func(s *melody.Session, msg []byte) {

		// Generar identificador único para la partida que no sea ninguna clave existente
		var code string
		for {
			code = strconv.Itoa(rand.Intn(9999))
			if _, ok := partidas[code]; !ok {
				break
			}
		}

		// Crear canal para la partida y almacenarlo en el mapa
		partidas[code] = make(chan string)

		// Llamar a la función partida con el canal correspondiente
		go partida.IniciarPartida(code, partidas[code])
		fmt.Println("Se ha iniciado la partida: ", code)

		// Responder al jugador con el identificador de la partida
		s.Write([]byte(code))
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
