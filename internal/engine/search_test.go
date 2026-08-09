package engine

import (
	"testing"

	"github.com/NaisuG/Golden-Team-Buddy/internal/board"
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

func TestIsReachableIgnoresOwnedRegardlessOfCost(t *testing.T) {
	c := catalog.New()
	c.Champions["Caro"] = catalog.Champion{Key: "Caro", Cost: 5}

	if isReachable(c, "Caro", nil, 1) {
		t.Error("Caro (costo 5) a nivel 1 sin estar owned debería ser inalcanzable")
	}
	owned := map[string]bool{"Caro": true}
	if !isReachable(c, "Caro", owned, 1) {
		t.Error("Caro ya conseguido (owned) debería ser alcanzable sin importar el costo")
	}
}

func TestBestNextChampionPicksDoubleBreakpoint(t *testing.T) {
	c := catalog.New()
	c.Traits["X"] = catalog.Trait{Key: "X", Breakpoints: []catalog.Breakpoint{{Style: "bronze", Min: 2}}}
	c.Traits["Y"] = catalog.Trait{Key: "Y", Breakpoints: []catalog.Breakpoint{{Style: "bronze", Min: 1}}}

	c.Champions["A"] = catalog.Champion{Key: "A", Traits: []string{"X"}}
	c.Champions["B"] = catalog.Champion{Key: "B", Traits: []string{"X", "Y"}}
	c.Champions["C"] = catalog.Champion{Key: "C", Traits: []string{"X"}}

	got, score := BestNextChampion(c, []string{"A"}, nil, 5, nil)

	if got != "B" {
		t.Errorf("BestNextChampion() = %q, quería %q", got, "B")
	}
	if score.TraitStrength <= 1 {
		t.Errorf("TraitStrength de agregar %q = %v, esperaba más de 1", got, score.TraitStrength)
	}
}

// Caso concreto de la conversación: cruzar un breakpoint nuevo (Fresh,
// bronce en 1) tiene que ganarle a apilar una copia de más en un trait
// que ya cruzó su bronce y no llega al siguiente (Shared, de 2 a 3, con
// plata recién en 4) -- sin importar que la copia de más comparta trait
// con toda la línea y la nueva no comparta nada todavía.
func TestBreakpointCrossingBeatsRedundantStacking(t *testing.T) {
	c := catalog.New()
	c.Traits["Shared"] = catalog.Trait{Key: "Shared", Breakpoints: []catalog.Breakpoint{{Style: "bronze", Min: 2}, {Style: "silver", Min: 4}}}
	c.Traits["Fresh"] = catalog.Trait{Key: "Fresh", Breakpoints: []catalog.Breakpoint{{Style: "bronze", Min: 1}}}

	c.Champions["Base1"] = catalog.Champion{Key: "Base1", Traits: []string{"Shared"}}
	c.Champions["Base2"] = catalog.Champion{Key: "Base2", Traits: []string{"Shared"}}
	c.Champions["Redundant"] = catalog.Champion{Key: "Redundant", Traits: []string{"Shared"}}
	c.Champions["Fresh1"] = catalog.Champion{Key: "Fresh1", Traits: []string{"Fresh"}}

	got, _ := BestNextChampion(c, []string{"Base1", "Base2"}, nil, 5, nil)

	if got != "Fresh1" {
		t.Errorf("BestNextChampion() = %q, quería %q -- cruzar un breakpoint nuevo debe ganarle a una copia que no cruza nada, sin importar la sinergia", got, "Fresh1")
	}
}

// Caso confirmado en la conversación: empujar un trait activo de bronce
// a plata (breakpoint real) vale más que abrir dos traits nuevos sin
// completar ninguno.
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

	pushScore := ComputeScore(c, []string{"Base1", "Base2", "ActivePush1", "ActivePush2"}, nil, 5)
	startScore := ComputeScore(c, []string{"Base1", "Base2", "NewStart1", "NewStart2"}, nil, 5)

	if pushScore.TraitStrength <= startScore.TraitStrength {
		t.Errorf("empujar bronce->plata (TraitStrength=%v) debería superar abrir dos traits nuevos sin completar (TraitStrength=%v)",
			pushScore.TraitStrength, startScore.TraitStrength)
	}
}

