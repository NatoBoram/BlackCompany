package bot

import (
	"github.com/aiseeq/s2l/lib/point"
	"github.com/aiseeq/s2l/lib/scl"
)

// AttackWave is a single attack wave, its units and its state. It should
// dictate the intent of its units, but the micro should be performed elsewhere.
type AttackWave struct {
	Tags   scl.Tags
	Target point.Point
}

type AttackWaves []AttackWave

func (a AttackWaves) Units(b *Bot) scl.Units {
	count := 0
	for _, wave := range a {
		count += wave.Tags.Len()
	}

	units := make(scl.Units, 0, count)
	for _, wave := range a {
		units = append(units, wave.Units(b)...)
	}

	return units
}

// Units gets the units in an attack wave
func (a *AttackWave) Units(b *Bot) scl.Units {
	if b.Units.MyAll.Empty() {
		b.Units.MyAll = make(scl.Units, 0)
	}

	return b.Units.MyAll.ByTags(a.Tags)
}
