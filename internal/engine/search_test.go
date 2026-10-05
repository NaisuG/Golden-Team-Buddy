package engine

import (
	"testing"

	"github.com/NaisuG/Golden-Team-Buddy/internal/board"
	"github.com/NaisuG/Golden-Team-Buddy/internal/catalog"
)

func bestCandidate(t *testing.T, c *catalog.Catalog, line []string) candidateScore {
	t.Helper()
	ranked := rankedCandidates(c, line, Owned{}, 5, nil)
	if len(ranked) == 0 {
		t.Fatal("rankedCandidates no devolvió candidatos")
	}
	return ranked[0]
}

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

func TestIsReachableIgnoresOwnedRegardlessOfCost(t *testing.T) {
	c := catalog.New()
	c.Champions["Caro"] = catalog.Champion{Key: "Caro", Cost: 5}

	if isReachable(c, "Caro", Owned{}, 1) {
		t.Error("Caro (costo 5) a nivel 1 sin estar owned debería ser inalcanzable")
	}
	owned := Owned{Bench: map[string]bool{"Caro": true}}
	if !isReachable(c, "Caro", owned, 1) {
		t.Error("Caro ya conseguido (owned) debería ser alcanzable sin importar el costo")
	}
}

func TestRankedCandidatesPicksDoubleBreakpoint(t *testing.T) {
	c := catalog.New()
	c.Traits["X"] = catalog.Trait{Key: "X", Breakpoints: []catalog.Breakpoint{{Style: "bronze", Min: 2}}}
	c.Traits["Y"] = catalog.Trait{Key: "Y", Breakpoints: []catalog.Breakpoint{{Style: "bronze", Min: 1}}}

	c.Champions["A"] = catalog.Champion{Key: "A", Traits: []string{"X"}}
	c.Champions["B"] = catalog.Champion{Key: "B", Traits: []string{"X", "Y"}}
	c.Champions["C"] = catalog.Champion{Key: "C", Traits: []string{"X"}}

	best := bestCandidate(t, c, []string{"A"})
	got, score := best.key, best.score

	if got != "B" {
		t.Errorf("mejor candidato = %q, quería %q", got, "B")
	}
	if score.TraitStrength <= 1 {
		t.Errorf("TraitStrength de agregar %q = %v, esperaba más de 1", got, score.TraitStrength)
	}
}

func TestBreakpointCrossingBeatsRedundantStacking(t *testing.T) {
	c := catalog.New()
	c.Traits["Shared"] = catalog.Trait{Key: "Shared", Breakpoints: []catalog.Breakpoint{{Style: "bronze", Min: 2}, {Style: "silver", Min: 4}}}
	c.Traits["Fresh"] = catalog.Trait{Key: "Fresh", Breakpoints: []catalog.Breakpoint{{Style: "bronze", Min: 1}}}

	c.Champions["Base1"] = catalog.Champion{Key: "Base1", Traits: []string{"Shared"}}
	c.Champions["Base2"] = catalog.Champion{Key: "Base2", Traits: []string{"Shared"}}
	c.Champions["Redundant"] = catalog.Champion{Key: "Redundant", Traits: []string{"Shared"}}
	c.Champions["Fresh1"] = catalog.Champion{Key: "Fresh1", Traits: []string{"Fresh"}}

	got := bestCandidate(t, c, []string{"Base1", "Base2"}).key

	if got != "Fresh1" {
		t.Errorf("mejor candidato = %q, quería %q", got, "Fresh1")
	}
}

func TestPushingToNextTierBeatsStartingTwoNewTraits(t *testing.T) {
	c := catalog.New()
	c.Traits["Active"] = catalog.Trait{Key: "Active", Breakpoints: []catalog.Breakpoint{{Style: "bronze", Min: 2}, {Style: "silver", Min: 4}}}
	c.Traits["New1"] = catalog.Trait{Key: "New1", Breakpoints: []catalog.Breakpoint{{Style: "bronze", Min: 2}}}
	c.Traits["New2"] = catalog.Trait{Key: "New2", Breakpoints: []catalog.Breakpoint{{Style: "bronze", Min: 2}}}

	c.Champions["Base1"] = catalog.Champion{Key: "Base1", Traits: []string{"Active"}}
	c.Champions["Base2"] = catalog.Champion{Key: "Base2", Traits: []string{"Active"}}
	c.Champions["ActivePush1"] = catalog.Champion{Key: "ActivePush1", Traits: []string{"Active"}}
	c.Champions["ActivePush2"] = catalog.Champion{Key: "ActivePush2", Traits: []string{"Active"}}
	c.Champions["NewStart1"] = catalog.Champion{Key: "NewStart1", Traits: []string{"New1"}}
	c.Champions["NewStart2"] = catalog.Champion{Key: "NewStart2", Traits: []string{"New2"}}

	pushScore := ComputeScore(c, []string{"Base1", "Base2", "ActivePush1", "ActivePush2"}, Owned{})
	startScore := ComputeScore(c, []string{"Base1", "Base2", "NewStart1", "NewStart2"}, Owned{})

	if pushScore.TraitStrength <= startScore.TraitStrength {
		t.Errorf("subir a plata (TraitStrength=%v) debería superar a abrir dos traits incompletos (TraitStrength=%v)",
			pushScore.TraitStrength, startScore.TraitStrength)
	}
}

