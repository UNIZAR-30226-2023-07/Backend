package main

import (
	"DB/DAO"
	"DB/VO"
	"fmt"
)

func main() {
	jDAO := DAO.JugadoresDAO{}
	aDAO := DAO.AmistadDAO{}
	cDAO := DAO.CartasDAO{}

	j1VO := VO.NewJugadorVO("Adrian1", "1234", make([]byte, 1), "Hola esto es una prueba1", 0, 0, "#abc")
	j3VO := VO.NewJugadorVO("Adrian3", "1234", make([]byte, 1), "Hola esto es una prueba3", 0, 0, "#efg")
	j2VO := VO.NewJugadorVO("Adrian2", "1234", make([]byte, 1), "Hola esto es una prueba2", 0, 0, "#hij")
	j4VO := VO.NewJugadorVO("Adrian4", "1234", make([]byte, 1), "Hola esto es una prueba4", 0, 0, "#klm")
	c1VO := VO.NewCartasVO(1, "espadas", make([]byte, 1))
	c2VO := VO.NewCartasVO(1, "oros", make([]byte, 1))
	c3VO := VO.NewCartasVO(1, "bastos", make([]byte, 1))

	jDAO.AddJugador(*j3VO)            //Añadimos un jugador
	jDAO.DelJugador(j3VO.GetCodigo()) //Borramos un jugdor

	jDAO.AddJugador(*j1VO) //Añadimos un jugador
	jDAO.AddJugador(*j2VO) //Añadimos un jugador
	jDAO.AddJugador(*j4VO) //Añadimos un jugador

	//Vamos a hacer que j1VO le mande una solicitud a j2VO y este la acepte
	aDAO.PeticionAmistad(j1VO.GetCodigo(), j2VO.GetCodigo())
	aDAO.AceptarPeticion(j2VO.GetCodigo(), j1VO.GetCodigo())

	//Ahora vamos a hacer que j2VO le mande una solicitud a j4Vo y este la rechace
	aDAO.PeticionAmistad(j2VO.GetCodigo(), j4VO.GetCodigo())
	aDAO.RechazarPeticion(j4VO.GetCodigo(), j2VO.GetCodigo())

	j22VO := VO.NewJugadorVO("ABCD", "1234", make([]byte, 1), "Hola esto es una prueba5", 0, 0, j2VO.GetCodigo())
	jDAO.ModJugador(*j22VO)

	//Ahora mostramos los amigos de J2V0
	amiguis := jDAO.ListarAmigos(j2VO.GetCodigo())
	for i := 0; i < len(amiguis); i++ {
		fmt.Printf(amiguis[i].GetCodigo()) //Solo deberia mostrar Adrian1
	}

	cDAO.AddCarta(*c1VO) //Añadir Carta
	cDAO.DelCarta(*c1VO) //Borrar Carta
	cDAO.AddCarta(*c2VO) //Añadir Carta
	cDAO.AddCarta(*c3VO) //Añadir Carta

}
