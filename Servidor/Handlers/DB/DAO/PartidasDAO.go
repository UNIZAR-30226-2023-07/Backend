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
	fullp := "COUNT(*) FROM PARTICIPAR WHERE partida = $1 HAVING COUNT(*) > 5"
	rows, err := db.Query(fullp, clave)
	CheckError(err)

	defer rows.Close()
	if rows.Next() {
		res = true
	}

	return res

}
