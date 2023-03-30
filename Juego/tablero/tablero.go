package tablero

import (
	"container/list"
	"fmt"
	"math/rand"
	"time"
	"github.com/emirpasic/gods/lists/doublylinkedlist"
	"juego/cartas"
)

type Tablero struct {
	Mazo          *doublylinkedlist.List
	Descartes     *doublylinkedlist.List
	Combinaciones *list.List //Es una lista de doublylinkedlist donde se guardan las cartas jugadas(trios y escaleras en cada lista)
}


func RobarCarta(list *doublylinkedlist.List, mano *doublylinkedlist.List) { //Función encargada de robar una carta del mazo
	r := rand.Intn(list.Size()) + 1 //Obtiene un número aleatorio de la lista
	value, ok := list.Get(r)        //Obtiene el valor de la carta de la lista
	fmt.Println("Has robado la carta ",value)
	if ok {
		mano.Add(value) //Añade el valor a la mano
		list.Remove(r)  //Elimina el valor del mazo
	}
}

func RobarDescartes(list *doublylinkedlist.List, mano *doublylinkedlist.List) {
	value, ok := list.Get(0)
	fmt.Println("Has robado la carta de descartes",value)
	if ok {
		mano.Add(value) //Añade el valor a la mano
	}
}

func AnyadirCombinaciones(t Tablero,comb *list.List){
	aux := list.New();
	for i := comb.Front(); i != nil; i = i.Next(){
		aux.PushBack(i.Value)
	}

	for e := aux.Front(); e != nil; e = e.Next(){
		fmt.Println(e.Value)
		t.Combinaciones.PushBack(e.Value)
	}
	/*for e := t.Combinaciones.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
		fmt.Println("HOLI2")
	}*/
}

func FinTurno(mazo *doublylinkedlist.List, mano *doublylinkedlist.List, descarte *doublylinkedlist.List, i int) {
	value, _ := mano.Get(i) //Obtiene el valor de la mano a descartar
	mano.Remove(i)          //Elimina el valor de la mano
	if descarte.Size() > 1 {
		fmt.Println(descarte, "DESCARTE METE A MAZO") //Si hay más de un valor en descartes lo añade a la lista de mazo
		value, _ = descarte.Get(0)
		mazo.Add(value)
		descarte.Remove(0)
	}
	descarte.Add(value)     //Añade el valor a descartes
}

func SumaCartas(jugada *doublylinkedlist.List) int { // cuenta los puntos de la primera jugada que se hace y devuelve true si llega a 51
	total := 0
	size := jugada.Size()
	fmt.Println(size)
	for i := 0; i < size; i++ {
		v1, _ := jugada.Get(i)
		carta, _ := v1.(cartas.Carta)
		if carta.Valor == 0{
			if TrioValido(jugada){
				aux_t, _ := jugada.Get((i%size) + 1);
				carta, _ := aux_t.(cartas.Carta)
				if carta.Valor == 1{
					fmt.Println("Hola1")
					total += 11
				}else if carta.Valor >= 10{
					fmt.Println("Hola2")
					total += 10
				}else{
					fmt.Println("Hola3")
					total += carta.Valor
				}
				
			}else if EscaleraValida(jugada){
				if(i == size - 1){
					aux_t, _ := jugada.Get((i-1));
					carta, _ := aux_t.(cartas.Carta)
					if carta.Valor == 13{
						fmt.Println("Hola4")
						total += 11
					}else if carta.Valor >= 10{
						fmt.Println("Hola5")
						total += 10
					}else{
						fmt.Println("Hola6")
						total += carta.Valor + 1
					}
					
				}else{
					aux_t, _ := jugada.Get((i+1));
					carta, _ := aux_t.(cartas.Carta)
					if carta.Valor == 1 || carta.Valor > 10{
						fmt.Println("Hola7")
						total += 10
					}else{
						fmt.Println("Hola8")
						total += (carta.Valor - 1)
					}
					
				}
				
			}
		}else if carta.Valor == 1 {
			total += 11
		}else if carta.Valor >= 10 {
			total += 10
		}else {
			total += carta.Valor
		}
		fmt.Println(total)
	}
	return total
}

