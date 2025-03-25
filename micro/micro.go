package micro

import "github.com/NatoBoram/BlackCompany/bot"

func Step(b *bot.Bot) {
	handleAttackWaves(b)
	handleTownHalls(b)
	handleWorkers(b)
}
