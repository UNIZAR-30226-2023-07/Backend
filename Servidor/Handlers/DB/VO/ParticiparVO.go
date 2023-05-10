package VO

type ParticiparVO struct {
	partida string
	jugador string
	puntos  int
	turno   int
	abierto string
	bot     int
}

func NewParticiparVO(partida string, jugador string, puntos int, turno int, abierto string, bot int) *ParticiparVO {
	p := ParticiparVO{partida: partida, jugador: jugador, puntos: puntos, turno: turno, abierto: abierto, bot: bot}
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

func (p *ParticiparVO) GetTurno() int {
	return p.turno
}

func (p *ParticiparVO) GetAbierto() string {
	return p.abierto
}

func (p *ParticiparVO) GetBot() int {
	return p.bot
}
