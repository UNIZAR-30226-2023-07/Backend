package VO

type JugadoresVO struct {
	nombre   string
	contra   string
	perfil   []byte
	descrp   string
	pjugadas int
	pganadas int
	codigo   string
}

func NewJugadorVO(nombre string, contra string, perfil []byte, descrp string, pjugadas int, pganadas int, codigo string) *JugadoresVO {
	j := JugadoresVO{nombre: nombre, contra: contra, perfil: perfil, descrp: descrp, pjugadas: pjugadas, pganadas: pganadas, codigo: codigo}
	return &j
}

func (j *JugadoresVO) GetNombre() string {
	return j.nombre
}

func (j *JugadoresVO) GetContra() string {
	return j.contra
}

func (j *JugadoresVO) GetPerfil() []byte {
	return j.perfil
}

func (j *JugadoresVO) GetDescrip() string {
	return j.descrp
}

func (j *JugadoresVO) GetCodigo() string {
	return j.codigo
}

type AmistadVO struct {
	estado string
	usr1   string
	usr2   string
}

func NewAmistadVO(estado string, usr1 string, usr2 string) *AmistadVO {
	a := AmistadVO{estado: estado, usr1: usr1, usr2: usr2}
	return &a
}

func (a *AmistadVO) GetEstado() string {
	return a.estado
}

func (a *AmistadVO) GetUsr1() string {
	return a.usr1
}

func (a *AmistadVO) GetUsr2() string {
	return a.usr2
}

type CartasVO struct {
	numero int
	palo   string
	foto   []byte
}

func NewCartasVO(numero int, palo string, foto []byte) *CartasVO {
	c := CartasVO{numero: numero, palo: palo, foto: foto}
	return &c
}

func (c *CartasVO) GetNumero() int {
	return c.numero
}

func (c *CartasVO) GetPalo() string {
	return c.palo
}

func (c *CartasVO) GetFoto() []byte {
	return c.foto
}
