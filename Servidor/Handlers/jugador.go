package Handlers

import (
	"DB/DAO"
	"DB/VO"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ModJug struct {
	Email  string `json:"email" binding:"required"`
	Nombre string `json:"nombre" binding:"required"`
	Foto   int    `json:"foto"` // binding:"required"`
	Descp  string `json:"descp" binding:"required"`
}

func PostModJug(c *gin.Context) {
	m := ModJug{}
	//Con el binding guardamos el json de la petición en u que es de tipo login
	if err := c.BindJSON(&m); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	jDAO := DAO.JugadoresDAO{}
	jVO := VO.NewJugadorVO(m.Nombre, "", m.Foto, m.Descp, 0, 0, 0, m.Email, "")

	jDAO.ModJugador(*jVO)

	c.JSON(http.StatusAccepted, gin.H{
		"res": "ok",
	})

}

func GetInfoUsuario(c *gin.Context) {
	email := c.Param("email")

	jDAO := DAO.JugadoresDAO{}

	jVO := jDAO.GetJugador(email)

	if jVO != nil {
		c.JSON(http.StatusOK, gin.H{
			"nombre":   jVO.GetNombre(),
			"foto":     jVO.GetFoto(),
			"descrp":   jVO.GetDescrip(),
			"pjugadas": jVO.GetPJugadas(),
			"pganadas": jVO.GetPGanadas(),
			"codigo":   jVO.GetCodigo(),
			"puntos":   jVO.GetPuntos(),
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{})
	}

}

func GetInfoUsuario2(c *gin.Context) {
	code := c.Param("code")

	jDAO := DAO.JugadoresDAO{}

	jVO := jDAO.GetJugador2(code)

	if jVO != nil {
		c.JSON(http.StatusOK, gin.H{
			"nombre":   jVO.GetNombre(),
			"foto":     jVO.GetFoto(),
			"descrp":   jVO.GetDescrip(),
			"pjugadas": jVO.GetPJugadas(),
			"pganadas": jVO.GetPGanadas(),
			"codigo":   jVO.GetCodigo(),
			"puntos":   jVO.GetPuntos(),
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{})
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

func DelJugador(c *gin.Context) {

	code := c.Param("code")

	jDAO := DAO.JugadoresDAO{}

	if jDAO.EstaJugador(code) {
		jDAO.DelJugador(code)
		c.JSON(http.StatusOK, gin.H{
			"res": "usuario eliminado",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"res": "usuario no existente",
		})
	}
}

func HistorialJugador(c *gin.Context) {

	code := c.Param("code")

	jDAO := DAO.JugadoresDAO{}

	type Msg struct {
		Tipo    string
		Creador string
		Clave   string
		Ganador bool
		Puntos  int
	}

	var msgs []Msg

	partidas, particip := jDAO.HisotialPartidas(code)
	for i := 0; i < len(partidas); i++ {

		var gana bool

		if partidas[i].GetTipo() == "amistosa" {
			gana = particip[i].GetPuntos() == 1
		} else {
			gana = particip[i].GetPuntos() < 100
		}

		m := Msg{
			Tipo:    partidas[i].GetTipo(),
			Creador: partidas[i].GetCreador(),
			Clave:   partidas[i].GetClave(),
			Ganador: gana,
			Puntos:  particip[i].GetPuntos(),
		}

		msgs = append(msgs, m)
	}

	c.JSON(http.StatusOK, gin.H{
		"partidas": msgs,
	})
}
