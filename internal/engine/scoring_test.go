package engine

import (
	"testing"

	"github.com/NaisuG/Golden-Team-Buddy/internal/catalog"
)

func buildTestCatalog() *catalog.Catalog {
	c := catalog.New()
	c.Traits["Anima"] = catalog.Trait{
		Key:  "Anima",
		Name: "Anima",
		Breakpoints: []catalog.Breakpoint{
			{Style: "bronze", Min: 3},
			{Style: "gold", Min: 5},
		},
	}
	c.Champions["Illaoi"] = catalog.Champion{Key: "Illaoi", Traits: []string{"Anima"}}
	c.Champions["Jinx"] = catalog.Champion{Key: "Jinx", Traits: []string{"Anima"}}
	c.Champions["Briar"] = catalog.Champion{Key: "Briar", Traits: []string{"Anima"}}
	c.Champions["Aurora"] = catalog.Champion{Key: "Aurora", Traits: []string{"Anima"}}
	c.Champions["Meepsie"] = catalog.Champion{Key: "Meepsie", Traits: []string{"Anima"}}
	return c
}

func TestTraitStrengthBreakpoints(t *testing.T) {
	c := buildTestCatalog()

	cases := []struct {
		name     string
		team     []string
		expected float64
	}{
		{name: "2 Anima, sin breakpoint", team: []string{"Illaoi", "Jinx"}, expected: 0},
		{name: "3 Anima, bronce (tier 1)", team: []string{"Illaoi", "Jinx", "Briar"}, expected: 1},
		{name: "5 Anima, oro (tier 2)", team: []string{"Illaoi", "Jinx", "Briar", "Aurora", "Meepsie"}, expected: 4},
	}

	for _, tc := range cases {
		got := TraitStrength(c, tc.team)
		if got != tc.expected {
			t.Errorf("%s: TraitStrength() = %v; want %v", tc.name, got, tc.expected)
		}
	}
}

func TestAverageSynergy(t *testing.T) {
	c := buildTestCatalog()

	cases := []struct {
		name     string
		team     []string
		expected float64
	}{
		{name: "menos de 2 campeones", team: []string{"Illaoi"}, expected: 0},
		{name: "todos comparten Anima", team: []string{"Illaoi", "Jinx", "Briar"}, expected: 1},
	}

	for _, tc := range cases {
		got := AverageSynergy(c, tc.team)
		if !floatsEqual(got, tc.expected) {
			t.Errorf("%s: AverageSynergy() = %v; want %v", tc.name, got, tc.expected)
		}
	}
}

func TestComputeScoreStrongerBeatsWeaker(t *testing.T) {
	c := buildTestCatalog()

	weak := ComputeScore(c, []string{"Illaoi", "Jinx"}, 5)
	strong := ComputeScore(c, []string{"Illaoi", "Jinx", "Briar", "Aurora", "Meepsie"}, 5)

	if strong.Total <= weak.Total {
		t.Errorf("composición de 5 Anima (Total=%v) debería puntuar más que una de 2 (Total=%v)",
			strong.Total, weak.Total)
	}
}
