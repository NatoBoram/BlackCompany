package filter_test

import (
	"testing"

	"github.com/NatoBoram/BlackCompany/filter"
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/api"
)

func TestNotIn_IsIn(t *testing.T) {
	unit := &scl.Unit{Unit: api.Unit{Tag: 1}}
	units := scl.Units{unit}

	got := filter.NotIn(units)(unit)
	expected := false
	if got != false {
		t.Errorf("NotIn(units)(unit) = %v, expected %v", got, expected)
	}
}

func TestNotIn_IsOut(t *testing.T) {
	isIn := &scl.Unit{Unit: api.Unit{Tag: 1}}
	isOut := &scl.Unit{Unit: api.Unit{Tag: 2}}
	units := scl.Units{isIn}

	got := filter.NotIn(units)(isOut)
	expected := true
	if got != expected {
		t.Errorf("NotIn(units)(isOut) = %v, expected %v", got, expected)
	}
}

func TestNotIn_Empty(t *testing.T) {
	unit := &scl.Unit{Unit: api.Unit{Tag: 1}}
	empty := scl.Units{}

	got := filter.NotIn(empty)(unit)
	expected := true
	if got != expected {
		t.Errorf("NotIn(empty)(unit) = %v, expected %v", got, expected)
	}
}
