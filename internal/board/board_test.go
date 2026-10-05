package board

import (
	"testing"

	"github.com/NaisuG/Golden-Team-Buddy/internal/catalog"
)

func testCatalog(keys ...string) *catalog.Catalog {
	c := catalog.New()
	for _, key := range keys {
		c.Champions[key] = catalog.Champion{Key: key}
	}
	return c
}

func TestValidateAcceptsValidBoard(t *testing.T) {
	c := testCatalog("A", "B", "C")
	if err := Validate(c, Board{Champions: []string{"A", "B"}}, Bench{Champions: []string{"C"}}); err != nil {
		t.Errorf("Validate() = %v, quería nil", err)
	}
}

func TestValidateRejectsInvalidInput(t *testing.T) {
	c := testCatalog("A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K")

	cases := []struct {
		name string
		brd  Board
		bch  Bench
	}{
		{name: "tablero lleno", brd: Board{Champions: []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K"}}},
		{name: "banca llena", bch: Bench{Champions: []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}}},
		{name: "campeón desconocido", brd: Board{Champions: []string{"Z"}}},
		{name: "repetido en tablero", brd: Board{Champions: []string{"A", "A"}}},
		{name: "repetido entre tablero y banca", brd: Board{Champions: []string{"A"}}, bch: Bench{Champions: []string{"A"}}},
	}

	for _, tc := range cases {
		if err := Validate(c, tc.brd, tc.bch); err == nil {
			t.Errorf("%s: Validate() = nil, quería un error", tc.name)
		}
	}
}
