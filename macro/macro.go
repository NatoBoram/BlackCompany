package macro

import (
	"math/rand"

	"github.com/NatoBoram/BlackCompany/bot"
	"github.com/NatoBoram/BlackCompany/filter"
	"github.com/NatoBoram/BlackCompany/log"
	"github.com/aiseeq/s2l/lib/point"
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/api"
	"github.com/aiseeq/s2l/protocol/enums/ability"
	"github.com/aiseeq/s2l/protocol/enums/terran"
)

// Macro executes a strategy.
func Macro(b *bot.Bot, s *bot.Strategy) {
	for _, step := range s.Steps {
		if step.Predicate(b) {
			step.Execute(b)
		}

		if !step.Next(b) {
			break
		}
	}
}

func deductMarines(b *bot.Bot, barracks *scl.Unit) int {
	if !b.CanBuy(ability.Train_Marine) {
		return 0
	}

	// Confirmed that we're about to train one marine.
	b.DeductResources(ability.Train_Marine)

	if barracks.AddOnTag == 0 {
		return 1
	}

	addon := b.Units.ByTag[barracks.AddOnTag]
	if addon == nil || !addon.Is(terran.BarracksReactor) || !b.CanBuy(ability.Train_Marine) {
		return 1
	}

	// Confirmed that we're about to train two marines.
	b.DeductResources(ability.Train_Marine)
	return 2
}

func rallyPoint(b *bot.Bot) *point.Point {
	townHalls := b.FindTownHalls().Filter(filter.IsCcAtExpansion(b.State.CcForExp))
	if townHalls.Empty() {
		return nil
	}

	closest := townHalls.ClosestTo(b.Locs.EnemyStart)
	rally := closest.Towards(b.Locs.EnemyStart, closest.SightRange())
	return &rally
}

func build(b *bot.Bot, name string, buildingId api.UnitTypeID, abilityId api.AbilityID, size scl.BuildingSize) {
	if !b.CanBuy(abilityId) {
		return
	}

	townHalls := b.FindTownHalls().Filter(filter.IsCcAtExpansion(b.State.CcForExp))
	if townHalls.Empty() {
		return
	}

	randomTownHall := townHalls[rand.Intn(len(townHalls))]

	pos := b.WhereToBuild(randomTownHall.Point(), size, buildingId, abilityId)
	if pos == nil {
		return
	}

	builder := b.FindIdleOrGatheringWorkers().ClosestTo(pos)
	if builder == nil {
		return
	}

	if resource := b.FindResourcesNearTownHalls(townHalls).ClosestTo(pos); resource != nil {
		log.Info("Building %s at %v and queuing to gather at %v", name, *pos, resource.Point())

		builder.CommandPos(abilityId, pos)
		builder.CommandTagQueue(ability.Smart, resource.Tag)

		if resource.IsMineral() {
			b.Miners.MineralForMiner[builder.Tag] = resource.Tag
		}

		if resource.IsGeyser() {
			b.Miners.GasForMiner[builder.Tag] = resource.Tag
		}
	} else {
		log.Info("Building %s at %v", name, *pos)
		builder.CommandPos(abilityId, pos)
	}

	b.DeductResources(abilityId)
}
