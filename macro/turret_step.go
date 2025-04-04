package macro

import (
	"github.com/NatoBoram/BlackCompany/bot"
	"github.com/NatoBoram/BlackCompany/filter"
	"github.com/NatoBoram/BlackCompany/log"
	"github.com/NatoBoram/BlackCompany/wheel"
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/enums/ability"
	"github.com/aiseeq/s2l/protocol/enums/terran"
)

// turretStep builds missile turrets in mineral lines then around buildings when
// flying enemies are detected
var turretStep = bot.BuildStep{
	Name: "Missile Turret",
	Predicate: func(b *bot.Bot) bool {
		if b.Units.My.OfType(terran.EngineeringBay).Filter(scl.Ready).Empty() {
			return false
		}

		// Wait for all ordered turrets to be at least started before going for the
		// next one
		ordered := b.Units.My.All().Filter(filter.IsOrderedToTag(ability.Build_MissileTurret, 0))
		if ordered.Exists() {
			return false
		}

		return b.CanBuy(ability.Build_MissileTurret)
	},

	Execute: func(b *bot.Bot) {
		turrets := b.Units.My.OfType(terran.MissileTurret)

		townHalls := b.FindTownHalls().Filter(filter.IsCcAtExpansion(b.State.CcForExp))
		if townHalls.Empty() {
			return
		}

		// Handle mineral lines first
		for _, th := range townHalls {
			unprotected := b.Units.Minerals.All().
				CloserThan(scl.ResourceSpreadDistance, th).
				Filter(filter.HasMinerals, filter.NotInTurretRange(turrets))
			if unprotected.Empty() {
				continue
			}

			center := unprotected.Center()
			pos := b.WhereToBuild(center, scl.S2x2, terran.MissileTurret, ability.Build_MissileTurret)
			if pos == nil {
				continue
			}

			worker := b.FindIdleOrGatheringWorkers().ClosestTo(pos)
			if worker == nil {
				continue
			}

			log.Info("Building missile turret at %v", pos)
			worker.CommandPos(ability.Build_MissileTurret, pos)
			b.DeductResources(ability.Build_MissileTurret)
			return
		}

		if !b.State.DetectedEnemyAirArmy {
			return
		}

		// At this point, there should be turrets in mineral lines. Otherwise, we're
		// lacking in minerals and we might not want turrets anyway.
		if turrets.Empty() {
			return
		}
		turret := turrets.First()

		unprotected := b.Units.My.All().Filter(scl.Structure, filter.NotInTurretRange(turrets))
		if unprotected.Empty() {
			return
		}

		random := wheel.RandomIn(unprotected)
		nearby := unprotected.CloserThan(turret.AirRange(), random)
		center := nearby.Center()
		pos := b.WhereToBuild(center, scl.S2x2, terran.MissileTurret, ability.Build_MissileTurret)
		if pos == nil {
			return
		}

		worker := b.FindIdleOrGatheringWorkers().ClosestTo(pos)
		if worker == nil {
			return
		}

		log.Info("Building missile turret at %v", pos)
		worker.CommandPos(ability.Build_MissileTurret, pos)
		b.DeductResources(ability.Build_MissileTurret)
	},

	Next: func(b *bot.Bot) bool {
		return true
	},
}
