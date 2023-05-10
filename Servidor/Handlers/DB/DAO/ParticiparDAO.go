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
	addp := "INSERT INTO PARTICIPAR (partida, jugador, puntos_resultado, enlobby, turno, abierto, bot) VALUES ($1, $2, 0, 1, $3, 'no', $4)"
	_, e := db.Exec(addp, pVO.GetPartida(), pVO.GetJugador(), pVO.GetTurno(), pVO.GetBot())
	CheckError(e)

}

// Asigna enlobby de un jugador o bot para una partida
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

// Asigna enlobby para todos jugadores de una partida que no son bots
func (pDAO *ParticiparDAO) ModLobby(p string, l int) {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Modifica enlobby de todos los jugadores de una partida
	modL := "UPDATE PARTICIPAR SET enlobby = $2 WHERE partida = $1 AND bot = 0"
	_, e := db.Exec(modL, p, l)
	CheckError(e)

}

// Actualiza si un jugador ha abierto
func (pDAO *ParticiparDAO) UpdateAbierto(p string, turno int, a string) {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Modifica si un jugador ha abierto
	moda := "UPDATE PARTICIPAR SET abierto = $3 WHERE partida = $2 AND turno = $1"
	_, e := db.Exec(moda, turno, p, a)
	CheckError(e)

}

// Devuelve true si esta participando en la partida
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

// Devuelve los jugadores junto a su turno
func (pDAO *ParticiparDAO) GetJugadoresTurnos(p string) [][]string {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	var res [][]string

	//Buscamos los jugadores junto a su turno
	isp := "SELECT jugador, turno FROM PARTICIPAR WHERE partida = $1"
	rows, err := db.Query(isp, p)
	CheckError(err)

	defer rows.Close()
	for rows.Next() {
		var jugador string
		var turno string

		err := rows.Scan(&jugador, &turno)
		CheckError(err)

		var tmp []string
		tmp = append(tmp, jugador)
		tmp = append(tmp, turno)

		res = append(res, tmp)
	}

	return res

}

// Devuelve los jugadores que están en el lobby
func (pDAO *ParticiparDAO) GetJugadoresEnLobby(p string) []string {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	var res []string

	//Buscamos los jugadores junto a su turno
	isp := "SELECT jugador FROM PARTICIPAR WHERE partida = $1 AND enlobby = 1"
	rows, err := db.Query(isp, p)
	CheckError(err)

	defer rows.Close()
	for rows.Next() {
		var jugador string

		err := rows.Scan(&jugador)
		CheckError(err)

		res = append(res, jugador)
	}

	return res

}

// Actualiza los puntos del torneo en curso
func (pDAO *ParticiparDAO) UpdatePuntosJug(turno int, partida string, puntos string) {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Incrementa los puntos de un jugador real en la tabla de jugadores
	upp := "UPDATE JUGADORES AS j " +
		"SET puntos = puntos + $3 " +
		"FROM PARTICIPAR AS p " +
		"WHERE p.jugador = j.codigo AND " +
		"p.bot = 0 AND p.partida = $2 " +
		"AND p.turno = $1 "
	_, e := db.Exec(upp, turno, partida, puntos)
	CheckError(e)
}

// Actualiza los puntos del torneo en curso
func (pDAO *ParticiparDAO) UpdatePartidasJug(turno int, partida string, ganador bool) {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Incrementamos una partida jugada
	uppj := "UPDATE JUGADORES AS j " +
		"SET pjugadas = pjugadas + 1 " +
		"FROM PARTICIPAR AS p " +
		"WHERE p.jugador = j.codigo AND " +
		"p.bot = 0 AND p.partida = $2 " +
		"AND p.turno = $1 "
	_, e := db.Exec(uppj, turno, partida)
	CheckError(e)

	//Si ha ganado incrementamos un partida ganada
	if ganador {
		uppg := "UPDATE JUGADORES AS j " +
			"SET pganadas = pganadas + 1 " +
			"FROM PARTICIPAR AS p " +
			"WHERE p.jugador = j.codigo AND " +
			"p.bot = 0 AND p.partida = $2 " +
			"AND p.turno = $1 "
		_, e := db.Exec(uppg, turno, partida)
		CheckError(e)
	}

}

// Actualiza los puntos del torneo en curso
func (pDAO *ParticiparDAO) UpdatePuntos(p string, j string, puntos string) {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Modifica los puntos de un jugador
	modp := "UPDATE PARTICIPAR SET puntos_resultado = $3 WHERE partida = $2 AND jugador = $1"
	_, e := db.Exec(modp, j, p, puntos)
	CheckError(e)

}

// Devuelve los puntos del torneo en curso
func (pDAO *ParticiparDAO) GetPuntos(p string) []string {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	var res []string

	//Devuelve los puntos de cada jugador
	punt := "SELECT turno, puntos_resultado FROM PARTICIPAR WHERE partida = $1 ORDER BY turno ASC"
	rows, e := db.Query(punt, p)
	CheckError(e)

	defer rows.Close()
	for rows.Next() {
		var turno int
		var puntos string

		err := rows.Scan(&turno, &puntos)
		CheckError(err)

		res = append(res, puntos)
	}

	return res

}

// Recupera si un jugados ha abierto o no ordenados por turno
func (pDAO *ParticiparDAO) GetAbierto(p string) []string {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	var res []string

	//Seleccionamos si ha abierto o no ordenador por turno
	a := "SELECT turno, abierto FROM PARTICIPAR WHERE partida = $1 ORDER BY turno ASC"
	rows, err := db.Query(a, p)
	CheckError(err)

	defer rows.Close()
	for rows.Next() {
		var turno int
		var abierto string

		err := rows.Scan(&turno, &abierto)
		CheckError(err)

		res = append(res, abierto)
	}

	return res

}
