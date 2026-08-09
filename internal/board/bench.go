package board

import "fmt"

// Bench es la banca: campeones ya comprados que todavía no están en
// tablero. Confirmado en la conversación: banca no aporta rasgos ni
// sinergia -- eso lo respeta engine.GenerateVariants al no forzar la
// inclusión de nadie de banca en la línea, solo tratarlo como candidato
// ya conseguido (gratis) igual que a un campeón de tablero.
type Bench struct {
	Champions []string
}

func (b *Bench) Add(championKey string) error {
	if len(b.Champions) >= MaxBenchSize {
		return fmt.Errorf("banca llena: máximo %d campeones", MaxBenchSize)
	}
	b.Champions = append(b.Champions, championKey)
	return nil
}
