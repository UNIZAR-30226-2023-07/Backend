package DAO

import (
	"DB/VO"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type CartasDAO struct{}

// Las cartas ya estarán añadidas en la base de datos
// Recupera las cartas de la base de datos junto con su id
func (cDAO *CartasDAO) GetCartas() map[int]*VO.CartasVO {

	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Recuperamos todas las cartas de la BD
	getc := "SELECT id, valor, palo, reverso FROM CARTAS"
	rows, err := db.Query(getc)
	CheckError(err)

	var res = make(map[int]*VO.CartasVO)

	defer rows.Close()
	for rows.Next() {

		var id int
		var valor int
		var palo int
		var reverso int

		err := rows.Scan(&id, &valor, &palo, &reverso)
		CheckError(err)

		res[id] = VO.NewCartasVO(valor, palo, reverso)

	}

	return res

}
