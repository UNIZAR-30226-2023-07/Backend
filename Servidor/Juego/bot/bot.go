package bot

import (
	"Juego/cartas"
	"Juego/jugadores"
	"Juego/tablero"
	"fmt"
	"time"

	//	"encoding/json"

	//	"github.com/olahol/melody"
	"container/list"

	"github.com/emirpasic/gods/lists/doublylinkedlist"
)

type Carta struct { //Struct utilizado para definir la estructura de datos que representa las cartas
	Valor int
	Palo  int
	Color int
}

// Función encargada de encontrar una escalera en la mano, devuelve los puntos del trio, las
// cartas que lo forman y si se ha encontrado trio
func calcularEscalerasJoker(mano *doublylinkedlist.List, joker *doublylinkedlist.List) (int,
	*doublylinkedlist.List, *doublylinkedlist.List, bool) {
	puntos := 0
	comb := doublylinkedlist.New()
	// ordenar la mano por palos de menor a mayor
	mano = cartas.SortStart(mano, 1)
	nuevoPalo := true
	hay_as := false
	ind_as := 0
	esc := false
	num_j := joker.Size()

	no_elim := -1
	if num_j > 0 {
		// bucle hasta que recorre toda la mano o encuentra una escalera
		for i := 0; i < mano.Size() && !esc; i++ {
			num_j_anyadidos := 0
			num_c := 1
			puntos_t := 0
			v1, _ := mano.Get(i)
			carta1, _ := v1.(Carta)
			// comprobar si hay as en el palo
			if nuevoPalo {
				hay_as = carta1.Valor == 1
				ind_as = i
			}
			if carta1.Valor >= 10 {
				puntos_t = puntos_t + 10
			} else {
				puntos_t = puntos_t + carta1.Valor
			}
			// lista temporal donde añadir las cartas que se van encontrando de la escalera
			l := *doublylinkedlist.New()
			l.Add(carta1)
			i_inf := i
			hay_esc := true
			mirar_j := false
			for hay_esc {
				v2, _ := mano.Get(i + 1)
				carta2, _ := v2.(Carta)
				if carta1.Palo == carta2.Palo {
					nuevoPalo = false
				} else {
					nuevoPalo = true
				}
				// comprobar si las dos cartas son escalera
				if i > mano.Size() {
					hay_esc = false
				} else if carta1.Valor+1 == carta2.Valor && carta1.Palo == carta2.Palo {
					//añadir la nueva carta a l
					l.Add(carta2)
					if carta2.Valor >= 10 {
						puntos_t = puntos_t + 10
					} else {
						puntos_t = puntos_t + carta2.Valor
					}
					num_c += 1
					i++
					carta1 = carta2
				} else if carta1.Valor == 13 && hay_as && !mirar_j {
					// hay escalera valida de la forma 11 12 AS
					// recupero la carta del as
					as, _ := mano.Get(ind_as)
					as_c, _ := as.(Carta)
					l.Add(as_c)
					puntos_t = puntos_t + 11
					num_c += 1
					mirar_j = true
				} else if carta1.Valor == carta2.Valor && carta1.Palo == carta2.Palo {
					// dos cartas con el mismo numero seguidas, avanzo indice
					i++
					if no_elim == -1 {
						no_elim = i
					} else {
						no_elim = no_elim*100 + i
					}

				} else if num_j > 0 { // mirar si puedo añadir el joker para hacer escalera
					v_joker, _ := joker.Get(num_j - 1)
					c_joker, _ := v_joker.(Carta)
					l.Add(c_joker)
					num_j_anyadidos++
					if carta1.Valor == 13 && !hay_as { //joker como as
						puntos_t = puntos_t + 11
					} else if carta1.Valor >= 10 {
						puntos_t = puntos_t + 10
						if mirar_j {
							l.Swap(l.Size()-1, l.Size()-3)
							l.Swap(l.Size()-1, l.Size()-2)
						}
					} else {
						puntos_t = puntos_t + carta1.Valor + 1
					}
					num_j--
					num_c++
					carta1 = Carta{carta1.Valor + 1, carta1.Palo, carta1.Color}
					if hay_as {
						hay_esc = false
					}
				} else {
					hay_esc = false
				}

			}
			if num_c >= 3 && num_c-num_j_anyadidos >= 2 {
				// si el numero de cartas seguidas ha sido >=3, escalera valida
				puntos += puntos_t
				// añado l a la combinación a devolver
				comb.Add(l)
				if !mirar_j {
					// si no hay AS, borro de la mano las cartas de los indices seguidos que correspondan
					k := no_elim % 100
					no_elim = no_elim / 100
					for j := i; j >= i_inf; j-- {
						if j != k {
							mano.Remove(j)
							if j < k {
								k = no_elim % 100
								no_elim = no_elim / 100
								if k == 0 {
									k = -1
								}
							}
						}

					}
				} else {
					// si hay AS, borro las cartas de la mano de los indices seguidos, ADEMAS del indice del AS
					k := no_elim % 100
					no_elim = no_elim / 100
					for j := i; j >= i_inf; j-- {
						if j != k {
							mano.Remove(j)
							if j < k {
								k = no_elim % 100
								no_elim = no_elim / 100
								if k == 0 {
									k = -1
								}
							}
						}
					}
					mano.Remove(ind_as)
				}
				for j := 0; j < num_j_anyadidos; j++ {
					joker.Remove(0)
				}
				esc = true
			} else {
				num_j = num_j + num_j_anyadidos
				mirar_j = false
			}
		}
	}

	return puntos, comb, joker, esc
}

