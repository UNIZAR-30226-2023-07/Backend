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
	addp := "INSERT INTO PARTIDAS VALUES ($1, $2, $3, 'creando')"
	_, e := db.Exec(addp, pVO.GetClave(), pVO.GetCreador(), pVO.GetTipo())
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

// Asigna como pausada aquella partida cuya clave sea la misma que la pVO
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

// Asigna como terminada aquella partida cuya clave sea la misma que la pVO
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

// Asigna como iniciada aquella partida cuya clave sea la misma que la pVO
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

// Devuelve true si la partida esta llena
func (pDAO *PartidasDAO) EstaLlena(clave string) bool {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	res := false

	//Buscamos si hay mas de cinco jugadores
	fullp := "SELECT COUNT(*) FROM PARTICIPAR WHERE partida = $1 HAVING COUNT(*) > 5"
	rows, err := db.Query(fullp, clave)
	CheckError(err)

	defer rows.Close()
	if rows.Next() {
		res = true
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
	esp := "SELECT * FROM PARTIDAS WHERE clave = $1 AND estado = pausada"
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
	esc := "SELECT COUNT(*) FROM PARTIDAS WHERE clave = $1 AND creador = $2"
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

// Añade una carta del mazo de un jugador en la BD
func (pDAO *PartidasDAO) AddCartaMazo(m VO.MazosVO) {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Añadimos una nueva carta al mazo de un jugador en la BD
	addm := "INSERT INTO MAZOS VALUES (DEFAULT, $1, $2, $3)"
	_, err = db.Exec(addm, m.GetPartida(), m.GetCarta(), m.GetJugador())
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

// Recupera el mazo de un jugador
func (pDAO *PartidasDAO) GetMazo(p string, j string) []*VO.CartasVO {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Obtenemos todas las combinaciones de una partida
	qMaz := "SELECT carta FROM MAZOS WHERE partida = $1 and jugador = $2"
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

// Elimina las combinaciones, descartes y mazos de una partida
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

	//Eliminamos mazos de una partida
	delm := "DELETE FROM MAZOS WHERE partida = $1"
	_, err = db.Exec(delm, p)
	CheckError(err)

	//Eliminamos combinaciones de una partida
	deld := "DELETE FROM DESCARTES WHERE partida = $1"
	_, err = db.Exec(deld, p)
	CheckError(err)
}
