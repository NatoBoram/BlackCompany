package main

import (
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/aiseeq/s2l/protocol/api"
	"github.com/aiseeq/s2l/protocol/client"
)

// launchProton launches StarCraft II using Proton.
//
//   - https://levelup.gitconnected.com/guide-to-starcraft-ii-proto-api-264811da8a50
//   - https://gist.github.com/michaelbutler/f364276f4030c5f449252f2c4d960bd2
func launchProton(paths *Sc2Paths) error {
	cmd := exec.Command(paths.Proton, "run", paths.Exe,
		"-displayMode", "0",
		"-listen", "127.0.0.1",
		"-port", "8168",
	)

	client.SetExecutable(paths.Exe)

	cmd.Dir = path.Join(paths.Sc2, "Support64")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

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
