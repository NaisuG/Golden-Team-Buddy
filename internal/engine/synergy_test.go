package engine

import (
	"testing"

	"github.com/NaisuG/Golden-Team-Buddy/internal/catalog"
)

func TestSynergySharedTraitIsCloserThanNone(t *testing.T) {
	c := catalog.New()
	c.Traits["Anima"] = catalog.Trait{Key: "Anima", Name: "Anima"}
	c.Traits["Vanguard"] = catalog.Trait{Key: "Vanguard", Name: "Vanguard"}
	c.Traits["Sniper"] = catalog.Trait{Key: "Sniper", Name: "Sniper"}

	c.Champions["Illaoi"] = catalog.Champion{Key: "Illaoi", Traits: []string{"Anima", "Vanguard"}}
	c.Champions["Jinx"] = catalog.Champion{Key: "Jinx", Traits: []string{"Anima"}}
	c.Champions["Caitlyn"] = catalog.Champion{Key: "Caitlyn", Traits: []string{"Sniper"}}

	illaoi := ChampionVector(c, "Illaoi")
	jinx := ChampionVector(c, "Jinx")
	caitlyn := ChampionVector(c, "Caitlyn")

	sharedTrait := Synergy(illaoi, jinx)
	noSharedTrait := Synergy(illaoi, caitlyn)

	if sharedTrait <= noSharedTrait {
		t.Errorf("Illaoi-Jinx (comparten Anima) = %.2f, debería ser mayor que Illaoi-Caitlyn (nada en común) = %.2f",
			sharedTrait, noSharedTrait)
	}
}