// Función encargada de encontrar una escalera en la mano, devuelve los puntos del trio, las
// cartas que lo forman y si se ha encontrado trio
func calcularEscaleras(mano *doublylinkedlist.List) (int, *doublylinkedlist.List, bool) {
	puntos := 0
	comb := doublylinkedlist.New()
	// ordenar la mano por palos de menor a mayor
	mano = cartas.SortStart(mano, 1)
	cartas.MostrarMano(mano)
	nuevoPalo := true
	hay_as := false
	ind_as := 0
	esc := false
	no_elim := -1
	// bucle hasta que recorre toda la mano o encuentra una escalera
	for i := 0; i < mano.Size() && !esc; i++ {
		fmt.Println(mano.Size())
		fmt.Println(mano.Get(i))
		num_c := 1
		puntos_t := 0
		v1, _ := mano.Get(i)
		carta1, _ := v1.(Carta)
		fmt.Println(carta1)
		// comprobar si hay as en el palo
		if nuevoPalo {
			hay_as = carta1.Valor == 1
			ind_as = i
		}
		if carta1.Valor >= 10 {
			puntos_t = puntos_t + 10
		} else {
			puntos_t = puntos_t + carta1.Valor
		}
		// lista temporal donde añadir las cartas que se van encontrando de la escalera
		l := *doublylinkedlist.New()
		l.Add(carta1)
		i_inf := i
		hay_esc := true
		borrar_as := false
		fmt.Println("Holi")
		for hay_esc {
			v2, _ := mano.Get(i + 1)
			carta2, _ := v2.(Carta)
			fmt.Println(carta2)
			if carta1.Palo == carta2.Palo {
				nuevoPalo = false
				fmt.Println("Holi2")
			} else {
				nuevoPalo = true
				fmt.Println("Holi3")
			}
			// comprobar si las dos cartas son escalera
			if carta1.Valor+1 == carta2.Valor && carta1.Palo == carta2.Palo {
				//añadir la nueva carta a
				fmt.Println("Holi4")
				l.Add(carta2)
				if carta2.Valor >= 10 {
					puntos_t = puntos_t + 10
				} else {
					puntos_t = puntos_t + carta2.Valor
				}
				num_c += 1
				i++
			} else if i > mano.Size() {
				hay_esc = false
			} else if num_c >= 2 && carta1.Valor == 13 && hay_as {
				// hay escalera valida de la forma 11 12 AS
				// recupero la carta del as
				fmt.Println("Holi5")
				as, _ := mano.Get(ind_as)
				as_c, _ := as.(Carta)
				l.Add(as_c)
				puntos_t = puntos_t + 11
				num_c += 1
				hay_esc = false
				borrar_as = true
			} else if carta1.Valor == carta2.Valor && carta1.Palo == carta2.Palo {
				// dos cartas con el mismo numero seguidas, avanzo indice
				fmt.Println("Holi6")
				i++
				fmt.Println(i)
				if no_elim == -1 {
					no_elim = i
				} else {
					no_elim = no_elim*100 + i
				}
			} else {
				fmt.Println("Holi7")
				// no hay escalera
				hay_esc = false
			}
			carta1 = carta2
		}
		fmt.Println("Holi")
		if num_c >= 3 {
			// si el numero de cartas seguidas ha sido >=3, escalera valida
			puntos += puntos_t
			// añado l a la combinación a devolver
			comb.Add(l)
			if !borrar_as {
				// si no hay AS, borro de la mano las cartas de los indices seguidos que correspondan
				k := no_elim % 100
				no_elim = no_elim / 100
				for j := i; j >= i_inf; j-- {
					if j != k {
						mano.Remove(j)
						if j < k {
							k = no_elim % 100
							no_elim = no_elim / 100
							if k == 0 {
								k = -1
							}
						}
					}
				}
			} else {
				// si hay AS, borro las cartas de la mano de los indices seguidos, ADEMAS del indice del AS
				k := no_elim % 100
				no_elim = no_elim / 100
				for j := i; j >= i_inf; j-- {
					if j != k {
						mano.Remove(j)
						if j < k {
							k = no_elim % 100
							no_elim = no_elim / 100
							if k == 0 {
								k = -1
							}
						}
					}
				}
				mano.Remove(ind_as)
			}
			esc = true
		}
	}
	iterator := comb.Iterator()
	i := 0
	fmt.Println("Mostramos nueva combinacion")
	for iterator.Next() {
		i++
		l := iterator.Value()
		lista := l.(doublylinkedlist.List)
		iterator2 := lista.Iterator()
		for iterator2.Next() {
			c := iterator2.Value()
			/*cartas := c.(doublylinkedlist.List)
			iterator_c := cartas.Iterator()
			for iterator_c.Next() {
				v := iterator_c.Value()
				valor := v.(Carta)
				fmt.Println(valor)
				//mano.Add(valor)
			}*/
			valor := c.(Carta)
			fmt.Println(valor)
		}
	}
	return puntos, comb, esc
}

