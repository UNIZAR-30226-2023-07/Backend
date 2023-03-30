package DAO

import (
	"DB/VO"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "52.166.36.105"
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
	addj := "INSERT INTO JUGADORES VALUES ($1, $2, 0, $3, 0, 0, $4, $5)"
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

// Devuelve true si existe un jugador con ese email y devuelve su información
func (JDAO *JugadoresDAO) GetJugador(JVO *VO.JugadoresVO) bool {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	res := false
	//Busca el usuario cuyo email coincide con el de JVO
	getj := "SELECT nombre, descrp, foto, pjugadas, pganadas, codigo FROM JUGADORES WHERE email = $1"
	rows, err := db.Query(getj, JVO.GetEmail())
	CheckError(err)

	defer rows.Close()
	if rows.Next() {
		res = true
		var nombre string
		var descripcion string
		var foto int
		var pjugadas int
		var pganadas int
		var codigo string

		err := rows.Scan(&nombre, &descripcion, &foto, &pjugadas, &pganadas, &codigo)
		CheckError(err)

		j := VO.NewJugadorVO(nombre, "", foto, descripcion, pjugadas, pganadas, JVO.GetEmail(), codigo)
		JVO = j

	} else {
		res = false
	}

	return res
}

// Modifica descripción, nombre de un jugador y su foto cuyo email coincida con el de jvo
func (jDAO *JugadoresDAO) ModJugador(jVO VO.JugadoresVO) {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Modifica descripción, nombre y foto
	modj := "UPDATE JUGADORES SET (nombre, foto, descrp) = ($1, $2, $3) WHERE email = $4"
	_, e := db.Exec(modj, jVO.GetNombre(), jVO.GetFoto(), jVO.GetDescrip(), jVO.GetEmail())
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
	qAmis := "SELECT r.nombre, r.foto, r.descrp, r.email, r.codigo " +
		"FROM AMISTAD AS a JOIN JUGADORES AS r ON a.usr2 = r.codigo " +
		"WHERE a.usr1 = $1 AND a.estado = 'confirmada'"
	rows, err := db.Query(qAmis, j)
	CheckError(err)

	var res []*VO.JugadoresVO

	defer rows.Close()
	for rows.Next() {
		var nombre string
		var foto int
		var descripcion string
		var email string
		var codigo string

		err := rows.Scan(&nombre, &foto, &descripcion, &email, &codigo)
		CheckError(err)

		j := VO.NewJugadorVO(nombre, "", foto, descripcion, 0, 0, email, codigo)
		res = append(res, j)

	}

	return res

}

// Devuelve true si el jugador y contraseña son validos
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
	CheckError(err)

	defer rows.Close()
	if rows.Next() {
		res = true
	}

	return res
}

// Devuelve true si el nombre del jugador esta en uso
func (JDAO *JugadoresDAO) EstaJugador(nombre string) bool {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	res := false

	//Buscamos si existe ese usuario
	isj := "SELECT * FROM JUGADORES WHERE nombre = $1"
	rows, err := db.Query(isj, nombre)
	CheckError(err)

	defer rows.Close()
	if rows.Next() {
		res = true
	}

	return res
}

// Modifica unicamenta la contraseña de un jugador dado el email del jugador y la contraseña
func (JDAO *JugadoresDAO) CambiarContra(email string, ncontra string) {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Actualizamos contraseña de un usuario
	cc := "UPDATE JUGADORES SET contra = ($1) WHERE email = $2"
	_, e := db.Exec(cc, ncontra, email)
	CheckError(e)
}

// Devuelve partidas pausadas por un jugador
func (jDAO *JugadoresDAO) PartidasPausadas(j string) []*VO.PartidasVO {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Obtenemos todos las partidas pausadas
	qPausa := "SELECT p.clave, p.creador, p.tipo " +
		"FROM PARTIDAS AS p JOIN PARTICIPAR AS pr ON p.clave = pr.partida " +
		"WHERE p.estado = 'pausada' AND pr.jugador = $1 "
	rows, err := db.Query(qPausa, j)
	CheckError(err)

	var res []*VO.PartidasVO

	defer rows.Close()
	for rows.Next() {
		var clave string
		var creador string
		var tipo string

		err := rows.Scan(&clave, &creador, &tipo)
		CheckError(err)

		p := VO.NewPartidasVO(clave, clave, tipo, "pausada")
		res = append(res, p)

	}

	return res

}
