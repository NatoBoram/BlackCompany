package bot

import "github.com/aiseeq/s2l/lib/scl"

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
