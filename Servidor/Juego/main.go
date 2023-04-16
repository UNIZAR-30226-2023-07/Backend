package main

import (
	"fmt"

	"Servidor/Juego/partida"
	"Servidor/Juego/torneo"
)

func main(){
	var input string
	fmt.Println("¿Que quieres jugar?")
	fmt.Scanln(&input)
	if (input == "partida") {
		partida.IniciarPartida()
	} else if (input == "torneo") {
		torneo.IniciarTorneo()
	}
}