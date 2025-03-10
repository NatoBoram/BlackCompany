package main

import (
	"fmt"
	"log"
	"math"

	"github.com/aiseeq/s2l/lib/point"
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/api"
	"github.com/aiseeq/s2l/protocol/client"
)

func main() {
	env, err := loadEnv()
	if err != nil {
		log.Fatalf("failed to load environment variables: %v", err)
	}

	cfg, err := launch(env)
	if err != nil {
		log.Fatalf("failed to launch the game: %v", err)
	}

	runAgent(cfg.Client)
}

// launch launches the game. It'll check for PROTON_PATH before launching the
// game in Proton or fallback to s2l's default behaviour.
func launch(env *Env) (*client.GameConfig, error) {
	bot := client.NewParticipant(api.Race_Terran, "BlackCompany")
	cpu := client.NewComputer(api.Race_Random, api.Difficulty_Easy, api.AIBuild_RandomBuild)

	if env.PROTON_PATH != "" && env.STEAM_COMPAT_DATA_PATH != "" {
		paths, err := sc2Paths(env)
		if err != nil {
			return nil, fmt.Errorf("failed to get StarCraft II paths: %w", err)
		}

		flags := loadFlags()

		if err = launchProton(paths, flags); err != nil {
			return nil, fmt.Errorf("failed to launch StarCraft II using Proton: %w", err)
		}

		return protonConfig(bot, cpu), nil
	}

	return client.LaunchAndJoin(bot, cpu), nil
}

// runAgent creates a bot and runs it.
func runAgent(c *client.Client) {
	bot := &Bot{
		Bot: scl.New(c, OnUnitCreated),
		state: BotState{
			CcForExp:            make(map[api.UnitTag]point.Point),
			CcForOrbitalCommand: 0,
		},
	}

	bot.FramesPerOrder = 16
	bot.LastLoop = -math.MaxInt

	stop := make(chan struct{})
	bot.Init(stop)

	bot.Observe()
	for bot.Client.Status == api.Status_in_game {
		bot.Step()

		step := api.RequestStep{Count: uint32(bot.FramesPerOrder)}
		if _, err := bot.Client.Step(step); err != nil {
			if err.Error() == "Not in a game" {
				break
			}

			log.Printf("An unknown error occurred while stepping: %v", err)
			break
		}

		bot.Observe()
	}

	stop <- struct{}{}
}
