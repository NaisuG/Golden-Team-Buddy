package engine

import "github.com/NaisuG/Golden-Team-Buddy/internal/catalog"

type Score struct {
	TraitStrength float64
	Synergy       float64
	Feasibility   float64
	Total         float64
}

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

func Feasibility(c *catalog.Catalog, championKeys []string, level int) float64 {
	if len(championKeys) == 0 {
		return 0
	}

	var total float64
	var known int
	for _, key := range championKeys {
		champ, ok := c.Champions[key]
		if !ok {
			continue
		}
		odds, ok := OddsOf(level, champ.Cost)
		if !ok {
			continue
		}
		total += odds
		known++
	}
	if known == 0 {
		return 0
	}
	return total / float64(known)
}

const (
	traitWeight       = 1.0
	synergyWeight     = 3.0
	feasibilityWeight = 2.0
)

func ComputeScore(c *catalog.Catalog, championKeys []string, level int) Score {
	s := Score{
		TraitStrength: TraitStrength(c, championKeys),
		Synergy:       AverageSynergy(c, championKeys),
		Feasibility:   Feasibility(c, championKeys, level),
	}
	s.Total = traitWeight*s.TraitStrength + synergyWeight*s.Synergy + feasibilityWeight*s.Feasibility
	return s
}