// Función encargada de encontrar un trío con joker en la mano, devuelve los puntos del trio, las
// cartas que lo forman, si se ha encontrado trio y los jokers que quedan
func calcularTriosJoker(mano *doublylinkedlist.List, joker *doublylinkedlist.List) (int,
	*doublylinkedlist.List, *doublylinkedlist.List, bool) {
	puntos := 0
	mano = cartas.SortStart(mano, 0)
	comb := doublylinkedlist.New()
	trio := false
	if !joker.Empty() {
		// bucle hasta que recorre toda la mano o encuentra un trio
		for i := 0; i < mano.Size()-2 && !trio; i++ {
			i_inf := i
			v1, _ := mano.Get(i)
			carta1, _ := v1.(Carta)
			v2, _ := mano.Get(i + 1)
			carta2, _ := v2.(Carta)
			if carta1.Valor == carta2.Valor {
				// las tres cartas tienen el mismo numero
				if carta1.Palo != carta2.Palo {
					// las tres cartas son de distinto palo
					trio = true
					// lista donde añadir las cartas del trio
					l := *doublylinkedlist.New()
					l.Add(carta1)
					l.Add(carta2)
					v_joker, _ := joker.Get(0)
					c_joker, _ := v_joker.(Carta)
					l.Add(c_joker)
					if carta1.Valor == 1 {
						puntos = puntos + 11*3
					} else if carta1.Valor >= 10 {
						puntos = puntos + 10*3
					} else {
						puntos = puntos + carta1.Valor*3
					}
					i += 1

					for j := i; j >= i_inf; j-- {
						// se eliminan de la mano las cartas que hemos cojido
						mano.Remove(j)
					}
					joker.Remove(0) // borro joker
					comb.Add(l)
				}
			}
		}
	}
	return puntos, comb, joker, trio
}

