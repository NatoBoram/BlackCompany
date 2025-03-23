package bot

import (
	"math/rand"

	"github.com/NatoBoram/BlackCompany/filter"
	"github.com/NatoBoram/BlackCompany/log"
	"github.com/NatoBoram/BlackCompany/sight"
	"github.com/aiseeq/s2l/lib/point"
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/enums/ability"
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

// Trim removes units that are no longer alive from the attack wave.
func (a *AttackWave) Trim(b *Bot) scl.Units {
	units := a.Units(b)
	a.Tags = units.Tags()
	return units
}

// Step is called at every step of the game. This is the main loop of an attack
// wave.
func (a *AttackWave) Step(b *Bot) {
	units := a.Trim(b)
	if units.Empty() {
		return
	}

	units = a.Recenter(b, units)
	a.Advance(b, units)

	a.UpdateTarget(b)
}

// Recenter moves units that are too far from the wave towards the center of the
// unit group.
func (a *AttackWave) Recenter(b *Bot, units scl.Units) scl.Units {
	if units.Empty() {
		return units
	}

	center := units.Center()
	for _, u := range units {
		if u.Dist(center) > sight.LineOfSightScannerSweep.Float64() {
			dest := center.Towards(u, sight.LineOfSightScannerSweep.Float64()-1)

			if filter.IsNotOrderedToTarget(ability.Move, dest)(u) {
				u.CommandPos(ability.Move, dest)
			}

			units.Remove(u)
		}
	}

	return units
}

// Advance moves the attack wave towards the target.
func (a *AttackWave) Advance(b *Bot, units scl.Units) scl.Units {
	if units.Empty() {
		return units
	}

	for _, u := range units {
		if filter.IsNotOrderedToTarget(ability.Attack, a.Target)(u) {
			u.CommandPos(ability.Attack, a.Target)
		}

		units.Remove(u)
	}

	return units
}

// UpdateTarget changes the focus of this attack wave to something else.
func (a *AttackWave) UpdateTarget(b *Bot) {
	center := a.Units(b).Center()
	dist := a.Target.Dist(center)
	if dist > sight.LineOfSightScannerSweep.Float64() {
		return
	}

	// It's too close to the target, let's find an enemy unit then target it
	enemies := b.Units.Enemy.All()

	// Destroy all buildings to win the game
	buildings := enemies.Filter(scl.Structure)
	if buildings.Exists() {
		target := buildings.ClosestTo(center).Point()

		if target.Dist(a.Target) > sight.LineOfSightScannerSweep.Float64() {
			log.Info("Switching target to a building at %v", target)
		}

		a.Target = target
		return
	}

	// Units are generally closer to buildings or they prevent from destroying
	// buildings
	units := enemies.Filter(scl.NotStructure)
	if units.Exists() {
		target := units.ClosestTo(center).Point()

		if target.Dist(a.Target) > sight.LineOfSightScannerSweep.Float64() {
			log.Info("Switching target to a unit at %v", target)
		}

		a.Target = target
		return
	}

	// Expansions are obvious choices for building locations
	target := b.Locs.EnemyExps[rand.Intn(len(b.Locs.EnemyExps))]

	if target.Dist(a.Target) > sight.LineOfSightScannerSweep.Float64() {
		log.Info("Switching target to an enemy expansion at %v", target)
	}

	a.Target = target
}

// AttackWaves handles attack waves.
func (b *Bot) AttackWaves() {
	keepWaves := make(AttackWaves, 0, len(b.State.AttackWaves))

	for _, wave := range b.State.AttackWaves {
		units := wave.Trim(b)
		if units.Empty() {
			continue
		}

		wave.Step(b)
		keepWaves = append(keepWaves, wave)
	}

	b.State.AttackWaves = keepWaves
}
