package VO

type JugadoresVO struct {
	nombre string
	contra string
	//	perfil   []byte
	descrp   string
	pjugadas int
	pganadas int
	email    string
	codigo   string
}

func NewJugadorVO(nombre string, contra string, descrp string, pjugadas int, pganadas int, email string, codigo string) *JugadoresVO {
	j := JugadoresVO{nombre: nombre, contra: contra, descrp: descrp, pjugadas: pjugadas, pganadas: pganadas, email: email, codigo: codigo}
	return &j
}

func (j *JugadoresVO) GetNombre() string {
	return j.nombre
}

func (j *JugadoresVO) GetContra() string {
	return j.contra
}

/*func (j *JugadoresVO) GetPerfil() []byte {
	return j.perfil
}
*/

func (j *JugadoresVO) GetDescrip() string {
	return j.descrp
}

func (j *JugadoresVO) GetEmail() string {
	return j.email
}

func (j *JugadoresVO) GetCodigo() string {
	return j.codigo
}