// Función encargada de encontrar un trío en la mano, devuelve los puntos del trio, las
// cartas que lo forman y si se ha encontrado trio
func calcularTrios(mano *doublylinkedlist.List) (int, *doublylinkedlist.List, bool) {
	puntos := 0
	mano = cartas.SortStart(mano, 0)
	comb := doublylinkedlist.New()
	trio := false
	// bucle hasta que recorre toda la mano o encuentra un trio
	for i := 0; i < mano.Size()-2 && !trio; i++ {
		i_inf := i
		palo := 0
		v1, _ := mano.Get(i)
		carta1, _ := v1.(Carta)
		v2, _ := mano.Get(i + 1)
		carta2, _ := v2.(Carta)
		v3, _ := mano.Get(i + 2)
		carta3, _ := v3.(Carta)
		if carta1.Valor == carta2.Valor && carta2.Valor == carta3.Valor {
			// las tres cartas tienen el mismo numero
			if carta1.Palo != carta2.Palo && carta2.Palo != carta3.Palo && carta1.Palo != carta3.Palo {
				// las tres cartas son de distinto palo
				trio = true
				// lista donde añadir las cartas del trio
				l := *doublylinkedlist.New()
				l.Add(carta1)
				l.Add(carta2)
				l.Add(carta3)
				// sumo los palos de las cartas, luego se explica porqué
				palo = palo + carta1.Palo + carta2.Palo + carta3.Palo
				if carta1.Valor == 1 {
					puntos = puntos + 11*3
				} else if carta1.Valor >= 10 {
					puntos = puntos + 10*3
				} else {
					puntos = puntos + carta1.Valor*3
				}
				i += 2
				v4, _ := mano.Get(i + 1)
				carta4, _ := v4.(Carta)
				palo += carta4.Palo
				// la suma de los cuatro palos 1+2+3+4 = 10
				// si al añadir la cuarta carta el valor que teniamos en palo + el palo de la nueva carta
				// es == 10, entonces significa que las 4 cartas tienen palo diferente, por eso puede
				// formar el cuarteto
				if carta1.Valor == carta4.Valor && palo == 10 {
					l.Add(carta4)
					if carta1.Valor == 1 {
						puntos = puntos + 11
					} else if carta1.Valor >= 10 {
						puntos = puntos + 10
					} else {
						puntos = puntos + carta1.Valor
					}
					i += 1
				}
				for j := i; j >= i_inf; j-- {
					// se eliminan de la mano las cartas que hemos cogido
					mano.Remove(j)
				}
				comb.Add(l)
			}
		}
	}
	return puntos, comb, trio
}

func separarJokers(mano *doublylinkedlist.List) (*doublylinkedlist.List, *doublylinkedlist.List) {
	mano = cartas.SortStart(mano, 0)
	joker := doublylinkedlist.New()
	hay_j := true
	cartas.MostrarMano(mano)
	fmt.Println("Me quedo aquí1")
	for hay_j {
		v, _ := mano.Get(mano.Size() - 1)
		carta, _ := v.(Carta)
		if carta.Valor == 0 {
			v2, _ := mano.Get(mano.Size() - 2)
			carta2, _ := v2.(Carta)
			if carta2.Valor == 0 {
				joker.Add(carta2)
				mano.Remove(mano.Size() - 2)
			}
			joker.Add(carta)
			mano.Remove(mano.Size() - 1)
		}
		hay_j = false
	}
	cartas.MostrarMano(mano)
	fmt.Println("Me quedo aquí2")
	return mano, joker
}

