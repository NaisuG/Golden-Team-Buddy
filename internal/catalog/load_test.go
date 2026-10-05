package catalog

import "testing"

func TestLoadFromBytes(t *testing.T) {
	champJSON := `{"season":"set17","champions":[{"key":"Illaoi","name":"Illaoi","cost":3,"traits":["Anima","Shepherd","Vanguard"]}]}`

	traitJSON := `{"season":"set17","traits":[{"key":"Anima","name":"Anima","styles":[{"style":"bronze","min":3},{"style":"gold","min":5}]}]}`

	c, err := loadFromBytes([]byte(champJSON), []byte(traitJSON))
	if err != nil {
		t.Fatalf("loadFromBytes() error = %v", err)
	}

	illaoi, ok := c.Champions["Illaoi"]
	if !ok {
		t.Fatal("Illaoi no se cargó")
	}
	if illaoi.Cost != 3 {
		t.Errorf("Illaoi.Cost = %d, quería 3", illaoi.Cost)
	}

	anima, ok := c.Traits["Anima"]
	if !ok {
		t.Fatal("Anima no se cargó")
	}
	if len(anima.Breakpoints) != 2 {
		t.Fatalf("Anima tiene %d breakpoints, quería 2", len(anima.Breakpoints))
	}
	if anima.Breakpoints[1].Min != 5 {
		t.Errorf("segundo breakpoint de Anima: Min = %d, quería 5", anima.Breakpoints[1].Min)
	}
}

func TestLoadEmbedded(t *testing.T) {
	c, err := LoadEmbedded()
	if err != nil {
		t.Fatalf("LoadEmbedded() error = %v", err)
	}
	if len(c.Champions) == 0 {
		t.Error("LoadEmbedded() no cargó ningún campeón")
	}
	if len(c.Traits) == 0 {
		t.Error("LoadEmbedded() no cargó ningún trait")
	}
}