func TestBuildLinesLevelZeroAddsNothing(t *testing.T) {
	c := catalog.New()
	c.Champions["A"] = catalog.Champion{Key: "A"}

	result := buildLines(c, nil, Owned{}, 0, nil)[0]

	if len(result.Champions) != 0 {
		t.Errorf("con nivel 0 (target 0), buildLines agregó %d campeones; esperaba 0", len(result.Champions))
	}
}

func TestBuildLinesReachesTargetSize(t *testing.T) {
	c := catalog.New()
	c.Traits["X"] = catalog.Trait{Key: "X", Breakpoints: []catalog.Breakpoint{{Style: "bronze", Min: 2}}}
	names := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	for _, n := range names {
		c.Champions[n] = catalog.Champion{Key: n, Traits: []string{"X"}}
	}

	result := buildLines(c, nil, Owned{}, 1, nil)[0]

	want := TargetSize(1)
	if len(result.Champions) != want {
		t.Errorf("buildLines con nivel 1 devolvió %d campeones; esperaba %d", len(result.Champions), want)
	}
}

func TestGenerateVariantsDoesNotForceUselessBenchMember(t *testing.T) {
	c := catalog.New()
	c.Traits["X"] = catalog.Trait{Key: "X", Breakpoints: []catalog.Breakpoint{{Style: "bronze", Min: 1}}}
	c.Champions["Useless"] = catalog.Champion{Key: "Useless"}
	names := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	for _, n := range names {
		c.Champions[n] = catalog.Champion{Key: n, Traits: []string{"X"}}
	}

	bench := board.Bench{Champions: []string{"Useless"}}
	variants := GenerateVariants(c, board.Board{}, bench, 1)

	if len(variants) == 0 {
		t.Fatal("GenerateVariants no devolvió ninguna variante")
	}
	for _, champ := range variants[0].Champions {
		if champ == "Useless" {
			t.Error("GenerateVariants incluyó a Useless desde la banca sin que aporte nada")
		}
	}
}

func TestGenerateVariantsCanDropAWeakBoardMember(t *testing.T) {
	c := catalog.New()
	c.Traits["X"] = catalog.Trait{Key: "X", Breakpoints: []catalog.Breakpoint{{Style: "bronze", Min: 1}}}
	c.Champions["Useless"] = catalog.Champion{Key: "Useless"}
	names := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	for _, n := range names {
		c.Champions[n] = catalog.Champion{Key: n, Traits: []string{"X"}}
	}

	brd := board.Board{Champions: []string{"Useless"}}
	variants := GenerateVariants(c, brd, board.Bench{}, 1)

	if len(variants) == 0 {
		t.Fatal("GenerateVariants no devolvió ninguna variante")
	}
	for _, champ := range variants[0].Champions {
		if champ == "Useless" {
			t.Error("GenerateVariants mantuvo a Useless del tablero sin que aporte nada")
		}
	}
}

func TestGenerateVariantsKeepsContributingOwnedChampion(t *testing.T) {
	c := catalog.New()
	c.Traits["X"] = catalog.Trait{Key: "X", Breakpoints: []catalog.Breakpoint{{Style: "bronze", Min: 1}}}
	names := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	for _, n := range names {
		c.Champions[n] = catalog.Champion{Key: n, Traits: []string{"X"}}
	}

	brd := board.Board{Champions: []string{"H"}}
	variants := GenerateVariants(c, brd, board.Bench{}, 1)

	if len(variants) == 0 {
		t.Fatal("GenerateVariants no devolvió ninguna variante")
	}
	for _, champ := range variants[0].Champions {
		if champ == "H" {
			return
		}
	}
	t.Errorf("GenerateVariants descartó a H, que está en tablero y aporta a X: %v", variants[0].Champions)
}

