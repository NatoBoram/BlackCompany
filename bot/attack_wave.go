package bot

import (
	"github.com/aiseeq/s2l/lib/point"
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/enums/ability"
)

type AttackWaves []AttackWave

type AttackWave struct {
	Tags   scl.Tags
	Target point.Point
}

func (a AttackWaves) Units(b *Bot) scl.Units {
	if a == nil {
		return nil
	}

	var units scl.Units
	for _, wave := range a {
		units = append(units, wave.Units(b)...)
	}

	return units
}

func (a *AttackWave) Units(b *Bot) scl.Units {
	if a == nil || a.Tags == nil {
		return nil
	}

	return b.Units.MyAll.ByTags(a.Tags)
}

func (a *AttackWave) Trim(b *Bot) scl.Units {
	units := a.Units(b)
	if units.Empty() {
		return nil
	}

	a.Tags = units.Tags()
	return units
}

func (a *AttackWave) Step(b *Bot) {
	units := a.Units(b)
	if units.Empty() {
		return
	}

	units.CommandPos(ability.Attack_Attack, a.Target)
}

// AttackWaves handles attack waves.
func (b *Bot) AttackWaves() {
	// Initialize attack waves if nil
	if b.State.AttackWaves == nil {
		b.State.AttackWaves = make(AttackWaves, 0)
		return
	}

	newWaves := make(AttackWaves, 0, len(b.State.AttackWaves))

	for _, wave := range b.State.AttackWaves {
		units := wave.Trim(b)
		if units.Empty() {
			continue
		}

		wave.Step(b)
		newWaves = append(newWaves, wave)
	}

	b.State.AttackWaves = newWaves
}
