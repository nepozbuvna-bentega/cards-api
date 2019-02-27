package main

import (
	"log"

	"github.com/nepozbuvna-bentega/cards-api/game"
)

func main() {
	blacks := []*game.BlackCard{game.NewBlackCard("black card <x>")}
	whites := []*game.WhiteCard{game.NewWhiteCard(map[string]string{"x": "Value"})}

	g := game.New(1, blacks, whites)
	game.ActionJoin("user-0")(g)

	log.Printf("%s", g)
}
