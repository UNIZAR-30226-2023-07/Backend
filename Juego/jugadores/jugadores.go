package jugadores

import (
	"github.com/emirpasic/gods/lists/doublylinkedlist"
	"juego/cartas"
)

type Jugador struct{
	Id int
	Mano *doublylinkedlist.List 
	P_tor int
};

func CrearJugador(id int,mazo *doublylinkedlist.List) (Jugador){
	j := Jugador{id,cartas.RepartirMano(mazo),0}
	return j
}


