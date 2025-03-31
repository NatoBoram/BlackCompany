package macro

import "github.com/NatoBoram/BlackCompany/bot"

func attackWaveStep(config *AttackWaveConfig) *bot.BuildStep {
	return &bot.BuildStep{
		Name:      config.Name,
		Predicate: config.Predicate,
		Execute:   config.Execute,
		Next: func(b *bot.Bot) bool {
			return true
		},
	}
}
