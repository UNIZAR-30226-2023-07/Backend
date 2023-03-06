package DAO

import (
	"DB/VO"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = "5432"
	user     = "frances"
	password = "1234"
	dbname   = "Pro_Soft"
)

/*-------------------------------------------------------------------------------------------------------------------*/

type JugadoresDAO struct{}

// Añade el jugador jVO a la base de datos
func (jDAO *JugadoresDAO) AddJugador(jVO VO.JugadoresVO) {

	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Añadir jugador j
	addj := "INSERT INTO JUGADORES VALUES ($1, $2, $3, 0, 0 ,$4)"
	_, e := db.Exec(addj, jVO.GetNombre(), jVO.GetContra(), jVO.GetDescrip(), jVO.GetCodigo())
	CheckError(e)

}

// Borra un jugador de la base de datos con codigo j
func (jDAO *JugadoresDAO) DelJugador(j string) {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Borrar jugador
	delj := "DELETE FROM JUGADORES WHERE codigo = $1"
	_, e := db.Exec(delj, j)
	CheckError(e)

}

// Modifica la información de un jugador, excepto la foto y los puntos
func (jDAO *JugadoresDAO) ModJugador(jVO VO.JugadoresVO) {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Modifica nombre, descripción, y contraseña
	modj := "UPDATE JUGADORES SET(nombre, contra, descrp) = ($1, $2, $3) WHERE codigo = $4"
	_, e := db.Exec(modj, jVO.GetNombre(), jVO.GetContra(), jVO.GetDescrip(), jVO.GetCodigo())
	CheckError(e)

}

// Devuelve todos los amigos confirmados de un jugador j
func (jDAO *JugadoresDAO) ListarAmigos(j string) []*VO.JugadoresVO {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Obtenemos todos los parametros para cada amigo del jugador
	qAmis := "SELECT r.nombre, r.descrp, r.pjugadas, r.pganadas, r.codigo " +
		"FROM AMISTAD AS a JOIN JUGADORES AS r ON a.usr2 = r.codigo " +
		"WHERE a.usr1 = $1 AND a.estado = 'confirmada'"
	rows, err := db.Query(qAmis, j)
	CheckError(err)

	var res []*VO.JugadoresVO

	defer rows.Close()
	for rows.Next() {
		var nombre string
		var descripcion string
		var pjugadas int
		var pganadas int
		var codigo string
		var contra string
		var foto []byte

		err := rows.Scan(&nombre, &descripcion, &pjugadas, &pganadas, &codigo)
		CheckError(err)

		j := VO.NewJugadorVO(nombre, contra, foto, descripcion, pjugadas, pganadas, codigo)
		res = append(res, j)

	}

	return res

}

/*-------------------------------------------------------------------------------------------------------------------*/

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

/*-------------------------------------------------------------------------------------------------------------------*/

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

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
