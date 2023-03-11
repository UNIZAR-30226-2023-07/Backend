package Handlers

import (
	"DB/DAO"
	"DB/VO"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Register struct {
	Nombre string `json:"nombre" binding:"required"`
	Email  string `json:"email" binding:"required"`
	Contra string `json:"contra" binding:"required"`
}

func PostRegister(c *gin.Context) {

	u := Register{}
	//Con el binding guardamos el json de la petición en u que es de tipo login
	if err := c.BindJSON(&u); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	id := rand.Intn(1000)
	codigo := u.Nombre + "#" + strconv.Itoa(id)

	jDAO := DAO.JugadoresDAO{}
	jVO := VO.NewJugadorVO(u.Nombre, u.Contra, "", 0, 0, u.Email, codigo)
	if jDAO.AddJugador(*jVO) {
		c.JSON(http.StatusAccepted, gin.H{
			"res": "ok",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"res": "email no valido",
		})
	}

}
