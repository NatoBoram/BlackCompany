package macro

import (
	"github.com/NatoBoram/BlackCompany/bot"
	"github.com/NatoBoram/BlackCompany/filter"
	"github.com/NatoBoram/BlackCompany/log"
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/enums/terran"
)

// defenseWaveStep assigns an attack wave to the defense of the bases.
var defenseWaveStep = bot.BuildStep{
	Name: "Defense Wave",
	Predicate: func(b *bot.Bot) bool {
		return true
	},
	Execute: func(b *bot.Bot) {
		enemyInBases := b.FindEnemiesInBases()
		if len(enemyInBases) == 0 {
			return
		}

		base, enemies := b.MostThreatenedBase(enemyInBases)
		if base == nil || enemies.Empty() {
			return
		}

		// Where is the enemy cluster at that base?
		cluster := bot.FindClusterAtBase(base, enemies)
		if cluster.Empty() {
			log.Error("No cluster found at base %v", base.Point())
			return
		}

		inWaves := b.State.AttackWaves.Units(b)
		marines := b.Units.My.OfType(terran.Marine).Filter(scl.Ready, filter.NotIn(inWaves))
		if marines.Empty() {
			return
		}

		wave := bot.AttackWave{
			Tags:   marines.Tags(),
			Target: cluster.Center(),
		}
		b.State.AttackWaves = append(b.State.AttackWaves, wave)
		log.Info("Sending %d marines defend base at %v", marines.Len(), base.Point())
	},
	Next: func(b *bot.Bot) bool {
		return true
	},
}
