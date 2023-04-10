package Handlers

import (
	"DB/DAO"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Todas las peticiones de post de amistad usaran este struct
type Amistad struct {
	Emisor   string `json:"emisor" binding:"required"`
	Receptor string `json:"receptor" binding:"required"`
}

func PostAmistadRm(c *gin.Context) {

	a := Amistad{}

	if err := c.BindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"res": "error",
		})
		return
	}

	aDAO := DAO.AmistadDAO{}

	aDAO.RechazarPeticion(a.Emisor, a.Receptor)

	c.JSON(http.StatusAccepted, gin.H{
		"res": "ok",
	})

}

func GetAmistadList(c *gin.Context) {

	jug := c.Param("code") //Sacamos el parametro que llega en :code

	jDAO := DAO.JugadoresDAO{}

	friends := jDAO.ListarAmigos(jug)

	//Solo necesitamos estos parametros de los amigos

	type Amigo struct {
		Nombre string
		Foto   int
		Descp  string
	}

	var amiguis []Amigo

	for i := 0; i < len(friends); i++ {
		a := Amigo{
			Nombre: friends[i].GetNombre(),
			Foto:   friends[i].GetFoto(),
			Descp:  friends[i].GetDescrip(),
		}

		amiguis = append(amiguis, a)
	}

	c.JSON(http.StatusOK, gin.H{
		"amistad": amiguis,
	})

}

func PostAmistadAdd(c *gin.Context) {

	a := Amistad{}

	if err := c.BindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"res": "error",
		})
		return
	}

	aDAO := DAO.AmistadDAO{}

	aDAO.PeticionAmistad(a.Emisor, a.Receptor)

	c.JSON(http.StatusAccepted, gin.H{
		"res": "ok",
	})

}

func PostAmistadAccept(c *gin.Context) {

	a := Amistad{}

	if err := c.BindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"res": "error",
		})
		return
	}

	aDAO := DAO.AmistadDAO{}

	aDAO.AceptarPeticion(a.Emisor, a.Receptor)

	c.JSON(http.StatusAccepted, gin.H{
		"res": "ok",
	})

}

func PostAmistadDeny(c *gin.Context) {

	a := Amistad{}

	if err := c.BindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"res": "error",
		})
		return
	}

	aDAO := DAO.AmistadDAO{}

	aDAO.RechazarPeticion(a.Emisor, a.Receptor)

	c.JSON(http.StatusAccepted, gin.H{
		"res": "ok",
	})

}

func GetPendientesList(c *gin.Context) {

	jug := c.Param("code") //Sacamos el parametro que llega en :code

	jDAO := DAO.JugadoresDAO{}

	friends, state := jDAO.ListarPendientes(jug)

	//Solo necesitamos estos parametros de los amigos

	type Amigo struct {
		Nombre string
		Foto   int
		Descp  string
		Estado string
	}

	var amiguis []Amigo

	for i := 0; i < len(friends); i++ {
		a := Amigo{
			Nombre: friends[i].GetNombre(),
			Foto:   friends[i].GetFoto(),
			Descp:  friends[i].GetDescrip(),
			Estado: state[i].GetEstado(),
		}

		amiguis = append(amiguis, a)
	}

	c.JSON(http.StatusOK, gin.H{
		"amistad": amiguis,
	})

}
