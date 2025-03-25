package main

import (
	"fmt"
	"math"
	"runtime/debug"

	"github.com/NatoBoram/BlackCompany/bot"
	"github.com/NatoBoram/BlackCompany/log"
	"github.com/NatoBoram/BlackCompany/micro"
	"github.com/aiseeq/s2l/lib/point"
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/api"
	"github.com/aiseeq/s2l/protocol/client"
)

func main() {
	env, err := loadEnv()
	if err != nil {
		log.Fatal("failed to load environment variables: %v", err)
	}

	cfg, err := launch(env)
	if err != nil {
		log.Fatal("failed to launch the game: %v", err)
	}

	flags := loadFlags()
	if flags.Replay == "" {
		runAgent(cfg.Client)
	}
}

// launch launches the game. It'll check for PROTON_PATH before launching the
// game in Proton or fallback to s2l's default behaviour.
func launch(env *Env) (*client.GameConfig, error) {
	bot := client.NewParticipant(api.Race_Terran, "BlackCompany")
	cpu := client.NewComputer(api.Race_Random, api.Difficulty_Hard, api.AIBuild_RandomBuild)

	if env.PROTON_PATH != "" && env.STEAM_COMPAT_DATA_PATH != "" {
		flags := loadFlags()

		paths, err := sc2Paths(env)
		if err != nil {
			return nil, fmt.Errorf("failed to get StarCraft II paths: %w", err)
		}

		if err = launchProton(paths, flags); err != nil {
			return nil, fmt.Errorf("failed to launch StarCraft II using Proton: %w", err)
		}

		if flags.Replay != "" && flags.Map != "" {
			return replayConfig(flags)
		}

		return protonConfig(bot, cpu), nil
	}

	return client.LaunchAndJoin(bot, cpu), nil
}

// runAgent creates a bot and runs it.
func runAgent(c *client.Client) {
	bot := &bot.Bot{
		Bot: scl.New(c, bot.OnUnitCreated),
		State: bot.BotState{
			CcForExp:            make(map[api.UnitTag]point.Point),
			CcForOrbitalCommand: 0,
			AttackWaves:         bot.AttackWaves{},
		},
	}

	bot.FramesPerOrder = 16
	bot.LastLoop = -math.MaxInt

	stop := make(chan struct{})
	bot.Init(stop)

	// Print version information for everyone to enjoy
	info, ok := debug.ReadBuildInfo()
	if ok {
		message := fmt.Sprintf("BlackCompany %s", info.Main.Version)
		bot.Actions.ChatSend(message, api.ActionChat_Team)
	}

	bot.Observe()
	for bot.Client.Status == api.Status_in_game {
		bot.Step()
		micro.Step(bot)

		step := api.RequestStep{Count: uint32(bot.FramesPerOrder)}
		if _, err := bot.Client.Step(step); err != nil {
			if err.Error() == "Not in a game" {
				break
			}

			log.Error("An unknown error occurred while stepping: %v", err)
			break
		}

		bot.Observe()
	}

	stop <- struct{}{}
}
