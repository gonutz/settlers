package main

import (
	"fmt"
	"github.com/gonutz/settlers/game"
	"math/rand"
)

func main() {
	fmt.Println("Hello World!")
	game := game.New([]game.Color{game.Red, game.Blue, game.White}, rand.Int)
	fmt.Printf("initial: %#v\n", game)
	game.DealResources(5)
	fmt.Printf("rolled 5: %#v\n", game)
}
