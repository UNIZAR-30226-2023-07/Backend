package DAO

import (
	"DB/VO"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type CartasDAO struct{}

// Añade una carta a la base de datos
func (cDAO *CartasDAO) AddCarta(cVO VO.CartasVO) {

	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Añadir carta
	addc := "INSERT INTO CARTAS VALUES ($1, $2)"
	_, e := db.Exec(addc, cVO.GetNumero(), cVO.GetPalo())
	CheckError(e)

}

// Elimina una carta de la base de datos
func (cDAO *CartasDAO) DelCarta(cVO VO.CartasVO) {

	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Borrar carta
	delc := "DELETE FROM CARTAS WHERE numero = $1 AND palo = $2"
	_, e := db.Exec(delc, cVO.GetNumero(), cVO.GetPalo())
	CheckError(e)

}
