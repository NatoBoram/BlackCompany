package main

import (
	"log"
	"os/exec"
	"path"
	"strconv"

	"github.com/aiseeq/s2l/protocol/api"
	"github.com/aiseeq/s2l/protocol/client"
)

// launchProton launches StarCraft II using Proton.
//
//   - https://levelup.gitconnected.com/guide-to-starcraft-ii-proto-api-264811da8a50
//   - https://gist.github.com/michaelbutler/f364276f4030c5f449252f2c4d960bd2
func launchProton(paths *Sc2Paths, flags Flags) error {
	client.SetExecutable(paths.Exe)

	args := []string{
		"run", paths.Exe,
		"-displaymode", strconv.Itoa(flags.DisplayMode),
		"-listen", flags.Listen,
		"-port", strconv.Itoa(flags.Port),
	}

	if flags.Realtime {
		log.Println("Running StarCraft II in real-time")
		args = append(args, "-realtime", strconv.FormatBool(flags.Realtime))
		client.SetRealtime()
	}
	if flags.Timeout > 0 {
		args = append(args, "-timeout", flags.Timeout.String())
		client.SetConnectTimeout(flags.Timeout)
	}

	if flags.WindowWidth > 0 {
		args = append(args, "-windowwidth", strconv.Itoa(flags.WindowWidth))
	}
	if flags.WindowHeight > 0 {
		args = append(args, "-windowheight", strconv.Itoa(flags.WindowHeight))
	}

	if flags.WindowX > 0 {
		args = append(args, "-windowx", strconv.Itoa(flags.WindowX))
	}
	if flags.WindowY > 0 {
		args = append(args, "-windowy", strconv.Itoa(flags.WindowY))
	}

	cmd := exec.Command(paths.Proton, args...)

	cmd.Dir = path.Join(paths.Sc2, "Support64")

	go func() {
		if err := cmd.Run(); err != nil {
			log.Fatalf("failed to run StarCraft II: %v", err)
		}
	}()

	return nil
}

func protonConfig(bot *api.PlayerSetup, cpu *api.PlayerSetup) *client.GameConfig {
	mapPath := random1v1Map()

	log.Printf("Using map: %q\n", mapPath)

	client.SetMap(mapPath)

	config := client.NewGameConfig(bot, cpu)
	config.Connect(8168)
	config.StartGame(mapPath)

	return config
}
