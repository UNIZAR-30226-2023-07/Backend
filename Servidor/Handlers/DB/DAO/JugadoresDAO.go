package DAO

import (
	"DB/VO"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "52.174.124.24"
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
	addj := "INSERT INTO JUGADORES VALUES ($1, $2, 0, $3, 0, 0, 0, $4, $5)"
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

// Devuelve la información del jugador que coincida con el email
func (JDAO *JugadoresDAO) GetJugador(email string) *VO.JugadoresVO {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Busca el usuario cuyo email coincide con el de JVO
	getj := "SELECT nombre, descrp, foto, pjugadas, pganadas, puntos, codigo FROM JUGADORES WHERE email = $1"
	rows, err := db.Query(getj, email)
	CheckError(err)

	defer rows.Close()
	if rows.Next() {
		var nombre string
		var descripcion string
		var foto int
		var pjugadas int
		var pganadas int
		var puntos int
		var codigo string

		err := rows.Scan(&nombre, &descripcion, &foto, &pjugadas, &pganadas, &puntos, &codigo)
		CheckError(err)

		jVO := VO.NewJugadorVO(nombre, "", foto, descripcion, pjugadas, pganadas, puntos, email, codigo)
		return jVO

	} else {
		return nil
	}

}

// Devuelve la información del jugador que coincida con el codigo
func (JDAO *JugadoresDAO) GetJugador2(code string) *VO.JugadoresVO {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Busca el usuario cuyo email coincide con el de JVO
	getj := "SELECT nombre, descrp, foto, pjugadas, pganadas, puntos, email FROM JUGADORES WHERE codigo = $1"
	rows, err := db.Query(getj, code)
	CheckError(err)

	defer rows.Close()
	if rows.Next() {
		var nombre string
		var descripcion string
		var foto int
		var pjugadas int
		var pganadas int
		var puntos int
		var email string

		err := rows.Scan(&nombre, &descripcion, &foto, &pjugadas, &pganadas, &puntos, &email)
		CheckError(err)

		jVO := VO.NewJugadorVO(nombre, "", foto, descripcion, pjugadas, pganadas, puntos, email, code)
		return jVO

	} else {
		return nil
	}

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
	qAmis := "SELECT r.nombre, r.foto, r.descrp, r.email, r.codigo, r.puntos " +
		"FROM AMISTAD AS a JOIN JUGADORES AS r ON a.usr2 = r.codigo " +
		"WHERE a.usr1 = $1 AND a.estado = 'confirmada' " +
		"ORDER BY r.puntos DESC"
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
		var puntos int

		err := rows.Scan(&nombre, &foto, &descripcion, &email, &codigo, &puntos)
		CheckError(err)

		j := VO.NewJugadorVO(nombre, "", foto, descripcion, 0, 0, puntos, email, codigo)
		res = append(res, j)

	}

	return res

}

// Lista las peticiones no confirmadas del jugador j
func (jDAO *JugadoresDAO) ListarPendientes(j string) ([]*VO.JugadoresVO, []*VO.AmistadVO) {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Obtenemos todos los parametros para cada jugador del que se espera confirmación
	qAmis := "SELECT r.nombre, r.foto, r.descrp, r.email, r.codigo " +
		"FROM AMISTAD AS a JOIN JUGADORES AS r ON a.usr2 = r.codigo " +
		"WHERE a.usr1 = $1 AND a.estado = 'esp_confirmacion'"
	rows, err := db.Query(qAmis, j)
	CheckError(err)

	var resj []*VO.JugadoresVO
	var resa []*VO.AmistadVO

	defer rows.Close()
	for rows.Next() {
		var nombre string
		var foto int
		var descripcion string
		var email string
		var codigo string

		err := rows.Scan(&nombre, &foto, &descripcion, &email, &codigo)
		CheckError(err)

		jVO := VO.NewJugadorVO(nombre, "", foto, descripcion, 0, 0, 0, email, codigo)
		resj = append(resj, jVO)

		aVO := VO.NewAmistadVO("esp_confirmacion", j, codigo)
		resa = append(resa, aVO)

	}

	//Obtenemos todos los parametros para cada jugador que tiene pendiente por responder
	qAmis = "SELECT r.nombre, r.foto, r.descrp, r.email, r.codigo " +
		"FROM AMISTAD AS a JOIN JUGADORES AS r ON a.usr2 = r.codigo " +
		"WHERE a.usr1 = $1 AND a.estado = 'pendiente'"
	rows, err = db.Query(qAmis, j)
	CheckError(err)

	defer rows.Close()
	for rows.Next() {
		var nombre string
		var foto int
		var descripcion string
		var email string
		var codigo string

		err := rows.Scan(&nombre, &foto, &descripcion, &email, &codigo)
		CheckError(err)

		jVO := VO.NewJugadorVO(nombre, "", foto, descripcion, 0, 0, 0, email, codigo)
		resj = append(resj, jVO)

		aVO := VO.NewAmistadVO("pendiente", codigo, j)
		resa = append(resa, aVO)

	}

	return resj, resa

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
	qPausa := "SELECT p.clave, p.creador, p.tipo, p.torneo " +
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
		var torneo string

		err := rows.Scan(&clave, &creador, &tipo, &torneo)
		CheckError(err)

		p := VO.NewPartidasVO(clave, clave, tipo, "pausada", torneo)
		res = append(res, p)

	}

	return res

}

// Actualiza los puntos
func (JDAO *JugadoresDAO) AddPuntos(code string, puntos int) {
	//String para la conexión
	psqlcon := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//abrir base de datos
	db, err := sql.Open("postgres", psqlcon)
	CheckError(err)

	//cerrar base de datos
	defer db.Close()

	//Actualizamos los puntos de un usuario
	ap := "UPDATE JUGADORES SET puntos = ($1) WHERE codigo = $2"
	_, e := db.Exec(ap, puntos, code)
	CheckError(e)
}
