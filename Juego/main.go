package main

import (
	//"bufio"
	//"container/list"
	//"math/rand"
	//"time"
	//"fmt"
	//"net"
	//"strings"
	//"strconv"
	//"os"
	"github.com/emirpasic/gods/lists/doublylinkedlist"
	"juego/partida"
)

func main(){
	
	primeraPartida := true
	ganador := false
	listaJtotal := doublylinkedlist.New()
	// realizar partidas hasta que todos los jugadores se pasen de 100 puntos menos uno (el ganador) 
	for !ganador {
		// nueva partida
		listaJ := partida.IniciarPartida()
		if primeraPartida { // se inicializa la lista de jugadores
			for i:= 0; i < listaJ.Size(); i++ {
				jugador,_ := listaJ.Get(i)
				jugador.(jugadores.Jugador).P_tor = 0
				listaJtotal.Add(jugador.(jugadores.Jugador))
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

		// si algún jugador se pasa de 100 puntos, obtiene la mayor puntuación (que no se haya pasado)
		listaAux := doublylinkedlist.New()
		for j:= 0; j < jMax.Size(); j++ {
			jugadorMax,_ := jMax.Get(i)
			for i:= 0; i < listaJtotal.Size(); i++ {
				jugadorPartida,_ := listaJtotal.Get(i)
				if jugadorMax.(jugadores.Jugador).Id == jugadorPartida.(jugadores.Jugador).Id {
					jugadorPartida.(jugadores.Jugador).P_tor += pMax
					break
				}
			}
		}

		// comprobar si hay ganador
		numPerdedores := 0
		for i:= 0; i < listaJtotal.Size(); i++ {
			jugador,_ := listaJtotal.Get(i)
			if jugador.(jugadores.Jugador).P_tor > 100 {
				numPerdedores++
			}
		}
		if numPerdedores == listaJtotal.Size() - 2 {
			ganador = true
		}
	}
}