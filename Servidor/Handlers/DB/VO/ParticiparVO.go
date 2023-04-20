package VO

type ParticiparVO struct {
	partida string
	jugador string
	puntos  int
}

func NewParticiparVO(partida string, jugador string, puntos int) *ParticiparVO {
	p := ParticiparVO{partida: partida, jugador: jugador, puntos: puntos}
	return &p
}

func (p *ParticiparVO) GetPartida() string {
	return p.partida
}

func (p *ParticiparVO) GetJugador() string {
	return p.jugador
}

func (p *ParticiparVO) GetPuntos() int {
	return p.puntos
}
