package Handlers

import (
	"DB/DAO"
	"DB/VO"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CrearPart struct {
	Tipo      string `json:"tipo" binding:"required"`
	Anfitrion string `json:"anfitrion" binding:"required"`
}

type JoinPart struct {
	Codigo string `json:"codigo" binding:"required"`
	Clave  string `json:"clave" binding:"required"`
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

func JoinPartida(c *gin.Context) {

	p := JoinPart{}

	if err := c.BindJSON(&p); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	pDAO := DAO.PartidasDAO{}

	if !pDAO.EstaLlena(p.Clave) {

		parVO := VO.NewParticiparVO(p.Clave, p.Codigo, 0)
		parDAO := DAO.ParticiparDAO{}
		parDAO.AddParticipar(*parVO)

		c.JSON(http.StatusOK, gin.H{
			"res": "ok",
		})

	} else {

		c.JSON(http.StatusExpectationFailed, gin.H{
			"res": "Sala llena",
		})
	}

}

func IniciarPartida(c *gin.Context) {

	clave := c.Param("clave")

	pDAO := DAO.PartidasDAO{}

	pDAO.IniciarPartida(clave)

	c.JSON(http.StatusOK, gin.H{
		"res": "ok",
	})
}
