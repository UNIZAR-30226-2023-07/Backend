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
	Foto   int    `json:"foto" binding:"required"`
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
	jVO := VO.NewJugadorVO(m.Nombre, "", m.Foto, m.Descp, 0, 0, m.Email, "")

	jDAO.ModJugador(*jVO)

	c.JSON(http.StatusAccepted, gin.H{
		"res": "ok",
	})

}

func GetInfoUsuario(c *gin.Context){
	email := c.Param("email")

	jDAO := DAO.JugadoresDAO{}

	jVO := VO.NewJugadorVO("", "", "", "", "", "", email, "")

	jDAO.GetJugador(jVO)

	struct Usuario{
		Nombre string
		Foto int
		Desc string
		PGanadas int
		PJugadas int
		Codigo string
		//Puntos int
	}

	user := Usuario{
		Nombre: jVO[i].GetCodigo(),
		Foto:   jVO[i].GetFoto(),
		Descp:  jVO[i].GetDescrip(),
		PGanadas : jVO[i].GetPGanadas(),
		PJugadas:   jVO[i].GetPJugadas(),
		Codigo:  jVO[i].GetCodigo(),
		//Puntos : jVO[i].Get(),
	}

	c.JSON(http.StatusOK, gin.H{
		"usuario": user,
	})

}
