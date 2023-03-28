package DAO

import (
	"DB/VO"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type ParticiparDAO struct{}

// Crea una participación en la BD
func (pDAO *ParticiparDAO) AddParticipar(pVO VO.ParticiparVO) {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Añadir participación
	addp := "INSERT INTO PARTICIPAR (partida, jugador, puntos_resultado) VALUES ($1, $2, 0)"
	_, e := db.Exec(addp, pVO.GetPartida(), pVO.GetJugador())
	CheckError(e)

}

// Crea una participación en la BD
func (pDAO *ParticiparDAO) AsignarPuntos(pVO VO.ParticiparVO) {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Actualiazmos los puntos del usuario
	puntos := "UPDATE PARTICIPAR SET puntos_resultado = ($1) WHERE partida = $2 AND jugador = $3"
	_, e := db.Exec(puntos, pVO.GetPuntos(), pVO.GetPartida(), pVO.GetJugador())
	CheckError(e)

}
