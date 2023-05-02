package tablero

import (
	"Juego/cartas"
	"container/list"
	"fmt"
	"math/rand"
	"time"

	"github.com/emirpasic/gods/lists/doublylinkedlist"
)

type Tablero struct {
	Mazo          *doublylinkedlist.List
	Descartes     *doublylinkedlist.List
	Combinaciones *list.List //Es una lista de doublylinkedlist donde se guardan las cartas jugadas(trios y escaleras en cada lista)
}

func RobarCarta(list *doublylinkedlist.List, mano *doublylinkedlist.List) cartas.Carta{ //Función encargada de robar una carta del mazo
	fmt.Println(list.Size())
	r := rand.Intn(list.Size()) + 1 //Obtiene un número aleatorio de la lista
	value, ok := list.Get(r)        //Obtiene el valor de la carta de la lista
	cart := value.(cartas.Carta)
	fmt.Println("Has robado la carta ", cart)
	if ok {
		mano.Add(value) //Añade el valor a la mano
		list.Remove(r)  //Elimina el valor del mazo
	}
	return cart
}

func RobarDescartes(list *doublylinkedlist.List, mano *doublylinkedlist.List) cartas.Carta{
	value, ok := list.Get(0)
	cart := value.(cartas.Carta)
	fmt.Println("Has robado la carta de descartes", cart)
	if ok {
		mano.Add(value) //Añade el valor a la mano
	}
	return cart
}

func AnyadirCombinaciones(t Tablero, comb *list.List) {
	aux := list.New()
	for i := comb.Front(); i != nil; i = i.Next() {
		aux.PushBack(i.Value)
	}

	for e := aux.Front(); e != nil; e = e.Next() {
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
	if descarte.Size() > 0 {
		fmt.Println(descarte, "DESCARTE METE A MAZO") //Si hay más de un valor en descartes lo añade a la lista de mazo
		valueDesc, _ := descarte.Get(0)
		mazo.Add(valueDesc)
		descarte.Remove(0)
	}
	descarte.Add(value) //Añade el valor a descartes
}

func SumaCartas(jugada *doublylinkedlist.List) int { // cuenta los puntos de la primera jugada que se hace y devuelve true si llega a 51
	total := 0
	size := jugada.Size()
	fmt.Println(size)
	for i := 0; i < size; i++ {
		v1, _ := jugada.Get(i)
		carta, _ := v1.(cartas.Carta)
		if carta.Valor == 0 {
			if TrioValido(jugada) {
				aux_t, _ := jugada.Get((i % size) + 1)
				carta, _ := aux_t.(cartas.Carta)
				if carta.Valor == 1 {
					fmt.Println("Hola1")
					total += 11
				} else if carta.Valor >= 10 {
					fmt.Println("Hola2")
					total += 10
				} else {
					fmt.Println("Hola3")
					total += carta.Valor
				}

			} else if EscaleraValida(jugada) {
				if i == size-1 {
					aux_t, _ := jugada.Get((i - 1))
					carta, _ := aux_t.(cartas.Carta)
					if carta.Valor == 13 {
						fmt.Println("Hola4")
						total += 11
					} else if carta.Valor >= 10 {
						fmt.Println("Hola5")
						total += 10
					} else {
						fmt.Println("Hola6")
						total += carta.Valor + 1
					}

				} else {
					aux_t, _ := jugada.Get((i + 1))
					carta, _ := aux_t.(cartas.Carta)
					if carta.Valor == 1 || carta.Valor > 10 {
						fmt.Println("Hola7")
						total += 10
					} else {
						fmt.Println("Hola8")
						total += (carta.Valor - 1)
					}

				}

			}
		} else if carta.Valor == 1 {
			total += 11
		} else if carta.Valor >= 10 {
			total += 10
		} else {
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
			if cartas.CompararCartasN(carta1, carta2) == 1 {
				mano.Swap(i, j)
				i++
			}
		} else if tipo == 1 {
			if cartas.CompararCartasE(carta1, carta2) == 1 {
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
func PosicionJoker(jugada *doublylinkedlist.List) *doublylinkedlist.List {
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
					break
				}
			}
		}
	}
	fmt.Println(jokers)
	return jokers
}

// función que añade los jokers que estaban en un principio a la lista original, en la misma posición que antes
func AnyadirJokers(posicionJokers *doublylinkedlist.List, listaJokers *doublylinkedlist.List, jugada *doublylinkedlist.List) *doublylinkedlist.List {
	aux := doublylinkedlist.New()
	j := 0
	indice := 0
	primerJocker := true
	for i := 0; i < jugada.Size(); i++ {
		fmt.Println(i)
		if j < posicionJokers.Size() { //añade el joker en la posicion en la que estaba
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
				v1, _ := jugada.Get(i) //añade las demás cartas de la jugada
				carta, _ := v1.(cartas.Carta)
				aux.Add(carta)
			}
		} else {
			i = indice
			v1, _ := jugada.Get(i) //añade las demás cartas de la jugada
			carta, _ := v1.(cartas.Carta)
			aux.Add(carta) //CAMBIO: cartas.carta
			indice++
		}
	}
	return aux
}

