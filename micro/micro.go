package micro

import (
	"github.com/NatoBoram/BlackCompany/bot"
)

func Step(b *bot.Bot) {
	// Skip repeated frames
	if b.LastLoop != b.Loop {
		return
	}

	handleAttackWaves(b)
	handleTownHalls(b)
	handleWorkers(b)
}
