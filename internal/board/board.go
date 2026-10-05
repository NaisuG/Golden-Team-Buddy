package board

import (
	"fmt"

	"github.com/NaisuG/Golden-Team-Buddy/internal/catalog"
)

const (
	MaxBoardSize = 10
	MaxBenchSize = 9
)

type Board struct {
	Champions []string
}

func (b Board) Level() int {
	return len(b.Champions)
}

type Bench struct {
	Champions []string
}

// Validate revisa los límites de tablero y banca, y que no haya campeones desconocidos ni repetidos.
func Validate(c *catalog.Catalog, brd Board, bch Bench) error {
	if len(brd.Champions) > MaxBoardSize {
		return fmt.Errorf("el tablero tiene %d campeones, máximo %d", len(brd.Champions), MaxBoardSize)
	}
	if len(bch.Champions) > MaxBenchSize {
		return fmt.Errorf("la banca tiene %d campeones, máximo %d", len(bch.Champions), MaxBenchSize)
	}

	seen := make(map[string]bool, len(brd.Champions)+len(bch.Champions))
	for _, key := range append(append([]string{}, brd.Champions...), bch.Champions...) {
		if _, ok := c.Champions[key]; !ok {
			return fmt.Errorf("campeón desconocido: %q", key)
		}
		if seen[key] {
			return fmt.Errorf("campeón repetido: %q", key)
		}
		seen[key] = true
	}
	return nil
}
