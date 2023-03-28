package DAO

import (
	"DB/VO"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type MensajesDAO struct{}

// Añade un mensaje en la base de datos
func (mDAO *MensajesDAO) AddMensaje(mVO VO.MensajesVO) {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Añadir partida
	addm := "INSERT INTO MENSAJES (jug_emi, jug_rcp, timestamp, contenido, leido) VALUES ($1, $2, $3, $4, $5) "
	_, e := db.Exec(addm, mVO.GetEmisor(), mVO.GetReceptor(), time.Now(), mVO.GetContenido(), mVO.GetLeido())
	CheckError(e)

}

// Pone todos los mensajes recibidos por el usuario como recibidos
func (mDAO *MensajesDAO) LeerMensajes(usr string) {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Ponemos a leido todos los mensajes anteriores a ahora
	rdm := "UPDATE MENSAJES SET leido = 1 WHERE jug_rcp = $1 AND timestamp < $2 "
	_, e := db.Exec(rdm, usr, time.Now())
	CheckError(e)
}

// Devuelve todos los mensajes en los que participa el usuario
func (mDAO *MensajesDAO) ObtenerMensajes(usr string) []*VO.MensajesVO {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Obtenemos todos los mensajes en los que participa el usuario ordenados por timestamp
	qMsg := "SELECT jug_emi, jug_rcp, contenido " +
		"FROM MENSAJES " +
		"WHERE jug_emi = $1 OR jug_rcp = $1 "
	rows, err := db.Query(qMsg, usr)
	CheckError(err)

	var res []*VO.MensajesVO

	defer rows.Close()
	for rows.Next() {
		var jug_emi string
		var jug_rcp string
		var contenido string

		err := rows.Scan(&jug_emi, &jug_rcp, &contenido)
		CheckError(err)

		m := VO.NewMensajesVO(jug_emi, jug_rcp, contenido, 1)
		res = append(res, m)

	}

	return res
}

// Devuelve el numero de mensajes no leidos por el usuario
func (mDAO *MensajesDAO) MensajesSinLerr(usr string) int {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Obtenemos todos los mensajes en los que participa el usuario ordenados por timestamp
	qMsg := "COUNT(*) " +
		"FROM MENSAJES " +
		"WHERE jug_rcp = $1 AND leido = 0 "

	rows, err := db.Query(qMsg, usr)
	CheckError(err)

	defer rows.Close()

	rows.Next()
	var noleidos int

	e := rows.Scan(&noleidos)
	CheckError(e)

	return noleidos
}