func TestBuildLineLevelZeroAddsNothing(t *testing.T) {
	c := catalog.New()
	c.Champions["A"] = catalog.Champion{Key: "A"}

	result := BuildLine(c, nil, nil, 0, nil)

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

	result := BuildLine(c, nil, nil, 1, nil)

	want := TargetSize(1)
	if len(result.Champions) != want {
		t.Errorf("BuildLine con nivel 1 devolvió %d campeones; esperaba %d", len(result.Champions), want)
	}
}

// El bug del día: banca contaba para rasgos/sinergia sin ganárselo.
// Acá, "Useless" no aporta ningún trait -- estar en banca (ya
// conseguido, gratis) no debería alcanzar para que la línea final lo
// incluya si no suma nada.
func TestGenerateVariantsDoesNotForceUselessBenchMember(t *testing.T) {
	c := catalog.New()
	// bronce en 1: un solo portador de X ya cruza el breakpoint, así que
	// "un X" y "nada" nunca pueden empatar en 0-0 en el primer paso --
	// eso evitaría que el test dependa del orden de mapa.
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
			t.Error("GenerateVariants incluyó a un campeón de banca que no aporta nada, solo por estar ya conseguido -- banca no debería forzar inclusión")
		}
	}
}

// Recombinación: tablero tampoco debería forzar inclusión. Si lo que
// tenés puesto no aporta nada, la búsqueda tiene que poder dejarlo
// afuera igual que a cualquier candidato débil.
func TestGenerateVariantsCanDropAWeakBoardMember(t *testing.T) {
	c := catalog.New()
	// mismo motivo que en el test de banca: bronce en 1 para que no haya
	// empate 0-0 posible en el primer paso.
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
			t.Error("GenerateVariants mantuvo en la línea a un campeón de tablero que no aporta nada -- tablero tampoco debería forzar inclusión")
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
		t.Errorf("con A y B empatados pero solo una composición posible, esperaba 0 sub-variantes (duplicada); tiene %d",
			len(variants[0].Children))
	}
}

// Reconstrucción en miniatura del hallazgo real de la conversación
// (Leona vs. Veigar+Lissandra): apilar de más un trait ya activo
// (Early) puede ganarle, PASO A PASO, a empezar a construir un segundo
// trait (Late1) que recién vale algo una vez completo. Pero la
// composición óptima de 6 en este catálogo SÍ activa los dos
// (Seed+Early+Late1 = TraitStrength 3) -- y ese óptimo solo aparece si
// la búsqueda compara líneas completas, no si se compromete paso a paso
// con lo que se ve mejor en el momento.
func TestBuildLineFindsOptimumGreedyWouldMiss(t *testing.T) {
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

	result := BuildLine(c, nil, nil, 1, nil)

	if result.Score.TraitStrength != 3 {
		t.Errorf("BuildLine encontró TraitStrength=%v (%v); el óptimo real de este catálogo es 3 (Seed+Early+Late1) -- cayó en la trampa de apilar Early de más en vez de completar Late1",
			result.Score.TraitStrength, result.Champions)
	}
}

// El chequeo viejo de sub-variante solo miraba el primer paso -- acá,
// Base1 es un primer pick único y sin ningún empate en ese momento (el
// chequeo viejo nunca hubiera mostrado ninguna sub-variante). Recién
// varios pasos después aparece un empate real: LateA1+LateA2 y
// LateB1+LateB2 dan exactamente el mismo TraitStrength final, siendo
// composiciones genuinamente distintas. Solo se detecta mirando el
// resultado final de la búsqueda.
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
		t.Fatalf("esperaba 1 sub-variante (empate profundo entre Late1 y Late2); encontró %d", len(variants[0].Children))
	}
	if variants[0].Score.TraitStrength != variants[0].Children[0].Score.TraitStrength {
		t.Errorf("la sub-variante debería empatar en TraitStrength con la principal: %v vs %v",
			variants[0].Score.TraitStrength, variants[0].Children[0].Score.TraitStrength)
	}
	if signature(variants[0].Champions) == signature(variants[0].Children[0].Champions) {
		t.Error("la sub-variante no debería ser exactamente la misma composición que la principal")
	}
}
func TestSignatureIsOrderIndependent(t *testing.T) {
	a := signature([]string{"Illaoi", "Jinx", "Briar"})
	b := signature([]string{"Briar", "Illaoi", "Jinx"})
	if a != b {
		t.Errorf("signature() debería dar lo mismo sin importar el orden: %q vs %q", a, b)
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