func PartitionMenorMayor(mano *doublylinkedlist.List, low, high int, tipo int) (*doublylinkedlist.List, int) { //Función del sort encargada de particionar los datos
	v1, _ := mano.Get(high)
	carta1, _ := v1.(cartas.Carta)
	i := low
	for j := low; j < high; j++ {
		v2, _ := mano.Get(j)
		carta2, _ := v2.(cartas.Carta)
		
		if tipo == 0 {
			if cartas.CompararCartasN(carta1, carta2) == 1{
				mano.Swap(i, j)
				i++
			}
		} else if tipo == 1 {
			if cartas.CompararCartasE(carta1, carta2) == 1{
				mano.Swap(i, j)
				i++
			}
		}
	}
	
	mano.Swap(i, high)
	return mano, i
}

func SortMenorMayor(mano *doublylinkedlist.List, low, high int, tipo int) *doublylinkedlist.List { //Función inicial del sort
	if low < high {
		var p int
		mano, p = PartitionMenorMayor(mano, low, high, tipo)
		mano = SortMenorMayor(mano, low, p-1, tipo)
		mano = SortMenorMayor(mano, p+1, high, tipo)
	}
	return mano
}

// igual que sort pero ordenando de menor a mayor (para que sea mas facil comprobar las escaleras)
func SortStartMenorMayor(mano *doublylinkedlist.List, tipo int) *doublylinkedlist.List { //Función inicial del sort
	return SortMenorMayor(mano, 0, mano.Size()-1, tipo)
}

// Devuelve la posición de los jokers que hay en una combinación 
func PosicionJoker(jugada *doublylinkedlist.List) *doublylinkedlist.List{
	jokers := doublylinkedlist.New()
	NumJokers := NumComodines(jugada)
	fmt.Println(NumJokers)
	index := 0
	if NumJokers > 0 {
		for numJ := 1; numJ <= NumJokers; numJ++ {
			for i := index; i < jugada.Size(); i++ {
				fmt.Println(i)
				v1, _ := jugada.Get(i)
				carta, _ := v1.(cartas.Carta)
				if carta.Valor == 0 {
					jokers.Add(i)
					index = i + 1
					break;
				}
			}
		}
	}
	fmt.Println(jokers)
	return jokers
}

// función que añade los jokers que estaban en un principio a la lista original, en la misma posición que antes
func AnyadirJokers(posicionJokers *doublylinkedlist.List, listaJokers *doublylinkedlist.List, jugada *doublylinkedlist.List) *doublylinkedlist.List{
	aux := doublylinkedlist.New()
	j := 0
	indice := 0
	primerJocker := true
	for i := 0; i < jugada.Size(); i++ {
		fmt.Println(i)
		if j < posicionJokers.Size() {	//añade el joker en la posicion en la que estaba
			v1, _ := posicionJokers.Get(j)
			pos, _ := v1.(int) 
			if i == pos {
				v1, _ = listaJokers.Get(j)
				carta, _ := v1.(cartas.Carta) 
				aux.Add(carta)
				j++
				if primerJocker {
					indice = i
					primerJocker = false
				} 
			} else {
				primerJocker = true
				v1, _ := jugada.Get(i)	//añade las demás cartas de la jugada
				carta, _ := v1.(cartas.Carta) 
				aux.Add(carta)
			}
		} else {
			i = indice
			v1, _ := jugada.Get(i)	//añade las demás cartas de la jugada
			carta, _ := v1.(cartas.Carta) 
			aux.Add(carta)	//CAMBIO: cartas.carta
			indice++
		}
	}
	return aux
}

