package catalog

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadJSON(t *testing.T) {
	dir := t.TempDir()

	champJSON := `{"season":"set17","champions":[{"key":"Illaoi","name":"Illaoi","cost":3,"traits":["Anima","Shepherd","Vanguard"]}]}`
	if err := os.WriteFile(filepath.Join(dir, "champions.json"), []byte(champJSON), 0644); err != nil {
		t.Fatal(err)
	}

	traitJSON := `{"season":"set17","traits":[{"key":"Anima","name":"Anima","styles":[{"style":"bronze","min":3},{"style":"gold","min":5}]}]}`
	if err := os.WriteFile(filepath.Join(dir, "traits.json"), []byte(traitJSON), 0644); err != nil {
		t.Fatal(err)
	}

	c, err := LoadJSON(dir)
	if err != nil {
		t.Fatalf("LoadJSON() error = %v", err)
	}

	illaoi, ok := c.Champions["Illaoi"]
	if !ok {
		t.Fatal("Illaoi no se cargo")
	}
	if illaoi.Cost != 3 {
		t.Errorf("Illaoi.Cost = %d, queria 3", illaoi.Cost)
	}

	anima, ok := c.Traits["Anima"]
	if !ok {
		t.Fatal("Anima no se cargo")
	}
	if len(anima.Breakpoints) != 2 {
		t.Fatalf("Anima tiene %d breakpoints, queria 2", len(anima.Breakpoints))
	}
	if anima.Breakpoints[1].Min != 5 {
		t.Errorf("segundo breakpoint de Anima .Min = %d, queria 5", anima.Breakpoints[1].Min)
	}
}
