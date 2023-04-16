package partida

import (
	//"bufio"
	"container/list"
	//"math/rand"
	//"time"
	"fmt"
	//"net"
	"strings"
	"strconv"
	//"os"
	"github.com/emirpasic/gods/lists/doublylinkedlist"
	"Servidor/Juego/cartas"
	"Servidor/Juego/jugadores"
	//"juego/partida"
	"Servidor/Juego/tablero"
	"Servidor/Juego/bot"
)

type Partida struct {
	Jug *doublylinkedlist.List
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

func IniciarPartida(idPartida string, canalPartida chan string) *doublylinkedlist.List{
	//jugad, err := strconv.Atoi(os.Args[1])
	//torn, err := strconv.Atoi(os.Args[2])
	//bots, err := strconv.Atoi(os.Args[3])
	var parametrosPartida string
	parametrosPartida = <-canalPartida
	// separar los parametros por el caracter ","
	param := strings.Split(parametrosPartida, ",")
	numJugad, _ := strconv.Atoi(param[0])

	input := ""

	t := tablero.IniciarTablero()				//función de inicio de tablero para la partida

	listaJ := doublylinkedlist.New()

	var ab [3]bool

	for i := 0; i < numJugad; i++{								//Inicio de los jugadores
		jugador := jugadores.CrearJugador(i,t.Mazo)
		listaJ.Add(jugador)
		ab[i] = false
	}

	espera := make(chan string)
	wait := make(chan bool)
	go inicio_turno(espera,wait)						//Inicio de la escucha a la terminal

	partida := true
	turno := true
	carta_robada := false
	id := 0

	for partida{							//Mientras sigamos en la partida
		fmt.Println("Turno del jugador ",id)
		jugador,err := listaJ.Get(id)		//Inicio de los turnos
		turno = true						//Ponemos turno a true porque seguimos en un turno
		carta_robada = false				//Y la carta robada a false para limitar las acciones hasta que robe una carta
		if err{
			ab[0] = false;					//PRUEBAS
			carta_robada = true;
			jugador.(jugadores.Jugador).Mano.Clear()
			carta := cartas.Carta{13,1,1}
			jugador.(jugadores.Jugador).Mano.Add(carta)
			carta = cartas.Carta{12,1,1}
			jugador.(jugadores.Jugador).Mano.Add(carta)
			carta = cartas.Carta{11,1,1}
			jugador.(jugadores.Jugador).Mano.Add(carta)
			carta = cartas.Carta{10,1,1}
			jugador.(jugadores.Jugador).Mano.Add(carta)
			carta = cartas.Carta{9,1,1}
			jugador.(jugadores.Jugador).Mano.Add(carta)
			carta = cartas.Carta{8,1,1}
			jugador.(jugadores.Jugador).Mano.Add(carta)
			carta = cartas.Carta{5,3,1}
			jugador.(jugadores.Jugador).Mano.Add(carta)
			carta = cartas.Carta{13,1,1}
			jugador.(jugadores.Jugador).Mano.Add(carta)
			carta = cartas.Carta{4,3,1}
			jugador.(jugadores.Jugador).Mano.Add(carta)
			carta = cartas.Carta{3,3,1}
			jugador.(jugadores.Jugador).Mano.Add(carta)
			carta = cartas.Carta{7,2,1}
			jugador.(jugadores.Jugador).Mano.Add(carta)
			for turno{						//Mientras nos encontremos en un turno
				if(id == 0){
					fmt.Println("El bot va a operar")

					tablero.RobarCarta(t.Mazo,jugador.(jugadores.Jugador).Mano)
				
					cartas.MostrarMano(jugador.(jugadores.Jugador).Mano)

					fmt.Println("Vamos a ver las combinaciones posibles de la mano del bot")
				
					p,comb := bot.CalcularPuntosPosibles(jugador.(jugadores.Jugador).Mano)
				
					fmt.Println("Tenemos " , p , " puntos con las combinaciones")
					fmt.Println(comb)
				
					bot.ComprobarColocarCarta(jugador.(jugadores.Jugador).Mano,&t)
				
					tablero.FinTurno(t.Mazo,jugador.(jugadores.Jugador).Mano,t.Descartes,0)
					wait <- false
					if(id >= 3){
						id = 0
					}else{
						id = id + 1
					}
					turno = false
				}else{
					for !carta_robada{			//Mientras no hayan robado una carta
						resp := <- espera
						if(resp == "Robar_carta"){			//Accion de robar una carta
							tablero.RobarCarta(t.Mazo,jugador.(jugadores.Jugador).Mano)		//Obtenemos la carta del mazo y se la damos al jugador
							carta_robada = true
							wait <- false
						}else if(resp == "Robar_carta_descartes"){
							if(t.Descartes.Size() > 0){
								tablero.RobarDescartes(t.Descartes,jugador.(jugadores.Jugador).Mano)	//En caso de que haya, robamos la carta del mazo de descartes y se la damos al jugador
								carta_robada = true
							}else{
								fmt.Println("Error, no hay cartas en el descarte")
							}
							wait <- false
						}else if(resp == "Fin_partida"){			//Final de partida por si fuera necesario
							wait <- true
							partida = false
							turno = false
							carta_robada = true
							goto SALIR
						}else if(resp == "Mostrar_mano"){			//Comando para mostrar la mano
							fmt.Println("Mostrando mano: ")
							cartas.MostrarMano(jugador.(jugadores.Jugador).Mano)		//Función que muestra la mano del jugador actual
							wait <- false
						}else if(resp == "Mostrar_tablero"){							//Comando para mostrar el tablero
							tablero.MostrarTablero(t);									//Función para mostrar el tablero
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
						fmt.Println("¿Que carta desea devolver?")						//En caso de querer devolver una carta
						fmt.Scanln(&input)												//El usuario deberá de introducir el ID necesario
						i_input,_:= strconv.Atoi(input)
						fmt.Println("Has introducido: ", i_input)
						aux := jugador.(jugadores.Jugador).Mano.Size()
						for i_input > aux {
							fmt.Println("Valor no valido, introduzca una carta correcta")
							fmt.Scanln(&input)												//El usuario deberá de introducir el ID necesario
							i_input,_ = strconv.Atoi(input)
							fmt.Println("Has introducido: ", i_input)
							aux := jugador.(jugadores.Jugador).Mano.Size()
							fmt.Println(aux)
						}
						tablero.FinTurno(t.Mazo,jugador.(jugadores.Jugador).Mano,t.Descartes,i_input)	//Y esta función colocará esa carta en el mazo de descartes
						if(jugador.(jugadores.Jugador).Mano.Size() == 0){								//En caso de no contar con más cartas terminará la partida
							wait <- true
							partida = false
							turno = false
						}else{																			//Y en caso contrario pasaremos al turno del siguiente jugador
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
					}else if (resp == "Colocar_combinacion"){				//Comando para colocar una nueva combinación en el tablero
						if(ab[id] == false){
							fmt.Println("No puedes colocar una carta porque no has abierto")		//En caso de no abrir da error
						}else{
							lista := list.New()
							combinacion := doublylinkedlist.New()
							var vector [13]int;
							for i:=0; i < 13; i++{
								vector[i] = -1;
							}
							cont := 0
							cont_a := 0
							fmt.Println("Indice los trios a probar")					//Hace una petición de los id de los trios que queremos comprobar
							fmt.Scanln(&input)
							input_V , _ := strconv.Atoi(input)
							for input != "FIN"{											//Hasta que no introduzca FIN no termina de añadir nuevas combinaciones
								if(input == "END"){										//Hasta que no introduzca END no termina de añadir valores a las nuevas combinaciones
									fmt.Println("Comprobando los valores...")
									for j := cont_a; j < cont; j++{						//En caso de ser haber añadido las cartas, ha sido guardado el id en un vector j que recorremos gracias a los
										carta, _ := jugador.(jugadores.Jugador).Mano.Get(vector[j])	//contadores cont_a y cont
										combinacion.Add(carta.(cartas.Carta))			//Añadiendo las combinaciones necesarias para a nuestra lista de nuevas combinaciones
									}
									//fmt.Printf("Tipo de combinacion: %T\n", combinacion)
									cartas.MostrarMano(combinacion)
									if(tablero.TrioValido(combinacion) || tablero.EscaleraValida(combinacion)){			//Si la combinación es valida, la añadimos a la lista definitiva
											// crea una copia de la lista original
										copia := doublylinkedlist.New()
										for e:= 0; e < combinacion.Size(); e++{								
											valor,_ := combinacion.Get(e)
											copia.Add(valor)															//Creamos copia
										}				
										lista.PushBack(copia)															//Y la añadimos a la lista
										
										
										/*for e := lista.Front(); e != nil; e = e.Next() {
											miLista := e.Value.(*doublylinkedlist.List)
											//fmt.Println(e.Value)
											//p,_ := miLista.Get(0)
											//fmt.Println(p)
											cartas.MostrarMano(miLista)
										}*/
										combinacion.Clear()
									}else{																		//Sino se elimina de la copia de la lista anterior por si se busca introducir una nueva combinación
										fmt.Println("Combinacion no valida,intentelo de nuevo")
										combinacion.Clear()
										cont = cont_a
									}
								}else if input_V >= 0 && input_V < 15{								//En caso de que se trate de un examen de 0 a 15 lo consideramos una carta, para ello la añadimos al vector
									i_input,_ := strconv.Atoi(input)
										for i := 0; i < cont; i++{
											if vector[i] == i_input{
												fmt.Println("Carta ya introducida")
												cont = cont_a										//Si la carta ya ha sido introducida dará error
												goto COMP_VALOR_COL
											}
										} 
										fmt.Println("Valor ", i_input, "guardado para comprobacion")
										vector[cont] = i_input
										cont++														//Y en caso de ser necesario, añadimos el ID y modificamos los contadores
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
								for l := 0 ; l < cont ; l++{							//Al recibir la señal de FIN, añadimos las combinaciones que sean necesarias
									fmt.Println(cont, " ", l, " ", vector[l])
									fmt.Println(jugador.(jugadores.Jugador).Mano.Get(vector[l]))
									cartas.MostrarMano(jugador.(jugadores.Jugador).Mano)
									jugador.(jugadores.Jugador).Mano.Remove(vector[l])
									cartas.MostrarMano(jugador.(jugadores.Jugador).Mano)
									for k := l+1 ; k < cont; k++{
										if(vector[k] > vector[l]){
											vector[k] = vector[k] - 1
										}
									}
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
						
					}else if (resp == "Abrir"){						//El funcionamiento de abrir es similar al comentado antes, pero comprobando el total de puntos.
						lista := list.New()
						combinacion := doublylinkedlist.New()
						puntos := 0
						var vector [14]int;
						for i:=0; i < 14; i++{
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
									for i := cont_a; i < cont; i++{
										vector[i] = -1
									}
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
									fmt.Println(vector[l])
									fmt.Println(jugador.(jugadores.Jugador).Mano.Get(vector[l]))
									jugador.(jugadores.Jugador).Mano.Remove(vector[l])
									for k := l+1 ; k < cont; k++{
										if(vector[k] > vector[l]){
											vector[k] = vector[k] - 1
										}
									}
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
						if(jugador.(jugadores.Jugador).Mano.Size() == 0){
							wait <- true
							partida = false
							turno = false
						}else{
							wait <- false
						}
					}else if (resp == "Colocar_carta"){											//Si buscamos colocar una carta
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
							r := tablero.AnyadirCarta(l_aux,jugador.(jugadores.Jugador).Mano,&t,t_combinacion)		//Comprobamos si no es valida, lo es, o si lo es y nos devuelve joker
							if r != -1 {															//Si lo es la colocamos
								fmt.Println("Carta colocada con exito")
								if r == 1{															//Y si es necesario recibimos el Joker
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
								fmt.Println("no valido")											//Sino no hacemos nada
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
		
	}
	SALIR: 															//Al acabar la partida terminamos contando los puntos de los jugadores que no han cerrado
	
	fmt.Println("fin, recuento de puntos de la partida actual...")
	listaJFinal := doublylinkedlist.New()
	for i:= 0; i < listaJ.Size(); i++ {
		jugador,_ := listaJ.Get(i)
		j := jugador.(jugadores.Jugador)
		manoJugador := j.Mano
		puntos := j.P_tor
		for j:= 0; j < manoJugador.Size(); j++ {
			carta,_ := manoJugador.Get(j)
			numero := carta.(cartas.Carta).Valor
			if numero == 1 {
				puntos += 11
			} else if numero >= 2 && numero <= 9 {
				puntos += numero
			} else if numero == 0 {
				puntos += 20
			} else {
				puntos += 10
			}
		}
		j.P_tor = puntos
		listaJFinal.Add(j)
		fmt.Println("Puntos del jugador",jugador.(jugadores.Jugador).Id,":",j.P_tor)
	}
	return listaJFinal

}