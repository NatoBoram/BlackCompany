package micro

import (
	"math/rand"

	"github.com/NatoBoram/BlackCompany/bot"
	"github.com/NatoBoram/BlackCompany/filter"
	"github.com/NatoBoram/BlackCompany/log"
	"github.com/NatoBoram/BlackCompany/sight"
	"github.com/aiseeq/s2l/lib/point"
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/enums/ability"
)

func handleAttackWaves(b *bot.Bot) {
	keep := make(bot.AttackWaves, 0, len(b.State.AttackWaves))

	for _, wave := range b.State.AttackWaves {
		units := trimWave(b, &wave)
		if units.Empty() {
			continue
		}

		handleAttackWave(b, &wave)
		keep = append(keep, wave)
	}

	b.State.AttackWaves = keep
}

// trimWave removes units that are no longer alive from the attack wave.
func trimWave(b *bot.Bot, a *bot.AttackWave) scl.Units {
	units := a.Units(b)
	a.Tags = units.Tags()
	return units
}

func handleAttackWave(b *bot.Bot, a *bot.AttackWave) {
	units := trimWave(b, a)
	if units.Empty() {
		return
	}

	units = recenterWave(units, a.Target)
	advanceWave(a, units)

	updateWaveTarget(b, a)
}

// recenterWave moves units that are too far from the wave towards the center of
// the unit group.
func recenterWave(units scl.Units, target point.Point) scl.Units {
	if units.Empty() {
		return units
	}

	center := units.Center()
	decentered := units.FurtherThan(sight.LineOfSightScannerSweep.Float64(), center)
	if decentered.Empty() {
		return units
	}

	if float64(decentered.Len()) <= float64(units.Len())*0.2 {
		return units
	}

	towards := center.Towards(target, 1)

	log.Info("Recentering %d units", decentered.Len())
	for _, unit := range units {
		if filter.IsNotOrderedToTarget(ability.Move, towards)(unit) {
			unit.CommandPos(ability.Move, towards)
		}

		units.Remove(unit)
	}

	return units
}

// advanceWave moves the attack wave towards the target.
func advanceWave(a *bot.AttackWave, units scl.Units) scl.Units {
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

// updateWaveTarget changes the focus of this attack wave to something else.
func updateWaveTarget(b *bot.Bot, a *bot.AttackWave) {
	center := a.Units(b).Center()
	dist := a.Target.Dist(center)
	if dist > sight.LineOfSightScannerSweep.Float64() {
		return
	}

	// It's too close to the target, let's find an enemy unit then target it
	enemies := b.Units.Enemy.All()

	// If someone's at home, gotta defend it!
	cluster := b.FindEnemyClusterAtHome()
	if cluster.Exists() {
		target := cluster.Center()
		if target.Dist(a.Target) > sight.LineOfSightScannerSweep.Float64() {
			log.Info("Switching target to defend base at %v", target)
		}

		a.Target = target
		return
	}

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
