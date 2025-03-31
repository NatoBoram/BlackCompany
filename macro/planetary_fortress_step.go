package macro

import (
	"math/rand"

	"github.com/NatoBoram/BlackCompany/bot"
	"github.com/NatoBoram/BlackCompany/filter"
	"github.com/NatoBoram/BlackCompany/log"
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/enums/ability"
	"github.com/aiseeq/s2l/protocol/enums/terran"
)

// planetaryFortressStep builds planetary fortresses as long as we have money
// for that.
var planetaryFortressStep = bot.BuildStep{
	Name: "Planetary Fortress",
	Predicate: func(b *bot.Bot) bool {
		if !b.CanBuy(ability.Morph_PlanetaryFortress) {
			return false
		}

		if b.Units.My.OfType(terran.CommandCenter).Filter(scl.Ready).Empty() {
			return false
		}

		if b.Units.My.OfType(terran.EngineeringBay).Empty() {
			return false
		}

		return true
	},

	Execute: func(b *bot.Bot) {
		if b.State.CcForPlanetaryFortress == 0 {

			commandCenters := b.Units.My.OfType(terran.CommandCenter).Filter(
				scl.Ready, scl.Ground,
				filter.IsCcAtExpansion(b.State.CcForExp),
				filter.IsNotOrderedToAny(ability.Morph_OrbitalCommand, ability.Morph_PlanetaryFortress),
				filter.IsNotTags(b.State.CcForOrbitalCommand, b.State.CcForPlanetaryFortress),
			)
			if commandCenters.Empty() {
				return
			}

			randomCommandCenter := commandCenters[rand.Intn(len(commandCenters))]
			b.State.CcForPlanetaryFortress = randomCommandCenter.Tag
		}

		// Check if the marked command center is still valid
		reserved := b.Units.ByTag[b.State.CcForPlanetaryFortress]
		if reserved == nil || !reserved.Is(terran.CommandCenter, terran.CommandCenterFlying) {
			b.State.CcForPlanetaryFortress = 0
			return
		}

		// If it's ordered to do anything else, cancel it
		ordered := filter.IsOrderedTo(ability.Morph_PlanetaryFortress)(reserved)
		if len(reserved.Orders) > 0 && !ordered {
			reserved.Command(ability.Cancel_Last)
			return
		}

		// If it's not morphing yet, morph it
		if !ordered {
			log.Info("Morphing planetary fortress at %v", reserved.Point())
			reserved.Command(ability.Morph_PlanetaryFortress)
			b.DeductResources(ability.Morph_PlanetaryFortress)
			b.State.CcForOrbitalCommand = 0
		}
	},

	Next: func(b *bot.Bot) bool {
		return true
	},
}
