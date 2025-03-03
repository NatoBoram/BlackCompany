package main

import (
	"fmt"
	"log"

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
	bot := client.NewParticipant(api.Race_Terran, "sc2bot")
	cpu := client.NewComputer(api.Race_Random, api.Difficulty_Easy, api.AIBuild_RandomBuild)

	if env.PROTON_PATH != "" && env.STEAM_COMPAT_DATA_PATH != "" {
		paths, err := sc2Paths(env)
		if err != nil {
			return nil, fmt.Errorf("failed to get StarCraft II paths: %w", err)
		}

		if err = launchProton(paths); err != nil {
			return nil, fmt.Errorf("failed to launch StarCraft II using Proton: %w", err)
		}

		return protonConfig(bot, cpu), nil
	}

	return client.LaunchAndJoin(bot, cpu), nil
}

// runAgent creates a bot and runs it.
func runAgent(c *client.Client) {
	bot := &Bot{c, nil}

	for bot.c.Status == api.Status_in_game {
		bot.Step()

		if _, err := c.Step(api.RequestStep{Count: uint32(3)}); err != nil {
			if err.Error() == "Not in a game" {
				break
			}

			log.Printf("An unknown error occurred while stepping: %v", err)
			break
		}

		bot.Observe()
	}
}
