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
		{name: "2 Anima, no activa nada", team: []string{"Illaoi", "Jinx"}, expected: 0},
		{name: "3 Anima, bronce (tier 1)", team: []string{"Illaoi", "Jinx", "Briar"}, expected: 1},
		{name: "5 Anima, oro (tier 2)", team: []string{"Illaoi", "Jinx", "Briar", "Aurora", "Meepsie"}, expected: 4},
	}

	for _, tc := range cases {
		got := TraitStrength(c, tc.team)
		if got != tc.expected {
			t.Errorf("%s: TraitStrength() = %v, quería %v", tc.name, got, tc.expected)
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
			t.Errorf("%s: AverageSynergy() = %v, quería %v", tc.name, got, tc.expected)
		}
	}
}

func TestComputeScoreStrongerBeatsWeaker(t *testing.T) {
	c := buildTestCatalog()

	weak := ComputeScore(c, []string{"Illaoi", "Jinx"}, Owned{})
	strong := ComputeScore(c, []string{"Illaoi", "Jinx", "Briar", "Aurora", "Meepsie"}, Owned{})

	if !lexicographicBetter(strong, weak) {
		t.Errorf("composición de 5 Anima (%+v) debería ganarle a una de 2 (%+v)", strong, weak)
	}
}

func TestActiveTraits(t *testing.T) {
	c := buildTestCatalog()

	if got := ActiveTraits(c, []string{"Illaoi", "Jinx"}); len(got) != 0 {
		t.Errorf("con 2 Anima no debería haber traits activos, hay %v", got)
	}

	got := ActiveTraits(c, []string{"Illaoi", "Jinx", "Briar", "Aurora", "Meepsie"})
	want := ActiveTrait{Name: "Anima", Count: 5, Style: "gold"}
	if len(got) != 1 || got[0] != want {
		t.Errorf("ActiveTraits() = %v, quería [%v]", got, want)
	}
}
