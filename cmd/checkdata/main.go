package main

import (
	"fmt"
	"log"

	"github.com/NaisuG/Golden-Team-Buddy/internal/catalog"
	"github.com/NaisuG/Golden-Team-Buddy/internal/engine"
)

func main() {
	c, err := catalog.LoadJSON("scripts/etl")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Campeones cargados:", len(c.Champions))
	fmt.Println("Traits cargados:", len(c.Traits))

	current := []string{"Illaoi"}
	level := 1

	variants := engine.GenerateVariants(c, current, level)

	for i, v := range variants {
		fmt.Printf("\n%d) Total=%.2f  %v\n", i+1, v.Score.Total, v.Champions)
		for j, sub := range v.Children {
			fmt.Printf("   %d%c) Total=%.2f  %v\n", i+1, rune('a'+j), sub.Score.Total, sub.Champions)
		}
	}
}
