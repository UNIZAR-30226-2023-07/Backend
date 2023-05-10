package jugadores

import (
	"Juego/cartas"

	"github.com/emirpasic/gods/lists/doublylinkedlist"
)

type Jugador struct {
	Id    int
	Mano  *doublylinkedlist.List
	P_tor int
}

func CrearJugador(id int, mazo *doublylinkedlist.List) Jugador {
	j := Jugador{id, cartas.RepartirMano(mazo), 0}
	return j
}

func CartaMasAlta(mano *doublylinkedlist.List) int{
	r := 0;
	for i := 0; i < mano.Size(); i++{
		carta, _ := mano.Get(i)
		carta_r, _ := mano.Get(r)
		if carta.(cartas.Carta).Valor >= carta_r.(cartas.Carta).Valor{
			r = i
		}
	}
	return r;
}

func CartaMasBaja(mano *doublylinkedlist.List) int{
	r := 0;
	for i := 0; i < mano.Size(); i++{
		carta, _ := mano.Get(i)
		carta_r, _ := mano.Get(r)
		if carta.(cartas.Carta).Valor < carta_r.(cartas.Carta).Valor{
			r = i
		}
	}
	return r;
}