func TestGenerateVariantsGivesNoPreferenceToBench(t *testing.T) {
	c := catalog.New()
	c.Traits["X"] = catalog.Trait{Key: "X", Breakpoints: []catalog.Breakpoint{{Style: "bronze", Min: 1}}}
	names := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	for _, n := range names {
		c.Champions[n] = catalog.Champion{Key: n, Traits: []string{"X"}}
	}

	bench := board.Bench{Champions: []string{"H"}}
	variants := GenerateVariants(c, board.Board{}, bench, 1)

	if len(variants) == 0 {
		t.Fatal("GenerateVariants no devolvió ninguna variante")
	}
	for _, champ := range variants[0].Champions {
		if champ == "H" {
			t.Errorf("H está en banca y no debería tener preferencia: %v", variants[0].Champions)
		}
	}
}

func TestGenerateVariantsAreDistinct(t *testing.T) {
	c := catalog.New()
	c.Traits["X"] = catalog.Trait{Key: "X", Breakpoints: []catalog.Breakpoint{{Style: "bronze", Min: 2}}}
	names := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}
	for _, n := range names {
		c.Champions[n] = catalog.Champion{Key: n, Traits: []string{"X"}}
	}

	variants := GenerateVariants(c, board.Board{}, board.Bench{}, 1)

	if len(variants) == 0 {
		t.Fatal("GenerateVariants no devolvió ninguna variante")
	}
	if len(variants) > 3 {
		t.Fatalf("GenerateVariants devolvió %d variantes; el máximo es 3", len(variants))
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

func TestGenerateVariantsSuppressesDuplicateSubVariant(t *testing.T) {
	c := catalog.New()
	c.Traits["X"] = catalog.Trait{Key: "X", Breakpoints: []catalog.Breakpoint{{Style: "bronze", Min: 1}}}
	c.Champions["A"] = catalog.Champion{Key: "A", Traits: []string{"X"}}
	c.Champions["B"] = catalog.Champion{Key: "B", Traits: []string{"X"}}

	variants := GenerateVariants(c, board.Board{}, board.Bench{}, 1)

	if len(variants) == 0 {
		t.Fatal("GenerateVariants no devolvió ninguna variante")
	}
	if len(variants[0].Children) != 0 {
		t.Errorf("esperaba 0 sub-variantes, hay %d",
			len(variants[0].Children))
	}
}

func TestBuildLinesFindsOptimumGreedyWouldMiss(t *testing.T) {
	c := catalog.New()
	c.Traits["Seed"] = catalog.Trait{Key: "Seed", Breakpoints: []catalog.Breakpoint{{Style: "bronze", Min: 1}}}
	c.Traits["Early"] = catalog.Trait{Key: "Early", Breakpoints: []catalog.Breakpoint{{Style: "bronze", Min: 2}}}
	c.Traits["Late1"] = catalog.Trait{Key: "Late1", Breakpoints: []catalog.Breakpoint{{Style: "bronze", Min: 2}}}

	c.Champions["Base1"] = catalog.Champion{Key: "Base1", Traits: []string{"Seed", "Early"}}
	c.Champions["Base2"] = catalog.Champion{Key: "Base2", Traits: []string{"Early"}}
	c.Champions["Decoy1"] = catalog.Champion{Key: "Decoy1", Traits: []string{"Early"}}
	c.Champions["Decoy2"] = catalog.Champion{Key: "Decoy2", Traits: []string{"Early"}}
	c.Champions["Decoy3"] = catalog.Champion{Key: "Decoy3", Traits: []string{"Early"}}
	c.Champions["Decoy4"] = catalog.Champion{Key: "Decoy4", Traits: []string{"Early"}}
	c.Champions["LateSeed"] = catalog.Champion{Key: "LateSeed", Traits: []string{"Late1"}}
	c.Champions["LatePartner"] = catalog.Champion{Key: "LatePartner", Traits: []string{"Late1"}}

	result := buildLines(c, nil, Owned{}, 1, nil)[0]

	if result.Score.TraitStrength != 3 {
		t.Errorf("buildLines encontró TraitStrength=%v (%v), quería 3",
			result.Score.TraitStrength, result.Champions)
	}
}

func TestGenerateVariantsFindsDeepTieSubVariant(t *testing.T) {
	c := catalog.New()
	c.Traits["Seed"] = catalog.Trait{Key: "Seed", Breakpoints: []catalog.Breakpoint{{Style: "bronze", Min: 1}}}
	c.Traits["Early"] = catalog.Trait{Key: "Early", Breakpoints: []catalog.Breakpoint{{Style: "bronze", Min: 4}}}
	c.Traits["Late1"] = catalog.Trait{Key: "Late1", Breakpoints: []catalog.Breakpoint{{Style: "bronze", Min: 2}}}
	c.Traits["Late2"] = catalog.Trait{Key: "Late2", Breakpoints: []catalog.Breakpoint{{Style: "bronze", Min: 2}}}

	c.Champions["Base1"] = catalog.Champion{Key: "Base1", Traits: []string{"Seed", "Early"}}
	c.Champions["Base2"] = catalog.Champion{Key: "Base2", Traits: []string{"Early"}}
	c.Champions["Base3"] = catalog.Champion{Key: "Base3", Traits: []string{"Early"}}
	c.Champions["Base4"] = catalog.Champion{Key: "Base4", Traits: []string{"Early"}}
	c.Champions["LateA1"] = catalog.Champion{Key: "LateA1", Traits: []string{"Late1"}}
	c.Champions["LateA2"] = catalog.Champion{Key: "LateA2", Traits: []string{"Late1"}}
	c.Champions["LateB1"] = catalog.Champion{Key: "LateB1", Traits: []string{"Late2"}}
	c.Champions["LateB2"] = catalog.Champion{Key: "LateB2", Traits: []string{"Late2"}}

	variants := GenerateVariants(c, board.Board{}, board.Bench{}, 1)

	if len(variants) == 0 {
		t.Fatal("GenerateVariants no devolvió ninguna variante")
	}
	if len(variants[0].Children) != 1 {
		t.Fatalf("esperaba 1 sub-variante, hay %d", len(variants[0].Children))
	}
	if variants[0].Score.TraitStrength != variants[0].Children[0].Score.TraitStrength {
		t.Errorf("la sub-variante debería empatar en TraitStrength con la principal: %v vs %v",
			variants[0].Score.TraitStrength, variants[0].Children[0].Score.TraitStrength)
	}
	if signature(variants[0].Champions) == signature(variants[0].Children[0].Champions) {
		t.Error("la sub-variante no debería ser exactamente la misma composición que la principal")
	}
}

func TestRankedCandidatesIsDeterministicOnTies(t *testing.T) {
	c := catalog.New()
	c.Traits["NeedsFour"] = catalog.Trait{Key: "NeedsFour", Breakpoints: []catalog.Breakpoint{{Style: "bronze", Min: 4}}}
	c.Champions["Amy"] = catalog.Champion{Key: "Amy", Traits: []string{"NeedsFour"}}
	c.Champions["Zed"] = catalog.Champion{Key: "Zed", Traits: []string{"NeedsFour"}}

	first := bestCandidate(t, c, nil).key

	for i := 0; i < 200; i++ {
		got := bestCandidate(t, c, nil).key
		if got != first {
			t.Fatalf("corrida %d: el mejor candidato fue %q, la primera corrida dio %q", i, got, first)
		}
	}
}

func TestSignatureIsOrderIndependent(t *testing.T) {
	a := signature([]string{"Illaoi", "Jinx", "Briar"})
	b := signature([]string{"Briar", "Illaoi", "Jinx"})
	if a != b {
		t.Errorf("signature() depende del orden: %q vs %q", a, b)
	}
}

func TestGenerateVariantsNeverRepeatsAComposition(t *testing.T) {
	c := catalog.New()
	c.Traits["X"] = catalog.Trait{Key: "X", Breakpoints: []catalog.Breakpoint{{Style: "bronze", Min: 2}}}
	names := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}
	for _, n := range names {
		c.Champions[n] = catalog.Champion{Key: n, Traits: []string{"X"}}
	}

	variants := GenerateVariants(c, board.Board{}, board.Bench{}, 1)

	seen := make(map[string]bool)
	for _, v := range variants {
		sig := signature(v.Champions)
		if seen[sig] {
			t.Errorf("composición repetida en el árbol: %v", v.Champions)
		}
		seen[sig] = true
		for _, child := range v.Children {
			childSig := signature(child.Champions)
			if seen[childSig] {
				t.Errorf("sub-variante repite una composición ya mostrada: %v", child.Champions)
			}
			seen[childSig] = true
		}
	}
}
