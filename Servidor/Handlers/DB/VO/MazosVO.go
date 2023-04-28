package VO

type MazosVO struct {
	partida string
	carta   int
}

func NewMazosVO(partida string, carta int) *MazosVO {
	c := MazosVO{partida: partida, carta: carta}
	return &c
}

func (c *MazosVO) GetPartida() string {
	return c.partida
}

func (c *MazosVO) GetCarta() int {
	return c.carta
}
