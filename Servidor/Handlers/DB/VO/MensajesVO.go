package VO

type MensajesVO struct {
	jug_emi   string
	jug_rcp   string
	contenido string
	leido     int
}

func NewMensajesVO(jug_emi string, jug_rcp string, contenido string, leido int) *MensajesVO {
	m := MensajesVO{jug_emi: jug_emi, jug_rcp: jug_rcp, contenido: contenido, leido: leido}
	return &m
}

func (m *MensajesVO) GetEmisor() string {
	return m.jug_emi
}

func (m *MensajesVO) GetReceptor() string {
	return m.jug_rcp
}

func (m *MensajesVO) GetContenido() string {
	return m.contenido
}

func (m *MensajesVO) GetLeido() int {
	return m.leido
}
