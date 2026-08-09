package engine

import (
	"sort"
	"strings"

	"github.com/NaisuG/Golden-Team-Buddy/internal/board"
	"github.com/NaisuG/Golden-Team-Buddy/internal/catalog"
)

// TargetSize devuelve a cuántos campeones apunta el árbol, según el nivel
// actual (cantidad de campeones en tablero):
//
//	0    -> 0 (módulo inactivo, mínimo 1 campeón en tablero para arrancar)
//	1-5  -> 6
//	6-7  -> 8
//	8    -> 9
//	9-10 -> 10
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

// candidateScore es el resultado de evaluar un campeón candidato.
type candidateScore struct {
	key   string
	score Score
}

// isReachable indica si un campeón puntual se puede conseguir ahora
// mismo: ya lo tenés (owned, sin importar si viene de tablero o banca),
// o su probabilidad de aparición a este nivel es mayor a 0. Un
// candidato no alcanzable ni se evalúa -- se descarta antes de competir
// por rasgos o sinergia.
func isReachable(c *catalog.Catalog, key string, owned map[string]bool, level int) bool {
	if owned[key] {
		return true
	}
	champ, ok := c.Champions[key]
	if !ok {
		return false
	}
	odds, ok := OddsOf(level, champ.Cost)
	if !ok {
		return true // sin dato conocido, no penalizamos
	}
	return odds > 0
}

// lexicographicBetter compara dos puntajes en cascada, no por suma:
// primero avance real de traits (breakpoints cruzados), y solo si
// empatan exacto ahí, decide la sinergia. Con esto, compartir un trait
// que ya no cruza nada nunca le gana a abrir uno que sí cruza, sin
// importar cuánta sinergia extra traiga la opción redundante.
func lexicographicBetter(a, b Score) bool {
	if a.TraitStrength != b.TraitStrength {
		return a.TraitStrength > b.TraitStrength
	}
	return a.Synergy > b.Synergy
}

// rankedCandidates evalúa todos los campeones elegibles (ni ya en line,
// ni en excluded, ni inalcanzables) y los devuelve ordenados de mejor a
// peor según lexicographicBetter.
func rankedCandidates(c *catalog.Catalog, line []string, owned map[string]bool, level int, excluded map[string]bool) []candidateScore {
	used := make(map[string]bool, len(line))
	for _, key := range line {
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
		if !isReachable(c, key, owned, level) {
			continue
		}
		candidate := append(append([]string{}, line...), key)
		results = append(results, candidateScore{key: key, score: ComputeScore(c, candidate, owned, level)})
	}

	sort.Slice(results, func(i, j int) bool {
		return lexicographicBetter(results[i].score, results[j].score)
	})

	return results
}

// BestNextChampion busca cuál candidato mejora más la línea, según
// lexicographicBetter (no según Total). excluded existe para forzar
// variantes distintas entre sí (ver GenerateVariants) -- pasá nil si no
// aplica.
//
// No la usa BuildLine (que hace beam search y necesita ver TODOS los
// candidatos de cada camino vivo, no solo el mejor de uno). Queda como
// utilidad de un solo paso, útil por separado y con sus propios tests.
func BestNextChampion(c *catalog.Catalog, line []string, owned map[string]bool, level int, excluded map[string]bool) (string, Score) {
	ranked := rankedCandidates(c, line, owned, level, excluded)
	if len(ranked) == 0 {
		return "", Score{}
	}
	return ranked[0].key, ranked[0].score
}

// Variant es una línea del árbol: una secuencia de campeones, con su
// puntaje y sus sub-variantes (1a, 1b, 1c).
type Variant struct {
	Champions []string
	Score     Score
	Children  []Variant
}

// beamWidth es cuántos caminos parciales se mantienen vivos en cada
// paso de la búsqueda, en vez de quedarse solo con el mejor. Con un
// catálogo real (~60-80 campeones) y líneas de a lo sumo 10, esto sigue
// siendo trivial de calcular incluso con un beam generoso -- no hace
// falta ser tacaño acá.
const beamWidth = 10

// beamEntry es un camino parcial vivo dentro de la búsqueda: la línea
// tal como va, y su puntaje.
type beamEntry struct {
	line  []string
	score Score
}

// dedupeBeam saca entradas que representan la MISMA composición (mismo
// conjunto de campeones, distinto orden de armado) -- sin esto, varios
// caminos podrían terminar ocupando lugares del beam con exactamente el
// mismo tablero, desperdiciando espacio que podría ir a una dirección
// genuinamente distinta.
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

