package catalog

type Breakpoint struct {
	Style string
	Min   int
}

type Trait struct {
	Key         string
	Name        string
	Breakpoints []Breakpoint
}

type Champion struct {
	Key    string
	Name   string
	Cost   int
	Traits []string
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
	var traits []Trait
	for _, key := range champ.Traits {
		if t, ok := c.Traits[key]; ok {
			traits = append(traits, t)
		}
	}
	return traits
}
