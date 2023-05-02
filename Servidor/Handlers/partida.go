package Handlers

import (
	"DB/DAO"
	"DB/VO"
	"Juego/partida"
	"Juego/torneo"
	"encoding/json"
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

type InitPart struct {
	Codigo string `json:"codigo" binding:"required"`
	Clave  string `json:"clave" binding:"required"`
	Bot	string `json:"bot" binding:"required"`
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

	// Generar identificador sumando 1 al id de la última partida guardada en el map
	id := len(partidas) + 1
	code := strconv.Itoa(id)
	_, ok := partidas[code]
	for ok || pDAO.HayPartida(code) {
		id = id + 1
		code = strconv.Itoa(id)
		_, ok = partidas[code]
	}

	// Crear canal para la partida y almacenarlo en el mapa
	partidas["/api/ws/partida/"+code] = make(chan string)

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

	parVO := VO.NewParticiparVO(code, p.Anfitrion, 1, 0, "no")

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

	j := parDAO.GetJugadoresEnLobby(p.Clave)

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
			"res":       "ok",
			"jugadores": j,
		})

	} else if !pDAO.EstaPausada(p.Clave) {

		n := pDAO.NJugadoresPartida(p.Clave)
		estor := pDAO.EsTorneo(p.Clave)
		if estor || n <= 4 {

			parVO := VO.NewParticiparVO(p.Clave, p.Codigo, 0, n, "no")
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
				"res":       "ok",
				"jugadores": j,
			})

		} else {

			c.JSON(http.StatusBadRequest, gin.H{
				"res": "Sala llena",
			})
		}
	}

}

