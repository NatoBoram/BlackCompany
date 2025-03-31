package main

import (
	"fmt"
	"math"

	"github.com/NatoBoram/BlackCompany/bot"
	"github.com/NatoBoram/BlackCompany/log"
	"github.com/NatoBoram/BlackCompany/macro"
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
	b := &bot.Bot{
		Bot: scl.New(c, bot.OnUnitCreated),
		State: bot.BotState{
			CcForExp:            make(map[api.UnitTag]point.Point),
			CcForOrbitalCommand: 0,
			AttackWaves:         bot.AttackWaves{},
		},
	}

	b.FramesPerOrder = 16
	b.LastLoop = -math.MaxInt

	stop := make(chan struct{})
	b.Init(stop)
	b.Observe()
	b.InitState()

	var lastStep string
	for b.Client.Status == api.Status_in_game {
		b.Step()
		lastStep = macro.Step(b, &macro.Standard, lastStep)
		micro.Step(b)

		// Once a step is done, send it to the game
		b.Cmds.Process(&b.Actions)
		if len(b.Actions) > 0 {
			if _, err := b.Client.Action(api.RequestAction{Actions: b.Actions}); err != nil {
				log.Warn("Failed to send actions: %v", err)
			}

			b.Actions = nil
		}

		step := api.RequestStep{Count: uint32(b.FramesPerOrder)}
		if _, err := b.Client.Step(step); err != nil {
			if err.Error() == "Not in a game" {
				break
			}

			log.Error("An unknown error occurred while stepping: %v", err)
			break
		}

		b.Observe()
	}

	stop <- struct{}{}
}
