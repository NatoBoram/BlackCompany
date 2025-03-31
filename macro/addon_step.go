package macro

import (
	"math/rand"

	"github.com/NatoBoram/BlackCompany/bot"
	"github.com/NatoBoram/BlackCompany/filter"
	"github.com/NatoBoram/BlackCompany/log"
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/api"
	"github.com/aiseeq/s2l/protocol/enums/ability"
)

// addonStep manages building add-ons for buildings.
//
// TODO: Check if there's enough space for the reactor, and if not, fly the
// building somewhere safe.
func addonStep(name string, buildingId api.UnitTypeID, addonId api.UnitTypeID, abilityId api.AbilityID, quantity int) *bot.BuildStep {
	return &bot.BuildStep{
		Name: stepName(name, quantity),
		Predicate: func(b *bot.Bot) bool {
			buildings := b.Units.My.OfType(buildingId).Filter(scl.Ready, scl.Ground, scl.NoAddon)
			if buildings.Empty() {
				return false
			}

			if !b.CanBuy(abilityId) {
				return false
			}

			if b.Units.My.OfType(addonId).Len() >= quantity {
				return false
			}

			return true
		},

		Execute: func(b *bot.Bot) {
			buildings := b.Units.My.OfType(buildingId).Filter(scl.Ready, scl.Ground, scl.NoAddon)
			if buildings.Empty() {
				return
			}

			// If there's no building marked for add-on, mark one
			if b.State.BuildingForAddOn == 0 {
				randomBuilding := buildings[rand.Intn(len(buildings))]
				b.State.BuildingForAddOn = randomBuilding.Tag
			}

			// Check if the marked building is still valid
			reserved := b.Units.ByTag[b.State.BuildingForAddOn]
			if reserved == nil || !reserved.Is(buildingId) || reserved.AddOnTag != 0 {
				b.State.BuildingForAddOn = 0
				return
			}

			// If it's ordered to do anything else, cancel it
			ordered := filter.IsOrderedTo(abilityId)(reserved)
			if len(reserved.Orders) > 0 && !ordered {
				reserved.Command(ability.Cancel_Last)
				return
			}

			log.Info("Building %s at %v", name, reserved.Point())
			reserved.Command(abilityId)

			// In case it fails, queue the add-on to a new location
			elsewhere := b.WhereToBuild(reserved.Point(), scl.S5x3, addonId, abilityId)
			reserved.CommandPosQueue(abilityId, elsewhere)

			b.DeductResources(abilityId)
			b.State.BuildingForAddOn = 0
		},

		Next: func(b *bot.Bot) bool {
			return b.Units.My.OfType(addonId).Len() >= quantity
		},
	}
}
