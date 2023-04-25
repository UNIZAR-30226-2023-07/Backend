package VO

type MazosVO struct {
	turno   int
	partida string
	carta   int
}

func NewMazosVO(turno int, partida string, carta int) *MazosVO {
	m := MazosVO{turno: turno, partida: partida, carta: carta}
	return &m
}

func (m *MazosVO) GetTurno() int {
	return m.turno
}

func (m *MazosVO) GetPartida() string {
	return m.partida
}

func (m *MazosVO) GetCarta() int {
	return m.carta
}
