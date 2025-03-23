package main

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"testing"

	"github.com/joho/godotenv"
)

// Env contains all the environment variables that are supported by this
// program.
type Env struct {
	// PROTON_PATH is the path to the Proton installation. When it's set, the game
	// will be launched from Steam instead of using the native client.
	//
	// Example: $HOME/.steam/root/steamapps/common/Proton - Experimental/proton
	PROTON_PATH string

	// SC2PATH is the path to the StarCraft II installation.
	//
	// Example: $HOME/.steam/debian-installation/steamapps/compatdata/3430940832/pfx/drive_c/Program Files (x86)/StarCraft II
	SC2PATH string

	// STEAM_COMPAT_CLIENT_INSTALL_PATH is the path to the installation of Steam
	// itself.
	//
	// Example: $HOME/.steam/debian-installation
	STEAM_COMPAT_CLIENT_INSTALL_PATH string

	// STEAM_COMPAT_DATA_PATH is the path to the non-Steam game's "compatdata"
	// directory inside of "steamapps".
	//
	// Example: $HOME/.steam/debian-installation/steamapps/compatdata/3430940832
	STEAM_COMPAT_DATA_PATH string
}

// Environment represents the current environment (development, test,
// production).
type Environment string

// Environment represents the current environment (development, test,
// production).
const (
	Development Environment = "development"
	Test        Environment = "test"
	Production  Environment = "production"
)

// toEnvironment converts a string to an Environment.
func toEnvironment(s string) Environment {
	switch s {
	case "development":
		return Development
	case "test":
		return Test
	case "production":
		return Production
	default:
		return "development"
	}
}

// getEnvironment returns the current environment.
func getEnvironment() Environment {
	environment := os.Getenv("GO_ENV")
	if environment != "" {
		return toEnvironment(environment)
	}

	if testing.Testing() {
		os.Setenv("GO_ENV", "test")
		return "test"
	}

	os.Setenv("GO_ENV", "development")
	return "development"
}

// loadEnv loads the environment variables from the .env files.
func loadEnv() (*Env, error) {
	environment := getEnvironment()

	files := []string{
		".env." + string(environment) + ".local",
		".env." + string(environment),
		".env.local",
		".env",
	}

	for _, file := range files {
		if err := godotenv.Load(file); err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load environment variables from %q: %w", file, err)
		}
	}

	return &Env{
		PROTON_PATH:                      os.Getenv("PROTON_PATH"),
		STEAM_COMPAT_CLIENT_INSTALL_PATH: os.Getenv("STEAM_COMPAT_CLIENT_INSTALL_PATH"),
		STEAM_COMPAT_DATA_PATH:           os.Getenv("STEAM_COMPAT_DATA_PATH"),
		SC2PATH:                          os.Getenv("SC2PATH"),
	}, nil
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

// sc2Paths returns the paths of StarCraft II and Proton as found in the
// environment variables.
func sc2Paths(env *Env) (*Sc2Paths, error) {
	if env.PROTON_PATH == "" {
		return nil, fmt.Errorf("PROTON_PATH is not set")
	}

	if env.STEAM_COMPAT_DATA_PATH == "" {
		return nil, fmt.Errorf("STEAM_COMPAT_DATA_PATH is not set")
	}

	pfx := path.Join(env.STEAM_COMPAT_DATA_PATH, "pfx")

	if err := os.Setenv("WINEPREFIX", pfx); err != nil {
		return nil, fmt.Errorf("failed to set WINEPREFIX: %w", err)
	}

	if err := os.Setenv("WINE_SIMULATE_WRITECOPY", "1"); err != nil {
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

// maps2024Season4 is the map pool for the 2024 Season 4 ladder.
var maps2024Season4 = []string{
	"AbyssalReefAIE",
	"AcropolisAIE",
	"AutomatonAIE",
	"EphemeronAIE",
	"InterloperAIE",
	"ThunderbirdAIE",
}

// random1v1Map returns a random map name from the current 1v1 ladder map pool.
func random1v1Map() string {
	currentMaps := maps2024Season4
	return currentMaps[rand.Intn(len(currentMaps))] + ".SC2Map"
}
