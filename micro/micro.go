package micro

import "github.com/NatoBoram/BlackCompany/bot"

func Step(b *bot.Bot) {
	handleAttackWaves(b)
	handleDefense(b)
	handleTownHalls(b)
	handleWorkers(b)
}
