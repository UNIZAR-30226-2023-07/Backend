package VO

type CartasVO struct {
	valor   int
	palo    int
	reverso int
}

func NewCartasVO(valor int, palo int, reverso int) *CartasVO {
	c := CartasVO{valor: valor, palo: palo, reverso: reverso}
	return &c
}

func (c *CartasVO) GetValor() int {
	return c.valor
}

func (c *CartasVO) GetPalo() int {
	return c.palo
}

func (c *CartasVO) GetReverso() int {
	return c.reverso
}