// Devuelve true cuabdo se ha podido abrir con exito y
// false en caso contrario
func Abrir(jugada *doublylinkedlist.List, mano *doublylinkedlist.List, t *Tablero) bool{ //falta comprobar trios y escaleras
	posJ := PosicionJoker(jugada)
	jugada,listaJokers := cartas.SepararJokers(jugada)
	jugada = SortStartMenorMayor(jugada, 0)
	fmt.Println(jugada)
	fmt.Println(listaJokers)
	jugada = AnyadirJokers(posJ, listaJokers, jugada)

	fmt.Println(jugada)
	if !EscaleraValida(jugada) && !TrioValido(jugada) {
		fmt.Println("no valido")
		return false
	}
	listaC := doublylinkedlist.New()
	for i := 0; i <= jugada.Size() - 1; i++ {
		v1, _ := jugada.Get(i)
		carta, _ := v1.(cartas.Carta)

		ind := mano.IndexOf(carta)
		mano.Remove(ind)
		fmt.Println("carta eliminada", carta)
		listaC.Add(carta)
	}
	t.Combinaciones.PushBack(listaC)
	return true
}

// función para añadir una carta a una combinación
// Devuelve -1 si es una jugada invalida, 0 si es valida y 1 si es valida y devuelve un comodin
func AnyadirCarta(jugada *doublylinkedlist.List, mano *doublylinkedlist.List, t *Tablero, idCombinacion int) int{
	fmt.Println("jugada:",jugada)
	if !jugada.Empty() {
		v1, _ := jugada.Get(0)
		carta, _ := v1.(cartas.Carta)
		devolverJoker := false
		id_comb := 0
		for e := t.Combinaciones.Front(); e != nil; e = e.Next() {
			fmt.Println(id_comb, " HGOLA")
			if id_comb == idCombinacion {
				fmt.Println("hola ," ,t.Combinaciones.Front())
				fmt.Println(e.Value)
				fmt.Println(t.Combinaciones)
				//listaC := e.Value.(*doublylinkedlist.List)
				//crea una copia de la combinacion
				listaC := doublylinkedlist.New()
				for a:= 0; a < (e.Value.(*doublylinkedlist.List)).Size(); a++{
					valor,_ := (e.Value.(*doublylinkedlist.List)).Get(a)
					listaC.Add(valor)
				}
				fmt.Println("listaC:",listaC)
				if NumComodines(listaC) > 0 {
					posJ := PosicionJoker(listaC)
					listaC,listaJokers := cartas.SepararJokers(listaC)
					listaC = SortStartMenorMayor(listaC, 0)
					fmt.Println("listaC:",listaC)
					listaC.Add(carta)
					fmt.Println("listaC:",listaC)
					listaC = SortStartMenorMayor(listaC, 0)
					fmt.Println("listaC:",listaC)
					fmt.Println(listaJokers)
					fmt.Println(posJ)
					indice := listaC.IndexOf(carta)

					
					for i:= 0; i < posJ.Size(); i++ {
						v1, _ := posJ.Get(i)
						j, _ := v1.(int)
						if indice == j {
							devolverJoker = true
							posJ.Remove(j)
							v1, _ = listaJokers.Get(i)
							joker, _ := v1.(int)
							listaJokers.Remove(joker)
						} 
					}

					listaC = AnyadirJokers(posJ, listaJokers, listaC)
					if !EscaleraValida(listaC) && !TrioValido(listaC) {
						return -1
					}
					t.Combinaciones.Remove(e)
					t.Combinaciones.PushBack(listaC)
					ind := mano.IndexOf(carta)
					mano.Remove(ind)
					if devolverJoker {
						return 1
					} else {
						return 0
					}


				} else {
					listaC.Add(carta)
					listaC = SortStartMenorMayor(listaC, 0)
					if !EscaleraValida(listaC) && !TrioValido(listaC) {
						return -1
					}
					t.Combinaciones.Remove(e)
					t.Combinaciones.PushBack(listaC)
					ind := mano.IndexOf(carta)
					mano.Remove(ind)

					return 0
				}
			}
			id_comb++
		}
	}
	return -1
}
