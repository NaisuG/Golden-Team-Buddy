package engine

import (
	"testing"

	"github.com/NaisuG/Golden-Team-Buddy/internal/catalog"
)

func TestTargetSize(t *testing.T) {
	cases := map[int]int{
		0: 0, 1: 6, 5: 6, 6: 8, 7: 8, 8: 9, 9: 10, 10: 10,
	}
	for level, want := range cases {
		if got := TargetSize(level); got != want {
			t.Errorf("TargetSize(%d) = %d, quería %d", level, got, want)
		}
	}
}

func TestBestNextChampionPicksDoubleBreakpoint(t *testing.T) {
	c := catalog.New()
	c.Traits["X"] = catalog.Trait{Key: "X", Breakpoints: []catalog.Breakpoint{{Style: "bronze", Min: 2}}}
	c.Traits["Y"] = catalog.Trait{Key: "Y", Breakpoints: []catalog.Breakpoint{{Style: "bronze", Min: 1}}}

	c.Champions["A"] = catalog.Champion{Key: "A", Traits: []string{"X"}}
	c.Champions["B"] = catalog.Champion{Key: "B", Traits: []string{"X", "Y"}}
	c.Champions["C"] = catalog.Champion{Key: "C", Traits: []string{"X"}}

	got, score := BestNextChampion(c, []string{"A"}, 5, nil)

	if got != "B" {
		t.Errorf("BestNextChampion() = %q, quería %q", got, "B")
	}
	if score.TraitStrength <= 1 {
		t.Errorf("TraitStrength de agregar %q = %v, esperaba más de 1", got, score.TraitStrength)
	}
}

func TestBuildLineLevelZeroAddsNothing(t *testing.T) {
	c := catalog.New()
	c.Champions["A"] = catalog.Champion{Key: "A"}

	result := BuildLine(c, []string{}, 0, nil)

	if len(result.Champions) != 0 {
		t.Errorf("con nivel 0 (target 0), BuildLine agregó %d campeones; esperaba 0", len(result.Champions))
	}
}

func TestBuildLineReachesTargetSize(t *testing.T) {
	c := catalog.New()
	c.Traits["X"] = catalog.Trait{Key: "X", Breakpoints: []catalog.Breakpoint{{Style: "bronze", Min: 2}}}
	names := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	for _, n := range names {
		c.Champions[n] = catalog.Champion{Key: n, Traits: []string{"X"}}
	}

	result := BuildLine(c, []string{}, 1, nil)

	want := TargetSize(1)
	if len(result.Champions) != want {
		t.Errorf("BuildLine con nivel 1 devolvió %d campeones; esperaba %d", len(result.Champions), want)
	}
}

func TestGenerateVariantsAreDistinct(t *testing.T) {
	c := catalog.New()
	c.Traits["X"] = catalog.Trait{Key: "X", Breakpoints: []catalog.Breakpoint{{Style: "bronze", Min: 2}}}
	names := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}
	for _, n := range names {
		c.Champions[n] = catalog.Champion{Key: n, Traits: []string{"X"}}
	}

	variants := GenerateVariants(c, []string{}, 1)

	if len(variants) != 3 {
		t.Fatalf("GenerateVariants devolvió %d variantes; esperaba 3", len(variants))
	}

	seenFirstPick := make(map[string]bool)
	for i, v := range variants {
		if len(v.Champions) != TargetSize(1) {
			t.Errorf("variante %d tiene %d campeones; esperaba %d", i, len(v.Champions), TargetSize(1))
		}
		first := v.Champions[0]
		if seenFirstPick[first] {
			t.Errorf("variante %d repite %q como primer campeón, ya lo había usado otra variante", i, first)
		}
		seenFirstPick[first] = true
	}
}

func TestGenerateVariantsCreatesSubVariantOnCloseTie(t *testing.T) {
	c := catalog.New()
	c.Traits["X"] = catalog.Trait{Key: "X", Breakpoints: []catalog.Breakpoint{{Style: "bronze", Min: 1}}}
	c.Champions["A"] = catalog.Champion{Key: "A", Traits: []string{"X"}}
	c.Champions["B"] = catalog.Champion{Key: "B", Traits: []string{"X"}}

	variants := GenerateVariants(c, []string{}, 1)

	if len(variants) == 0 {
		t.Fatal("GenerateVariants no devolvió ninguna variante")
	}
	if len(variants[0].Children) != 1 {
		t.Errorf("con A y B empatados exacto, esperaba 1 sub-variante en la primera línea; tiene %d",
			len(variants[0].Children))
	}
}

func TestGenerateVariantsNoSubVariantWhenGapIsLarge(t *testing.T) {
	c := catalog.New()
	c.Traits["X"] = catalog.Trait{Key: "X", Breakpoints: []catalog.Breakpoint{{Style: "bronze", Min: 1}}}
	c.Traits["Y"] = catalog.Trait{Key: "Y", Breakpoints: []catalog.Breakpoint{{Style: "bronze", Min: 1}}}
	c.Champions["A"] = catalog.Champion{Key: "A", Traits: []string{"X"}}
	c.Champions["B"] = catalog.Champion{Key: "B", Traits: []string{"X"}}
	c.Champions["C"] = catalog.Champion{Key: "C", Traits: []string{"X", "Y"}}
	variants := GenerateVariants(c, []string{}, 1)

	if len(variants) == 0 {
		t.Fatal("GenerateVariants no devolvió ninguna variante")
	}
	if len(variants[0].Children) != 0 {
		t.Errorf("con una brecha grande entre 1er y 2do lugar, no esperaba sub-variante; tiene %d",
			len(variants[0].Children))
	}
}
