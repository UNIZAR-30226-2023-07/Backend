package main

import (
	//"bufio"
	//"container/list"
	//"math/rand"
	//"time"
	"fmt"
	//"net"
	//"strings"
	//"strconv"
	//"os"

	"github.com/emirpasic/gods/lists/doublylinkedlist"
	"juego/partida"
	"juego/jugadores"
)

func main(){
	
	primeraPartida := true
	ganador := false
	listaJtotal := doublylinkedlist.New()
	// realizar partidas hasta que todos los jugadores se pasen de 100 puntos menos uno (el ganador) 
	for !ganador {
		// nueva partida
		fmt.Println("Nueva Partida")
		listaJ := partida.IniciarPartida()
		
		if primeraPartida { // se inicializa la lista de jugadores
			for i:= 0; i < listaJ.Size(); i++ {
				jugador,_ := listaJ.Get(i)
				j := jugador.(jugadores.Jugador)
				j.P_tor = 0
				listaJtotal.Add(j)
			}
			primeraPartida = false
		}

		// contar puntos
		pMax := 0
		jMax := doublylinkedlist.New()
		for i:= 0; i < listaJ.Size(); i++{
			jugador,_ := listaJ.Get(i)
			puntos := jugador.(jugadores.Jugador).P_tor
			
			if puntos > pMax && puntos <= 100 { // quedarse con la mayor puntuacion que no se pasa de 100
				pMax = puntos
			} else if puntos > 100 { // quedarse con los jugadores que se pasan de 100
				jMax.Add(jugador.(jugadores.Jugador))
			}
		}
		fmt.Println("la mayor puntuacion que no llega a 100 es",pMax)
		fmt.Println("jugadores que se pasan de 100:")
		for i:= 0; i < jMax.Size(); i++ {
			jugador,_ := jMax.Get(i)
			fmt.Println("id:",jugador.(jugadores.Jugador).Id,"puntos:",jugador.(jugadores.Jugador).P_tor)
		}

		// si algún jugador se pasa de 100 puntos, obtiene la mayor puntuación (que no se haya pasado)
		listaAux := doublylinkedlist.New()
		for i:= 0; i < listaJtotal.Size(); i++ {
			jugadorPartida,_ := listaJtotal.Get(i)
			jug := jugadorPartida.(jugadores.Jugador)
			sePasa := false
			for j:= 0; j < jMax.Size(); j++ {
				jugadorMax,_ := jMax.Get(j)
				if jugadorMax.(jugadores.Jugador).Id == jug.Id {
					jug.P_tor += pMax
					sePasa = true
					break
				}
			}
			if !sePasa { // si no se pasa de 100 puntos se suman solo los puntos de la partida
				jugador,_ := listaJ.Get(i)
				j := jugador.(jugadores.Jugador)
				jug.P_tor += j.P_tor
			}
			listaAux.Add(jug)
		}

		listaJtotal = listaAux
		fmt.Println("Recuento de puntos:")
		for i:= 0; i < listaJtotal.Size(); i++ {
			jugador,_ := listaJtotal.Get(i)
			fmt.Println("id:",jugador.(jugadores.Jugador).Id,"puntos:",jugador.(jugadores.Jugador).P_tor)
		}

		// comprobar si hay ganador
		numPerdedores := 0
		for i:= 0; i < listaJtotal.Size(); i++ {
			jugador,_ := listaJtotal.Get(i)
			if jugador.(jugadores.Jugador).P_tor > 100 {
				numPerdedores++
			}
		}
		if numPerdedores == listaJtotal.Size() - 1 {
			ganador = true
			fmt.Println("Hay ganador")
		}
	}
}