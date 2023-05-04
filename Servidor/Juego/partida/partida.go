package partida

import (
	//"bufio"
	"container/list"
	"encoding/json"

	//"math/rand"
	//"time"
	"fmt"
	//"net"
	"strconv"
	"strings" //DESCOMENTAR

	//"os"
	"Juego/cartas"
	"Juego/jugadores"

	"github.com/emirpasic/gods/lists/doublylinkedlist"

	//"juego/partida"
	"Juego/bot"
	"Juego/tablero"

	"github.com/olahol/melody"
)

type Partida struct {
	Jug *doublylinkedlist.List
}

func Add_jug(j jugadores.Jugador, p Partida) {
	p.Jug.Add(j)
}

type RespuestaDescarte struct {
	Emisor        string     `json:"emisor"`
	Receptor      string     `json:"receptor"`
	Tipo          string     `json:"tipo"`
	Info          string     `json:"info"`
	Descartes     []string   `json:"descartes"`
	Combinaciones [][]string `json:"combinaciones"`
	Turno         string     `json:"turno"`
	Abrir         string     `json:"abrir"`
	Ganador       string     `json:"ganador"`
}

// func inicio_turno(espera chan string, wait chan bool) { //COMENTADO
func inicio_turno(espera chan string, wait chan bool, canalPartida chan string) { //DESCOMENTAR
	fin := false
	for !fin {
		var input string
		fmt.Println("¿Que acción desea hacer?")
		//fmt.Scanln(&input) //COMENTADO
		input = <-canalPartida //DESCOMENTAR
		if input == "Fin_partida" {
			fin = true
			fmt.Println("FINAL")
		} else if input == "Pausar" {
			fin = true
			fmt.Println("FINAL")
		}
		espera <- input
		fin = <-wait

	}

}

