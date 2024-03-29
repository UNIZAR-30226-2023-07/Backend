package VO

type PartidasVO struct {
	clave   string
	creador string
	tipo    string
	estado  string
}

func NewPartidasVO(clave string, creador string, tipo string, estado string) *PartidasVO {
	p := PartidasVO{clave: clave, creador: creador, tipo: tipo, estado: estado}
	return &p
}

func (p *PartidasVO) GetClave() string {
	return p.clave
}

func (p *PartidasVO) GetCreador() string {
	return p.creador
}

func (p *PartidasVO) GetTipo() string {
	return p.tipo
}

func (p *PartidasVO) GetEstado() string {
	return p.estado
}
