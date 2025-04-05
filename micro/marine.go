package micro

import (
	"github.com/NatoBoram/BlackCompany/bot"
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/enums/terran"
	"github.com/aiseeq/s2l/protocol/enums/zerg"
)

func handleMarines(b *bot.Bot) {
	marines := b.Units.My.OfType(terran.Marine)
	if marines.Empty() {
		return
	}

	killChangelingsOnSight(b, marines)
}

func killChangelingsOnSight(b *bot.Bot, army scl.Units) {
	changelings := b.Enemies.All.OfType(zerg.Changeling)
	if army.Empty() || changelings.Empty() {
		return
	}

	for _, unit := range army {
		inSight := changelings.CloserThan(unit.SightRange(), unit.Point())
		if inSight.Empty() {
			continue
		}

		unit.Attack(inSight)
	}
}
