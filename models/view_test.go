package models

import "testing"

func TestLcm(t *testing.T) {
	a := lcm([]int{1, 2, 3})

	if a != 6 {
		t.Errorf("Expected 6 but got %d", a)
	}

	a = lcm([]int{2, 3, 4})

	if a != 12 {
		t.Errorf("Expected 12 but got %d", a)
	}
}

func TestGridColumns(t *testing.T) {
	a := gridColumns([]int{1, 2, 3, 1, 1, 1, 1})

	if a != "1fr 6fr 3fr 3fr 2fr 2fr 2fr 6fr 6fr 6fr 6fr" {
		t.Errorf("Expected '1fr 6fr 3fr 3fr 2fr 2fr 2fr 6fr 6fr 6fr 6fr' but got %s", a)
	}
}
