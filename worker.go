package main

import (
	"github.com/aiseeq/s2l/lib/scl"
)

// Workers handles idle workers
func (b *Bot) Workers() {
	idle := b.findWorkers().Filter(scl.Idle)

	townHalls := b.findTownHalls()
	if townHalls.Empty() {
		return
	}

	mineralFields := b.findUnsaturatedMineralFieldsNearTownHalls(townHalls)
	if idle.Exists() && mineralFields.Exists() {
		b.FillMineralsUpTo2(&idle, townHalls, mineralFields)
	}

	vespeneGeysers := b.findUnsaturatedVespeneGeysersNearTownHalls(townHalls)
	if idle.Exists() && vespeneGeysers.Exists() {
		b.FillGases(&idle, townHalls, vespeneGeysers)
	}

	if idle.Exists() && mineralFields.Exists() {
		b.FillMineralsUpTo3(&idle, townHalls, mineralFields)
	}
}