// func IniciarPartida() *doublylinkedlist.List { //COMENTADO
func IniciarPartida(idPartida string, canalPartida chan string, estabaPausada bool, es_bot []bool, ws *melody.Melody) *doublylinkedlist.List { //DESCOMENTAR
	//jugad, err := strconv.Atoi(os.Args[1])
	//torn, err := strconv.Atoi(os.Args[2])
	//bots, err := strconv.Atoi(os.Args[3])
	fmt.Println("Iniciando")
	input := ""
	parametrosPartida := <-canalPartida //DESCOMENTAR

	//separar los parametros por el caracter ","
	param := strings.Split(parametrosPartida, ",") //DESCOMENTAR
	numJugad, _ := strconv.Atoi(param[0])          //DESCOMENTAR
	listaJ := doublylinkedlist.New()
	//fmt.Println("Numero de jugadores: ", numJugad)
	//var ab [3]bool //COMENTADO

	for i := 0; i < len(es_bot); i++ { //DESCOMENTAR
		if es_bot[i] {
			numJugad = 4
			break
		}
	}
	ab := make([]bool, numJugad) //DESCOMENTAR

	t := tablero.Tablero{doublylinkedlist.New(), doublylinkedlist.New(), list.New()} //DESCOMENTAR

	if !estabaPausada {
		fmt.Println("Partida creada")
		t = tablero.IniciarTablero() //función de inicio de tablero para la partida	//DESCOMENTAR

		for i := 0; i < numJugad; i++ { //Inicio de los jugadores DESCOMENTAR
			//for i := 0; i < 3; i++ { //Inicio de los jugadores //COMENTADO
			jugador := jugadores.CrearJugador(i, t.Mazo)
			listaJ.Add(jugador)
			ab[i] = false
		}
	} else { //DESCOMENTAR
		fmt.Println("Partida reanudada")

		// llegan las cartas del mazo del tablero por el canal, hay que guardarlas en t.Mazo
		respuesta := <-canalPartida
		for respuesta != "Fin_mazo" {
			// separar el string por el caracter ","
			V_P_C := strings.Split(respuesta, ",")
			valor, _ := strconv.Atoi(V_P_C[0])
			palo, _ := strconv.Atoi(V_P_C[1])
			color, _ := strconv.Atoi(V_P_C[2])
			carta := cartas.Carta{valor, palo, color}
			t.Mazo.Add(carta)
			respuesta = <-canalPartida
		}
		fmt.Println("Mazo recibido")

		// llegan las cartas de descartes del tablero por el canal, hay que guardarlas en t.Descartes
		respuesta = <-canalPartida
		for respuesta != "Fin_descartes" {
			// separar el string por el caracter ","
			V_P_C := strings.Split(respuesta, ",")
			valor, _ := strconv.Atoi(V_P_C[0])
			palo, _ := strconv.Atoi(V_P_C[1])
			color, _ := strconv.Atoi(V_P_C[2])
			carta := cartas.Carta{valor, palo, color}
			t.Descartes.Add(carta)
			respuesta = <-canalPartida
		}
		fmt.Println("Descartes recibidos")

		// llegan las cartas de las combinaciones del tablero por el canal, hay que guardarlas en t.Combinaciones (lista de listas)
		respuesta = <-canalPartida
		for respuesta != "Fin_combinaciones" {
			comb := doublylinkedlist.New()
			for respuesta != "Fin_combinacion" {
				// separar el string por el caracter ","
				V_P_C := strings.Split(respuesta, ",")
				valor, _ := strconv.Atoi(V_P_C[0])
				palo, _ := strconv.Atoi(V_P_C[1])
				color, _ := strconv.Atoi(V_P_C[2])
				carta := cartas.Carta{valor, palo, color}
				comb.Add(carta)
				respuesta = <-canalPartida
			}
			t.Combinaciones.PushBack(comb)
			respuesta = <-canalPartida
		}
		fmt.Println("Combinaciones recibidas")

		// llegan las cartas de los jugadores por el canal, hay que guardarlas en jugadores.Jugador.Mano
		for i := 0; i < numJugad; i++ {
			respuesta = <-canalPartida
			mano := doublylinkedlist.New()
			for respuesta != "Fin_mano" {
				// separar el string por el caracter ","
				V_P_C := strings.Split(respuesta, ",")
				valor, _ := strconv.Atoi(V_P_C[0])
				palo, _ := strconv.Atoi(V_P_C[1])
				color, _ := strconv.Atoi(V_P_C[2])
				carta := cartas.Carta{valor, palo, color}
				mano.Add(carta)
				respuesta = <-canalPartida
			}
			j := jugadores.Jugador{i, mano, 0}
			listaJ.Add(j)
		}
		fmt.Println("Manos recibidas")

		// llega si han abierto o no los jugadores por el canal, hay que guardarlos en ab
		for i := 0; i < numJugad; i++ {
			respuesta = <-canalPartida
			if respuesta == "si" {
				ab[i] = true
			} else {
				ab[i] = false
			}
		}
		fmt.Println("Ya empieza la partida")
	} //DESCOMENTAR

	espera := make(chan string)
	wait := make(chan bool)
	//go inicio_turno(espera, wait) //Inicio de la escucha a la terminal //COMENTADO
	go inicio_turno(espera, wait, canalPartida) //Inicio de la escucha a la terminal //DESCOMENTAR

	partida := true
	turno := true
	carta_robada := false
	id := 0

	for partida { //Mientras sigamos en la partida
		fmt.Println("Turno del jugador ", id)
		jugador, err := listaJ.Get(id) //Inicio de los turnos
		turno = true                   //Ponemos turno a true porque seguimos en un turno
		carta_robada = false           //Y la carta robada a false para limitar las acciones hasta que robe una carta
		if err {
			//COMENTADO
			/*ab[0] = false //PRUEBAS
			carta_robada = false
			jugador.(jugadores.Jugador).Mano.Clear()
			carta := cartas.Carta{0, 4, 1}
			jugador.(jugadores.Jugador).Mano.Add(carta)
			carta = cartas.Carta{12, 1, 1}
			jugador.(jugadores.Jugador).Mano.Add(carta)
			carta = cartas.Carta{11, 1, 1}
			jugador.(jugadores.Jugador).Mano.Add(carta)
			carta = cartas.Carta{10, 1, 1}
			jugador.(jugadores.Jugador).Mano.Add(carta)
			carta = cartas.Carta{9, 1, 1}
			jugador.(jugadores.Jugador).Mano.Add(carta)
			carta = cartas.Carta{9, 1, 1}
			jugador.(jugadores.Jugador).Mano.Add(carta)
			carta = cartas.Carta{8, 1, 1}
			jugador.(jugadores.Jugador).Mano.Add(carta)
			carta = cartas.Carta{13, 1, 1}
			jugador.(jugadores.Jugador).Mano.Add(carta)
			carta = cartas.Carta{5, 3, 1}
			jugador.(jugadores.Jugador).Mano.Add(carta)
			carta = cartas.Carta{4, 3, 1}
			jugador.(jugadores.Jugador).Mano.Add(carta)
			carta = cartas.Carta{3, 3, 1}
			jugador.(jugadores.Jugador).Mano.Add(carta)
			carta = cartas.Carta{6, 2, 1}
			jugador.(jugadores.Jugador).Mano.Add(carta)
			carta = cartas.Carta{7, 2, 1}
			jugador.(jugadores.Jugador).Mano.Add(carta)*/

			for turno { //Mientras nos encontremos en un turno

				//COMENTADO
				if es_bot[id] {
					bot.Bot_En_Funcionamiento(t, jugador, ab[id])

					var RD RespuestaDescarte
					RD.Emisor = "Servidor"
					RD.Receptor = "todos"
					RD.Tipo = "Descarte"

					// devolver descartes y combinaciones
					// recorrer el mazo de descartes y pasar cada componente a string
					for i := 0; i < t.Descartes.Size(); i++ { //DESCOMENTAR todo el for
						carta, _ := t.Descartes.Get(i)
						carta2 := carta.(cartas.Carta)
						cartaString := strconv.Itoa(carta2.Valor) + "," + strconv.Itoa(carta2.Palo) + "," + strconv.Itoa(carta2.Color)
						RD.Descartes = append(RD.Descartes, cartaString)
					}

					// recorrer las combinaciones y pasar cada componente a string
					for e := t.Combinaciones.Front(); e != nil; e = e.Next() { //DESCOMENTAR todo el for
						combinacion := e.Value.(*doublylinkedlist.List)
						var comb []string
						for j := 0; j < combinacion.Size(); j++ {
							carta, _ := combinacion.Get(j)
							carta2 := carta.(cartas.Carta)
							cartaString := strconv.Itoa(carta2.Valor) + "," + strconv.Itoa(carta2.Palo) + "," + strconv.Itoa(carta2.Color)
							comb = append(comb, cartaString)
						}
						RD.Combinaciones = append(RD.Combinaciones, comb)
					}

					// Devolver siguiente turno y si ha abierto, si hay ganador devolverlo
					if jugador.(jugadores.Jugador).Mano.Size() == 0 {
						RD.Ganador = strconv.Itoa(id)
						wait <- true
						partida = false
						turno = false
					} else { //Y en caso contrario pasaremos al turno del siguiente jugador
						//if id >= 3 { //COMENTADO
						canalPartida <- "no"
						if id >= numJugad-1 { //DESCOMENTAR
							id = 0
						} else {
							id = id + 1
						}
						RD.Turno = strconv.Itoa(id)
						if ab[id] {
							RD.Abrir = "si"
						} else {
							RD.Abrir = "no"
						}
						turno = false
						wait <- false
					}

					msg, _ := json.MarshalIndent(&RD, "", "\t")
					//basura := <-espera
					//fmt.Println(basura)
					ws.BroadcastFilter(msg, func(q *melody.Session) bool { //Envia la información a todos con la misma url
						return q.Request.URL.Path == "/api/ws/partida/"+idPartida
					})
				} else {
					for !carta_robada { //Mientras no hayan robado una carta
						resp := <-espera
						if resp == "Robar_carta" { //Accion de robar una carta
							//tablero.RobarCarta(t.Mazo, jugador.(jugadores.Jugador).Mano)  //COMENTADO
							c := tablero.RobarCarta(t.Mazo, jugador.(jugadores.Jugador).Mano) //Obtenemos la carta del mazo y se la damos al jugador //DESCOMENTAR
							carta_robada = true
							canalPartida <- strconv.Itoa(c.Valor) + "," + strconv.Itoa(c.Palo) + "," + strconv.Itoa(c.Color) //DESCOMENTAR
							wait <- false
						} else if resp == "Robar_carta_descartes" {
							if t.Descartes.Size() > 0 {
								//tablero.RobarDescartes(t.Descartes, jugador.(jugadores.Jugador).Mano) //COMENTADO
								c := tablero.RobarDescartes(t.Descartes, jugador.(jugadores.Jugador).Mano) //En caso de que haya, robamos la carta del mazo de descartes y se la damos al jugador //DESCOMENTAR
								carta_robada = true
								canalPartida <- strconv.Itoa(c.Valor) + "," + strconv.Itoa(c.Palo) + "," + strconv.Itoa(c.Color) //DESCOMENTAR
							} else {
								fmt.Println("Error, no hay cartas en el descarte")
								canalPartida <- "Error, no hay cartas en el descarte" //DESCOMENTAR
							}
							wait <- false
						} else if resp == "Fin_partida" { //Final de partida por si fuera necesario
							canalPartida <- "fin" //DESCOMENTAR
							wait <- true
							partida = false
							turno = false
							carta_robada = true
							goto SALIR
						} else if resp == "Pausar" {
							pausar(t, canalPartida, listaJ, ab) //DESCOMENTAR
							wait <- true
							partida = false
							turno = false
							carta_robada = true
							goto SALIR
						} else if resp == "Mostrar_mano" { //Comando para mostrar la mano
							fmt.Println("Mostrando mano: ")
							cartas.MostrarMano(jugador.(jugadores.Jugador).Mano) //Función que muestra la mano del jugador actual

							// recorrer jugador.(jugadores.Jugador).Mano y pasar cada componente a string
							for i := 0; i < jugador.(jugadores.Jugador).Mano.Size(); i++ { //DESCOMENTAR todo el for
								carta, _ := jugador.(jugadores.Jugador).Mano.Get(i)
								carta2 := carta.(cartas.Carta)
								cartaString := strconv.Itoa(carta2.Valor) + "," + strconv.Itoa(carta2.Palo) + "," + strconv.Itoa(carta2.Color)
								canalPartida <- cartaString
							}
							canalPartida <- "fin" //DESCOMENTAR

							wait <- false
						} else if resp == "Mostrar_manos" { //Comando para mostrar las manos de todos los jugadores
							// recorrer la lista de jugadores y pasar cada componente a string
							for i := 0; i < listaJ.Size(); i++ {
								jugador, _ := listaJ.Get(i)
								jug := jugador.(jugadores.Jugador)
								for j := 0; j < jug.Mano.Size(); j++ {
									carta, _ := jug.Mano.Get(j)
									carta2 := carta.(cartas.Carta)
									cartaString := strconv.Itoa(carta2.Valor) + "," + strconv.Itoa(carta2.Palo) + "," + strconv.Itoa(carta2.Color)
									canalPartida <- cartaString
								}
								canalPartida <- "finJ"
							}
							canalPartida <- "fin"

							wait <- false
						} else if resp == "Mostrar_tablero" { //Comando para mostrar el tablero
							tablero.MostrarTablero(t) //Función para mostrar el tablero
							// recorrer el mazo y pasar cada componente a string
							for i := 0; i < t.Mazo.Size(); i++ { //DESCOMENTAR todo el for
								carta, _ := t.Mazo.Get(i)
								carta2 := carta.(cartas.Carta)
								cartaString := strconv.Itoa(carta2.Valor) + "," + strconv.Itoa(carta2.Palo) + "," + strconv.Itoa(carta2.Color)
								canalPartida <- cartaString
							}
							canalPartida <- "fin" //DESCOMENTAR

							// recorrer el mazo de descartes y pasar cada componente a string
							for i := 0; i < t.Descartes.Size(); i++ { //DESCOMENTAR todo el for
								carta, _ := t.Descartes.Get(i)
								carta2 := carta.(cartas.Carta)
								cartaString := strconv.Itoa(carta2.Valor) + "," + strconv.Itoa(carta2.Palo) + "," + strconv.Itoa(carta2.Color)
								canalPartida <- cartaString
							}
							canalPartida <- "fin" //DESCOMENTAR

							// recorrer las combinaciones y pasar cada componente a string
							for e := t.Combinaciones.Front(); e != nil; e = e.Next() { //DESCOMENTAR todo el for
								combinacion := e.Value.(*doublylinkedlist.List)
								for j := 0; j < combinacion.Size(); j++ {
									carta, _ := combinacion.Get(j)
									carta2 := carta.(cartas.Carta)
									cartaString := strconv.Itoa(carta2.Valor) + "," + strconv.Itoa(carta2.Palo) + "," + strconv.Itoa(carta2.Color)
									canalPartida <- cartaString
								}
								canalPartida <- "finC" //DESCOMENTAR
							}
							canalPartida <- "fin" //DESCOMENTAR

							wait <- false
						} else {
							fmt.Println("Error, primero tienes que robar una carta")
							canalPartida <- "Error, primero tienes que robar una carta" //DESCOMENTAR
							wait <- false
						}
					}
					resp := <-espera
					if resp == "Mostrar_mano" {
						fmt.Println("Mostrando mano: ")
						cartas.MostrarMano(jugador.(jugadores.Jugador).Mano)

						// recorrer jugador.(jugadores.Jugador).Mano y pasar cada componente a string
						for i := 0; i < jugador.(jugadores.Jugador).Mano.Size(); i++ { //DESCOMENTAR todo el for
							carta, _ := jugador.(jugadores.Jugador).Mano.Get(i)
							carta2 := carta.(cartas.Carta)
							cartaString := strconv.Itoa(carta2.Valor) + "," + strconv.Itoa(carta2.Palo) + "," + strconv.Itoa(carta2.Color)
							canalPartida <- cartaString
						}
						canalPartida <- "fin" //DESCOMENTAR
						wait <- false
					} else if resp == "Descarte" {
						canalPartida <- "Ok" //DESCOMENTAR
						var input string
						fmt.Println("¿Que carta desea devolver?") //En caso de querer devolver una carta
						//fmt.Scanln(&input)                        //El usuario deberá de introducir el ID necesario //COMENTADO
						input = <-canalPartida //DESCOMENTAR
						i_input, _ := strconv.Atoi(input)
						fmt.Println("Has introducido: ", i_input)
						aux := jugador.(jugadores.Jugador).Mano.Size()

						if i_input > aux { //DESCOMENTAR el if entero
							canalPartida <- "Valor no valido, introduzca una carta correcta"
						} else {
							canalPartida <- "Ok"
						}
						/*for i_input > aux {
							fmt.Println("Valor no valido, introduzca una carta correcta")
							//fmt.Scanln(&input) //El usuario deberá de introducir el ID necesario //COMENTADO
							input = <- canalPartida
							i_input, _ = strconv.Atoi(input)
							fmt.Println("Has introducido: ", i_input)
							aux := jugador.(jugadores.Jugador).Mano.Size()
							fmt.Println(aux)

							if i_input > aux {
								canalPartida <- "Valor no valido, introduzca una carta correcta"
							} else {
								canalPartida <- "Ok"
							}
						}*/

						tablero.FinTurno(t.Mazo, jugador.(jugadores.Jugador).Mano, t.Descartes, i_input) //Y esta función colocará esa carta en el mazo de descartes
						// devolver descartes y combinaciones
						// recorrer el mazo de descartes y pasar cada componente a string
						for i := 0; i < t.Descartes.Size(); i++ { //DESCOMENTAR todo el for
							carta, _ := t.Descartes.Get(i)
							carta2 := carta.(cartas.Carta)
							cartaString := strconv.Itoa(carta2.Valor) + "," + strconv.Itoa(carta2.Palo) + "," + strconv.Itoa(carta2.Color)
							canalPartida <- cartaString
						}
						canalPartida <- "fin" //DESCOMENTAR

						// recorrer las combinaciones y pasar cada componente a string
						for e := t.Combinaciones.Front(); e != nil; e = e.Next() { //DESCOMENTAR todo el for
							combinacion := e.Value.(*doublylinkedlist.List)
							for j := 0; j < combinacion.Size(); j++ {
								carta, _ := combinacion.Get(j)
								carta2 := carta.(cartas.Carta)
								cartaString := strconv.Itoa(carta2.Valor) + "," + strconv.Itoa(carta2.Palo) + "," + strconv.Itoa(carta2.Color)
								canalPartida <- cartaString
							}
							canalPartida <- "finC" //DESCOMENTAR
						}
						canalPartida <- "fin" //DESCOMENTAR

						// Devolver siguiente turno y si ha abierto, si hay ganador devolverlo
						if jugador.(jugadores.Jugador).Mano.Size() == 0 {
							canalPartida <- "ganador" //En caso de no contar con más cartas terminará la partida
							canalPartida <- strconv.Itoa(id)
							wait <- true
							partida = false
							turno = false
						} else { //Y en caso contrario pasaremos al turno del siguiente jugador
							//if id >= 3 { //COMENTADO
							canalPartida <- "no"
							if id >= numJugad-1 { //DESCOMENTAR
								id = 0
							} else {
								id = id + 1
							}
							canalPartida <- strconv.Itoa(id)
							if ab[id] {
								canalPartida <- "si"
							} else {
								canalPartida <- "no"
							}
							turno = false
							wait <- false
						}
					} else if resp == "Mostrar_manos" { //Comando para mostrar las manos de todos los jugadores
						// recorrer la lista de jugadores y pasar cada componente a string
						for i := 0; i < listaJ.Size(); i++ {
							jugador, _ := listaJ.Get(i)
							jug := jugador.(jugadores.Jugador)
							for j := 0; j < jug.Mano.Size(); j++ {
								carta, _ := jug.Mano.Get(j)
								carta2 := carta.(cartas.Carta)
								cartaString := strconv.Itoa(carta2.Valor) + "," + strconv.Itoa(carta2.Palo) + "," + strconv.Itoa(carta2.Color)
								canalPartida <- cartaString
							}
							canalPartida <- "finJ"
						}
						canalPartida <- "fin"

						wait <- false
					} else if resp == "Mostrar_tablero" {
						tablero.MostrarTablero(t)

						// recorrer el mazo y pasar cada componente a string
						for i := 0; i < t.Mazo.Size(); i++ { //DESCOMENTAR todo el for
							carta, _ := t.Mazo.Get(i)
							carta2 := carta.(cartas.Carta)
							cartaString := strconv.Itoa(carta2.Valor) + "," + strconv.Itoa(carta2.Palo) + "," + strconv.Itoa(carta2.Color)
							canalPartida <- cartaString
						}
						canalPartida <- "fin" //DESCOMENTAR

						// recorrer el mazo de descartes y pasar cada componente a string
						for i := 0; i < t.Descartes.Size(); i++ { //DESCOMENTAR todo el for
							carta, _ := t.Descartes.Get(i)
							carta2 := carta.(cartas.Carta)
							cartaString := strconv.Itoa(carta2.Valor) + "," + strconv.Itoa(carta2.Palo) + "," + strconv.Itoa(carta2.Color)
							canalPartida <- cartaString
						}
						canalPartida <- "fin" //DESCOMENTAR

						// recorrer las combinaciones y pasar cada componente a string
						for e := t.Combinaciones.Front(); e != nil; e = e.Next() { //DESCOMENTAR todo el for
							combinacion := e.Value.(*doublylinkedlist.List)
							for j := 0; j < combinacion.Size(); j++ {
								carta, _ := combinacion.Get(j)
								carta2 := carta.(cartas.Carta)
								cartaString := strconv.Itoa(carta2.Valor) + "," + strconv.Itoa(carta2.Palo) + "," + strconv.Itoa(carta2.Color)
								canalPartida <- cartaString
							}
							canalPartida <- "finC" //DESCOMENTAR
						}
						canalPartida <- "fin" //DESCOMENTAR

						wait <- false
					} else if resp == "Colocar_combinacion" { //Comando para colocar una nueva combinación en el tablero
						if ab[id] == false {
							fmt.Println("No puedes colocar una carta porque no has abierto")    //En caso de no abrir da error
							canalPartida <- "No puedes colocar una carta porque no has abierto" //DESCOMENTAR
						} else {
							canalPartida <- "Ok" //DESCOMENTAR
							lista := list.New()
							combinacion := doublylinkedlist.New()
							var vector [13]int
							for i := 0; i < 13; i++ {
								vector[i] = -1
							}
							cont := 0
							cont_a := 0
							fmt.Println("Indice los trios a probar") //Hace una petición de los id de los trios que queremos comprobar
							//fmt.Scanln(&input) //COMENTADO
							input := <-canalPartida //DESCOMENTAR
							input_V, _ := strconv.Atoi(input)
							for input != "FIN" { //Hasta que no introduzca FIN no termina de añadir nuevas combinaciones
								if input == "END" { //Hasta que no introduzca END no termina de añadir valores a las nuevas combinaciones
									fmt.Println("Comprobando los valores...")
									for j := cont_a; j < cont; j++ { //En caso de ser haber añadido las cartas, ha sido guardado el id en un vector j que recorremos gracias a los
										carta, _ := jugador.(jugadores.Jugador).Mano.Get(vector[j]) //contadores cont_a y cont
										combinacion.Add(carta.(cartas.Carta))                       //Añadiendo las combinaciones necesarias para a nuestra lista de nuevas combinaciones
									}
									//fmt.Printf("Tipo de combinacion: %T\n", combinacion)
									cartas.MostrarMano(combinacion)
									if tablero.TrioValido(combinacion) || tablero.EscaleraValida(combinacion) { //Si la combinación es valida, la añadimos a la lista definitiva
										// crea una copia de la lista original
										copia := doublylinkedlist.New()
										for e := 0; e < combinacion.Size(); e++ {
											valor, _ := combinacion.Get(e)
											copia.Add(valor) //Creamos copia
										}
										lista.PushBack(copia) //Y la añadimos a la lista

										/*for e := lista.Front(); e != nil; e = e.Next() {
											miLista := e.Value.(*doublylinkedlist.List)
											//fmt.Println(e.Value)
											//p,_ := miLista.Get(0)
											//fmt.Println(p)
											cartas.MostrarMano(miLista)
										}*/
										combinacion.Clear()
										canalPartida <- "Ok" //DESCOMENTAR
									} else { //Sino se elimina de la copia de la lista anterior por si se busca introducir una nueva combinación
										fmt.Println("Combinacion no valida,intentelo de nuevo")
										combinacion.Clear()
										cont = cont_a
										canalPartida <- "Combinacion no valida,intentelo de nuevo" //DESCOMENTAR
									}
								} else if input_V >= 0 && input_V < 15 { //En caso de que se trate de un examen de 0 a 15 lo consideramos una carta, para ello la añadimos al vector
									i_input, _ := strconv.Atoi(input)
									for i := 0; i < cont; i++ {
										if vector[i] == i_input {
											fmt.Println("Carta ya introducida")
											cont = cont_a //Si la carta ya ha sido introducida dará error
											goto COMP_VALOR_COL
										}
									}
									fmt.Println("Valor ", i_input, "guardado para comprobacion")
									vector[cont] = i_input
									cont++ //Y en caso de ser necesario, añadimos el ID y modificamos los contadores
									fmt.Println(jugador.(jugadores.Jugador).Mano.Get(i_input))
								COMP_VALOR_COL:
								} else {
									fmt.Println("Comando erroneo")
								}
								//fmt.Scanln(&input) //COMENTADO
								input = <-canalPartida //DESCOMENTAR
							}
							if lista.Len() > 0 {
								//fmt.Println("Contenido de la lista:")
								//fmt.Println(lista)
								/*for e := lista.Back(); e != nil; e = e.Next() {
									miLista := e.Value.(*doublylinkedlist.List)
									cartas.MostrarMano(miLista)
								}*/
								fmt.Println(cont)
								for l := 0; l < cont; l++ { //Al recibir la señal de FIN, añadimos las combinaciones que sean necesarias
									fmt.Println(cont, " ", l, " ", vector[l])
									fmt.Println(jugador.(jugadores.Jugador).Mano.Get(vector[l]))
									cartas.MostrarMano(jugador.(jugadores.Jugador).Mano)
									jugador.(jugadores.Jugador).Mano.Remove(vector[l])
									cartas.MostrarMano(jugador.(jugadores.Jugador).Mano)
									for k := l + 1; k < cont; k++ {
										if vector[k] > vector[l] {
											vector[k] = vector[k] - 1
										}
									}
								}

								fmt.Println("lista antes de añadir al tablero")
								aaa := lista
								for elem := aaa.Front(); elem != nil; elem = elem.Next() {
									fmt.Println(elem.Value)
								}

								tablero.AnyadirCombinaciones(t, lista)
								canalPartida <- "Ok" //DESCOMENTAR
							} else {
								canalPartida <- "No se ha introducido ninguna combinacion" //DESCOMENTAR
							}
						}
						if jugador.(jugadores.Jugador).Mano.Size() == 0 {
							canalPartida <- "ganador"        //DESCOMENTAR
							canalPartida <- strconv.Itoa(id) //DESCOMENTAR
							wait <- true
							partida = false
							turno = false
						} else {
							canalPartida <- "no ganador" //DESCOMENTAR
							wait <- false
						}

					} else if resp == "Abrir" { //El funcionamiento de abrir es similar al comentado antes, pero comprobando el total de puntos.
						canalPartida <- "Ok" //DESCOMENTAR
						lista := list.New()
						combinacion := doublylinkedlist.New()
						puntos := 0
						var vector [14]int
						for i := 0; i < 14; i++ {
							vector[i] = -1
						}
						cont := 0
						cont_a := 0
						fmt.Println("Indice los trios a probar")
						//fmt.Scanln(&input) //COMENTADO
						input := <-canalPartida //DESCOMENTAR
						input_V, _ := strconv.Atoi(input)
						for input != "FIN" {
							if input == "END" {
								fmt.Println("Comprobando los valores...")
								for j := cont_a; j < cont; j++ {
									carta, _ := jugador.(jugadores.Jugador).Mano.Get(vector[j])
									combinacion.Add(carta.(cartas.Carta))
								}
								//fmt.Printf("Tipo de combinacion: %T\n", combinacion)
								cartas.MostrarMano(combinacion)
								if tablero.TrioValido(combinacion) || tablero.EscaleraValida(combinacion) {
									puntos = puntos + tablero.SumaCartas(combinacion)
									fmt.Println("Combinación valida, tienes ", puntos, "puntos")
									// crea una copia de la lista original
									copia := doublylinkedlist.New()
									for e := 0; e < combinacion.Size(); e++ {
										valor, _ := combinacion.Get(e)
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
									canalPartida <- "Ok" //DESCOMENTAR
								} else {
									fmt.Println("Combinacion no valida,intentelo de nuevo")
									combinacion.Clear()
									for i := cont_a; i < cont; i++ {
										vector[i] = -1
									}
									cont = cont_a
									canalPartida <- "Combinacion no valida,intentelo de nuevo"
								}
							} else if input_V >= 0 && input_V < 15 {
								i_input, _ := strconv.Atoi(input)
								for i := 0; i < cont; i++ {
									if vector[i] == i_input {
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
							} else {
								fmt.Println("Comando erroneo")
							}
							//fmt.Scanln(&input) //COMENTADO
							input = <-canalPartida //DESCOMENTAR
						}
						if puntos >= 51 {
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
								for l := 0; l < cont; l++ {
									fmt.Println(vector[l])
									fmt.Println(jugador.(jugadores.Jugador).Mano.Get(vector[l]))
									jugador.(jugadores.Jugador).Mano.Remove(vector[l])
									for k := l + 1; k < cont; k++ {
										if vector[k] > vector[l] {
											vector[k] = vector[k] - 1
										}
									}
								}

								fmt.Println("lista antes de añadir al tablero")
								aaa := lista
								for elem := aaa.Front(); elem != nil; elem = elem.Next() {
									fmt.Println(elem.Value)
								}

								tablero.AnyadirCombinaciones(t, lista)
								fmt.Println("Felicidades, has abierto")
								canalPartida <- "Ok" //DESCOMENTAR
								ab[id] = true
							} else {
								fmt.Println("La lista está vacía")
								canalPartida <- "La lista está vacía" //DESCOMENTAR
							}
						} else {
							//fmt.Println("No has conseguido suficientes puntos, y no has podido abrir") //COMENTADO
							canalPartida <- "No has conseguido suficientes puntos, y no has podido abrir" //DESCOMENTAR
							lista.Init()
						}
						if jugador.(jugadores.Jugador).Mano.Size() == 0 {
							canalPartida <- "ganador"
							canalPartida <- strconv.Itoa(id)
							wait <- true
							partida = false
							turno = false
						} else {
							canalPartida <- "no ganador"
							wait <- false
						}
					} else if resp == "Colocar_carta" { //Si buscamos colocar una carta
						if ab[id] == false {
							fmt.Println("No puedes colocar una carta porque no has abierto")
							canalPartida <- "No puedes colocar una carta porque no has abierto" //DESCOMENTAR
						} else {
							canalPartida <- "Ok" //DESCOMENTAR
							fmt.Println("¿En que combinación desea introducir su carta?")
							//fmt.Scanln(&input) //COMENTADO
							input = <-canalPartida //DESCOMENTAR
							t_combinacion, _ := strconv.Atoi(input)
							fmt.Println("¿Que carta desea introducir?")
							//fmt.Scanln(&input) //COMENTADO
							input = <-canalPartida //DESCOMENTAR
							i_carta, _ := strconv.Atoi(input)
							l_aux := doublylinkedlist.New()
							l_aux.Add(jugador.(jugadores.Jugador).Mano.Get(i_carta))
							r := tablero.AnyadirCarta(l_aux, jugador.(jugadores.Jugador).Mano, &t, t_combinacion) //Comprobamos si no es valida, lo es, o si lo es y nos devuelve joker
							if r != -1 {                                                                          //Si lo es la colocamos
								fmt.Println("Carta colocada con exito")
								if r == 1 { //Y si es necesario recibimos el Joker
									canalPartida <- "joker"
									value := cartas.Carta{0, 4, 1}
									jugador.(jugadores.Jugador).Mano.Add(value)
									canalPartida <- strconv.Itoa(value.Valor) + "," + strconv.Itoa(value.Palo) + "," + strconv.Itoa(value.Color)
								} else {
									if jugador.(jugadores.Jugador).Mano.Size() == 0 {
										canalPartida <- "ganador" //En caso de no contar con más cartas terminará la partida
										canalPartida <- strconv.Itoa(id)
										wait <- true
										partida = false
										turno = false
									} else {
										canalPartida <- "Ok" //DESCOMENTAR
									}
								}
							} else {
								fmt.Println("no valido")    //Sino no hacemos nada
								canalPartida <- "no valido" //DESCOMENTAR
							}

						}
						if !(jugador.(jugadores.Jugador).Mano.Size() == 0) {
							wait <- false
						}
					} else if resp == "Fin_partida" {
						canalPartida <- "fin" //DESCOMENTAR
						wait <- true
						partida = false
						turno = false
					} else if resp == "Pausar" {
						pausar(t, canalPartida, listaJ, ab) //DESCOMENTAR
						wait <- true
						partida = false
						turno = false
					} else {
						fmt.Println("Operacion no valida")
						wait <- false
					}
				} //COMENTADO
			}
		}

	}
SALIR: //Al acabar la partida terminamos contando los puntos de los jugadores que no han cerrado

	fmt.Println("fin, recuento de puntos de la partida actual...")
	listaJFinal := doublylinkedlist.New()
	for i := 0; i < listaJ.Size(); i++ {
		jugador, _ := listaJ.Get(i)
		j := jugador.(jugadores.Jugador)
		manoJugador := j.Mano
		puntos := j.P_tor
		for j := 0; j < manoJugador.Size(); j++ {
			carta, _ := manoJugador.Get(j)
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
		fmt.Println("Puntos del jugador", jugador.(jugadores.Jugador).Id, ":", j.P_tor)
	}
	return listaJFinal
}

func pausar(t tablero.Tablero, canalPartida chan string, listaJ *doublylinkedlist.List, ab []bool) {
	// recorrer el mazo y pasar cada componente a string
	for i := 0; i < t.Mazo.Size(); i++ {
		carta, _ := t.Mazo.Get(i)
		carta2 := carta.(cartas.Carta)
		cartaString := strconv.Itoa(carta2.Valor) + "," + strconv.Itoa(carta2.Palo) + "," + strconv.Itoa(carta2.Color)
		canalPartida <- cartaString
	}
	canalPartida <- "fin"

	// recorrer el mazo de descartes y pasar cada componente a string
	for i := 0; i < t.Descartes.Size(); i++ {
		carta, _ := t.Descartes.Get(i)
		carta2 := carta.(cartas.Carta)
		cartaString := strconv.Itoa(carta2.Valor) + "," + strconv.Itoa(carta2.Palo) + "," + strconv.Itoa(carta2.Color)
		canalPartida <- cartaString
	}
	canalPartida <- "fin"

	// recorrer las combinaciones y pasar cada componente a string
	for e := t.Combinaciones.Front(); e != nil; e = e.Next() {
		combinacion := e.Value.(*doublylinkedlist.List)
		for j := 0; j < combinacion.Size(); j++ {
			carta, _ := combinacion.Get(j)
			carta2 := carta.(cartas.Carta)
			cartaString := strconv.Itoa(carta2.Valor) + "," + strconv.Itoa(carta2.Palo) + "," + strconv.Itoa(carta2.Color)
			canalPartida <- cartaString
		}
		canalPartida <- "finC"
	}
	canalPartida <- "fin"

	// recorrer la lista de jugadores y pasar cada componente a string
	for i := 0; i < listaJ.Size(); i++ {
		jugador, _ := listaJ.Get(i)
		jug := jugador.(jugadores.Jugador)
		for j := 0; j < jug.Mano.Size(); j++ {
			carta, _ := jug.Mano.Get(j)
			carta2 := carta.(cartas.Carta)
			cartaString := strconv.Itoa(carta2.Valor) + "," + strconv.Itoa(carta2.Palo) + "," + strconv.Itoa(carta2.Color)
			canalPartida <- cartaString
		}
		canalPartida <- "finJ"
	}
	canalPartida <- "fin"

	// recorrer vector ab
	for i := 0; i < len(ab); i++ {
		if ab[i] == true {
			canalPartida <- "si"
		} else {
			canalPartida <- "no"
		}
	}
	canalPartida <- "fin"
}
