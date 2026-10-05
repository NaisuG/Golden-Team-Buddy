package catalog

type Breakpoint struct {
	Style string `json:"style"`
	Min   int    `json:"min"`
}

// Trait es un rasgo del set con sus breakpoints, ordenados de menor a mayor.
type Trait struct {
	Key         string       `json:"key"`
	Name        string       `json:"name"`
	Breakpoints []Breakpoint `json:"styles"`
}

// Champion es un campeón del set. Cost es su costo en tienda (1 a 5).
type Champion struct {
	Key    string   `json:"key"`
	Name   string   `json:"name"`
	Cost   int      `json:"cost"`
	Traits []string `json:"traits"`
}

// Catalog contiene los campeones y traits del set, indexados por su key.
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
