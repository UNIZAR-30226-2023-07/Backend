package DAO

import (
	"DB/VO"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "192.168.56.5"
	port     = "5432"
	user     = "frances"
	password = "1234"
	dbname   = "Pro_Soft"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

/*-------------------------------------------------------------------------------------------------------------------*/

type JugadoresDAO struct{}

// Añade el jugador jVO a la base de datos
func (jDAO *JugadoresDAO) AddJugador(jVO VO.JugadoresVO) bool {

	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Añadir jugador j
	addj := "INSERT INTO JUGADORES VALUES ($1, $2, $3, 0, 0, $4, $5)"
	_, e := db.Exec(addj, jVO.GetNombre(), jVO.GetContra(), jVO.GetDescrip(), jVO.GetEmail(), jVO.GetCodigo())

	if e == nil {
		return true
	} else {
		return false
	}

}

// Borra un jugador de la base de datos con codigo j
func (jDAO *JugadoresDAO) DelJugador(j string) bool {
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
	if e == nil {
		return true
	} else {
		return false
	}

}

// Devuelve true si el jugador existe y devuelve su información
func (JDAO *JugadoresDAO) GetJugador(JVO *VO.JugadoresVO) bool {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	res := false
	//Busca el usuario cuyo codigo coincide con el de JVO
	getj := "SELECT nombre, descrp, pjugadas, pganadas, email  FROM JUGADORES WHERE codigo = $1"
	rows, err := db.Query(getj, JVO.GetCodigo())

	if rows.Next() == true {
		res = true
		var nombre string
		var descripcion string
		var pjugadas int
		var pganadas int
		var codigo string

		err := rows.Scan(&nombre, &descripcion, &pjugadas, &pganadas, &codigo)
		CheckError(err)

		j := VO.NewJugadorVO(nombre, "", descripcion, pjugadas, pganadas, JVO.GetEmail(), codigo)
		JVO = j

	} else {
		res = false
	}

	return res
}

// Modifica contraseña, descripción, nombre de un jugador
func (jDAO *JugadoresDAO) ModJugador(jVO VO.JugadoresVO) {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Modifica descripción, y contraseña
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
	qAmis := "SELECT r.nombre, r.descp, r.email, r.codigo " +
		"FROM AMISTAD AS a JOIN JUGADORES AS r ON a.usr2 = r.codigo " +
		"WHERE a.usr1 = $1 AND a.estado = 'confirmada'"
	rows, err := db.Query(qAmis, j)
	CheckError(err)

	var res []*VO.JugadoresVO

	defer rows.Close()
	for rows.Next() {
		var nombre string
		var descripcion string
		var email string
		var codigo string

		err := rows.Scan(&nombre, &descripcion, &email, &codigo)
		CheckError(err)

		j := VO.NewJugadorVO(nombre, "", descripcion, 0, 0, email, codigo)
		res = append(res, j)

	}

	return res

}

// Devuelve true si el jugador existe y devuelve su información
func (JDAO *JugadoresDAO) ValJugador(JVO *VO.JugadoresVO) bool {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	res := false
	//Busca el usuario cuyo codigo coincide con el de JVO
	getj := "SELECT * FROM JUGADORES WHERE email = $1 AND contra = $2"
	rows, err := db.Query(getj, JVO.GetEmail(), JVO.GetContra())

	if rows.Next() == true {
		res = true
	} else {
		res = false
	}

	return res
}
