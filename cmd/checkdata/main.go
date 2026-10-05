package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/NaisuG/Golden-Team-Buddy/internal/board"
	"github.com/NaisuG/Golden-Team-Buddy/internal/catalog"
	"github.com/NaisuG/Golden-Team-Buddy/internal/engine"
)

func main() {
	c, err := catalog.LoadEmbedded()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Campeones cargados:", len(c.Champions))
	fmt.Println("Traits cargados:", len(c.Traits))

	brd := board.Board{Champions: []string{"Illaoi"}}
	variants := engine.GenerateVariants(c, brd, board.Bench{}, brd.Level())

	for i, v := range variants {
		fmt.Printf("\n%d) %v\n   %s\n", i+1, v.Champions, formatTraits(v.Traits))
		for j, sub := range v.Children {
			fmt.Printf("   %d%c) %v\n       %s\n", i+1, rune('a'+j), sub.Champions, formatTraits(sub.Traits))
		}
	}
}

func formatTraits(traits []engine.ActiveTrait) string {
	parts := make([]string, len(traits))
	for i, t := range traits {
		parts[i] = fmt.Sprintf("%d %s (%s)", t.Count, t.Name, t.Style)
	}
	return strings.Join(parts, ", ")
}
