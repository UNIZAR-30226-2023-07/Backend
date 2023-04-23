package VO

type DescartesVO struct {
	partida string
	carta   int
}

func NewDescartesVO(partida string, carta int) *DescartesVO {
	d := DescartesVO{partida: partida, carta: carta}
	return &d
}

func (d *DescartesVO) GetPartida() string {
	return d.partida
}

func (d *DescartesVO) GetCarta() int {
	return d.carta
}
