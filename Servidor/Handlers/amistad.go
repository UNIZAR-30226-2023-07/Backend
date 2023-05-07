package Handlers

import (
	"DB/DAO"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

// Todas las peticiones de post de amistad usaran este struct
type Amistad struct {
	Emisor   string `json:"emisor" binding:"required"`
	Receptor string `json:"receptor" binding:"required"`
}

func PostAmistadRm(c *gin.Context, chat *melody.Melody) {
	type M_rcp struct {
		Emisor    string `json:"emisor"`
		Receptor  string `json:"receptor"`
		Contenido string `json:"contenido"`
	}

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

	var M M_rcp
	M.Emisor = "Servidor"
	M.Receptor = a.Receptor
	M.Contenido = "Remove"
	msg, _ := json.MarshalIndent(&M, "", "\t")

	//Retransmitir el mensaje al receptor
	chat.BroadcastFilter(msg, func(q *melody.Session) bool {
		return q.Request.URL.Path == ("/api/ws/chat/" + a.Receptor)
	})

}

func GetAmistadList(c *gin.Context) {

	jug := c.Param("code") //Sacamos el parametro que llega en :code

	jDAO := DAO.JugadoresDAO{}

	friends := jDAO.ListarAmigos(jug)

	//Solo necesitamos estos parametros de los amigos

	type Amigo struct {
		Nombre string
		Codigo string
		Foto   int
		Descp  string
		Puntos int
	}

	var amiguis []Amigo

	for i := 0; i < len(friends); i++ {
		a := Amigo{
			Nombre: friends[i].GetNombre(),
			Codigo: friends[i].GetCodigo(),
			Foto:   friends[i].GetFoto(),
			Descp:  friends[i].GetDescrip(),
			Puntos: friends[i].GetPuntos(),
		}

		amiguis = append(amiguis, a)
	}

	c.JSON(http.StatusOK, gin.H{
		"amistad": amiguis,
	})

}

func PostAmistadAdd(c *gin.Context, chat *melody.Melody) {

	type M_rcp struct {
		Emisor    string `json:"emisor"`
		Receptor  string `json:"receptor"`
		Contenido string `json:"contenido"`
	}

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

	var M M_rcp
	M.Emisor = "Servidor"
	M.Receptor = a.Receptor
	M.Contenido = "Add"
	msg, _ := json.MarshalIndent(&M, "", "\t")

	//Retransmitir el mensaje al receptor
	chat.BroadcastFilter(msg, func(q *melody.Session) bool {
		return q.Request.URL.Path == ("/api/ws/chat/" + a.Receptor)
	})

}

func PostAmistadAccept(c *gin.Context, chat *melody.Melody) {
	type M_rcp struct {
		Emisor    string `json:"emisor"`
		Receptor  string `json:"receptor"`
		Contenido string `json:"contenido"`
	}

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

	var M M_rcp
	M.Emisor = "Servidor"
	M.Receptor = a.Receptor
	M.Contenido = "Accept"
	msg, _ := json.MarshalIndent(&M, "", "\t")

	//Retransmitir el mensaje al receptor
	chat.BroadcastFilter(msg, func(q *melody.Session) bool {
		return q.Request.URL.Path == ("/api/ws/chat/" + a.Receptor)
	})

}

func PostAmistadDeny(c *gin.Context, chat *melody.Melody) {
	type M_rcp struct {
		Emisor    string `json:"emisor"`
		Receptor  string `json:"receptor"`
		Contenido string `json:"contenido"`
	}

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

	var M M_rcp
	M.Emisor = "Servidor"
	M.Receptor = a.Receptor
	M.Contenido = "Deny"
	msg, _ := json.MarshalIndent(&M, "", "\t")

	//Retransmitir el mensaje al receptor
	chat.BroadcastFilter(msg, func(q *melody.Session) bool {
		return q.Request.URL.Path == ("/api/ws/chat/" + a.Receptor)
	})

}

func GetPendientesList(c *gin.Context) {

	jug := c.Param("code") //Sacamos el parametro que llega en :code

	jDAO := DAO.JugadoresDAO{}

	friends, state := jDAO.ListarPendientes(jug)

	//Solo necesitamos estos parametros de los amigos

	type Amigo struct {
		Nombre string
		Codigo string
		Foto   int
		Descp  string
		Estado string
	}

	var amiguis []Amigo

	for i := 0; i < len(friends); i++ {
		a := Amigo{
			Nombre: friends[i].GetNombre(),
			Codigo: friends[i].GetCodigo(),
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
