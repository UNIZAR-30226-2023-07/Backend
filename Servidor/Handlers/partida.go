package Handlers

import (
	"DB/DAO"
	"DB/VO"
	"encoding/json"
	"net/http"

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

	pVO := VO.NewPartidasVO(clave, p.Anfitrion, p.Tipo, "")

	pDAO := DAO.PartidasDAO{}

	pDAO.AddPartida(*pVO)

	parDAO := DAO.ParticiparDAO{}

	parVO := VO.NewParticiparVO(clave, p.Anfitrion, 0)

	parDAO.AddParticipar(*parVO)

	c.JSON(http.StatusOK, gin.H{
		"clave": clave,
	})

}

func JoinPartida(c *gin.Context, partidaNueva *melody.Melody) {

	p := JoinPart{}

	if err := c.BindJSON(&p); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	pDAO := DAO.PartidasDAO{}
	parDAO := DAO.ParticiparDAO{}

	if !pDAO.EstaLlena(p.Clave) {

		if pDAO.EstaPausada(p.Clave) {
			parDAO.ModLobbyJug(p.Codigo, p.Clave, 1)
		} else {
			parVO := VO.NewParticiparVO(p.Clave, p.Codigo, 1)
			parDAO.AddParticipar(*parVO)
		}

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

	} else {

		c.JSON(http.StatusBadRequest, gin.H{
			"res": "Sala llena",
		})
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

		c.JSON(http.StatusOK, gin.H{
			"res": "ok",
		})

	} else {

		c.JSON(http.StatusBadRequest, gin.H{
			"res": "error",
		})

	}

}

func PausarPartida(c *gin.Context, partidaNueva *melody.Melody) {

	p := JoinPart{}
	//Con el binding guardamos el json de la petición en u que es de tipo login
	if err := c.BindJSON(&p); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	pDAO := DAO.PartidasDAO{}
	parDAO := DAO.ParticiparDAO{}

	if pDAO.HayPartida(p.Clave) && pDAO.EsCreador(p.Clave, p.Codigo) {
		pDAO.PausarPartida(p.Clave) //Marcamos partida como pausada

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

func GetPausadas(c *gin.Context) {
	code := c.Param("code")

	jDAO := DAO.JugadoresDAO{}

	part := jDAO.PartidasPausadas(code)

	type Partida struct {
		Tipo    string
		Creador string
		Clave   string
	}

	var partidas []Partida

	for i := 0; i < len(part); i++ {
		p := Partida{
			Tipo:    part[i].GetTipo(),
			Creador: part[i].GetCreador(),
			Clave:   part[i].GetClave(),
		}
		partidas = append(partidas, p)
	}

	c.JSON(http.StatusOK, gin.H{
		"partidas": partidas,
	})
}
