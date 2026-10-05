package main

import (
	"context"
	"log"
	"sort"

	"github.com/NaisuG/Golden-Team-Buddy/internal/board"
	"github.com/NaisuG/Golden-Team-Buddy/internal/catalog"
	"github.com/NaisuG/Golden-Team-Buddy/internal/engine"
)

type App struct {
	ctx     context.Context
	catalog *catalog.Catalog
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	c, err := catalog.LoadEmbedded()
	if err != nil {
		log.Fatalf("no se pudo cargar el catálogo: %v", err)
	}
	a.catalog = c
}

// GenerateVariants devuelve hasta tres composiciones sugeridas a partir del tablero y la banca.
func (a *App) GenerateVariants(boardChampions []string, benchChampions []string) ([]engine.Variant, error) {
	brd := board.Board{Champions: boardChampions}
	bch := board.Bench{Champions: benchChampions}
	if err := board.Validate(a.catalog, brd, bch); err != nil {
		return nil, err
	}
	return engine.GenerateVariants(a.catalog, brd, bch, brd.Level()), nil
}

// ListChampions devuelve todos los campeones del set, ordenados por nombre.
func (a *App) ListChampions() []catalog.Champion {
	champions := make([]catalog.Champion, 0, len(a.catalog.Champions))
	for _, champ := range a.catalog.Champions {
		champions = append(champions, champ)
	}
	sort.Slice(champions, func(i, j int) bool {
		return champions[i].Name < champions[j].Name
	})
	return champions
}
