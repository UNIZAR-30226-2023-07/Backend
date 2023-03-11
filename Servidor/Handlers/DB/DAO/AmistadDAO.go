package DAO

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type AmistadDAO struct{}

// Añade una nueva relación de amistad y la establece como esperando confirmación para el que la manda y pendiente para el que la recibe
func (aDAO *AmistadDAO) PeticionAmistad(j1 string, j2 string) {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Añadir petición
	adda1 := "INSERT INTO AMISTAD VALUES ('esp_confirmacion', $1, $2)"
	_, e1 := db.Exec(adda1, j1, j2)
	CheckError(e1)

	adda2 := "INSERT INTO AMISTAD VALUES ('pendiente', $1, $2)"
	_, e2 := db.Exec(adda2, j2, j1)
	CheckError(e2)
}

// Confirma una relación de amistad pendiente
func (aDAO *AmistadDAO) AceptarPeticion(j1 string, j2 string) {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Confirmar ambas peticiones
	adda := "UPDATE AMISTAD SET estado = 'confirmada' WHERE usr1 = $1 AND usr2 = $2"
	_, e1 := db.Exec(adda, j1, j2)
	CheckError(e1)

	_, e2 := db.Exec(adda, j2, j1)
	CheckError(e2)
}

// Rechaza una relación de amistad pendiente
func (aDAO *AmistadDAO) RechazarPeticion(j1 string, j2 string) {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Confirmar ambas peticiones
	adda := "DELETE FROM AMISTAD WHERE usr1 = $1 AND usr2 = $2"
	_, e1 := db.Exec(adda, j1, j2)
	CheckError(e1)

	_, e2 := db.Exec(adda, j2, j1)
	CheckError(e2)
}
