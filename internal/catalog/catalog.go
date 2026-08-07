package catalog

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Breakpoint struct {
	Style string `json:"style"`
	Min   int    `json:"min"`
}

type Trait struct {
	Key         string       `json:"key"`
	Name        string       `json:"name"`
	Breakpoints []Breakpoint `json:"styles"`
}

type Champion struct {
	Key    string   `json:"key"`
	Name   string   `json:"name"`
	Cost   int      `json:"cost"`
	Traits []string `json:"traits"`
}

type Catalog struct {
	Champions map[string]Champion
	Traits    map[string]Trait
}

func New() *Catalog {
	return &Catalog{
		Champions: make(map[string]Champion),
		Traits:    make(map[string]Trait),
	}
}

func (c *Catalog) TraitsOf(championKey string) []Trait {
	champ, ok := c.Champions[championKey]
	if !ok {
		return nil
	}
	traits := make([]Trait, 0, len(champ.Traits))
	for _, key := range champ.Traits {
		if t, ok := c.Traits[key]; ok {
			traits = append(traits, t)
		}
	}
	return traits
}

type championsFile struct {
	Champions []Champion `json:"champions"`
}

type traitsFile struct {
	Traits []Trait `json:"traits"`
}

func LoadJSON(dataDir string) (*Catalog, error) {
	c := New()

	champData, err := os.ReadFile(filepath.Join(dataDir, "champions.json"))
	if err != nil {
		return nil, err
	}
	var cf championsFile
	if err := json.Unmarshal(champData, &cf); err != nil {
		return nil, err
	}
	for _, champ := range cf.Champions {
		c.Champions[champ.Key] = champ
	}

	traitData, err := os.ReadFile(filepath.Join(dataDir, "traits.json"))
	if err != nil {
		return nil, err
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
