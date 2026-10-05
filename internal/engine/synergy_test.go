package engine

import (
	"testing"

	"github.com/NaisuG/Golden-Team-Buddy/internal/catalog"
)

func TestSynergySharedTraitIsCloserThanNone(t *testing.T) {
	c := catalog.New()
	c.Traits["Anima"] = catalog.Trait{Key: "Anima", Name: "Anima"}
	c.Traits["Vanguard"] = catalog.Trait{Key: "Vanguard", Name: "Vanguard"}
	c.Traits["Sniper"] = catalog.Trait{Key: "Sniper", Name: "Sniper"}

	c.Champions["Illaoi"] = catalog.Champion{Key: "Illaoi", Traits: []string{"Anima", "Vanguard"}}
	c.Champions["Jinx"] = catalog.Champion{Key: "Jinx", Traits: []string{"Anima"}}
	c.Champions["Caitlyn"] = catalog.Champion{Key: "Caitlyn", Traits: []string{"Sniper"}}

	illaoi := ChampionVector(c, "Illaoi")
	jinx := ChampionVector(c, "Jinx")
	caitlyn := ChampionVector(c, "Caitlyn")

	sharedTrait := Synergy(illaoi, jinx)
	noSharedTrait := Synergy(illaoi, caitlyn)

	if sharedTrait <= noSharedTrait {
		t.Errorf("Illaoi-Jinx (comparten Anima) = %.2f, debería ser mayor que Illaoi-Caitlyn (nada en común) = %.2f",
			sharedTrait, noSharedTrait)
	}
}

const epsilon = 1e-9

func floatsEqual(a, b float64) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff < epsilon
}

func TestSynergyCases(t *testing.T) {
	identical := Vector{"Anima": 1, "Vanguard": 1}
	noOverlap := Vector{"Sniper": 1}
	empty := Vector{}

	cases := []struct {
		name     string
		a, b     Vector
		expected float64
	}{
		{name: "vectores idénticos", a: identical, b: identical, expected: 1.0},
		{name: "sin traits en común", a: identical, b: noOverlap, expected: 0.0},
		{name: "vector vacío", a: empty, b: identical, expected: 0.0},
	}

	for _, c := range cases {
		got := Synergy(c.a, c.b)
		if !floatsEqual(got, c.expected) {
			t.Errorf("%s: Synergy() = %f, quería %f", c.name, got, c.expected)
		}
	}
}
