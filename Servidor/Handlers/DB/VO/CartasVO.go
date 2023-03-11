package VO

type CartasVO struct {
	numero int
	palo   string
	foto   []byte
}

func NewCartasVO(numero int, palo string, foto []byte) *CartasVO {
	c := CartasVO{numero: numero, palo: palo, foto: foto}
	return &c
}

func (c *CartasVO) GetNumero() int {
	return c.numero
}

func (c *CartasVO) GetPalo() string {
	return c.palo
}

func (c *CartasVO) GetFoto() []byte {
	return c.foto
}
