package wheel_test

import (
	"testing"

	wheel "github.com/NatoBoram/BlackCompany/wheel"
)

func TestRandomIn(t *testing.T) {
	strings := []string{"a"}

	got := wheel.RandomIn(strings)
	expected := "a"
	if got != expected {
		t.Errorf("wheel.RandomIn(strings) = %v, expected %v", got, expected)
	}
}
