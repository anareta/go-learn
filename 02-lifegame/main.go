package main

import (
	"fmt"
	"github.com/anareta/go-learn/02-lifegame/game"
	"time"
)

func main() {
	fmt.Println("Life Game !!")

	var lifeMap = game.CreateMap([][]bool{
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, true, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, true, false, false, false, false, false, false, false},
		{false, false, false, false, true, true, false, false, true, true, true, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false}})

	// ターン経過
	turn := 0
	for true {
		// main loop
		turn++
		fmt.Printf("【世代 : %v】\n", turn)
		fmt.Println(lifeMap.Display())
		fmt.Printf("\n\n")

		time.Sleep(time.Second * 1)
		lifeMap = lifeMap.Next()
	}
}
