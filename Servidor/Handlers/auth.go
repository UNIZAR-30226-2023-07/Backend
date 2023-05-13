package Handlers

import (
	"DB/DAO"
	"DB/VO"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Login struct {
	Email  string `json:"email" binding:"required"`
	Contra string `json:"contra" binding:"required"`
}

type Register struct {
	Nombre string `json:"nombre" binding:"required"`
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
	jVO := VO.NewJugadorVO("", u.Contra, 0, "", 0, 0, 0, u.Email, "")
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

func PostRegister(c *gin.Context) {

	u := Register{}
	//Con el binding guardamos el json de la petición en u que es de tipo login
	if err := c.BindJSON(&u); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	jDAO := DAO.JugadoresDAO{}

	id := rand.Intn(1000)
	codigo := u.Nombre + "_" + strconv.Itoa(id)

	for jDAO.EstaJugador(codigo) {
		id = rand.Intn(1000)
		codigo = u.Nombre + "_" + strconv.Itoa(id)
	}

	jVO := VO.NewJugadorVO(u.Nombre, u.Contra, 0, "", 0, 0, 0, u.Email, codigo)

	if jDAO.AddJugador(*jVO) {
		c.JSON(http.StatusAccepted, gin.H{
			"res":  "ok",
			"code": codigo,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"res": "email no valido",
		})
	}
}

func PostModLogin(c *gin.Context) {

	log := Login{}
	//Con el binding guardamos el json de la petición en log que es de tipo login
	if err := c.BindJSON(&log); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Cambia la contraseña del jugador dado el email
	jDAO := DAO.JugadoresDAO{}
	jDAO.CambiarContra(log.Email, log.Contra)

	c.JSON(http.StatusAccepted, gin.H{
		"res": "ok",
	})
}
