package main

import (
	"fmt"
	"log"
	"math/rand"
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
		err := cmd.Run()
		if err != nil {
			log.Fatalf("failed to run StarCraft II: %v", err)
		}
	}()

	return nil
}

// Sc2Paths contains the paths of StarCraft II and Proton.
type Sc2Paths struct {
	// Exe is the path of the game's executable
	Exe string
	// Maps is the paths of the game's maps folder
	Maps string
	// Prefix is the path of the game's wine prefix
	Prefix string
	// Proton is the path of Proton's executable
	Proton string
	// Sc2 is the path of the game's main installation folder
	Sc2 string
}

func sc2Paths(env *Env) (*Sc2Paths, error) {
	if env.PROTON_PATH == "" {
		return nil, fmt.Errorf("PROTON_PATH is not set")
	}
	if env.STEAM_COMPAT_DATA_PATH == "" {
		return nil, fmt.Errorf("STEAM_COMPAT_DATA_PATH is not set")
	}

	pfx := path.Join(env.STEAM_COMPAT_DATA_PATH, "pfx")
	err := os.Setenv("WINEPREFIX", pfx)
	if err != nil {
		return nil, fmt.Errorf("failed to set WINEPREFIX: %w", err)
	}

	err = os.Setenv("WINE_SIMULATE_WRITECOPY", "1")
	if err != nil {
		return nil, fmt.Errorf("failed to set WINE_SIMULATE_WRITECOPY: %w", err)
	}

	sc2 := path.Join(pfx, "drive_c", "Program Files (x86)", "StarCraft II")
	exe := path.Join(sc2, "Versions", "Base93333", "SC2_x64.exe")

	return &Sc2Paths{
		Exe:    exe,
		Maps:   path.Join(sc2, "Maps"),
		Prefix: pfx,
		Proton: env.PROTON_PATH,
		Sc2:    sc2,
	}, nil
}

var Maps2024Season4 = []string{
	"AbyssalReefAIE",
	"AcropolisAIE",
	"AutomatonAIE",
	"EphemeronAIE",
	"InterloperAIE",
	"ThunderbirdAIE",
}

// Random1v1Map returns a random map name from the current 1v1 ladder map pool.
func Random1v1Map() string {
	currentMaps := Maps2024Season4
	return currentMaps[rand.Intn(len(currentMaps))] + ".SC2Map"
}

func protonConfig(bot *api.PlayerSetup, cpu *api.PlayerSetup) *client.GameConfig {
	mapPath := Random1v1Map()

	log.Printf("Using map: %q\n", mapPath)

	client.SetMap(mapPath)

	config := client.NewGameConfig(bot, cpu)
	config.Connect(8168)
	config.StartGame(mapPath)

	return config
}
