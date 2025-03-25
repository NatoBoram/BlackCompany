package bot

import (
	"github.com/NatoBoram/BlackCompany/filter"
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/api"
)

// acknowledgeMiners saves metadata about miners.
func (b *Bot) acknowledgeMiners() {
	townHalls := b.FindTownHalls().Filter(filter.IsCcAtExpansion(b.State.CcForExp))
	if townHalls.Empty() {
		return
	}

	miners := b.FindMiners()
	if miners.Empty() {
		return
	}

	seen := make(scl.Units, 0, miners.Len())

	// Let's do the CCForMiner first. For each town hall, if the SCV is gathering
	// or returning, then we put it in `CCForMiner`. This is simple enough that
	// we can afford to delete the previous list.
	b.Miners.CCForMiner = map[api.UnitTag]api.UnitTag{}
	for _, th := range townHalls {
		resources := b.findResourcesNearTownHalls(scl.Units{th})
		targets := filter.ToTags(resources)
		targets = append(targets, th.Tag)

		for _, miner := range miners {
			maybeAssigned := filter.HasAnyTargetTag(targets)(miner)
			workingNearby := miner.IsCloserThan(th.SightRange(), th) && filter.IsGatheringOrReturning(miner)

			if maybeAssigned || workingNearby {
				b.Miners.CCForMiner[miner.Tag] = th.Tag
				b.Miners.LastSeen[miner.Tag] = b.Loop

				seen = append(seen, miner)
			}
		}
	}

	// Now let's do the MineralForMiner. For each mineral field, if there's a
	// miner gathering from it, then we put it in `MineralForMiner` and we delete
	// it from `GasForMiner`.
	mineralFields := b.FindMineralFieldsNearTownHalls(townHalls)
	for _, mf := range mineralFields {
		for _, miner := range miners.Filter(filter.HasTargetTag(mf.Tag)) {
			b.Miners.MineralForMiner[miner.Tag] = mf.Tag
			b.Miners.LastSeen[miner.Tag] = b.Loop

			delete(b.Miners.GasForMiner, miner.Tag)
			seen = append(seen, miner)
		}
	}

	vespeneGeysers := b.FindClaimedVespeneGeysersNearTownHalls(townHalls)
	for _, gas := range vespeneGeysers {
		for _, miner := range miners.Filter(filter.HasTargetTag(gas.Tag)) {
			b.Miners.GasForMiner[miner.Tag] = gas.Tag
			b.Miners.LastSeen[miner.Tag] = b.Loop

			delete(b.Miners.MineralForMiner, miner.Tag)
			seen = append(seen, miner)
		}
	}

	// Finally, we need to cleanup miners that we didn't see this time by deleting
	// them from all the maps.
	//
	// Some miners are somehow not tagged to a mineral field, vespene geyser, or
	// town hall, but are still in the process of gathering or returning. Let's
	// keep them.
	for _, miner := range miners.Filter(scl.NotReady) {
		if seen.ByTag(miner.Tag) == nil {
			delete(b.Miners.CCForMiner, miner.Tag)
			delete(b.Miners.GasForMiner, miner.Tag)
			delete(b.Miners.MineralForMiner, miner.Tag)
		}
	}

	// logger.Info(
	// 	"Miners: %d, CCForMiner: %d, GasForMiner: %d, MineralForMiner: %d",
	// 	miners.Len(), len(b.Miners.CCForMiner), len(b.Miners.GasForMiner), len(b.Miners.MineralForMiner),
	// )
}

func (b *Bot) resetMinerAcknowledgements() {
	// CCForMiner[miner.Tag] = cc.Tag
	b.Miners.CCForMiner = map[api.UnitTag]api.UnitTag{}
	// GasForMiner[miner.Tag] = gas.Tag
	b.Miners.GasForMiner = map[api.UnitTag]api.UnitTag{}
	// MineralForMiner[miner.Tag] = mf.Tag
	b.Miners.MineralForMiner = map[api.UnitTag]api.UnitTag{}
	// Miners.LastSeen[miner.Tag] = b.Loop
	b.Miners.LastSeen = map[api.UnitTag]int{}
}
