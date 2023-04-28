package VO

type ManosVO struct {
	turno   int
	partida string
	carta   int
}

func NewManosVO(turno int, partida string, carta int) *ManosVO {
	m := ManosVO{turno: turno, partida: partida, carta: carta}
	return &m
}

func (m *ManosVO) GetTurno() int {
	return m.turno
}

func (m *ManosVO) GetPartida() string {
	return m.partida
}

func (m *ManosVO) GetCarta() int {
	return m.carta
}
