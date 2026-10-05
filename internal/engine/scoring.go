package engine

import (
	"sort"

	"github.com/NaisuG/Golden-Team-Buddy/internal/catalog"
)

// Score es el puntaje de una composición, desglosado por criterio.
type Score struct {
	TraitStrength float64
	BoardCount    int
	Synergy       float64
}

func traitCounts(c *catalog.Catalog, championKeys []string) map[string]int {
	counts := make(map[string]int)
	for _, key := range championKeys {
		for _, trait := range c.TraitsOf(key) {
			counts[trait.Key]++
		}
	}
	return counts
}

// TraitStrength suma el valor de los breakpoints que alcanza una composición.
func TraitStrength(c *catalog.Catalog, championKeys []string) float64 {
	var total float64
	for traitKey, count := range traitCounts(c, championKeys) {
		total += breakpointValue(c.Traits[traitKey], count)
	}
	return total
}

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

// BoardCount cuenta los campeones del tablero que aportan a algún trait activo.
func BoardCount(c *catalog.Catalog, championKeys []string, onBoard map[string]bool) int {
	counts := traitCounts(c, championKeys)
	count := 0
	for _, key := range championKeys {
		if !onBoard[key] {
			continue
		}
		for _, trait := range c.TraitsOf(key) {
			if breakpointValue(trait, counts[trait.Key]) > 0 {
				count++
				break
			}
		}
	}
	return count
}

// ComputeScore calcula todos los criterios de puntaje de una composición.
func ComputeScore(c *catalog.Catalog, championKeys []string, owned Owned) Score {
	return Score{
		TraitStrength: TraitStrength(c, championKeys),
		BoardCount:    BoardCount(c, championKeys, owned.Board),
		Synergy:       AverageSynergy(c, championKeys),
	}
}

// ActiveTrait es un trait que alcanzó al menos su primer breakpoint en una composición.
type ActiveTrait struct {
	Name  string
	Count int
	Style string
}

func ActiveTraits(c *catalog.Catalog, championKeys []string) []ActiveTrait {
	var active []ActiveTrait
	for traitKey, count := range traitCounts(c, championKeys) {
		trait := c.Traits[traitKey]
		style := ""
		for _, bp := range trait.Breakpoints {
			if count >= bp.Min {
				style = bp.Style
			}
		}
		if style == "" {
			continue
		}
		active = append(active, ActiveTrait{Name: trait.Name, Count: count, Style: style})
	}
	sort.Slice(active, func(i, j int) bool {
		if active[i].Count != active[j].Count {
			return active[i].Count > active[j].Count
		}
		return active[i].Name < active[j].Name
	})
	return active
}
