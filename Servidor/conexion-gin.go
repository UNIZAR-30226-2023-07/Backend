package main

import (
	"DB/DAO"
	"DB/VO"
	"Handlers"
	"encoding/json"
	"fmt"

	//"math/rand"
	"net/http"
	//"strconv"

	"strings"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

func main() {
	// canales de todas las partidas -> clave: codigo_partida, valor: canal de la partida
	partidas := make(map[string]chan string)
	torneos := make(map[string]string)

	// Set the router as the default one shipped with Gin
	router := gin.Default()
	router.Use(cors.Default())

	//Habrá que crear un melody por cada paquete de ws
	chat_lobby := melody.New()
	prueba := melody.New()
	chat := melody.New()
	partidaNueva := melody.New()
	torneoNuevo := melody.New()

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
		api.POST("/amistad/remove", func(c *gin.Context) {
			Handlers.PostAmistadRm(c, chat)
		})

		//Manda una solicitud de amistad
		api.POST("/amistad/add", func(c *gin.Context) {
			Handlers.PostAmistadAdd(c, chat)
		})

		//Acepta una solicitud de amistad
		api.POST("/amistad/accept", func(c *gin.Context) {
			Handlers.PostAmistadAccept(c, chat)
		})

		//Rechaza una solicitud de amistad
		api.POST("/amistad/deny", func(c *gin.Context) {
			Handlers.PostAmistadDeny(c, chat)
		})

		//Devuelve la lista de solicitudes pendientes
		api.GET("/amistad/get/pendientes/:code", Handlers.GetPendientesList)

		//Devuelve la información del usuario
		api.GET("/jugador/get/:email", Handlers.GetInfoUsuario)

		//Devuelve la información del usuario
		api.GET("/jugador/get2/:code", Handlers.GetInfoUsuario2)

		//Modifica nombre, foto y descripción de un jugador
		api.POST("/jugador/mod", Handlers.PostModJug)

		//Borrar jugador
		api.GET("/jugador/del/:code", Handlers.DelJugador)

		//Devulve las partidas pendientes en las que participa el usuario
		api.GET("/partidas/pausadas/get/:code", Handlers.GetPausadas)

		//Devulve los mensajes de un usuario
		api.GET("/msg/get/:code", Handlers.GetMsgList)

		//Pone a leidos los mensajes recibidos por el receptor del emidor
		api.POST("/msg/leer", Handlers.PostLeer)

		//WebSocket del chat entre amigos
		api.GET("/ws/chat/:code", func(c *gin.Context) {
			chat.HandleRequest(c.Writer, c.Request)
		})

		//Crea una nueva partida
		api.POST("/partida/crear", func(c *gin.Context) {
			Handlers.CreatePartida(c, partidas, torneos)
		})

		//Unirse a un partida existente
		api.POST("/partida/join", func(c *gin.Context) {
			Handlers.JoinPartida(c, partidaNueva, torneoNuevo)
		})

		//Inicia una partida creada
		api.POST("/partida/iniciar", func(c *gin.Context) {
			Handlers.IniciarPartida(c, partidaNueva, torneoNuevo, partidas, torneos)
		})

		//Pausa una partida inciada
		api.POST("/partida/pausar", func(c *gin.Context) {
			Handlers.PausarPartida(c, partidaNueva, partidas)
		})

		//ws para transmitir la información del juego en una partida normal
		api.GET("/ws/partida/:lobby", func(c *gin.Context) {
			//Pasa la petición al ws
			partidaNueva.HandleRequest(c.Writer, c.Request)
		})

		//ws para transmitir la información del juego en torneo
		api.GET("/ws/torneo/:lobby", func(c *gin.Context) {
			//Pasa la petición al ws
			torneoNuevo.HandleRequest(c.Writer, c.Request)
		})

		//ws para el chat de partida
		api.GET("/ws/chat/lobby/:lobby", func(c *gin.Context) {
			//Pasa la petición al ws
			chat_lobby.HandleRequest(c.Writer, c.Request)
		})

		//----------------Ejemplos-----------------------------//

		//Carga la página del chat/lobby
		api.GET("/channel/:lobby", func(c *gin.Context) {
			c.HTML(http.StatusOK, "chan.html", nil)
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
		msgs := string(msg)
		fmt.Println(msgs)

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
		msgs := string(msg)
		fmt.Println(msgs)

		type Mensaje struct {
			Emisor string   `json:"emisor"`
			Tipo   string   `json:"tipo"`
			Cartas []string `json:"cartas"` // que sea ["1,2,3", "4,5,6", "7,8,9""]
			Info   string   `json:"info"`
		}

		type Respuesta struct {
			Emisor   string   `json:"emisor"`
			Receptor string   `json:"receptor"`
			Tipo     string   `json:"tipo"`
			Cartas   []string `json:"cartas"`
			Info     string   `json:"info"`
		}

		type RespuestaTablero struct {
			Emisor        string     `json:"emisor"`
			Receptor      string     `json:"receptor"`
			Tipo          string     `json:"tipo"`
			Mazo          []string   `json:"mazo"`
			Descartes     []string   `json:"descartes"`
			Combinaciones [][]string `json:"combinaciones"`
		}

		type RespuestaDescarte struct {
			Emisor        string     `json:"emisor"`
			Receptor      string     `json:"receptor"`
			Tipo          string     `json:"tipo"`
			Info          string     `json:"info"`
			Descartes     []string   `json:"descartes"`
			Combinaciones [][]string `json:"combinaciones"`
			Turno         string     `json:"turno"`
			Abrir         string     `json:"abrir"`
			Ganador       string     `json:"ganador"`
		}

		type RespuestaManos struct {
			Emisor   string     `json:"emisor"`
			Receptor string     `json:"receptor"`
			Tipo     string     `json:"tipo"`
			Manos    [][]string `json:"manos"`
		}

		var M Mensaje
		var R Respuesta
		var RT RespuestaTablero
		var RD RespuestaDescarte
		var RM RespuestaManos

		json.Unmarshal(msg, &M)

		R.Emisor = "Servidor"
		R.Receptor = M.Emisor
		R.Tipo = M.Tipo
		R.Cartas = M.Cartas
		R.Info = M.Info

		par := strings.Split(s.Request.URL.Path, "/")
		idPartida := par[4]

		if M.Tipo == "Jugadores" {
			fmt.Println("Se envia el numero de jugadores")
			partidas[s.Request.URL.Path] <- M.Info
			respuesta := <-partidas[s.Request.URL.Path]
			fmt.Println("Respuesta:", respuesta)
		} else if M.Tipo == "Robar_carta" || M.Tipo == "Robar_carta_descartes" {
			partidas[s.Request.URL.Path] <- M.Tipo
			fmt.Println("1")
			respuesta := <-partidas[s.Request.URL.Path]
			fmt.Println(respuesta)
			R.Info = respuesta
		} else if M.Tipo == "Fin_partida" || M.Tipo == "FIN" || M.Tipo == "END" {
			partidas[s.Request.URL.Path] <- M.Tipo
			respuesta := <-partidas[s.Request.URL.Path]
			fmt.Println(respuesta)
		} else if M.Tipo == "Abrir" || M.Tipo == "Colocar_combinacion" {
			partidas[s.Request.URL.Path] <- M.Tipo
			respuesta := <-partidas[s.Request.URL.Path]
			if respuesta == "Ok" {
				for i := 0; i < len(M.Cartas); i++ {
					// separamos M.Cartas[i] por comas y enviamos cada numero por el canal
					nums := strings.Split(M.Cartas[i], ",")
					for j := 0; j < len(nums); j++ {
						partidas[s.Request.URL.Path] <- nums[j]
					}
					// al terminar esta componente se envia "END"
					partidas[s.Request.URL.Path] <- "END"
					respuesta = <-partidas[s.Request.URL.Path]
					if respuesta != "Ok" {
						fmt.Println("Error:", respuesta)
						goto SALIR
					}
				}
			SALIR:
				partidas[s.Request.URL.Path] <- "FIN"
				respuesta = <-partidas[s.Request.URL.Path]
				fmt.Println(respuesta)
				R.Info = respuesta
				respuesta = <-partidas[s.Request.URL.Path]
				if respuesta == "ganador" {
					respuesta = <-partidas[s.Request.URL.Path]
					R.Info = respuesta
					Handlers.FinPartida(respuesta, idPartida)
				}
			} else {
				fmt.Println(respuesta)
				R.Info = respuesta
			}
		} else if M.Tipo == "Colocar_carta" {
			partidas[s.Request.URL.Path] <- M.Tipo
			respuesta := <-partidas[s.Request.URL.Path]
			if respuesta == "Ok" {
				parametros := strings.Split(M.Info, ",")
				partidas[s.Request.URL.Path] <- parametros[0]
				partidas[s.Request.URL.Path] <- parametros[1]

				respuesta = <-partidas[s.Request.URL.Path]
				if respuesta == "joker" {
					respuesta = <-partidas[s.Request.URL.Path]
				} else if respuesta == "ganador" {
					respuesta = <-partidas[s.Request.URL.Path]
					Handlers.FinPartida(respuesta, idPartida)
				}
				fmt.Println("Respuesta:", respuesta)
				R.Info = respuesta

			} else {
				fmt.Println(respuesta)
				R.Info = respuesta
			}
		} else if M.Tipo == "Descarte" {
			partidas[s.Request.URL.Path] <- M.Tipo
			respuesta := <-partidas[s.Request.URL.Path]
			fmt.Println(respuesta)
			RD.Info = respuesta
			if respuesta == "Ok" {
				partidas[s.Request.URL.Path] <- M.Info
				respuesta = <-partidas[s.Request.URL.Path] // info
				fmt.Println(respuesta)
				RD.Info = respuesta
				if respuesta == "Ok" {
					respuesta = <-partidas[s.Request.URL.Path]
					for respuesta != "fin" { // Descartes
						RD.Descartes = append(RD.Descartes, respuesta)
						respuesta = <-partidas[s.Request.URL.Path]
					}
					respuesta = <-partidas[s.Request.URL.Path]
					for respuesta != "fin" { // Combinaciones
						var comb []string
						for respuesta != "finC" {
							comb = append(comb, respuesta)
							respuesta = <-partidas[s.Request.URL.Path]
						}
						RD.Combinaciones = append(RD.Combinaciones, comb)
						respuesta = <-partidas[s.Request.URL.Path]
					}
					respuesta = <-partidas[s.Request.URL.Path] // ganador
					if respuesta == "ganador" {
						respuesta = <-partidas[s.Request.URL.Path]
						RD.Ganador = respuesta
						Handlers.FinPartida(respuesta, idPartida)
					} else {
						respuesta = <-partidas[s.Request.URL.Path] // turno
						RD.Turno = respuesta
						respuesta = <-partidas[s.Request.URL.Path] // abrir
						RD.Abrir = respuesta
					}

					fmt.Println(RD.Info)
				}
			}
		} else if M.Tipo == "Mostrar_mano" {
			partidas[s.Request.URL.Path] <- M.Tipo
			respuesta := <-partidas[s.Request.URL.Path]
			for respuesta != "fin" {
				R.Cartas = append(R.Cartas, respuesta)
				respuesta = <-partidas[s.Request.URL.Path]
			}
			fmt.Println(respuesta)
		} else if M.Tipo == "Mostrar_tablero" { // este no se usa
			partidas[s.Request.URL.Path] <- M.Tipo
			respuesta := <-partidas[s.Request.URL.Path]
			for respuesta != "fin" {
				RT.Mazo = append(RT.Mazo, respuesta)
				respuesta = <-partidas[s.Request.URL.Path]
			}
			respuesta = <-partidas[s.Request.URL.Path]
			for respuesta != "fin" {
				RT.Descartes = append(RT.Descartes, respuesta)
				respuesta = <-partidas[s.Request.URL.Path]
			}
			respuesta = <-partidas[s.Request.URL.Path]
			for respuesta != "fin" {
				var comb []string
				for respuesta != "finC" {
					comb = append(comb, respuesta)
					respuesta = <-partidas[s.Request.URL.Path]
				}
				RT.Combinaciones = append(RT.Combinaciones, comb)
				respuesta = <-partidas[s.Request.URL.Path]
			}
		} else if M.Tipo == "Mostrar_manos" { //CAMBIADO
			partidas[s.Request.URL.Path] <- M.Tipo
			respuesta := <-partidas[s.Request.URL.Path]
			for respuesta != "fin" {
				var mano []string
				for respuesta != "finJ" {
					mano = append(mano, respuesta)
					respuesta = <-partidas[s.Request.URL.Path]
				}
				RM.Manos = append(RM.Manos, mano)
				respuesta = <-partidas[s.Request.URL.Path]
			}
		}

		if M.Tipo == "Mostrar_tablero" {
			RT.Emisor = "Servidor"
			RT.Receptor = M.Emisor
			RT.Tipo = M.Tipo
			msg, _ = json.MarshalIndent(&RT, "", "\t")
		} else if M.Tipo == "Descarte" {
			RD.Emisor = "Servidor"
			RD.Receptor = "todos"
			RD.Tipo = M.Tipo
			msg, _ = json.MarshalIndent(&RD, "", "\t")
		} else if M.Tipo == "Mostrar_manos" {
			RM.Emisor = "Servidor"
			RM.Receptor = "todos"
			RM.Tipo = M.Tipo
			msg, _ = json.MarshalIndent(&RM, "", "\t")
		} else {
			msg, _ = json.MarshalIndent(&R, "", "\t")
		}

		partidaNueva.BroadcastFilter(msg, func(q *melody.Session) bool { //Envia la información a todos con la misma url
			return q.Request.URL.Path == s.Request.URL.Path
		})

	})

	torneoNuevo.HandleMessage(func(s *melody.Session, msg []byte) {
		msgs := string(msg)
		fmt.Println(msgs)

		type Mensaje struct {
			Emisor string   `json:"emisor"`
			Tipo   string   `json:"tipo"`
			Cartas []string `json:"cartas"` // que sea ["1,2,3", "4,5,6", "7,8,9""]
			Info   string   `json:"info"`
		}

		type Respuesta struct {
			Emisor   string   `json:"emisor"`
			Receptor string   `json:"receptor"`
			Tipo     string   `json:"tipo"`
			Cartas   []string `json:"cartas"`
			Info     string   `json:"info"`
		}

		type RespuestaTablero struct {
			Emisor        string     `json:"emisor"`
			Receptor      string     `json:"receptor"`
			Tipo          string     `json:"tipo"`
			Mazo          []string   `json:"mazo"`
			Descartes     []string   `json:"descartes"`
			Combinaciones [][]string `json:"combinaciones"`
		}

		type RespuestaDescarte struct {
			Emisor        string     `json:"emisor"`
			Receptor      string     `json:"receptor"`
			Tipo          string     `json:"tipo"`
			Info          string     `json:"info"`
			Descartes     []string   `json:"descartes"`
			Combinaciones [][]string `json:"combinaciones"`
			Turno         string     `json:"turno"`
			Abrir         string     `json:"abrir"`
			Ganador       string     `json:"ganador"`
		}

		type RespuestaManos struct {
			Emisor   string     `json:"emisor"`
			Receptor string     `json:"receptor"`
			Tipo     string     `json:"tipo"`
			Manos    [][]string `json:"manos"`
		}

		var M Mensaje
		var R Respuesta
		var RT RespuestaTablero
		var RD RespuestaDescarte
		var RM RespuestaManos

		json.Unmarshal(msg, &M)

		R.Emisor = "Servidor"
		R.Receptor = M.Emisor
		R.Tipo = M.Tipo
		R.Cartas = M.Cartas
		R.Info = M.Info

		if M.Tipo == "Jugadores" {
			fmt.Println("Se envia el numero de jugadores")
			partidas[torneos[s.Request.URL.Path]] <- M.Info
			respuesta := <-partidas[s.Request.URL.Path]
			fmt.Println("Respuesta:", respuesta)
		} else if M.Tipo == "Robar_carta" || M.Tipo == "Robar_carta_descartes" {
			partidas[s.Request.URL.Path] <- M.Tipo
			respuesta := <-partidas[s.Request.URL.Path]
			fmt.Println(respuesta)
			R.Info = respuesta
		} else if M.Tipo == "Fin_partida" || M.Tipo == "FIN" || M.Tipo == "END" {
			partidas[s.Request.URL.Path] <- M.Tipo
			respuesta := <-partidas[s.Request.URL.Path]
			fmt.Println(respuesta)
		} else if M.Tipo == "Abrir" || M.Tipo == "Colocar_combinacion" {
			partidas[s.Request.URL.Path] <- M.Tipo
			respuesta := <-partidas[s.Request.URL.Path]
			if respuesta == "Ok" {
				for i := 0; i < len(M.Cartas); i++ {
					// separamos M.Cartas[i] por comas y enviamos cada numero por el canal
					nums := strings.Split(M.Cartas[i], ",")
					for j := 0; j < len(nums); j++ {
						partidas[s.Request.URL.Path] <- nums[j]
					}
					// si quedan mas componentes se envia "END"
					if i < len(M.Cartas)-1 {
						partidas[s.Request.URL.Path] <- "END"
						respuesta = <-partidas[s.Request.URL.Path]
						if respuesta != "Ok" {
							fmt.Println("Error:", respuesta)
							goto SALIR
						}
					}
				}
			SALIR:
				partidas[s.Request.URL.Path] <- "FIN"
				respuesta = <-partidas[s.Request.URL.Path]
				fmt.Println(respuesta)
				respuesta = <-partidas[s.Request.URL.Path]
				if respuesta == "ganador" {
					respuesta = <-partidas[s.Request.URL.Path]
				}
				R.Info = respuesta
			} else {
				fmt.Println(respuesta)
				R.Info = respuesta
			}
		} else if M.Tipo == "Colocar_carta" {
			partidas[s.Request.URL.Path] <- M.Tipo
			respuesta := <-partidas[s.Request.URL.Path]
			if respuesta == "Ok" {
				parametros := strings.Split(M.Info, ",")
				partidas[s.Request.URL.Path] <- parametros[0]
				for i := 1; i < len(parametros); i++ {
					partidas[s.Request.URL.Path] <- parametros[i]
				}
				respuesta := <-partidas[s.Request.URL.Path]
				if respuesta == "joker" || respuesta == "ganador" {
					respuesta = <-partidas[s.Request.URL.Path]
				}
				fmt.Println(respuesta)
				R.Info = respuesta

			} else {
				fmt.Println(respuesta)
				R.Info = respuesta
			}
		} else if M.Tipo == "Descarte" {
			partidas[s.Request.URL.Path] <- M.Tipo
			respuesta := <-partidas[s.Request.URL.Path]
			fmt.Println(respuesta)
			RD.Info = respuesta
			if respuesta == "Ok" {
				partidas[s.Request.URL.Path] <- M.Info
				respuesta := <-partidas[s.Request.URL.Path] // info
				fmt.Println(respuesta)
				RD.Info = respuesta
				if respuesta == "Ok" {
					respuesta = <-partidas[s.Request.URL.Path]
					for respuesta != "fin" { // Descartes
						RD.Descartes = append(RD.Descartes, respuesta)
						respuesta = <-partidas[s.Request.URL.Path]
					}
					respuesta = <-partidas[s.Request.URL.Path]
					for respuesta != "fin" { // Combinaciones
						var comb []string
						for respuesta != "finC" {
							comb = append(comb, respuesta)
							respuesta = <-partidas[s.Request.URL.Path]
						}
						RD.Combinaciones = append(RD.Combinaciones, comb)
						respuesta = <-partidas[s.Request.URL.Path]
					}
					respuesta = <-partidas[s.Request.URL.Path] // ganador
					if respuesta == "ganador" {
						respuesta = <-partidas[s.Request.URL.Path]
						RD.Ganador = respuesta
					} else {
						respuesta = <-partidas[s.Request.URL.Path] // turno
						RD.Turno = respuesta
						respuesta = <-partidas[s.Request.URL.Path] // abrir
						RD.Abrir = respuesta
					}

					fmt.Println(RD.Info)
				}
			}
		} else if M.Tipo == "Mostrar_mano" {
			partidas[s.Request.URL.Path] <- M.Tipo
			respuesta := <-partidas[s.Request.URL.Path]
			for respuesta != "fin" {
				R.Cartas = append(R.Cartas, respuesta)
				respuesta = <-partidas[s.Request.URL.Path]
			}
			fmt.Println(respuesta)
		} else if M.Tipo == "Mostrar_tablero" { // este no se usa
			partidas[s.Request.URL.Path] <- M.Tipo
			respuesta := <-partidas[s.Request.URL.Path]
			for respuesta != "fin" {
				RT.Mazo = append(RT.Mazo, respuesta)
				respuesta = <-partidas[s.Request.URL.Path]
			}
			respuesta = <-partidas[s.Request.URL.Path]
			for respuesta != "fin" {
				RT.Descartes = append(RT.Descartes, respuesta)
				respuesta = <-partidas[s.Request.URL.Path]
			}
			respuesta = <-partidas[s.Request.URL.Path]
			for respuesta != "fin" {
				var comb []string
				for respuesta != "finC" {
					comb = append(comb, respuesta)
					respuesta = <-partidas[s.Request.URL.Path]
				}
				RT.Combinaciones = append(RT.Combinaciones, comb)
				respuesta = <-partidas[s.Request.URL.Path]
			}
		} else if M.Tipo == "Mostrar_manos" { //CAMBIADO
			partidas[s.Request.URL.Path] <- M.Tipo
			respuesta := <-partidas[s.Request.URL.Path]
			for respuesta != "fin" {
				var mano []string
				for respuesta != "finJ" {
					mano = append(mano, respuesta)
					respuesta = <-partidas[s.Request.URL.Path]
				}
				RM.Manos = append(RM.Manos, mano)
				respuesta = <-partidas[s.Request.URL.Path]
			}
		}

		if M.Tipo == "Mostrar_tablero" {
			RT.Emisor = "Servidor"
			RT.Receptor = M.Emisor
			RT.Tipo = M.Tipo
			msg, _ = json.MarshalIndent(&RT, "", "\t")
		} else if M.Tipo == "Descarte" {
			RD.Emisor = "Servidor"
			RD.Receptor = "todos"
			RD.Tipo = M.Tipo
			msg, _ = json.MarshalIndent(&RD, "", "\t")
		} else if M.Tipo == "Mostrar_manos" {
			RM.Emisor = "Servidor"
			RM.Receptor = "todos"
			RM.Tipo = M.Tipo
			msg, _ = json.MarshalIndent(&RM, "", "\t")
		} else {
			msg, _ = json.MarshalIndent(&R, "", "\t")
		}

		partidaNueva.BroadcastFilter(msg, func(q *melody.Session) bool { //Envia la información a todos con la misma url
			return q.Request.URL.Path == s.Request.URL.Path
		})

	})
	// Start and run the server
	router.Run(":3001")
}
