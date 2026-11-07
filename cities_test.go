package main

import "testing"

func TestParseCity(t *testing.T) {
	for _, name := range []string{"Farg'ona", "FARG'ONA", "fargʻona", "Farg’ona", "  fargona "} {
		id, ok := ParseCity(name)
		if id != CityFergana || !ok {
			t.Errorf("ParseCity(%q) = %d, %v; want %d, true", name, id, ok, CityFergana)
		}
	}
	if id, ok := ParseCity("Moskva"); id != 0 || ok {
		t.Errorf("ParseCity(Moskva) = %d, %v; want 0, false", id, ok)
	}
}

func TestEveryCityHasAName(t *testing.T) {
	for id := CityTashkent; id <= CityNamangan; id++ {
		name, ok := CityNames[id]
		if !ok {
			t.Fatalf("city %d has no name", id)
		}
		if back, ok := ParseCity(name); back != id || !ok {
			t.Errorf("ParseCity(%q) = %d, %v; want %d", name, back, ok, id)
		}
	}
}
