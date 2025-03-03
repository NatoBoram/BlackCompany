package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

type Env struct {
	// PROTON_PATH is the path to the Proton installation. When it's set, the game
	// will be launched from Steam instead of using the native client.
	//
	// Example: $HOME/.steam/root/steamapps/common/Proton - Experimental/proton
	PROTON_PATH string

	// STEAM_COMPAT_CLIENT_INSTALL_PATH is the path to the installation of Steam
	// itself.
	//
	// Exemple: $HOME/.steam/debian-installation
	STEAM_COMPAT_CLIENT_INSTALL_PATH string

	// STEAM_COMPAT_DATA_PATH is the path to the non-Steam game's "compatdata"
	// directory inside of "steamapps".
	//
	// Exemple: $HOME/.steam/debian-installation/steamapps/compatdata/3430940832
	STEAM_COMPAT_DATA_PATH string
}

func getEnvironment() string {
	environment := os.Getenv("GO_ENV")
	if environment != "" {
		return environment
	}

	if testing.Testing() {
		os.Setenv("GO_ENV", "test")
		return "test"
	}

	os.Setenv("GO_ENV", "development")
	return "development"
}

func loadEnv() (*Env, error) {
	environment := getEnvironment()

	files := []string{
		".env." + environment + ".local",
		".env." + environment,
		".env.local",
		".env",
	}

	for _, file := range files {
		err := godotenv.Load(file)
		if err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load environment variables from %q: %w", file, err)
		}
	}

	return &Env{
		PROTON_PATH:                      os.Getenv("PROTON_PATH"),
		STEAM_COMPAT_CLIENT_INSTALL_PATH: os.Getenv("STEAM_COMPAT_CLIENT_INSTALL_PATH"),
		STEAM_COMPAT_DATA_PATH:           os.Getenv("STEAM_COMPAT_DATA_PATH"),
	}, nil
}
