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

func orbitalCommandStep(quantity int) *bot.BuildStep {
	return &bot.BuildStep{
		Name: "Orbital Command",
		Predicate: func(b *bot.Bot) bool {
			barracks := b.Units.My.OfType(terran.Barracks).Filter(scl.Ready, scl.Ground)
			if barracks.Empty() {
				return false
			}

			if !b.CanBuy(ability.Morph_OrbitalCommand) {
				return false
			}

			orbitalCommands := b.Units.My.OfType(terran.OrbitalCommand, terran.OrbitalCommandFlying)
			inProgress := b.Units.My.OfType(terran.CommandCenter).Filter(filter.IsOrderedTo(ability.Morph_OrbitalCommand))
			if orbitalCommands.Len()+inProgress.Len() >= quantity {
				b.State.CcForOrbitalCommand = 0
				return false
			}

			return true
		},

		Execute: func(b *bot.Bot) {
			if b.State.CcForOrbitalCommand == 0 {
				// There's no command center marked for morphing into an orbital
				// command, so let's mark one
				commandCenters := b.Units.My.
					OfType(terran.CommandCenter).
					Filter(scl.Ready,
						filter.IsCcAtExpansion(b.State.CcForExp),
						filter.IsNotOrderedToAny(ability.Morph_OrbitalCommand, ability.Morph_PlanetaryFortress),
						filter.IsNotTags(b.State.CcForOrbitalCommand, b.State.CcForPlanetaryFortress),
					)
				if commandCenters.Empty() {
					return
				}

				randomCommandCenter := commandCenters[rand.Intn(len(commandCenters))]
				b.State.CcForOrbitalCommand = randomCommandCenter.Tag
			}

			// Check if the marked command center is still valid
			reserved := b.Units.ByTag[b.State.CcForOrbitalCommand]
			if reserved == nil || !reserved.Is(terran.CommandCenter, terran.CommandCenterFlying) {
				b.State.CcForOrbitalCommand = 0
				return
			}

			// If it's ordered to do anything else, cancel it
			ordered := filter.IsOrderedTo(ability.Morph_OrbitalCommand)(reserved)
			if len(reserved.Orders) > 0 && !ordered {
				reserved.Command(ability.Cancel_Last)
				return
			}

			// If it's not morphing yet, morph it
			if !ordered {
				log.Info("Morphing orbital command at %v", reserved.Point())
				reserved.Command(ability.Morph_OrbitalCommand)
				b.DeductResources(ability.Morph_OrbitalCommand)
				b.State.CcForOrbitalCommand = 0
			}
		},

		Next: func(b *bot.Bot) bool {
			orbitalCommands := b.Units.My.OfType(terran.OrbitalCommand, terran.OrbitalCommandFlying)
			inProgress := b.Units.My.OfType(terran.CommandCenter).Filter(filter.IsOrderedTo(ability.Morph_OrbitalCommand))
			return orbitalCommands.Len()+inProgress.Len() >= quantity
		},
	}
}
