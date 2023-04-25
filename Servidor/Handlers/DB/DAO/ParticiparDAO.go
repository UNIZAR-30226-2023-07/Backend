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
	addp := "INSERT INTO PARTICIPAR (partida, jugador, puntos_resultado, enlobby) VALUES ($1, $2, 0, 1, $3)"
	_, e := db.Exec(addp, pVO.GetPartida(), pVO.GetJugador(), pVO.GetTurno())
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
	modL := "UPDATE PARTICIPAR SET enlobby = $3 WHERE partida = $2 AND jugador = $1"
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
	modL := "UPDATE PARTICIPAR SET enlobby = $2 WHERE partida = $1"
	_, e := db.Exec(modL, p, l)
	CheckError(e)

}

// Devuelve true si esa partida es un torneo
func (pDAO *ParticiparDAO) EstaParticipando(p string, j string) bool {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	res := false

	//Buscamos si el jugador esta participando
	isp := "SELECT * FROM PARTICIPAR WHERE partida = $1 AND jugador = $2"
	rows, err := db.Query(isp, p, j)
	CheckError(err)

	defer rows.Close()
	if rows.Next() {
		res = true
	}

	return res

}
