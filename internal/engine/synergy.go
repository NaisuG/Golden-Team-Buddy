package engine

import (
	"math"

	"github.com/NaisuG/Golden-Team-Buddy/internal/catalog"
)

type Vector map[string]float64

func ChampionVector(c *catalog.Catalog, championKey string) Vector {
	traits := c.TraitsOf(championKey)
	v := make(Vector, len(traits))
	for _, t := range traits {
		v[t.Key] = 1
	}
	return v
}

// Synergy es la similitud coseno entre los traits de dos campeones.
func Synergy(a, b Vector) float64 {
	var dot, magA, magB float64
	for key, va := range a {
		magA += va * va
		if vb, ok := b[key]; ok {
			dot += va * vb
		}
	}
	for _, vb := range b {
		magB += vb * vb
	}
	if magA == 0 || magB == 0 {
		return 0
	}
	return dot / (math.Sqrt(magA) * math.Sqrt(magB))
}
