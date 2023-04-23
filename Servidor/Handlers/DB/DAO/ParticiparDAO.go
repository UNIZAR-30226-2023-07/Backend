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
	addp := "INSERT INTO PARTICIPAR (partida, jugador, puntos_resultado, enlobby) VALUES ($1, $2, 0, 1)"
	_, e := db.Exec(addp, pVO.GetPartida(), pVO.GetJugador())
	CheckError(e)

}

// Asigna enlobby de un jugador para una partida
func (pDAO *ParticiparDAO) ModLobbyJug(j string, p string, l int) {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Modifica enlobby de un jugador
	modL := "UPDATE PARTICIPAR SET (enlobby) = $3 WHERE partida = $2 AND jugador = $1"
	_, e := db.Exec(modL, j, p, l)
	CheckError(e)

}

// Asigna enlobby para todos jugadores de una partida
func (pDAO *ParticiparDAO) ModLobby(p string, l int) {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Modifica enlobby de todos los jugadores de una partida
	modL := "UPDATE PARTICIPAR SET (enlobby) = $2 WHERE partida = $1"
	_, e := db.Exec(modL, p, l)
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
