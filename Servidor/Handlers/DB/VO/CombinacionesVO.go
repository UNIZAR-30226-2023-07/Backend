package VO

type CombinacionesVO struct {
	partida string
	carta   int
	ncomb   int
}

func NewCombinacionesVO(partida string, carta int, ncomb int) *CombinacionesVO {
	c := CombinacionesVO{partida: partida, carta: carta, ncomb: ncomb}
	return &c
}

func (c *CombinacionesVO) GetPartida() string {
	return c.partida
}

func (c *CombinacionesVO) GetCarta() int {
	return c.carta
}

func (c *CombinacionesVO) GetNcomb() int {
	return c.ncomb
}