func descarteBot(mazo *doublylinkedlist.List, mano *doublylinkedlist.List, descarte *doublylinkedlist.List) {
	mano = cartas.SortStart(mano, 0)
	tablero.FinTurno(mazo, mano, descarte, mano.Size()-1)
}

func CalcularPuntosPosibles(mano *doublylinkedlist.List) (int, *doublylinkedlist.List) { //Función encargada de revisar los puntos posibles de una mano
	puntos := 0
	puntos_trio := 0
	esc := true
	comb := doublylinkedlist.New()
	comb_trio := doublylinkedlist.New()
	fmt.Println("Hola")
	mano, joker := separarJokers(mano)
	copia_mano_trio := doublylinkedlist.New()
	for e := 0; e < mano.Size(); e++ {
		valor, _ := mano.Get(e)
		copia_mano_trio.Add(valor) //Creamos copia
	}
	cartas.MostrarMano(mano)
	cartas.MostrarMano(copia_mano_trio)
	fmt.Println("LUCIA JOPEEE")
	trio := true
	for esc {
		// bucle para encontrar todas las escaleras
		puntos_m, combE, escR := calcularEscaleras(mano)
		puntos += puntos_m
		if escR {
			//añade a comb la nueva escalera encontrada
			comb.Add(combE)
		}
		esc = escR
		fmt.Println("Me quedo aquí")
		fmt.Println(esc)
	}
	esc_j := true
	for esc_j {
		puntos_m, combE, jokerR, escR := calcularEscalerasJoker(mano, joker)
		puntos += puntos_m
		if escR {
			//añade a comb el nuevo trio encontrado
			comb.Add(combE)
		}
		esc_j = escR
		joker = jokerR
	}
	for trio {
		// bucle para encontrar todos los trios
		puntos_m, combT, trioR := calcularTrios(copia_mano_trio)
		puntos_trio += puntos_m
		if trioR {
			//añade a comb el nuevo trio encontrado
			comb_trio.Add(combT)
		}
		trio = trioR
	}
	trio_j := true
	for trio_j {
		puntos_m, combT, jokerR, trioR := calcularTriosJoker(copia_mano_trio, joker)
		puntos_trio += puntos_m
		if trioR {
			//añade a comb el nuevo trio encontrado
			comb_trio.Add(combT)
		}
		trio_j = trioR
		joker = jokerR
	}

	if !(esc_j && trio_j) {
		for i := 0; i < joker.Size(); i++ {
			valor, _ := joker.Get(i)
			mano.Add(valor)
		}
	}
	if puntos > puntos_trio {
		fmt.Println("Mas puntos en combinacion")
		return puntos, comb
	}
	fmt.Println("Mas puntos en trio")
	return puntos_trio, comb_trio

}

func ComprobarColocarCarta(m *doublylinkedlist.List, t *tablero.Tablero) {
	aux := m.Size()
	enc := false
	for i := 0; i < aux; i++ {
		if enc {
			i = 0
			enc = false
			cartas.MostrarMano(m)
			fmt.Println("RESTO J", i, aux)
		}
		l_aux := doublylinkedlist.New()
		c_aux, _ := m.Get(i)
		cartas.MostrarMano(m)
		fmt.Println("Comprobamos la carta ", c_aux)
		l_aux.Add(c_aux)
		/*if(c_aux.Valor == 0){
			tablero.anyadirJoker(l_aux, m, t, i)
			break
		}*/
		for j := 0; j < t.Combinaciones.Len(); j++ {
			r := tablero.AnyadirCarta(l_aux, m, t, j)
			if r != -1 {
				fmt.Printf("Se ha añadido una carta a una combinacion")
				if r == 1 {
					fmt.Println(" Y se ha obtenido un Joker \n")
				}
				t_aux := *t
				tablero.MostrarTablero(t_aux)
				fmt.Println("RESTO I")
				enc = true
				aux = m.Size()
				goto SALIR

			}
		}
	SALIR:
		cartas.MostrarMano(m)
	}
}

