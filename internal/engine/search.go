package engine

import (
	"sort"
	"strings"

	"github.com/NaisuG/Golden-Team-Buddy/internal/board"
	"github.com/NaisuG/Golden-Team-Buddy/internal/catalog"
)

// TargetSize devuelve cuántos campeones debe tener una línea según el nivel.
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

func isReachable(c *catalog.Catalog, key string, owned Owned, level int) bool {
	if owned.Has(key) {
		return true
	}
	champ, ok := c.Champions[key]
	if !ok {
		return false
	}
	odds, ok := OddsOf(level, champ.Cost)
	if !ok {
		return true
	}
	return odds > 0
}

func lexicographicBetter(a, b Score) bool {
	if a.TraitStrength != b.TraitStrength {
		return a.TraitStrength > b.TraitStrength
	}
	if a.BoardCount != b.BoardCount {
		return a.BoardCount > b.BoardCount
	}
	return a.Synergy > b.Synergy
}

func rankedCandidates(c *catalog.Catalog, line []string, owned Owned, level int, excluded map[string]bool) []candidateScore {
	used := make(map[string]bool, len(line))
	for _, key := range line {
		used[key] = true
	}
	for key := range excluded {
		used[key] = true
	}

	keys := make([]string, 0, len(c.Champions))
	for key := range c.Champions {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var results []candidateScore
	for _, key := range keys {
		if used[key] {
			continue
		}
		if !isReachable(c, key, owned, level) {
			continue
		}
		candidate := append(append([]string{}, line...), key)
		results = append(results, candidateScore{key: key, score: ComputeScore(c, candidate, owned)})
	}

	sort.SliceStable(results, func(i, j int) bool {
		return lexicographicBetter(results[i].score, results[j].score)
	})

	return results
}

// Variant es una composición sugerida, con sus sub-variantes equivalentes.
type Variant struct {
	Champions []string
	Traits    []ActiveTrait
	Score     Score
	Children  []Variant
}

const beamWidth = 10

type beamEntry struct {
	line  []string
	score Score
}

func dedupeBeam(entries []beamEntry) []beamEntry {
	seen := make(map[string]bool, len(entries))
	result := make([]beamEntry, 0, len(entries))
	for _, e := range entries {
		sig := signature(e.line)
		if seen[sig] {
			continue
		}
		seen[sig] = true
		result = append(result, e)
	}
	return result
}

func buildLines(c *catalog.Catalog, seed []string, owned Owned, level int, excluded map[string]bool) []Variant {
	target := TargetSize(level)
	seedLine := append([]string{}, seed...)

	beam := []beamEntry{{
		line:  seedLine,
		score: ComputeScore(c, seedLine, owned),
	}}

	for len(beam[0].line) < target {
		var next []beamEntry
		for _, entry := range beam {
			for _, cand := range rankedCandidates(c, entry.line, owned, level, excluded) {
				extended := append(append([]string{}, entry.line...), cand.key)
				next = append(next, beamEntry{line: extended, score: cand.score})
			}
		}
		if len(next) == 0 {
			break
		}

		sort.SliceStable(next, func(i, j int) bool {
			return lexicographicBetter(next[i].score, next[j].score)
		})
		next = dedupeBeam(next)
		if len(next) > beamWidth {
			next = next[:beamWidth]
		}
		beam = next
	}

	variants := make([]Variant, len(beam))
	for i, entry := range beam {
		variants[i] = Variant{Champions: entry.line, Traits: ActiveTraits(c, entry.line), Score: entry.score}
	}
	return variants
}

func signature(championKeys []string) string {
	sorted := append([]string{}, championKeys...)
	sort.Strings(sorted)
	return strings.Join(sorted, "|")
}

func showSubVariant(best, runnerUp Score) bool {
	return best.TraitStrength == runnerUp.TraitStrength && best.BoardCount == runnerUp.BoardCount
}

// GenerateVariants arma hasta tres composiciones distintas con beam search.
func GenerateVariants(c *catalog.Catalog, brd board.Board, bch board.Bench, level int) []Variant {
	owned := NewOwned(brd, bch)

	var variants []Variant
	excluded := make(map[string]bool)
	seen := make(map[string]bool)

	for i := 0; i < 3; i++ {
		lines := buildLines(c, nil, owned, level, excluded)
		if len(lines) == 0 {
			break
		}

		best := lines[0]
		if len(best.Champions) == 0 {
			break
		}

		if seen[signature(best.Champions)] {
			excluded[best.Champions[0]] = true
			continue
		}
		seen[signature(best.Champions)] = true
		excluded[best.Champions[0]] = true

		for _, alt := range lines[1:] {
			if !showSubVariant(best.Score, alt.Score) {
				break
			}
			if seen[signature(alt.Champions)] {
				continue
			}
			seen[signature(alt.Champions)] = true
			best.Children = []Variant{alt}
			break
		}

		variants = append(variants, best)
	}

	return variants
}
