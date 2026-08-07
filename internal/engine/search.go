package engine

import (
	"sort"

	"github.com/NaisuG/Golden-Team-Buddy/internal/catalog"
)

func TargetSize(level int) int {
	switch {
	case level == 0:
		return 0
	case level <= 5:
		return 6
	case level <= 7:
		return 8
	case level == 8:
		return 9
	default:
		return 10
	}
}

type candidateScore struct {
	key   string
	score Score
}

func rankedCandidates(c *catalog.Catalog, current []string, level int, excluded map[string]bool) []candidateScore {
	used := make(map[string]bool, len(current))
	for _, key := range current {
		used[key] = true
	}
	for key := range excluded {
		used[key] = true
	}

	var results []candidateScore
	for key := range c.Champions {
		if used[key] {
			continue
		}
		candidate := append(append([]string{}, current...), key)
		results = append(results, candidateScore{key: key, score: ComputeScore(c, candidate, level)})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].score.Total > results[j].score.Total
	})

	return results
}

func BestNextChampion(c *catalog.Catalog, current []string, level int, excluded map[string]bool) (string, Score) {
	ranked := rankedCandidates(c, current, level, excluded)
	if len(ranked) == 0 {
		return "", Score{}
	}
	return ranked[0].key, ranked[0].score
}

type Variant struct {
	Champions []string
	Score     Score
	Children  []Variant
}

func BuildLine(c *catalog.Catalog, current []string, level int, excluded map[string]bool) Variant {
	target := TargetSize(level)
	line := append([]string{}, current...)

	for len(line) < target {
		next, _ := BestNextChampion(c, line, level, excluded)
		if next == "" {
			break
		}
		line = append(line, next)
	}

	return Variant{
		Champions: line,
		Score:     ComputeScore(c, line, level),
	}
}

const branchThreshold = 0.15

func GenerateVariants(c *catalog.Catalog, current []string, level int) []Variant {
	var variants []Variant
	excluded := make(map[string]bool)

	for i := 0; i < 3; i++ {
		ranked := rankedCandidates(c, current, level, excluded)
		if len(ranked) == 0 {
			break
		}

		best := ranked[0]
		line := BuildLine(c, append(append([]string{}, current...), best.key), level, nil)

		if len(ranked) > 1 {
			runnerUp := ranked[1]
			if runnerUp.score.Total >= best.score.Total*(1-branchThreshold) {
				subLine := BuildLine(c, append(append([]string{}, current...), runnerUp.key), level, nil)
				line.Children = []Variant{subLine}
			}
		}

		variants = append(variants, line)
		excluded[best.key] = true
	}

	return variants
}
