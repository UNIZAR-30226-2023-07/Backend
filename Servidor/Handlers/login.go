package Handlers

import (
	"DB/DAO"
	"DB/VO"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Login struct {
	Email  string `json:"email" binding:"required"`
	Contra string `json:"contra" binding:"required"`
}

func PostLogin(c *gin.Context) {

	u := Login{}
	//Con el binding guardamos el json de la petición en u que es de tipo login
	if err := c.BindJSON(&u); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	jDAO := DAO.JugadoresDAO{}
	jVO := VO.NewJugadorVO("", u.Contra, "", 0, 0, u.Email, "")
	if jDAO.ValJugador(jVO) {
		c.JSON(http.StatusAccepted, gin.H{
			"res": "ok",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"res": "email o contraseña no validos",
		})
	}

}
