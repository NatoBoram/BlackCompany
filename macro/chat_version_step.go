package macro

import (
	"fmt"
	"runtime/debug"

	"github.com/NatoBoram/BlackCompany/bot"
	"github.com/aiseeq/s2l/protocol/api"
)

// chatVersionStep announces the version number of the bot.
func chatVersionStep() *bot.BuildStep {
	announced := false

	return &bot.BuildStep{
		Name: "Announce version number",
		Predicate: func(b *bot.Bot) bool {
			return true
		},
		Execute: func(b *bot.Bot) {
			if announced {
				return
			}

			// Print version information for everyone to enjoy
			info, ok := debug.ReadBuildInfo()
			if ok {
				message := fmt.Sprintf("BlackCompany %s", info.Main.Version)
				b.Actions.ChatSend(message, api.ActionChat_Team)
			}

			announced = true
		},
		Next: func(b *bot.Bot) bool {
			return true
		},
	}
}
