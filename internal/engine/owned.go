package engine

import "github.com/NaisuG/Golden-Team-Buddy/internal/board"

// Owned son los campeones que el jugador ya tiene, separados en tablero y banca.
type Owned struct {
	Board map[string]bool
	Bench map[string]bool
}

func NewOwned(brd board.Board, bch board.Bench) Owned {
	o := Owned{
		Board: make(map[string]bool, len(brd.Champions)),
		Bench: make(map[string]bool, len(bch.Champions)),
	}
	for _, key := range brd.Champions {
		o.Board[key] = true
	}
	for _, key := range bch.Champions {
		o.Bench[key] = true
	}
	return o
}

func (o Owned) Has(key string) bool {
	return o.Board[key] || o.Bench[key]
}
