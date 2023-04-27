package Handlers

import (
	"DB/DAO"
	"DB/VO"
	"Juego/partida"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

type CrearPart struct {
	Tipo      string `json:"tipo" binding:"required"`
	Anfitrion string `json:"anfitrion" binding:"required"`
}

type JoinPart struct {
	Codigo string `json:"codigo" binding:"required"`
	Clave  string `json:"clave" binding:"required"`
}

// Retrasmitir mensaje en el ws de partida
type Mensaje struct {
	Emisor string   `json:"emisor"`
	Tipo   string   `json:"tipo"`
	Cartas []string `json:"cartas"`
	Info   string   `json:"info"`
}

type Turnos struct {
	Emisor string     `json:"emisor"`
	Tipo   string     `json:"tipo"`
	Turnos [][]string `json:"turnos"`
}

func CreatePartida(c *gin.Context, partidas map[string]chan string, torneos map[string]string) {

	pDAO := DAO.PartidasDAO{}
	parDAO := DAO.ParticiparDAO{}

	// Generar identificador único para la partida que no sea ninguna clave existente
	var code string
	for {
		code = strconv.Itoa(rand.Intn(9999))
		if _, ok := partidas[code]; !ok && !pDAO.HayPartida(code) {
			break
		}
	}

	// Crear canal para la partida y almacenarlo en el mapa
	partidas["/api/ws/partida/"+code] = make(chan string)

	// Llamar a la función partida con el canal correspondiente
	go partida.IniciarPartida(code, partidas["/api/ws/partida/"+code],false) // el bool indica que no se ha pausado

	p := CrearPart{}
	//Con el binding guardamos el json de la petición en u que es de tipo login
	if err := c.BindJSON(&p); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	pVO := VO.NewPartidasVO(code, p.Anfitrion, p.Tipo, "", "")

	if p.Tipo == "torneo" { //Guardamos la partida actual del torneo, en la BD puede ser la primera
		pVO = VO.NewPartidasVO(code, p.Anfitrion, p.Tipo, "", code)
		torneos[code] = "/api/ws/partida/" + code
	}

	pDAO.AddPartida(*pVO)

	parVO := VO.NewParticiparVO(code, p.Anfitrion, 1, 0)

	parDAO.AddParticipar(*parVO)

	c.JSON(http.StatusOK, gin.H{
		"clave": code,
	})

}

func JoinPartida(c *gin.Context, partidaNueva *melody.Melody, torneoNuevo *melody.Melody) {

	p := JoinPart{}

	if err := c.BindJSON(&p); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	pDAO := DAO.PartidasDAO{}
	parDAO := DAO.ParticiparDAO{}

	if pDAO.EstaPausada(p.Clave) && parDAO.EstaParticipando(p.Clave, p.Codigo) { //Si estaba pausada y el jugador estaba participando

		parDAO.ModLobbyJug(p.Codigo, p.Clave, 1) //Marcamos que ha llegado el jugador al lobby

		var M Mensaje

		M.Emisor = "Servidor"
		M.Tipo = "Nuevo_Jugador : " + p.Codigo

		msg, _ := json.MarshalIndent(&M, "", "\t")

		if pDAO.EsTorneo(p.Clave) {
			torneoNuevo.BroadcastFilter(msg, func(q *melody.Session) bool { //Envia la información a todos con la misma url
				return q.Request.URL.Path == "/api/ws/torneo/"+p.Clave
			})
		} else {
			partidaNueva.BroadcastFilter(msg, func(q *melody.Session) bool { //Envia la información a todos con la misma url
				return q.Request.URL.Path == "/api/ws/partida/"+p.Clave
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"res": "ok",
		})

	} else if !pDAO.EstaPausada(p.Clave) {

		n := pDAO.NJugadoresPartida(p.Clave)
		estor := pDAO.EsTorneo(p.Clave)
		if estor || n <= 4 {

			parVO := VO.NewParticiparVO(p.Clave, p.Codigo, 0, n)
			parDAO.AddParticipar(*parVO)

			var M Mensaje

			M.Emisor = "Servidor"
			M.Tipo = "Nuevo_Jugador: " + p.Codigo

			msg, _ := json.MarshalIndent(&M, "", "\t")

			if estor {
				torneoNuevo.BroadcastFilter(msg, func(q *melody.Session) bool { //Envia la información a todos con la misma url
					return q.Request.URL.Path == "/api/ws/torneo/"+p.Clave
				})
			} else {
				partidaNueva.BroadcastFilter(msg, func(q *melody.Session) bool { //Envia la información a todos con la misma url
					return q.Request.URL.Path == "/api/ws/partida/"+p.Clave
				})
			}

			c.JSON(http.StatusOK, gin.H{
				"res": "ok",
			})

		} else {

			c.JSON(http.StatusBadRequest, gin.H{
				"res": "Sala llena",
			})
		}
	}

}

func IniciarPartida(c *gin.Context, partidaNueva *melody.Melody, torneoNuevo *melody.Melody, partidas map[string]chan string, torneos map[string]string) {

	p := JoinPart{}

	if err := c.BindJSON(&p); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	pDAO := DAO.PartidasDAO{}
	parDAO := DAO.ParticiparDAO{}

	if pDAO.HayPartida(p.Clave) && pDAO.EsCreador(p.Clave, p.Codigo) && pDAO.JugadoresEnLobby(p.Clave) {

		turnos := parDAO.GetJugadoresTurnos(p.Clave)
		njug := pDAO.NJugadoresPartida(p.Clave)

		var M1 Turnos
		M1.Emisor = "Servidor"
		M1.Tipo = "Partida_Iniciada"
		M1.Turnos = turnos
		msg1, _ := json.MarshalIndent(&M1, "", "\t")

		pDAO.IniciarPartida(p.Clave)

		if pDAO.EsTorneo(p.Clave) {
			torneoNuevo.BroadcastFilter(msg1, func(q *melody.Session) bool { //Envia la información a todos con la misma url
				return q.Request.URL.Path == "/api/ws/torneo/"+p.Clave
			})

			partidas[torneos["/api/ws/torneo/"+p.Clave]] <- strconv.Itoa(njug)

		} else {
			partidaNueva.BroadcastFilter(msg1, func(q *melody.Session) bool { //Envia la información a todos con la misma url
				return q.Request.URL.Path == "/api/ws/partida/"+p.Clave
			})

			fmt.Println(partidas)

			partidas["/api/ws/partida/"+p.Clave] <- strconv.Itoa(njug)

			if pDAO.EstaPausada(p.Clave) {
				//falta mazo del tablero
				descarte := pDAO.GetDescarte(p.Clave)
				//combinaciones := pDAO.GetCombinaciones(p.Clave)

				/* DESCOMENTAR CUANDO ESTÉ EL MAZO
				for i := 0; i < len(mazo); i++ {
					partidas["/api/ws/partida/"+p.Clave] <- mazo[i]
				}
				partidas["/api/ws/partida/"+p.Clave] <- "Fin_mazo"
				*/

				cartaDescarte := strconv.Itoa(descarte.GetValor()) + strconv.Itoa(descarte.GetPalo()) + strconv.Itoa(descarte.GetReverso())
				partidas["/api/ws/partida/"+p.Clave] <- cartaDescarte
				partidas["/api/ws/partida/"+p.Clave] <- "Fin_descartes"

				/* No se como pasar del vo a un [][]string pero combinaciones tiene que ser asi:
					[["valor,palo,reverso","valor,palo,reverso","valor,palo,reverso"],["valor,palo,reverso","valor,palo,reverso","valor,palo,reverso"]]
					Ese ejemplo serían dos combinaciones.

					Y despues de pasarlo a eso hay que descomentar esto de abajo

				for i := 0; i < len(combinaciones); i++ {
					for j := 0; j < len(combinaciones[i]); j++ {
						partidas["/api/ws/partida/"+p.Clave] <- combinaciones[i][j]
					}
					partidas["/api/ws/partida/"+p.Clave] <- "Fin_combinacion"
				}
				partidas["/api/ws/partida/"+p.Clave] <- "Fin_combinaciones"
				*/


				/* Con los jugadores tampoco lo sé :(
					Habría que descomentar esto para enviarlo.

				for i := 0; i < len(jugadores); i++ {
					for j := 0; j < len(jugadores[i]); j++ {
						partidas["/api/ws/partida/"+p.Clave] <- jugadores[i][j]
					}
					partidas["/api/ws/partida/"+p.Clave] <- "Fin_mano"
				}
				*/

				/* DESCOMENTAR CUANDO ESTÉ EL VECTOR DE QUIÉN HA ABIERTO
				for i := 0; i < len(abiertos); i++ {
				partidas["/api/ws/partida/"+p.Clave] <- strconv.FormatBool(abiertos[i])
				}
				*/

			}
		}

		if pDAO.EstaPausada(p.Clave) {
			//Recuperar las cartas
			//pDAO.DelTableroGuardado(p.Clave) //Una vez recuperadas borramos la informacion
		} else {
			//Repartir las cartas
		}

		c.JSON(http.StatusOK, gin.H{
			"res": "ok",
		})

	} else {

		c.JSON(http.StatusBadRequest, gin.H{
			"res": "error",
		})

	}

}

func PausarPartida(c *gin.Context, partidaNueva *melody.Melody, partidas map[string]chan string) {

	p := JoinPart{}
	//Con el binding guardamos el json de la petición en u que es de tipo login
	if err := c.BindJSON(&p); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	partida := partidas["api/ws/partida/"+p.Clave]

	pDAO := DAO.PartidasDAO{}
	parDAO := DAO.ParticiparDAO{}

	pVO := pDAO.GetPartida(p.Clave)
	fmt.Println(pVO)

	if pVO.GetCreador() == p.Codigo {
		pDAO.PausarPartida(p.Clave) //Marcamos partida como pausada

		var Mazo []string
		var Descartes []string
		var Combinaciones [][]string
		var Manos [][]string
		partida <- "Pausar"
		// Devuelve el tablero: mazo, descartes y combinaciones
		respuesta := <-partida
		for respuesta != "fin" {
			Mazo = append(Mazo, respuesta)
			respuesta = <-partida
		}
		respuesta = <-partida
		for respuesta != "fin" {
			Descartes = append(Descartes, respuesta)
			respuesta = <-partida
		}
		respuesta = <-partida
		for respuesta != "fin" {
			var comb []string
			for respuesta != "finC" {
				comb = append(comb, respuesta)
				respuesta = <-partida
			}
			Combinaciones = append(Combinaciones, comb)
			respuesta = <-partida
		}

		// Lista con la mano de cada jugador --> [["1,2,3","4,5,6"],["1,2,3","4,5,6]] --> cada string es valor,palo,color y cada lista es una mano
		respuesta = <-partida
		for respuesta != "fin" {
			var mano []string
			for respuesta != "finJ" {
				mano = append(mano, respuesta)
				respuesta = <-partida
			}
			Manos = append(Manos, mano)
			respuesta = <-partida
		}

		// Lista de qué jugadores han abierto o no
		/* 
		var ab []string
		respuesta = <-partida
		for respuesta != "fin" {
			ab = append(ab,respuesta)
			respuesta = <-partida
		}
		*/



		//Guardamos combinaciones en la BD -> [["1,2,3","2,3,1","1,2,3"]]
		for i := 0; i < len(Combinaciones); i++ {
			for j := 0; j < len(Combinaciones[i]); j++ {
				comb := strings.Split(Combinaciones[i][j], ",")
				carta, _ := strconv.Atoi((comb[0] + comb[1] + comb[2]))
				c := VO.NewCombinacionesVO(p.Clave, carta, i)
				pDAO.AddCombinacion(*c)
			}
		}

		//Guardamos los descartes en la BD
		for j := 0; j < len(Descartes); j++ {
			desc := strings.Split(Descartes[j], ",")
			carta, _ := strconv.Atoi((desc[0] + desc[1] + desc[2]))
			d := VO.NewDescartesVO(p.Clave, carta)
			pDAO.AddDescarte(*d)
		}

		//Guardamos el mazo de cada jugador (Por ahora no guardamos turno)
		for i := 0; i < len(Manos); i++ {
			for j := 0; j < len(Manos[i]); j++ {
				man := strings.Split(Manos[i][j], ",")
				carta, _ := strconv.Atoi((man[0] + man[1] + man[2]))
				m := VO.NewMazosVO(i, p.Clave, carta)
				pDAO.AddCartaMazo(*m)
			}
		}

		parDAO.ModLobby(p.Clave, 0) //Guardamos que los jugadores ya no estan en lobby o jugando

		var M Mensaje

		M.Emisor = "Servidor"
		M.Tipo = "Partida_Pausada"

		msg, _ := json.MarshalIndent(&M, "", "\t")

		partidaNueva.BroadcastFilter(msg, func(q *melody.Session) bool { //Envia la información a todos con la misma url
			return q.Request.URL.Path == "/api/ws/partida/"+p.Clave
		})

		c.JSON(http.StatusOK, gin.H{
			"res": "ok",
		})

	} else {

		c.JSON(http.StatusBadRequest, gin.H{
			"res": "error",
		})
	}

}
