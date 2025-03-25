package micro

import (
	"github.com/NatoBoram/BlackCompany/bot"
	"github.com/NatoBoram/BlackCompany/log"
	"github.com/aiseeq/s2l/lib/scl"
)

// handleWorkers handles idle workers
func handleWorkers(b *bot.Bot) {
	idle := b.FindWorkers().Filter(scl.Idle)
	if idle.Empty() {
		return
	}

	townHalls := b.FindTownHalls()
	if townHalls.Empty() {
		return
	}

	log.Info("Sending %d workers back to work", idle.Len())

	mineralFields := b.FindUnsaturatedMineralFieldsNearTownHalls(townHalls)
	if idle.Exists() && mineralFields.Exists() {
		b.FillMineralsUpTo2(&idle, townHalls, mineralFields)
	}

	vespeneGeysers := b.FindUnsaturatedVespeneGeysersNearTownHalls(townHalls)
	if idle.Exists() && vespeneGeysers.Exists() {
		b.FillGases(&idle, townHalls, vespeneGeysers)
	}

	if idle.Exists() && mineralFields.Exists() {
		b.FillMineralsUpTo3(&idle, townHalls, mineralFields)
	}
}
