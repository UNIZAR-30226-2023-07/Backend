package VO

type AmistadVO struct {
	estado string
	usr1   string
	usr2   string
}

func NewAmistadVO(estado string, usr1 string, usr2 string) *AmistadVO {
	a := AmistadVO{estado: estado, usr1: usr1, usr2: usr2}
	return &a
}

func (a *AmistadVO) GetEstado() string {
	return a.estado
}

func (a *AmistadVO) GetUsr1() string {
	return a.usr1
}

func (a *AmistadVO) GetUsr2() string {
	return a.usr2
}
