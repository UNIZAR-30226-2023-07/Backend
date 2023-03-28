package VO

type JugadoresVO struct {
	nombre   string
	contra   string
	foto     int
	descrp   string
	pjugadas int
	pganadas int
	email    string
	codigo   string
}

func NewJugadorVO(nombre string, contra string, foto int, descrp string, pjugadas int, pganadas int, email string, codigo string) *JugadoresVO {
	j := JugadoresVO{nombre: nombre, contra: contra, foto: foto, descrp: descrp, pjugadas: pjugadas, pganadas: pganadas, email: email, codigo: codigo}
	return &j
}

func (j *JugadoresVO) GetNombre() string {
	return j.nombre
}

func (j *JugadoresVO) GetContra() string {
	return j.contra
}

func (j *JugadoresVO) GetFoto() int {
	return j.foto
}

func (j *JugadoresVO) GetDescrip() string {
	return j.descrp
}

func (j *JugadoresVO) GetEmail() string {
	return j.email
}

func (j *JugadoresVO) GetCodigo() string {
	return j.codigo
}

func (j *JugadoresVO) GetPGanadas() int {
	return j.pganadas
}

func (j *JugadoresVO) GetPJugadas() int {
	return j.pjugadas
}

/*func (j *JugadoresVO) GetPuntos() int {
	return j.Puntos
}*/
