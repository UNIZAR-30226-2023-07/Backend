package DAO

import (
	"DB/VO"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PartidasDAO struct{}

// Crea una nueva partida en la BD
func (pDAO *PartidasDAO) AddPartida(pVO VO.PartidasVO) {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Añadir partida
	addp := "INSERT INTO PARTIDAS VALUES ($1, $2, $3, 'creando', $4)"
	_, e := db.Exec(addp, pVO.GetClave(), pVO.GetCreador(), pVO.GetTipo(), pVO.GetPactual())
	CheckError(e)

}

// Devuelve true si ya existe una partida con esa clave
func (pDAO *PartidasDAO) HayPartida(clave string) bool {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	res := false

	//Buscamos si existe ya alguna partida con esa clave
	isp := "SELECT * FROM PARTIDAS WHERE clave = $1"
	rows, err := db.Query(isp, clave)
	CheckError(err)

	defer rows.Close()
	if rows.Next() {
		res = true
	}

	return res

}

// Devuelve true si esa partida es un torneo
func (pDAO *PartidasDAO) EsTorneo(clave string) bool {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	res := false

	//Buscamos si existe ya alguna partida con esa clave y es torneo
	isp := "SELECT * FROM PARTIDAS WHERE clave = $1 AND tipo = 'torneo'"
	rows, err := db.Query(isp, clave)
	CheckError(err)

	defer rows.Close()
	if rows.Next() {
		res = true
	}

	return res

}

// Devuelve la partida actual de un torneo
func (pDAO *PartidasDAO) PartidaActualTorneo(torneo string) string {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	pcatual := ""

	//Buscamos si hay alguna sala con gente libre en el torneo
	isp := "SELECT pactual FROM PARTIDAS WHERE clave = $1"
	rows, err := db.Query(isp, torneo)
	CheckError(err)

	defer rows.Close()
	if rows.Next() {

		err := rows.Scan(&pcatual)
		CheckError(err)
	}

	return pcatual

}

// Devuelve la partida con clave clave
func (pDAO *PartidasDAO) GetPartida(clave string) *VO.PartidasVO {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Buscamos el creador del torneo
	isp := "SELECT creador, tipo, estado, pactual FROM PARTIDAS WHERE clave = $1"
	rows, err := db.Query(isp, clave)
	CheckError(err)

	var partida *VO.PartidasVO
	defer rows.Close()
	if rows.Next() {

		var creador string
		var tipo string
		var estado string
		var pactual string

		err := rows.Scan(&creador, &tipo, &estado, &pactual)

		partida = VO.NewPartidasVO(clave, creador, tipo, estado, pactual)
		CheckError(err)

	}

	return partida

}

// Asigna como pausada aquella partida que tenga clave como clave
func (pDAO *PartidasDAO) PausarPartida(clave string) {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Marcamos la partida como pausada
	pausep := "UPDATE PARTIDAS SET estado = 'pausada' WHERE clave = $1"
	_, e := db.Exec(pausep, clave)
	CheckError(e)

}

// Asigna como terminada aquella partida que tenga clave como clave
func (pDAO *PartidasDAO) TerminarPartida(clave string) {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Marcamos la partida como terminada
	terminp := "UPDATE PARTIDAS SET estado = 'terminada' WHERE clave = $1"
	_, e := db.Exec(terminp, clave)
	CheckError(e)

}

// Asigna como iniciada aquella partida que tenga clave como clave
func (pDAO *PartidasDAO) IniciarPartida(clave string) {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Marcamos la partida como iniciada
	initp := "UPDATE PARTIDAS SET estado = 'iniciada' WHERE clave = $1"
	_, e := db.Exec(initp, clave)
	CheckError(e)
}

// Devuelve el numero de jugadores que estan en la partida
func (pDAO *PartidasDAO) NJugadoresPartida(clave string) int {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	res := 0

	//Numero de jugadores
	fullp := "SELECT COUNT(*) FROM PARTICIPAR WHERE partida = $1"
	rows, err := db.Query(fullp, clave)
	CheckError(err)

	defer rows.Close()
	if rows.Next() {
		err := rows.Scan(&res)
		CheckError(err)
	}

	return res

}

// Devulve true si la partida está pausada
func (pDAO *PartidasDAO) EstaPausada(p string) bool {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	res := false

	//Comprobamos si j es el creador de la partida
	esp := "SELECT * FROM PARTIDAS WHERE clave = $1 AND estado = 'pausada'"
	rows, err := db.Query(esp, p)
	CheckError(err)

	defer rows.Close()
	if rows.Next() {
		res = true
	}

	return res

}

// Devuelve true si j es el creador de la partida p
func (pDAO *PartidasDAO) EsCreador(p string, j string) bool {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	res := false

	//Comprobamos si j es el creador de la partida
	esc := "SELECT * FROM PARTIDAS WHERE clave = $1 AND creador = $2"
	rows, err := db.Query(esc, p, j)
	CheckError(err)

	defer rows.Close()
	if rows.Next() {
		res = true
	}

	return res

}

// Devuelve true si todos los jugadores están en el lobby
func (pDAO *PartidasDAO) JugadoresEnLobby(clave string) bool {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Contamos los jugadores de una partida
	countj := "SELECT COUNT(*) FROM PARTICIPAR WHERE partida = $1"
	rows, err := db.Query(countj, clave)
	CheckError(err)

	jug := -1

	defer rows.Close()
	if rows.Next() {
		err := rows.Scan(&jug)
		CheckError(err)
	}

	//Contamos los jugadores en lobby de una partida
	countl := "SELECT COUNT(*) FROM PARTICIPAR WHERE partida = $1 AND enlobby = 1"
	rows, err = db.Query(countl, clave)
	CheckError(err)

	lobby := -1
	defer rows.Close()
	if rows.Next() {
		err := rows.Scan(&lobby)
		CheckError(err)
	}

	return (jug != -1 && lobby != -1 && jug == lobby)

}

// Añade una combinación en la BD
func (pDAO *PartidasDAO) AddCombinacion(c VO.CombinacionesVO) {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Añadimos una nueva combinación a la BD
	addc := "INSERT INTO COMBINACIONES VALUES (DEFAULT, $1, $2, $3)"
	_, err = db.Exec(addc, c.GetPartida(), c.GetCarta(), c.GetNcomb())
	CheckError(err)
}

// Añade una combinación en la BD
func (pDAO *PartidasDAO) AddCartaMazo(m VO.MazosVO) {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Añadimos una nueva combinación a la BD
	addc := "INSERT INTO MAZOS VALUES (DEFAULT, $1, $2)"
	_, err = db.Exec(addc, m.GetPartida(), m.GetCarta())
	CheckError(err)
}

// Añade una carta a la mano de un jugador en la BD
func (pDAO *PartidasDAO) AddCartaMano(m VO.ManosVO) {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Añadimos una nueva carta a la mano de un jugador en la BD
	addm := "INSERT INTO MANOS VALUES (DEFAULT, $1, $2, $3)"
	_, err = db.Exec(addm, m.GetPartida(), m.GetCarta(), m.GetTurno())
	CheckError(err)
}

// Añade un descarte en la BD
func (pDAO *PartidasDAO) AddDescarte(d VO.DescartesVO) {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Añadimos un nuevo descarte a la BD
	addd := "INSERT INTO DESCARTES VALUES (DEFAULT, $1, $2)"
	_, err = db.Exec(addd, d.GetPartida(), d.GetCarta())
	CheckError(err)
}

// Recupera todas las combinaciones de una partida
func (pDAO *PartidasDAO) GetCombinaciones(p string) []*VO.CombinacionesVO {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Obtenemos todas las combinaciones de una partida
	qComb := "SELECT carta, ncomb " +
		"FROM COMBINACIONES " +
		"WHERE partida = $1 " +
		"ORDER BY ncomb, carta ASC"
	rows, err := db.Query(qComb, p)
	CheckError(err)

	var res []*VO.CombinacionesVO

	defer rows.Close()
	for rows.Next() {
		var carta int
		var ncomb int

		err := rows.Scan(&carta, &ncomb)
		CheckError(err)

		c := VO.NewCombinacionesVO(p, carta, ncomb)
		res = append(res, c)

	}

	return res

}

// Recupera el Mazo de la partida
func (pDAO *PartidasDAO) GetMazo(p string) []*VO.CartasVO {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Obtenemos todas las combinaciones de una partida
	qMaz := "SELECT carta FROM MAZOS WHERE partida = $1"
	rows, err := db.Query(qMaz, p)
	CheckError(err)

	var res []*VO.CartasVO

	defer rows.Close()
	for rows.Next() {
		var carta int

		err := rows.Scan(&carta)
		CheckError(err)

		c := VO.NewCartasVO((carta/10)/10, (carta/10)%10, carta%10)
		res = append(res, c)

	}

	return res

}

// Recupera la Mano de un jugador
func (pDAO *PartidasDAO) GetMano(p string, j string) []*VO.CartasVO {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Obtenemos todas las combinaciones de una partida
	qMaz := "SELECT m.carta " +
		"FROM MANOS AS m JOIN PARTICIPAR AS p ON m.turno = p.turno AND m.partida = p.partida " +
		"WHERE m.partida = $1 and p.jugador = $2"
	rows, err := db.Query(qMaz, p, j)
	CheckError(err)

	var res []*VO.CartasVO

	defer rows.Close()
	for rows.Next() {
		var carta int

		err := rows.Scan(&carta)
		CheckError(err)

		c := VO.NewCartasVO((carta/10)/10, (carta/10)%10, carta%10)
		res = append(res, c)

	}

	return res

}

// Recupera el descarte de una partida
func (pDAO *PartidasDAO) GetDescarte(p string) *VO.CartasVO {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Obtenemos todas las combinaciones de una partida
	qDes := "SELECT carta FROM DESCARTES WHERE partida = $1"
	rows, err := db.Query(qDes, p)
	CheckError(err)

	var res *VO.CartasVO

	defer rows.Close()
	if rows.Next() {
		var carta int

		err := rows.Scan(&carta)
		CheckError(err)

		res = VO.NewCartasVO((carta/10)/10, (carta/10)%10, carta%10)
	}

	return res

}

// Elimina las combinaciones, descartes, manos y mazo de una partida
func (pDAO *PartidasDAO) DelTableroGuardado(p string) {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Eliminamos combinaciones de una partida
	delc := "DELETE FROM COMBINACIONES WHERE partida = $1"
	_, err = db.Exec(delc, p)
	CheckError(err)

	//Eliminamos las manos de una partida
	delm := "DELETE FROM MANOS WHERE partida = $1"
	_, err = db.Exec(delm, p)
	CheckError(err)

	//Eliminamos combinaciones de una partida
	deld := "DELETE FROM DESCARTES WHERE partida = $1"
	_, err = db.Exec(deld, p)
	CheckError(err)

	//Eliminamos el mazo de una partida
	delmz := "DELETE FROM MAZOS WHERE partida = $1"
	_, err = db.Exec(delmz, p)
	CheckError(err)
}
