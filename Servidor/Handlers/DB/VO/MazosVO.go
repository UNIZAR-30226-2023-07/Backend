package VO

type MazosVO struct {
	jugador string
	partida string
	carta   int
}

func NewMazosVO(jugador string, partida string, carta int) *MazosVO {
	m := MazosVO{jugador: jugador, partida: partida, carta: carta}
	return &m
}

func (m *MazosVO) GetJugador() string {
	return m.jugador
}

func (m *MazosVO) GetPartida() string {
	return m.partida
}

func (m *MazosVO) GetCarta() int {
	return m.carta
}