func Bot_En_Funcionamiento(t tablero.Tablero, jugador interface{}, ab bool) {

	fmt.Println("El bot va a operar")

	// sleep de 5 segundos
	time.Sleep(5 * time.Second)

	tablero.RobarCarta(t.Mazo, jugador.(jugadores.Jugador).Mano)

	cartas.MostrarMano(jugador.(jugadores.Jugador).Mano)

	fmt.Println("Vamos a ver las combinaciones posibles de la mano del bot")

	p, comb := CalcularPuntosPosibles(jugador.(jugadores.Jugador).Mano)

	fmt.Println("Tenemos ", p, " puntos con las combinaciones")
	fmt.Println(comb)

	iterator := comb.Iterator()
	i := 0
	fmt.Println("Mostramos nueva combinacion")
	for iterator.Next() {
		i++
		l := iterator.Value()
		lista := l.(*doublylinkedlist.List)
		iterator2 := lista.Iterator()
		for iterator2.Next() {
			c := iterator2.Value()
			cartas := c.(doublylinkedlist.List)
			iterator_c := cartas.Iterator()
			for iterator_c.Next() {
				v := iterator_c.Value()
				fmt.Println(v)
				//mano.Add(valor)
			}
		}
	}
	if ab == false {
		if p < 51 {
			cartas.MostrarMano(jugador.(jugadores.Jugador).Mano)
			fmt.Println("No podemos abrir")
			iterator := comb.Iterator()
			i := 0
			for iterator.Next() {
				i++
				l := iterator.Value()
				lista := l.(*doublylinkedlist.List)
				iterator2 := lista.Iterator()
				for iterator2.Next() {
					c := iterator2.Value()
					cartas := c.(doublylinkedlist.List)
					iterator_c := cartas.Iterator()
					for iterator_c.Next() {
						v := iterator_c.Value()
						fmt.Println("Añadimos la carta ", v)
						jugador.(jugadores.Jugador).Mano.Add(v)
					}
				}

			}
			cartas.MostrarMano(jugador.(jugadores.Jugador).Mano)

			fmt.Println()

			des := jugadores.CartaMasAlta(jugador.(jugadores.Jugador).Mano)

			tablero.FinTurno(t.Mazo, jugador.(jugadores.Jugador).Mano, t.Descartes, des)

			cartas.MostrarMano(jugador.(jugadores.Jugador).Mano)
		} else {
			fmt.Println("Vamos a abrir")
			listaA := list.New()
			//cartas.MostrarMano(jugador.(jugadores.Jugador).Mano)
			iterator := comb.Iterator()
			for iterator.Next() {
				i++
				l := iterator.Value()
				lista := l.(*doublylinkedlist.List)
				iterator2 := lista.Iterator()
				for iterator2.Next() {
					c := iterator2.Value()
					fmt.Println(c)
					cartas := c.(doublylinkedlist.List)
					iterator_c := cartas.Iterator()
					copia := doublylinkedlist.New()
					for iterator_c.Next() {
						v := iterator_c.Value()
						copia.Add(v) //Creamos copia

					}
					listaA.PushBack(copia)
				}
			}
			fmt.Println(listaA)
			tablero.AnyadirCombinaciones(t, listaA)
			tablero.MostrarTablero(t)

			ComprobarColocarCarta(jugador.(jugadores.Jugador).Mano, &t)

			cartas.MostrarMano(jugador.(jugadores.Jugador).Mano)

			des := jugadores.CartaMasAlta(jugador.(jugadores.Jugador).Mano)

			tablero.FinTurno(t.Mazo, jugador.(jugadores.Jugador).Mano, t.Descartes, des)

			cartas.MostrarMano(jugador.(jugadores.Jugador).Mano)
		}
	} else {
		fmt.Println("VAMOS A HACER PRUEBAS")

		if p > 0 {

			listaA := list.New()
			//cartas.MostrarMano(jugador.(jugadores.Jugador).Mano)
			iterator := comb.Iterator()
			for iterator.Next() {
				i++
				l := iterator.Value()
				lista := l.(*doublylinkedlist.List)
				iterator2 := lista.Iterator()
				for iterator2.Next() {
					c := iterator2.Value()
					fmt.Println(c)
					cartas := c.(doublylinkedlist.List)
					iterator_c := cartas.Iterator()
					copia := doublylinkedlist.New()
					for iterator_c.Next() {
						v := iterator_c.Value()
						copia.Add(v) //Creamos copia

					}
					listaA.PushBack(copia)
				}
			}
			fmt.Println("AÑADIMOS LA COMBINACION")
			tablero.AnyadirCombinaciones(t, listaA)
		}

		fmt.Println("MI MANO ES:")
		cartas.MostrarMano(jugador.(jugadores.Jugador).Mano)

		fmt.Println("EL TABLERO ES:")
		tablero.MostrarTablero(t)

		ComprobarColocarCarta(jugador.(jugadores.Jugador).Mano, &t)

		fmt.Println("MI MANO ES:")
		cartas.MostrarMano(jugador.(jugadores.Jugador).Mano)

		fmt.Println("EL TABLERO ES:")
		tablero.MostrarTablero(t)

		des := jugadores.CartaMasAlta(jugador.(jugadores.Jugador).Mano)

		tablero.FinTurno(t.Mazo, jugador.(jugadores.Jugador).Mano, t.Descartes, des)

	}

}

