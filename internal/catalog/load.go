package catalog

import (
	"embed"
	"encoding/json"
)

type championsFile struct {
	Champions []Champion `json:"champions"`
}

type traitsFile struct {
	Traits []Trait `json:"traits"`
}

func loadFromBytes(champData, traitData []byte) (*Catalog, error) {
	c := New()

	var cf championsFile
	if err := json.Unmarshal(champData, &cf); err != nil {
		return nil, err
	}
	for _, champ := range cf.Champions {
		c.Champions[champ.Key] = champ
	}

	var tf traitsFile
	if err := json.Unmarshal(traitData, &tf); err != nil {
		return nil, err
	}
	for _, trait := range tf.Traits {
		c.Traits[trait.Key] = trait
	}

	return c, nil
}

//go:embed data/champions.json data/traits.json
var embeddedData embed.FS

// LoadEmbedded carga el catálogo incluido en el binario desde data/.
func LoadEmbedded() (*Catalog, error) {
	champData, err := embeddedData.ReadFile("data/champions.json")
	if err != nil {
		return nil, err
	}
	traitData, err := embeddedData.ReadFile("data/traits.json")
	if err != nil {
		return nil, err
	}
	return loadFromBytes(champData, traitData)
}