// buildLines corre el mismo beam search que BuildLine, pero devuelve
// las hasta beamWidth líneas completas que sobrevivieron al final -- no
// solo la mejor. GenerateVariants usa esto para elegir sub-variantes a
// partir de lo que la propia búsqueda ya encontró al terminar, en vez
// de un chequeo aparte y más pobre sobre el primer paso nada más.
func buildLines(c *catalog.Catalog, seed []string, owned map[string]bool, level int, excluded map[string]bool) []Variant {
	target := TargetSize(level)
	seedLine := append([]string{}, seed...)

	beam := []beamEntry{{
		line:  seedLine,
		score: ComputeScore(c, seedLine, owned, level),
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
			break // ningún camino vivo tiene más candidatos para extender
		}

		sort.Slice(next, func(i, j int) bool {
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
		variants[i] = Variant{Champions: entry.line, Score: entry.score}
	}
	return variants
}

// BuildLine arma la mejor línea completa desde `seed` (normalmente
// vacío, o un solo campeón cuando GenerateVariants quiere forzar el
// primer pick de una variante) usando beam search -- ver buildLines
// para el detalle del mecanismo. No hay ningún prefijo obligatorio más
// allá de seed -- ni tablero ni banca fuerzan su inclusión, compiten
// como cualquier otro candidato.
func BuildLine(c *catalog.Catalog, seed []string, owned map[string]bool, level int, excluded map[string]bool) Variant {
	return buildLines(c, seed, owned, level, excluded)[0]
}

// signature da una representación canónica de una composición, sin
// importar el orden de los campeones. Sirve para detectar cuándo dos
// caminos de búsqueda distintos terminaron en exactamente el mismo
// tablero final.
func signature(championKeys []string) string {
	sorted := append([]string{}, championKeys...)
	sort.Strings(sorted)
	return strings.Join(sorted, "|")
}

// showSubVariant decide si el segundo candidato amerita mostrarse como
// sub-variante: solo cuando empata EXACTO en avance de traits con el
// mejor -- misma fuerza real, distinta forma de lograrla. Ya no existe
// una noción de "casi tan bueno" en TraitStrength, porque los
// breakpoints son todo-o-nada; la sinergia sí puede diferir entre
// ambos, y es justamente lo que hace que valga la pena mostrar la
// alternativa.
func showSubVariant(best, runnerUp Score) bool {
	return best.TraitStrength == runnerUp.TraitStrength
}

// ownedSet junta tablero y banca en un solo set. Para factibilidad da
// lo mismo de cuál de los dos venga -- ya está conseguido, no hace
// falta que "aparezca". La diferencia real entre tablero y banca (que
// solo tablero aporta rasgos y sinergia) la decide sola la búsqueda:
// arranca de cero y arma la línea por mérito, así que un campeón de
// banca (o de tablero) solo termina en la línea final si de verdad se
// gana el lugar -- estar "ya conseguido" nunca fuerza su inclusión.
func ownedSet(brd board.Board, bch board.Bench) map[string]bool {
	owned := make(map[string]bool, len(brd.Champions)+len(bch.Champions))
	for _, key := range brd.Champions {
		owned[key] = true
	}
	for _, key := range bch.Champions {
		owned[key] = true
	}
	return owned
}

// GenerateVariants arma hasta 3 líneas de primer nivel (forzando que
// difieran en el primer campeón nuevo, como multi-PV de un motor de
// ajedrez). La sub-variante de cada una sale del propio beam search de
// buildLines: la siguiente línea del mismo beam que empate EXACTO en
// TraitStrength con la mejor y todavía no se haya mostrado -- no de un
// chequeo aparte sobre el primer paso. Así, un empate que solo aparece
// varios pasos adentro de la búsqueda (como una elección entre dos
// segundos rasgos igual de fuertes) también se detecta, no solo los
// empates que ya eran visibles de entrada.
//
// brd y bch son lo que el usuario ya tiene -- se juntan en owned
// (gratis, sin necesidad de comprar), pero NINGUNO de los dos se
// inserta a la fuerza en la línea. Cada campeón, esté en tablero, en
// banca, o solo en el catálogo, compite por mérito por cada lugar. Esto
// es lo que habilita recombinación sin un caso especial: si conviene
// sentar a alguien de banca, sacar a alguien de tablero, o comprar algo
// nuevo, la búsqueda lo encuentra sola.
//
// Ninguna composición final se muestra dos veces: se lleva un registro
// (seen) de cada tablero completo ya mostrado, y si una línea o
// sub-variante termina coincidiendo con una ya vista -- aunque haya
// llegado por un camino distinto -- se descarta en vez de repetirla.
func GenerateVariants(c *catalog.Catalog, brd board.Board, bch board.Bench, level int) []Variant {
	owned := ownedSet(brd, bch)

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
			break // nivel 0 u otro caso sin objetivo: nada que generar
		}

		if seen[signature(best.Champions)] {
			excluded[best.Champions[0]] = true
			continue
		}
		seen[signature(best.Champions)] = true
		excluded[best.Champions[0]] = true

		for _, alt := range lines[1:] {
			if !showSubVariant(best.Score, alt.Score) {
				break // el beam viene ordenado: si esta ya no empata, ninguna de las siguientes tampoco
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
