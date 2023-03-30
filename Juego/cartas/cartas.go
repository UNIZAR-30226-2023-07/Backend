package cartas 

import (

	"fmt"
	"math/rand"
	"time"
	"github.com/emirpasic/gods/lists/doublylinkedlist"
)

type Carta struct { //Struct utilizado para definir la estructura de datos que representa las cartas
	Valor int
	Palo  int
	Color int
}

func CreacionBaraja(list *doublylinkedlist.List) { //Función que inicializa la baraja de cartas del sistema
	rand.Seed(time.Now().UnixNano())
	carta := Carta{0, 0, 0}
	for i := 1; i <= 2; i++ {
		carta.Color = i
		for j := 1; j <= 4; j++ {
			carta.Palo = j
			for k := 1; k <= 13; k++ {
				carta.Valor = k
				list.Add(carta)
			}
		}
	}
	carta.Valor = 0 // Se añaden 2 comodines
	carta.Color = 1
	list.Add(carta)
	carta.Color = 2
	list.Add(carta)
}

func RepartirMano(list *doublylinkedlist.List) *doublylinkedlist.List { //Función encargada de, a partir de la creación de la baraja de cartas, repartir 14 de ellas
	listR := doublylinkedlist.New()
	for j := 0; j < 14; j++ {
		r := rand.Intn(list.Size()) + 1 //Crea aleatorio
		value, ok := list.Get(r)        //Obtiene el valor a repartir
		for !ok {
			r = rand.Intn(list.Size()) + 1
			value, ok = list.Get(r)
		}
		listR.Add(value) //Lo añade a la mano
		list.Remove(r)   //Lo borra

	}

	return listR
}



func MostrarMano(mano *doublylinkedlist.List) { //Función que muestra los valores de la mano repartida
	mano.Each(func(index int, value interface{}) {
		fmt.Printf("%d: %v\n", index, value)
	})
}

func CompararCartasN(a Carta, b Carta) int { //Función parte del sort encargada de filtrar las cartas por valor y color
	if a.Valor < b.Valor {
		return -1
	} else if a.Valor > b.Valor {
		return 1
	} else {
		if a.Color < b.Color {
			return -1
		} else if a.Color > b.Color {
			return 1
		} else {
			return 0
		}
	}
}

func CompararCartasE(a Carta, b Carta) int { //Función parte del sort encargada de filtrar las cartas por palo y valor
	if a.Palo < b.Palo {
		return 1
	} else if a.Palo > b.Palo {
		return -1
	} else {
		if a.Valor < b.Valor {
			return 1
		} else if a.Valor > b.Valor {
			return -1
		} else {
			return 0
		}
	}

}

func SepararJokers(mano *doublylinkedlist.List) (*doublylinkedlist.List, *doublylinkedlist.List) {
	fmt.Println("Separar jokers")
	mano = SortStart(mano, 0)
	joker := doublylinkedlist.New()
	MostrarMano(mano)
	hay_j := true
	for hay_j {
		v, _ := mano.Get(mano.Size() - 1)
		carta, _ := v.(Carta)
		fmt.Println("mirar joker ", carta)
		if carta.Valor == 0 {
			joker.Add(carta)
			mano.Remove(mano.Size() - 1)
		} else {
			hay_j = false
		}
	}
	return mano, joker
}

func Sort(mano *doublylinkedlist.List, low, high int, tipo int) *doublylinkedlist.List { //Función inicial del sort
	if low < high {
		var p int
		mano, p = partition(mano, low, high, tipo)
		mano = Sort(mano, low, p-1, tipo)
		mano = Sort(mano, p+1, high, tipo)
	}
	return mano
}

func SortStart(mano *doublylinkedlist.List, tipo int) *doublylinkedlist.List { //Función inicial del sort
	return Sort(mano, 0, mano.Size()-1, tipo)
}

func partition(mano *doublylinkedlist.List, low, high int, tipo int) (*doublylinkedlist.List, int) { //Función del sort encargada de particionar los datos
	v1, _ := mano.Get(high)
	carta1, _ := v1.(Carta)
	i := low
	for j := low; j < high; j++ {
		v2, _ := mano.Get(j)
		carta2, _ := v2.(Carta)
		if tipo == 0 {
			if CompararCartasN(carta1, carta2) == -1 {
				mano.Swap(i, j)
				i++
			}
		} else if tipo == 1 {
			if CompararCartasE(carta1, carta2) == -1 {
				mano.Swap(i, j)
				i++
			}
		}
	}
	mano.Swap(i, high)
	return mano, i
}
