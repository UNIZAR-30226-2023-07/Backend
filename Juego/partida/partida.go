package partida

import (
	"juego/jugadores"

	"github.com/emirpasic/gods/lists/doublylinkedlist"
)

type Partida struct {
	Jug *doublylinkedlist.List
}

func Crear_torneo(){
}

func Add_jug(j jugadores.Jugador,p Partida){
	p.Jug.Add(j)
}