// Devuelve true cuabdo se ha podido abrir con exito y
// false en caso contrario
func Abrir(jugada *doublylinkedlist.List, mano *doublylinkedlist.List, t *Tablero) bool { //falta comprobar trios y escaleras
	posJ := PosicionJoker(jugada)
	jugada, listaJokers := cartas.SepararJokers(jugada)
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
	for i := 0; i <= jugada.Size()-1; i++ {
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

func AnyadirJoker (jugada *doublylinkedlist.List, mano *doublylinkedlist.List, t *Tablero, idCombinacion int) int {
	if !jugada.Empty(){
		for e := t.Combinaciones.Front(); e != nil; e = e.Next() {
			/*if id_comb == idCombinacion {
				
			}*/
		}
		return 0
	}else{
		return -1
	}
}

func AnyadirCarta(jugada *doublylinkedlist.List, mano *doublylinkedlist.List, t *Tablero, idCombinacion int) int {
	if !jugada.Empty() {
		listaJokers := doublylinkedlist.New()
		v1, _ := jugada.Get(0)
		carta, _ := v1.(cartas.Carta)
		id_comb := 0
		for e := t.Combinaciones.Front(); e != nil; e = e.Next() {
			if id_comb == idCombinacion {				
				listaC := doublylinkedlist.New()
				for a := 0; a < (e.Value.(*doublylinkedlist.List)).Size(); a++ {
					valor, _ := (e.Value.(*doublylinkedlist.List)).Get(a)
					listaC.Add(valor)
				}
				fmt.Println("listaC:", listaC)
				if NumComodines(listaC) > 0 {
					listaC.Add(carta)
					posJ := PosicionJoker(listaC)
					for i := 0; i < posJ.Size(); i++ {
						v1, _ := posJ.Get(i)
						j, _ := v1.(int)
						listaJokers.Add(listaC.Get(j))
						listaC.Remove(j)
						
					}
					listaC = SortStartMenorMayor(listaC, 0)
					if EscaleraValida(listaC) || TrioValido(listaC){
						/*carta := cartas.Carta{0, 4, 1}
						mano.Add(carta)*/
						return 1
					}
					fmt.Println("listaCP:", listaC)
					for i := 0; i < posJ.Size(); i++ {
						v1, _ := posJ.Get(i)
						j, _ := v1.(int)
						fmt.Println(j,"HEEEY",posJ.Size())
						aux, _ := listaJokers.Get(i)
						aux_j, _ := aux.(cartas.Carta)
						listaC.Insert(j + 1,aux_j)
					}
					fmt.Println("listaCP:", listaC)
					fmt.Println("Seguimos")
					if !EscaleraValida(listaC) && !TrioValido(listaC) {
						fmt.Println("No valido")
						return -1
					}
					t.Combinaciones.Remove(e)
					t.Combinaciones.PushBack(listaC)
					ind := mano.IndexOf(carta)
					mano.Remove(ind)
					fmt.Println("Valido")
					return 0
				}else{
					listaC.Add(carta)
					listaC = SortStartMenorMayor(listaC, 0)
					fmt.Println("listaCA:", listaC)
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

/*// función para añadir una carta a una combinación
// Devuelve -1 si es una jugada invalida, 0 si es valida y 1 si es valida y devuelve un comodin
func AnyadirCarta(jugada *doublylinkedlist.List, mano *doublylinkedlist.List, t *Tablero, idCombinacion int) int {
	fmt.Println("jugada:", jugada)
	if !jugada.Empty() {
		v1, _ := jugada.Get(0)
		carta, _ := v1.(cartas.Carta)
		devolverJoker := false
		id_comb := 0
		for e := t.Combinaciones.Front(); e != nil; e = e.Next() {
			fmt.Println(id_comb, " HGOLA")
			if id_comb == idCombinacion {
				fmt.Println("hola ,", t.Combinaciones.Front())
				fmt.Println(e.Value)
				fmt.Println(t.Combinaciones)
				//listaC := e.Value.(*doublylinkedlist.List)
				//crea una copia de la combinacion
				listaC := doublylinkedlist.New()
				for a := 0; a < (e.Value.(*doublylinkedlist.List)).Size(); a++ {
					valor, _ := (e.Value.(*doublylinkedlist.List)).Get(a)
					listaC.Add(valor)
				}
				fmt.Println("listaC:", listaC)
				if NumComodines(listaC) > 0 {
					posJ := PosicionJoker(listaC)
					listaC, listaJokers := cartas.SepararJokers(listaC)
					listaC = SortStartMenorMayor(listaC, 0)
					fmt.Println("listaC:", listaC)
					listaC.Add(carta)
					fmt.Println("listaC:", listaC)
					listaC = SortStartMenorMayor(listaC, 0)
					fmt.Println("listaC:", listaC)
					t_aux := *t
					MostrarTablero(t_aux)
					fmt.Println(listaJokers)
					fmt.Println(posJ)
					indice := listaC.IndexOf(carta)

					for i := 0; i < posJ.Size(); i++ {
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
					fmt.Println("AAAAAAAAA")
					MostrarTablero(t_aux)
					//fmt.Println("listaC:", listaC)
					//listaC = AnyadirJokers(posJ, listaJokers, listaC)
					fmt.Println("listaC:", listaC)
					if !EscaleraValida(listaC) && !TrioValido(listaC) {
						return -1
					}
					t.Combinaciones.Remove(e)
					t.Combinaciones.PushBack(listaC)
					fmt.Println("AAAAAAAAA")
					MostrarTablero(t_aux)
					ind := mano.IndexOf(carta)
					mano.Remove(ind)
					if devolverJoker {
						carta := cartas.Carta{0, 4, 1}
						mano.Add(carta)
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
}*/

// Se muestra el Tablero por pantalla (para depuración)
func MostrarTablero(t Tablero) {
	fmt.Println("MAZO: ", t.Mazo)

	fmt.Println("DESCARTES: ", t.Descartes)

	fmt.Println("COMBINACIONES: ", t.Combinaciones)

	l := t.Combinaciones
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
	fmt.Println("---------------------------------------\n")

}

// Inicializa el Tablero y la mano del jugador (hay que cambiar lo de repartirMano cuando se hagan más jugadores)
func IniciarTablero() Tablero {

	rand.Seed(time.Now().UnixNano())
	mazo := doublylinkedlist.New()
	descarte := doublylinkedlist.New()

	cartas.CreacionBaraja(mazo)

	t := Tablero{mazo, descarte, list.New()}

	aux := doublylinkedlist.New()
	carta := cartas.Carta{0, 4, 1}
	aux.Add(carta)
	carta = cartas.Carta{11, 1, 1}
	aux.Add(carta)
	carta = cartas.Carta{12, 1, 1}
	aux.Add(carta)
	carta = cartas.Carta{13, 1, 1}
	aux.Add(carta)

	t.Combinaciones.PushBack(aux)

	aux2:= doublylinkedlist.New()
	carta = cartas.Carta{6, 4, 1}
	aux2.Add(carta)
	carta = cartas.Carta{6, 1, 1}
	aux2.Add(carta)
	carta = cartas.Carta{6, 3, 1}
	aux2.Add(carta)

	t.Combinaciones.PushBack(aux2)

	return t
}

//función que llama a la jugada que indique el jugador
/*func RealizarJugada(t *Tablero, mano *doublylinkedlist.List, jugada int, i int, cartasAjugar *doublylinkedlist.List) {
	switch jugada {
	case 0: //Descarte
		FinTurno(t.Mazo, mano, t.Descartes, i)
		return
	case 1: //Robar
		RobarCarta(t.Mazo, mano)
		return
	case 2: //Abrir
		if Suma51(cartasAjugar) {
			Abrir(cartasAjugar, mano, t)
		}
		cartasAjugar.Clear()
		return
	case 3: //Añadir 1 carta a una combinación existente
		AnyadirCarta(cartasAjugar, mano, t, 0)
		cartasAjugar.Clear()
		return
	default:
	}
}*/

/*
Pre: TRUE
Post: return true si es un comodin, es decir vale 0

	false en caso contrario
*/
func EsComodin(valor int) bool {
	return valor == 0
}

/*
Pre: TRUE
Post: devuelve el numero de comodines en la lista
*/
func NumComodines(jugada *doublylinkedlist.List) int {
	num_comodines := 0
	for j := 0; j < jugada.Size(); j++ {
		cart, _ := jugada.Get(j)
		carta, _ := cart.(cartas.Carta)
		ValorCarta := carta.Valor

		if EsComodin(ValorCarta) {
			num_comodines++
		}
	}
	return num_comodines
}

/*
Pre: lista ordenada en orden de jugada (Consideramos que el comodin es el 0)
Post: return true si es una escalera válida en el juego del Rabino, y

	false en caso contrario
*/
func EscaleraValida(jugada *doublylinkedlist.List) bool {

	if jugada.Empty() { //Si la lista de la jugada es vacia
		return false
	}

	num_cartas := jugada.Size()

	//COMPROBACION: NUMERO DE CARTAS VALIDO
	//Escalera maxima: 1,2,3,4,5,6,7,8,9,10,Sota(11),Caballo(12),Rey(13),As(1 o 0)
	if num_cartas > 14 { //Tamagno maximo de escalera 14
		return false
	}

	num_comodines := NumComodines(jugada)

	//COMPROBACION: NUMERO DE COMODINES VALIDO
	if num_comodines > (num_cartas - 2) { //Numero de comodines es como mucho num_cartas - 2
		return false
	}

	//COMPROBACION: TIENEN EL MISMO PALO
	index := 0                   //Indice inicial
	cart, _ := jugada.Get(index) //Cogemos la primera carta
	carta, _ := cart.(cartas.Carta)
	CartaValorRef := carta.Valor //Sacamos este valor por si es un comodin
	PaloCartaRef := carta.Palo

	for EsComodin(CartaValorRef) { //Tomamos de referencia una carta que no sea comodin
		index++                      //Miramos la siguiente carta
		cart, _ := jugada.Get(index) //Cogemos la primera carta
		carta, _ := cart.(cartas.Carta)
		CartaValorRef = carta.Valor //Sacamos este valor por si es un comodin
		PaloCartaRef = carta.Palo
	}

	for u := index + 1; u < jugada.Size(); u++ { //Miramos que tenga todas las cartas el mismo palo
		cart1, _ := jugada.Get(u)
		carta1, _ := cart1.(cartas.Carta)
		CartaValorMirar := carta1.Valor //Sacamos este valor por si es un comodin
		PaloCartaMirar := carta1.Palo

		fmt.Println("SAUL", PaloCartaRef,PaloCartaMirar,carta,carta1,!EsComodin(CartaValorRef),!EsComodin(CartaValorMirar), PaloCartaRef != PaloCartaMirar)
		if PaloCartaRef != PaloCartaMirar && !EsComodin(CartaValorRef) && !EsComodin(CartaValorMirar) { //Si tiene distinto palo, no valido
			fmt.Println("False 1")
			return false
		}
	}

	//COMPROBACION: VALOR DE CARTAS CRECIENTE
	index = 0 //Indice inicial
	cart, _ = jugada.Get(index)
	carta, _ = cart.(cartas.Carta)
	CartaValorRef = carta.Valor //Sacamos este valor de la primera carta

	for EsComodin(CartaValorRef) { //Tomamos de referencia una carta que no sea comodin
		index++ //Miramos la siguiente carta
		cart, _ = jugada.Get(index)
		carta, _ = cart.(cartas.Carta)
		CartaValorRef = carta.Valor
	}

	//El numero de comodines que hay delante de 1 debe ser cero, el numero de comodines que pueden ir delante
	// de 2 es 1, y asi sucesibamente. Index en este caso tambien tomaria el numero de comodines que hay delante
	if (CartaValorRef - index) <= 0 {
		fmt.Println("False 1")
		return false
	}

	for j := index + 1; j < jugada.Size(); j++ { //Miramos si el valor de las cartas es creciente
		//Empezamos a comparar con cartas posteriores a la de referencia
		CartaValorRef++
		cart1, _ := jugada.Get(j)
		carta1, _ := cart1.(cartas.Carta)
		CartaValorMirar := carta1.Valor //Sacamos este valor de la carta

		if !EsComodin(CartaValorMirar) { //Si es un comodin seguro que valdra para la escalera
			if CartaValorMirar != 1 || CartaValorRef != 14 { //Esta condicion no se cumple cuando despues del Rey(13), se pone un As(1)
				if CartaValorMirar != CartaValorRef { //Si concuerda el valor con lo que deberia dar(p.ejem: 2 != 15(valor despues de As))
					fmt.Println("False 2")
					return false
				}
			}
		} else if CartaValorRef > 14 { //Si la jugada continua despues de ..., Rey, As,...sera erronea
			fmt.Println("False 3")
			return false
		}
	}

	return true //Si cumple todas las condiciones
}

/*
Pre: TRUE
Post: return true si es una trio o cuarteto válida en el juego del Rabino, y

	false en caso contrario
*/
func TrioValido(jugada *doublylinkedlist.List) bool {

	if jugada.Empty() { //Si la lista de la jugada es vacia
		return false
	}

	//COMPROBACION: ES UN TRIO O UN CUARTETO
	if jugada.Size() < 3 || jugada.Size() > 4 { //El tamaño de la jugada puede ser 3 o 4
		return false
	}

	//COMPROBACION: NUMERO DE COMODINES VALIDO
	num_comodines := NumComodines(jugada)
	if num_comodines > 1 { //Numero de comodines es como mucho 1
		return false
	}

	//COMPROBACION: TIENEN EL MISMO VALOR
	index := 0
	cart, _ := jugada.Get(index)
	carta, _ := cart.(cartas.Carta)
	ValorCartaRef := carta.Valor

	for EsComodin(ValorCartaRef) { //Tomamos de referencia una carta que no sea comodin
		index++ //Miramos la siguiente carta
		cart, _ = jugada.Get(index)
		carta, _ = cart.(cartas.Carta)
		ValorCartaRef = carta.Valor
	}

	for i := index + 1; i < jugada.Size(); i++ { //Comprobamos que todas las cartas tengan el mismo valor

		cart1, _ := jugada.Get(i)
		carta1, _ := cart1.(cartas.Carta)
		ValorCarta := carta1.Valor

		if ValorCartaRef != ValorCarta && !EsComodin(ValorCarta) { //Si la carta a comparar es un comodin sera valida
			return false
		}
	}

	//COMPROBACION: TIENEN DISTINTO PALO
	for j := 0; j < jugada.Size(); j++ {

		cart, _ := jugada.Get(j)
		carta, _ := cart.(cartas.Carta)
		CartaValorRef := carta.Valor //Sacamos este valor por si es un comodin
		PaloCartaRef := carta.Palo

		if !EsComodin(CartaValorRef) { //Si es un comodin seguro que sera valido
			for u := j + 1; u < jugada.Size(); u++ { //Miramos que tenga todas las cartas el mismo palo
				cart1, _ := jugada.Get(u)
				carta1, _ := cart1.(cartas.Carta)
				CartaValorMirar := carta1.Valor //Sacamos este valor por si es un comodin
				PaloCartaMirar := carta1.Palo

				if PaloCartaRef == PaloCartaMirar && !EsComodin(CartaValorMirar) { //Si tiene distinto palo, no valido
					return false
				}
			}
		}
	}

	return true //Si cumple todas las condiciones

}
