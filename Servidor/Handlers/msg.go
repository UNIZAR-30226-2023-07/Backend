package Handlers

import (
	"DB/DAO"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMsgList(c *gin.Context) {

	m := c.Param("code")

	mDAO := DAO.MensajesDAO{}

	vmsg := mDAO.ObtenerMensajes(m)

	c.JSON(http.StatusAccepted, gin.H{
		"msg": vmsg,
	})

}

func PostLeer(c *gin.Context) {

	a := Amistad{}

	if err := c.BindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"res": "error",
		})
		return
	}

	mDAO := DAO.MensajesDAO{}

	mDAO.LeerMensajes(a.Receptor, a.Emisor)

	c.JSON(http.StatusAccepted, gin.H{
		"res": "ok",
	})

}