/*func partition(mano *doublylinkedlist.List, low, high int, tipo int) (*doublylinkedlist.List, int) { //Función del sort encargada de particionar los datos
	v1, _ := mano.Get(high)
	carta1, _ := v1.(Carta)
	i := low
	for j := low; j < high; j++ {
		v2, _ := mano.Get(j)
		carta2, _ := v2.(Carta)
		if tipo == 0 {
			if cartas.compararCartasN(carta1, carta2) == -1 {
				mano.Swap(i, j)
				i++
			}
		} else if tipo == 1 {
			if cartas.compararCartasE(carta1, carta2) == -1 {
				mano.Swap(i, j)
				i++
			}
		}
	}
	mano.Swap(i, high)
	return mano, i
}*/

/*func Sort(mano *doublylinkedlist.List, low, high int, tipo int) *doublylinkedlist.List { //Función inicial del sort
	if low < high {
		var p int
		mano, p = partition(mano, low, high, tipo)
		mano = Sort(mano, low, p-1, tipo)
		mano = Sort(mano, p+1, high, tipo)
	}
	return mano
}*/

/*func SortStart(mano *doublylinkedlist.List, tipo int) *doublylinkedlist.List { //Función inicial del sort
	return Sort(mano, 0, mano.Size()-1, tipo)
}*/

// cuando se empieza el juego y no se tiene 51 puntos hay que
// añadir a la mano las combinaciones que habiamos calculado y separado
// para ello este codigo que en archivo anterior esta en una de las
// pruebas del main

/*
	iterator := comb.Iterator()
	i := 0
	for iterator.Next() {
		i++
		l := iterator.Value()
		lista := l.(*doublylinkedlist.List)
		iterator2 := lista.Iterator()
		for iterator2.Next() {
			c := iterator2.Value()
			cartas := c.(doublylinkedlist.List)
			iterator_c := cartas.Iterator()
			for iterator_c.Next() {
				v := iterator_c.Value()
				valor := v.(Carta)
				mano.Add(valor)
			}
		}

	}
*/