func IniciarPartida(c *gin.Context, partidaNueva *melody.Melody, torneoNuevo *melody.Melody, partidas map[string]chan string, torneos map[string]string) {

	p := InitPart{}

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
		/*
		es_bot := make([]bool, 4) 
		b := 1
		if p.Bot == "si" {  
			meterBots = true
			for i := 0; i < len(es_bot); i++ {
				if i >= njug {
					es_bot[i] = true
					stringBot := ["bot" + strconv.Itoa(b),strconv.Itoa(i)]
					turnos = append(turnos, stringBot)
					b += 1
				} else {
					es_bot[i] = false
				}
			}
			// Indicar al DAO que hay bots en la partida
		}*/

		msg1, _ := json.MarshalIndent(&M1, "", "\t")

		pDAO.IniciarPartida(p.Clave)

		// Llamar a la función partida con el canal correspondiente
		if !pDAO.EsTorneo(p.Clave) {
			if pDAO.EstaPausada(p.Clave) {
				go partida.IniciarPartida(p.Clave, partidas["/api/ws/partida/"+p.Clave], true) // el bool indica que se ha pausado
			} else {
				go partida.IniciarPartida(p.Clave, partidas["/api/ws/partida/"+p.Clave], false) // el bool indica que no se ha pausado
			}
		} else {
			if pDAO.EstaPausada(p.Clave) {
				go torneo.IniciarTorneo(p.Clave, partidas[torneos["/api/ws/torneo/"+p.Clave]], true)	// el bool indica que se ha pausado
			} else{
				go torneo.IniciarTorneo(p.Clave, partidas[torneos["/api/ws/torneo/"+p.Clave]], false)	// el bool indica que no se ha pausado
			}
		}

		if pDAO.EsTorneo(p.Clave) {
			torneoNuevo.BroadcastFilter(msg1, func(q *melody.Session) bool { //Envia la información a todos con la misma url
				return q.Request.URL.Path == "/api/ws/torneo/"+p.Clave
			})

			partidas[torneos["/api/ws/torneo/"+p.Clave]] <- strconv.Itoa(njug)

		} else {
			partidaNueva.BroadcastFilter(msg1, func(q *melody.Session) bool { //Envia la información a todos con la misma url
				return q.Request.URL.Path == "/api/ws/partida/"+p.Clave
			})

			partidas["/api/ws/partida/"+p.Clave] <- strconv.Itoa(njug)

			if pDAO.EstaPausada(p.Clave) {
				mazo := pDAO.GetMazo(p.Clave)
				descarte := pDAO.GetDescarte(p.Clave)
				combinaciones := pDAO.GetCombinaciones(p.Clave)

				//Recuperar el mazo
				for i := 0; i < len(mazo); i++ {
					partidas["/api/ws/partida/"+p.Clave] <- strconv.Itoa(mazo[i].GetValor()) + "," + strconv.Itoa(mazo[i].GetPalo()) + "," + strconv.Itoa(mazo[i].GetReverso())
				}
				partidas["/api/ws/partida/"+p.Clave] <- "Fin_mazo"

				//Recuperar el descarte
				cartaDescarte := strconv.Itoa(descarte.GetValor()) + "," + strconv.Itoa(descarte.GetPalo()) + "," + strconv.Itoa(descarte.GetReverso())
				partidas["/api/ws/partida/"+p.Clave] <- cartaDescarte
				partidas["/api/ws/partida/"+p.Clave] <- "Fin_descartes"

				//Recuperar combinaciones
				ncomb := 0
				for i := 0; i < len(combinaciones); i++ {
					partidas["/api/ws/partida/"+p.Clave] <- strconv.Itoa((combinaciones[i].GetCarta()/10)/10) + "," + strconv.Itoa((combinaciones[i].GetCarta()/10)%10) + "," + strconv.Itoa(combinaciones[i].GetCarta()%10)
					if i+1 < len(combinaciones) && ncomb != combinaciones[i+1].GetNcomb() {
						ncomb = combinaciones[i+1].GetNcomb()
						partidas["/api/ws/partida/"+p.Clave] <- "Fin_combinacion"
					}
				}
				partidas["/api/ws/partida/"+p.Clave] <- "Fin_combinacion"

				//Recuperar manos
				//Como todos los jugadores de antes deben de estar en el lobby podemos usar esta funcion
				jugadores := parDAO.GetJugadoresEnLobby(p.Clave)
				for i := 0; i < len(jugadores); i++ {
					mano := pDAO.GetMano(p.Clave, jugadores[i])
					for j := 0; j < len(mano); j++ {
						partidas["/api/ws/partida/"+p.Clave] <- strconv.Itoa(mano[i].GetValor()) + "," + strconv.Itoa(mano[i].GetPalo()) + "," + strconv.Itoa(mano[i].GetReverso())
					}
					partidas["/api/ws/partida/"+p.Clave] <- "Fin_mano"
				}

				//Recuperamos abiertos
				abiertos := parDAO.GetAbierto(p.Clave)
				for i := 0; i < len(abiertos); i++ {
					partidas["/api/ws/partida/"+p.Clave] <- abiertos[i]
				}

				//pDAO.DelTableroGuardado(p.Clave)

			}
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

	partida := partidas["/api/ws/partida/"+p.Clave]

	pDAO := DAO.PartidasDAO{}
	parDAO := DAO.ParticiparDAO{}

	pVO := pDAO.GetPartida(p.Clave)

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

		// Lista de qué jugadores han abierto o no ordenados por turnos
		var ab []string
		respuesta = <-partida
		for respuesta != "fin" {
			ab = append(ab, respuesta)
			respuesta = <-partida
		}

		//Guardamos en la BD
		for i := 0; i < len(ab); i++ {
			parDAO.UpdateAbierto(p.Clave, i, ab[i])
		}

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

		//Guardamos las manos de cada jugador
		for i := 0; i < len(Manos); i++ {
			for j := 0; j < len(Manos[i]); j++ {
				man := strings.Split(Manos[i][j], ",")
				carta, _ := strconv.Atoi((man[0] + man[1] + man[2]))
				m := VO.NewManosVO(i, p.Clave, carta)
				pDAO.AddCartaMano(*m)
			}
		}

		//Guardamos el mazo de la partida
		for j := 0; j < len(Mazo); j++ {
			maz := strings.Split(Mazo[j], ",")
			carta, _ := strconv.Atoi((maz[0] + maz[1] + maz[2]))
			m := VO.NewMazosVO(p.Clave, carta)
			pDAO.AddCartaMazo(*m)
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