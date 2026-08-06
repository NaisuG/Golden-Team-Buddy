package board

import "fmt"

const MaxBenchSize = 9
const MaxBoardSize = 10

type Board struct {
	Champions []string
}

func (b *Board) Level() int {
	return len(b.Champions)
}

func (b *Board) Add(championKey string) error {
	if len(b.Champions) >= MaxBoardSize {
		return fmt.Errorf("tablero lleno: máximo %d campeones", MaxBoardSize)
	}
	b.Champions = append(b.Champions, championKey)
	return nil
}
