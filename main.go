package main

import (
	"fmt"
	"log"

	"github.com/aiseeq/s2l/protocol/api"
	"github.com/aiseeq/s2l/protocol/client"
)

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

func main() {
	env, err := loadEnv()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := launch(env)
	if err != nil {
		log.Fatal(err)
	}

	runAgent(cfg.Client)
}

func launch(env *Env) (*client.GameConfig, error) {
	bot := client.NewParticipant(api.Race_Terran, "sc2bot")
	cpu := client.NewComputer(api.Race_Random, api.Difficulty_Easy, api.AIBuild_RandomBuild)

	if env.PROTON_PATH != "" && env.STEAM_COMPAT_DATA_PATH != "" {
		paths, err := sc2Paths(env)
		if err != nil {
			return nil, fmt.Errorf("failed to get StarCraft II paths: %w", err)
		}

		err = launchProton(paths)
		if err != nil {
			return nil, fmt.Errorf("failed to launch StarCraft II using Proton: %w", err)
		}

		return protonConfig(bot, cpu), nil
	}

	return client.LaunchAndJoin(bot, cpu), nil
}
