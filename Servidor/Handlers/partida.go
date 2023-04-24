package Handlers

import (
	"DB/DAO"
	"DB/VO"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

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
}

func CreatePartida(c *gin.Context, clave string) {

	p := CrearPart{}
	//Con el binding guardamos el json de la petición en u que es de tipo login
	if err := c.BindJSON(&p); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	pDAO := DAO.PartidasDAO{}

	var pVO *VO.PartidasVO

	if p.Tipo == "torneo" {

		torneo := strconv.Itoa(rand.Intn(9999))
		for pDAO.HayPartida(torneo) {
			torneo = strconv.Itoa(rand.Intn(9999))
		}

		tVO := VO.NewPartidasVO(torneo, p.Anfitrion, p.Tipo, "", "")

		pDAO.AddPartida(*tVO)

		pVO = VO.NewPartidasVO(clave, p.Anfitrion, "amistosa", "", torneo)

		c.JSON(http.StatusOK, gin.H{
			"clave": torneo,
		})

	} else {

		pVO = VO.NewPartidasVO(clave, p.Anfitrion, p.Tipo, "", "")

		c.JSON(http.StatusOK, gin.H{
			"clave": clave,
		})
	}

	pDAO.AddPartida(*pVO)

	parDAO := DAO.ParticiparDAO{}

	parVO := VO.NewParticiparVO(clave, p.Anfitrion, 0)

	parDAO.AddParticipar(*parVO)

}

func JoinPartida(c *gin.Context, partidaNueva *melody.Melody, nuevoLobby string) bool {

	p := JoinPart{}

	if err := c.BindJSON(&p); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return false
	}

	pDAO := DAO.PartidasDAO{}
	parDAO := DAO.ParticiparDAO{}

	if pDAO.EsTorneo(p.Clave) {
		//Comprobamos si hay alguna sala vacía dentro del torneo
		hay, sala := pDAO.HayPartidaTorneo(p.Clave)

		if hay { //Si hay una sala disponible lo añadimos

			parVO := VO.NewParticiparVO(sala, p.Codigo, 1)
			parDAO.AddParticipar(*parVO)

			var M Mensaje

			M.Emisor = "Servidor"
			M.Tipo = "Nuevo_Jugador : " + p.Codigo

			msg, _ := json.MarshalIndent(&M, "", "\t")

			partidaNueva.BroadcastFilter(msg, func(q *melody.Session) bool { //Envia la información a todos con la misma url
				return q.Request.URL.Path == "/api/ws/partida/"+p.Clave
			})

			c.JSON(http.StatusOK, gin.H{
				"res": "ok",
			})

			return false

		} else { //Si no hay salas disponibles creamos una

			tVO := pDAO.GetPartida(p.Clave)

			pVO := VO.NewPartidasVO(nuevoLobby, tVO.GetCreador(), "amistosa", "", p.Clave)
			pDAO.AddPartida(*pVO)

			parVO := VO.NewParticiparVO(nuevoLobby, p.Codigo, 1)
			parDAO.AddParticipar(*parVO)

			//Como será el primero entrar no hay que avisar de que entra
			c.JSON(http.StatusOK, gin.H{
				"res": nuevoLobby,
			})

			return true
		}

	} else { //Esta parte será comun en las partidas normales y en las continuadas

		if pDAO.EstaPausada(p.Clave) && parDAO.EstaParticipando(p.Clave, p.Codigo) {
			parDAO.ModLobbyJug(p.Codigo, p.Clave, 1)

			var M Mensaje

			M.Emisor = "Servidor"
			M.Tipo = "Nuevo_Jugador : " + p.Codigo

			msg, _ := json.MarshalIndent(&M, "", "\t")

			partidaNueva.BroadcastFilter(msg, func(q *melody.Session) bool { //Envia la información a todos con la misma url
				return q.Request.URL.Path == "/api/ws/partida/"+p.Clave
			})

			c.JSON(http.StatusOK, gin.H{
				"res": "ok",
			})

		} else if !pDAO.EstaPausada(p.Clave) {

			if !pDAO.EstaLlena(p.Clave) {

				parVO := VO.NewParticiparVO(p.Clave, p.Codigo, 1)
				parDAO.AddParticipar(*parVO)

				var M Mensaje

				M.Emisor = "Servidor"
				M.Tipo = "Nuevo_Jugador: " + p.Codigo

				msg, _ := json.MarshalIndent(&M, "", "\t")

				partidaNueva.BroadcastFilter(msg, func(q *melody.Session) bool { //Envia la información a todos con la misma url
					return q.Request.URL.Path == "/api/ws/partida/"+p.Clave
				})

				c.JSON(http.StatusOK, gin.H{
					"res": "ok",
				})

			} else {

				c.JSON(http.StatusBadRequest, gin.H{
					"res": "Sala llena",
				})
			}

		} else {

			c.JSON(http.StatusBadRequest, gin.H{
				"res": "error",
			})
		}

		return false

	}

}

func IniciarPartida(c *gin.Context, partidaNueva *melody.Melody) {

	p := JoinPart{}

	if err := c.BindJSON(&p); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	pDAO := DAO.PartidasDAO{}

	if pDAO.HayPartida(p.Clave) && pDAO.EsCreador(p.Clave, p.Codigo) && pDAO.JugadoresEnLobby(p.Clave) {

		var M Mensaje

		M.Emisor = "Servidor"
		M.Tipo = "Partida_Iniciada"

		msg, _ := json.MarshalIndent(&M, "", "\t")

		if pDAO.EsTorneo(p.Clave) {

			partidas := pDAO.GetPartidasTorneo(p.Clave)

			fmt.Println(partidas)

			for i := 0; i < len(partidas); i++ {

				partidaNueva.BroadcastFilter(msg, func(q *melody.Session) bool { //Envia la información a todos con la misma url
					return q.Request.URL.Path == "/api/ws/partida/"+partidas[i]
				})

				pDAO.IniciarPartida(partidas[i])

			}

			//Repartir las cartas -> los torneos no se pueden parar

		} else {

			partidaNueva.BroadcastFilter(msg, func(q *melody.Session) bool { //Envia la información a todos con la misma url
				return q.Request.URL.Path == "/api/ws/partida/"+p.Clave
			})

			pDAO.IniciarPartida(p.Clave)

			if pDAO.EstaPausada(p.Clave) {
				//Recuperar las cartas
				//pDAO.DelTableroGuardado(p.Clave)	//Una vez recuperadas borramos la informacion
			} else {
				//Repartir las cartas
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

func PausarPartida(c *gin.Context, partidaNueva *melody.Melody, partida chan string) {

	p := JoinPart{}
	//Con el binding guardamos el json de la petición en u que es de tipo login
	if err := c.BindJSON(&p); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	pDAO := DAO.PartidasDAO{}
	parDAO := DAO.ParticiparDAO{}

	//POR AHORA LOS TONEOS NO SE PUEDEN PARAR!!

	pVO := pDAO.GetPartida(p.Clave)
	fmt.Println(pVO)

	if pVO.GetCreador() == p.Codigo && pVO.GetTorneo() == "" && pVO.GetTipo() != "torneo" {
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

		// no se que funcion del dao sería pero habria que guardar el mazo para robar(Mazo), el descarte(Descartes) y las combinaciones(Combinaciones), y tambien
		// las manos de cada jugador(Manos). Los jugadores están ordenados segun el turno que tienen.


		//Habrá que guardar las combinaciones, cada carta de cada jugador y el descarte
		//pDAO.AddCartaMazo(m)
		//pDAO.AddCombinacion(c)
		//pDAO.AddDescarte(d)

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
