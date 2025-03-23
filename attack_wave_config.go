package main

import (
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/enums/terran"
)

// AttackWaveConfig holds the configuration for launching an attack wave.
type AttackWaveConfig struct {
	// Name is for usage in a build step
	Name string

	// Predicate determines if this attack wave should be launched or not
	Predicate func(b *Bot) bool

	// Execute attributes units to an attack wave
	Execute func(b *Bot)
}

// firstWaveConfig puts marines into a group for launching a marine rush timing
// attack after combat shield is started. It's never executed again after the
// first wave is launched.
func firstWaveConfig() *AttackWaveConfig {
	launched := false

	return &AttackWaveConfig{
		Name: "First Attack Wave",
		Predicate: func(b *Bot) bool {
			return !launched
		},
		Execute: func(b *Bot) {
			if launched {
				return
			}

			inWaves := b.state.AttackWaves.Units(b)
			marines := b.Units.My.OfType(terran.Marine).Filter(scl.Ready, NotIn(inWaves))
			if marines.Empty() {
				return
			}

			wave := AttackWave{
				Tags:   marines.Tags(),
				Target: b.Locs.EnemyStart,
			}
			b.state.AttackWaves = append(b.state.AttackWaves, wave)

			launched = true
			logger.Info("Sending %d marines to enemy base %v", marines.Len(), b.Locs.EnemyStart)
		},
	}
}

func fullSupplyWaveConfig() *AttackWaveConfig {
	return &AttackWaveConfig{
		Name: "Full Supply Attack Wave",
		Predicate: func(b *Bot) bool {
			marines := b.Units.My.OfType(terran.Marine).Filter(scl.Ready, NotIn(b.state.AttackWaves.Units(b)))
			return b.Obs.PlayerCommon.FoodUsed >= 200 && marines.Len() >= 10
		},
		Execute: func(b *Bot) {
			marines := b.Units.My.OfType(terran.Marine).Filter(scl.Ready, NotIn(b.state.AttackWaves.Units(b)))
			if marines.Empty() {
				return
			}

			wave := AttackWave{
				Tags:   marines.Tags(),
				Target: marines.Center(),
			}

			b.state.AttackWaves = append(b.state.AttackWaves, wave)
			logger.Info("Preparing new attack wave with %d units", marines.Len())
		},
	}
}
