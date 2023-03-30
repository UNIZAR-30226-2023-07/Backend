package partida

import (
	//"bufio"
	"container/list"
	//"math/rand"
	//"time"
	"fmt"
	//"net"
	//"strings"
	"strconv"
	//"os"
	"github.com/emirpasic/gods/lists/doublylinkedlist"
	"juego/cartas"
	"juego/jugadores"
	//"juego/partida"
	"juego/tablero"
)

type Partida struct {
	Jug *doublylinkedlist.List
}

func Crear_torneo(){
}

func Add_jug(j jugadores.Jugador,p Partida){
	p.Jug.Add(j)
}

func inicio_turno(espera chan string,wait chan bool){
	fin := false
	for !fin{
		var input string
		fmt.Println("¿Que acción desea hacer?")
		fmt.Scanln(&input)
		if(input == "Fin_partida"){
			fin = true
			fmt.Println("FINAL")
		}
		espera <- input
		fin = <- wait

	}

}

func IniciarPartida(){
	//jugad, err := strconv.Atoi(os.Args[1])
	//torn, err := strconv.Atoi(os.Args[2])
	//bots, err := strconv.Atoi(os.Args[3])

	input := ""

	t := tablero.IniciarTablero()

	listaJ := doublylinkedlist.New()

	var ab [3]bool

	for i := 0; i < 3; i++{
		jugador := jugadores.CrearJugador(i,t.Mazo)
		listaJ.Add(jugador)
		ab[i] = false
	}

	espera := make(chan string)
	wait := make(chan bool)
	go inicio_turno(espera,wait)

	partida := true
	turno := true
	carta_robada := false
	id := 0

	for partida{
		fmt.Println("Turno del jugador ",id)
		jugador,err := listaJ.Get(id)
		turno = true
		carta_robada = false
		if err{
			/*ab[0] = true;		//Borrar
			carta_robada = true;
			jugador.(jugadores.Jugador).Mano.Clear()
			carta := cartas.Carta{13,1,1}
			jugador.(jugadores.Jugador).Mano.Add(carta)
			carta = cartas.Carta{2,1,1}
			jugador.(jugadores.Jugador).Mano.Add(carta)
			carta = cartas.Carta{3,1,1}
			jugador.(jugadores.Jugador).Mano.Add(carta)
			carta = cartas.Carta{4,1,1}
			jugador.(jugadores.Jugador).Mano.Add(carta)*/
			for turno{
				for !carta_robada{
					resp := <- espera
					if(resp == "Robar_carta"){
						tablero.RobarCarta(t.Mazo,jugador.(jugadores.Jugador).Mano)
						carta_robada = true
						wait <- false
					}else if(resp == "Robar_carta_descartes"){
						if(t.Descartes.Size() > 0){
							tablero.RobarDescartes(t.Descartes,jugador.(jugadores.Jugador).Mano)
							carta_robada = true
						}else{
							fmt.Println("Error, no hay cartas en el descarte")
						}
						wait <- false
					}else if(resp == "Fin_partida"){
						wait <- true
						partida = false
						turno = false
						carta_robada = true
						goto SALIR
					}else if(resp == "Mostrar_mano"){
						fmt.Println("Mostrando mano: ")
						cartas.MostrarMano(jugador.(jugadores.Jugador).Mano)
						wait <- false
					}else if(resp == "Mostrar_tablero"){
						tablero.MostrarTablero(t);
						wait <- false
					}else{
						fmt.Println("Error, primero tienes que robar una carta")
						wait <- false
					}
				}
				resp := <- espera
				if(resp == "Mostrar_mano"){
					fmt.Println("Mostrando mano: ")
					cartas.MostrarMano(jugador.(jugadores.Jugador).Mano)
					wait <- false
				}else if(resp == "Descarte"){
					var input string
					fmt.Println("¿Que carta desea devolver?")
					fmt.Scanln(&input)
					i_input,_:= strconv.Atoi(input)
					tablero.FinTurno(t.Mazo,jugador.(jugadores.Jugador).Mano,t.Descartes,i_input)
					if(jugador.(jugadores.Jugador).Mano.Size() == 0){
						wait <- true
						partida = false
						turno = false
					}else{
						if(id >= 3){
							id = 0
						}else{
							id = id + 1
						}
						turno = false
						wait <- false
					}
				}else if (resp == "Mostrar_tablero"){
					tablero.MostrarTablero(t);
					wait <- false
				}else if (resp == "Colocar_combinacion"){
					if(ab[id] == false){
						fmt.Println("No puedes colocar una carta porque no has abierto")
					}else{
						lista := list.New()
						combinacion := doublylinkedlist.New()
						var vector [13]int;
						for i:=0; i < 13; i++{
							vector[i] = -1;
						}
						cont := 0
						cont_a := 0
						fmt.Println("Indice los trios a probar")
						fmt.Scanln(&input)
						input_V , _ := strconv.Atoi(input)
						for input != "FIN"{
							if(input == "END"){
								fmt.Println("Comprobando los valores...")
								for j := cont_a; j < cont; j++{
									carta, _ := jugador.(jugadores.Jugador).Mano.Get(vector[j])
									combinacion.Add(carta.(cartas.Carta))
								}
								//fmt.Printf("Tipo de combinacion: %T\n", combinacion)
								cartas.MostrarMano(combinacion)
								if(tablero.TrioValido(combinacion) || tablero.EscaleraValida(combinacion)){
										// crea una copia de la lista original
									copia := doublylinkedlist.New()
									for e:= 0; e < combinacion.Size(); e++{
										valor,_ := combinacion.Get(e)
										copia.Add(valor)
									}
									lista.PushBack(copia)
									
									
									/*for e := lista.Front(); e != nil; e = e.Next() {
										miLista := e.Value.(*doublylinkedlist.List)
										//fmt.Println(e.Value)
										//p,_ := miLista.Get(0)
										//fmt.Println(p)
										cartas.MostrarMano(miLista)
									}*/
									combinacion.Clear()
								}else{
									fmt.Println("Combinacion no valida,intentelo de nuevo")
									combinacion.Clear()
									cont = cont_a
								}
							}else if input_V >= 0 && input_V < 15{
								i_input,_ := strconv.Atoi(input)
									for i := 0; i < cont; i++{
										if vector[i] == i_input{
											fmt.Println("Carta ya introducida")
											cont = cont_a
											goto COMP_VALOR_COL
										}
									} 
									fmt.Println("Valor ", i_input, "guardado para comprobacion")
									vector[cont] = i_input
									cont++
									fmt.Println(jugador.(jugadores.Jugador).Mano.Get(i_input))
									COMP_VALOR_COL:
							}else{
								fmt.Println("Comando erroneo")
							}
							fmt.Scanln(&input)
						}
						if lista.Len() > 0 {
							//fmt.Println("Contenido de la lista:")
							//fmt.Println(lista)
							/*for e := lista.Back(); e != nil; e = e.Next() {
								miLista := e.Value.(*doublylinkedlist.List)
								cartas.MostrarMano(miLista)
							}*/
							fmt.Println(cont)
							for l := 0 ; l < cont ; l++{
								fmt.Println(cont, " ", l, " ", vector[l]-l)
								fmt.Println(jugador.(jugadores.Jugador).Mano.Get(vector[l]-l))
								cartas.MostrarMano(jugador.(jugadores.Jugador).Mano)
								jugador.(jugadores.Jugador).Mano.Remove(vector[l]-l)
								cartas.MostrarMano(jugador.(jugadores.Jugador).Mano)
							}
	
							fmt.Println("lista antes de añadir al tablero")
							aaa := lista
							for elem := aaa.Front(); elem != nil; elem = elem.Next() {
								fmt.Println(elem.Value)
							}
	
	
	
							tablero.AnyadirCombinaciones(t,lista)
						}
					}
					if(jugador.(jugadores.Jugador).Mano.Size() == 0){
						wait <- true
						partida = false
						turno = false
					}else{
						wait <- false
					}
					
				}else if (resp == "Abrir"){
					lista := list.New()
					combinacion := doublylinkedlist.New()
					puntos := 0
					var vector [13]int;
					for i:=0; i < 13; i++{
						vector[i] = -1;
					}
					cont := 0
					cont_a := 0
					fmt.Println("Indice los trios a probar")
					fmt.Scanln(&input)
					input_V , _ := strconv.Atoi(input)
					for input != "FIN"{
						if(input == "END"){
							fmt.Println("Comprobando los valores...")
							for j := cont_a; j < cont; j++{
								carta, _ := jugador.(jugadores.Jugador).Mano.Get(vector[j])
								combinacion.Add(carta.(cartas.Carta))
							}
							//fmt.Printf("Tipo de combinacion: %T\n", combinacion)
							cartas.MostrarMano(combinacion)
							if(tablero.TrioValido(combinacion) || tablero.EscaleraValida(combinacion)){
								puntos = puntos + tablero.SumaCartas(combinacion)
								fmt.Println("Combinación valida, tienes ", puntos, "puntos")
								    // crea una copia de la lista original
								copia := doublylinkedlist.New()
								for e:= 0; e < combinacion.Size(); e++{
									valor,_ := combinacion.Get(e)
									copia.Add(valor)
								}
								lista.PushBack(copia)
								
								
								/*for e := lista.Front(); e != nil; e = e.Next() {
									miLista := e.Value.(*doublylinkedlist.List)
									//fmt.Println(e.Value)
									//p,_ := miLista.Get(0)
									//fmt.Println(p)
									cartas.MostrarMano(miLista)
								}*/
								combinacion.Clear()
								fmt.Println("Inserte más combinaciones")
								cont_a = cont
							}else{
								fmt.Println("Combinacion no valida,intentelo de nuevo")
								combinacion.Clear()
								cont = cont_a
							}
						}else if input_V >= 0 && input_V < 15{
							i_input,_:= strconv.Atoi(input)
							for i := 0; i < cont; i++{
								if vector[i] == i_input{
									fmt.Println("Carta ya introducida")
									cont = cont_a
									goto COMP_VALOR
								}
							} 
							fmt.Println("Valor ", i_input, "guardado para comprobacion")
							vector[cont] = i_input
							cont++
							fmt.Println(jugador.(jugadores.Jugador).Mano.Get(i_input))
							COMP_VALOR:
						}else{
							fmt.Println("Comando erroneo")
						}
						fmt.Scanln(&input)
					}
					if puntos >= 51{
						/*jugadores.Abrir(jugador.(jugadores.Jugador))
						for e := lista.Front(); e != nil; e = e.Next() {
							miLista := e.Value.(*doublylinkedlist.List)
							p,_ := miLista.Get(0)
							fmt.Println(p)
							cartas.MostrarMano(miLista)
							fmt.Println("Felicidades, has abierto")
						}*/
						if lista.Len() > 0 {
							//fmt.Println("Contenido de la lista:")
							//fmt.Println(lista)
							/*for e := lista.Back(); e != nil; e = e.Next() {
								miLista := e.Value.(*doublylinkedlist.List)
								cartas.MostrarMano(miLista)
							}*/
							for l := 0 ; l < cont ; l++{
								jugador.(jugadores.Jugador).Mano.Remove(vector[l]-l)
							}

							fmt.Println("lista antes de añadir al tablero")
							aaa := lista
							for elem := aaa.Front(); elem != nil; elem = elem.Next() {
								fmt.Println(elem.Value)
							}



							tablero.AnyadirCombinaciones(t,lista)
							fmt.Println("Felicidades, has abierto")
							ab[id] = true
						} else {
							fmt.Println("La lista está vacía")
						}
					}else{
						fmt.Println("No has conseguido suficientes puntos, y no has podido abrir")
						lista.Init()
					}
					wait <- false
				}else if (resp == "Colocar_carta"){
					if(ab[id] == false){
						fmt.Println("No puedes colocar una carta porque no has abierto")
					}else{
						fmt.Println("¿En que combinación desea introducir su carta?")
						fmt.Scanln(&input)
						t_combinacion,_ := strconv.Atoi(input)
						fmt.Println("¿Que carta desea introducir?")
						fmt.Scanln(&input)
						i_carta,_ := strconv.Atoi(input)
						l_aux := doublylinkedlist.New()
						l_aux.Add(jugador.(jugadores.Jugador).Mano.Get(i_carta))
						r := tablero.AnyadirCarta(l_aux,jugador.(jugadores.Jugador).Mano,&t,t_combinacion)
						if r != -1 {
							fmt.Println("Carta colocada con exito")
							if r == 1{
								value := cartas.Carta{0, 4, 1}
								jugador.(jugadores.Jugador).Mano.Add(value)
							}else{
								if(jugador.(jugadores.Jugador).Mano.Size() == 0){
									wait <- true
									partida = false
									turno = false
								}
							}
						} else {
							fmt.Println("no valido")
						}
						
					}
					if(!(jugador.(jugadores.Jugador).Mano.Size() == 0)){
						wait <- false
					}
				}else if (resp == "Fin_partida"){
					wait <- true
					partida = false
					turno = false
				}else{
					fmt.Println("Operacion no valida")
					wait <- false
				}
			}
		}
		
	}
	SALIR: 
	

}