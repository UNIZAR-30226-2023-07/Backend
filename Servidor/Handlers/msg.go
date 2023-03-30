package Handlers

import (
	"DB/DAO"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMsgList(c *gin.Context) {

	m := c.Param("code")

	mDAO := DAO.MensajesDAO{}

	msgs := mDAO.ObtenerMensajes(m)

	type Vmsg struct {
		Emisor    string
		Receptor  string
		Contenido string
		Leido     int
	}

	var vmsg []Vmsg
	for i := 0; i < len(msgs); i++ {
		vmsg = append(vmsg, Vmsg{msgs[i].GetEmisor(), msgs[i].GetReceptor(), msgs[i].GetContenido(), msgs[i].GetLeido()})
	}

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
