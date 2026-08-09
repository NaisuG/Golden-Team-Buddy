package engine

import "github.com/NaisuG/Golden-Team-Buddy/internal/catalog"

// Score es el puntaje de una composición, desglosado por componente.
type Score struct {
	TraitStrength float64
	Synergy       float64
	Feasibility   float64
	Total         float64
}

// TraitStrength suma el valor de los breakpoints que una composición
// alcanza. Cada Trait.Breakpoints debe venir ordenado de menor a mayor
// Min -- eso es lo que permite quedarnos con "el más alto alcanzado" con
// un solo recorrido.
func TraitStrength(c *catalog.Catalog, championKeys []string) float64 {
	counts := make(map[string]int)
	for _, key := range championKeys {
		for _, trait := range c.TraitsOf(key) {
			counts[trait.Key]++
		}
	}

	var total float64
	for traitKey, count := range counts {
		trait := c.Traits[traitKey]
		total += breakpointValue(trait, count)
	}
	return total
}

// breakpointValue es todo-o-nada: en TFT un trait recién existe al
// cruzar su primer breakpoint entero (2 Vanguardias activa Vanguardia,
// 1 sola no activa nada). No hay puntos intermedios por "ir
// construyendo hacia" un breakpoint que todavía no se cruzó.
func breakpointValue(trait catalog.Trait, count int) float64 {
	var value float64
	for i, bp := range trait.Breakpoints {
		if count >= bp.Min {
			tier := float64(i + 1)
			value = tier * tier
		}
	}
	return value
}

// AverageSynergy calcula el promedio de Synergy() entre todos los pares
// posibles de campeones de la composición. Con menos de 2 campeones no
// hay ningún par que comparar, así que devuelve 0.
func AverageSynergy(c *catalog.Catalog, championKeys []string) float64 {
	if len(championKeys) < 2 {
		return 0
	}

	vectors := make([]Vector, len(championKeys))
	for i, key := range championKeys {
		vectors[i] = ChampionVector(c, key)
	}

	var total float64
	var pairs int
	for i := 0; i < len(vectors); i++ {
		for j := i + 1; j < len(vectors); j++ {
			total += Synergy(vectors[i], vectors[j])
			pairs++
		}
	}
	return total / float64(pairs)
}

// Feasibility devuelve la menor probabilidad de aparición conocida entre
// los campeones de la composición que el usuario TODAVÍA NO TIENE.
// owned marca cuáles ya están conseguidos (tablero o banca, da lo mismo
// cuál) -- esos se saltan, porque no dependen de aparecer en la tienda.
//
// Con isReachable ya filtrando en rankedCandidates, cualquier línea que
// llega hasta acá solo contiene picks individualmente alcanzables --
// este valor queda como dato informativo del puntaje final (qué tan
// cómoda es la línea en conjunto), no como filtro adicional.
func Feasibility(c *catalog.Catalog, championKeys []string, owned map[string]bool, level int) float64 {
	min := 1.0
	known := false
	for _, key := range championKeys {
		if owned[key] {
			continue
		}
		champ, ok := c.Champions[key]
		if !ok {
			continue
		}
		odds, ok := OddsOf(level, champ.Cost)
		if !ok {
			continue
		}
		known = true
		if odds < min {
			min = odds
		}
	}
	if !known {
		return 1.0
	}
	return min
}

// Pesos de arranque para el Total informativo. No están calibrados
// contra datos reales todavía -- eso es trabajo de v2, con partidas
// reales via "Finalizar Tablero".
const (
	traitWeight   = 1.0
	synergyWeight = 3.0
)

// ComputeScore junta las tres señales. Total es un número único a modo
// informativo para mostrar en el árbol -- pero OJO: no es lo que decide
// qué candidato gana durante la búsqueda. Esa decisión la toma
// lexicographicBetter() en search.go, comparando primero TraitStrength
// y solo en empate exacto Synergy. Total nunca entra en esa
// comparación, solo se muestra al final.
func ComputeScore(c *catalog.Catalog, championKeys []string, owned map[string]bool, level int) Score {
	s := Score{
		TraitStrength: TraitStrength(c, championKeys),
		Synergy:       AverageSynergy(c, championKeys),
		Feasibility:   Feasibility(c, championKeys, owned, level),
	}
	s.Total = traitWeight*s.TraitStrength + synergyWeight*s.Synergy
	return s
}