package main

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Flags struct {
	// Blizzard flags

	// DisplayMode is the display mode to use.
	//
	//  * 0: windowed (default)
	//  * 1: fullscreen
	DisplayMode int

	// WindowWidth is the width of the window in pixels.
	WindowWidth int
	// WindowHeight is the height of the window in pixels
	WindowHeight int

	// WindowX is the horizontal placement of the window in pixels.
	WindowX int
	// WindowY is the vertical placement of the window in pixels.
	WindowY int

	// Listen is the IP address to listen on.
	//
	// Default: 127.0.0.1
	Listen string

	// Port is the port to listen on.
	//
	// Default: 8168
	Port int

	// Realtime is whether to run StarCraft II in real-time or not.
	//
	//  - 0: Do not run in real-time (default)
	//  - 1: Run in real-time
	Realtime bool

	// Timeout for how long the library will block for a response.
	Timeout time.Duration
}

func loadFlags() Flags {
	parsed := ParseFlags()

	flags := Flags{}
	flags.Realtime = flagBool(parsed, "realtime", false)
	flags.DisplayMode = flagInt(parsed, "displaymode", 0)
	flags.Port = flagInt(parsed, "port", 8168)
	flags.WindowHeight = flagInt(parsed, "windowheight", 0)
	flags.WindowWidth = flagInt(parsed, "windowwidth", 0)
	flags.WindowX = flagInt(parsed, "windowx", 0)
	flags.WindowY = flagInt(parsed, "windowy", 0)
	flags.Listen = flagString(parsed, "listen", "127.0.0.1")

	twoMinutes, _ := time.ParseDuration("2m")
	flags.Timeout = flagDuration(parsed, "timeout", twoMinutes)

	return flags
}

// flagBool gets a boolean flag value from the provided flag map
func flagBool(flags map[string]string, name string, fallback bool) bool {
	if val, ok := flags[name]; ok {
		switch strings.ToLower(val) {
		case "true", "1", "t", "yes", "y":
			return true
		case "false", "0", "f", "no", "n":
			return false
		}
	}
	return fallback
}

// flagInt gets an integer flag value from the provided flag map
func flagInt(flags map[string]string, name string, fallback int) int {
	if value, ok := flags[name]; ok {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}

// flagString gets a string flag value from the provided flag map
func flagString(flags map[string]string, name string, fallback string) string {
	if value, ok := flags[name]; ok {
		return value
	}
	return fallback
}

// flagDuration gets a duration flag value from the provided flag map
func flagDuration(flags map[string]string, name string, fallback time.Duration) time.Duration {
	if value, ok := flags[name]; ok {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return fallback
}

// ParseFlags parses command-line arguments directly and bypasses the flag
// package because some library decided it was a good idea to declare flags in
// `init` and keep them private.
func ParseFlags() map[string]string {
	flags := make(map[string]string)

	// Skip the program name
	args := os.Args[1:]

	for i := 0; i < len(args); i++ {
		advance, name, value := ParseArg(args, i)

		if name != "" {
			flags[name] = value
		}

		// Skip any used arguments
		i += advance
	}

	return flags
}

// ParseArg processes a single argument and returns:
//   - how many additional arguments were consumed other than the current one
//   - the flag name (if it's a flag)
//   - the flag value
//
// Returns empty strings for name and value if the argument is not a flag
func ParseArg(args []string, index int) (advance int, name string, value string) {
	if index >= len(args) {
		return 0, "", ""
	}

	arg := args[index]

	// Check if it's a flag (starts with -)
	if len(arg) <= 1 || arg[0] != '-' {
		// Not a flag
		return 0, "", ""
	}

	// Handle --flags
	if arg[1] == '-' && len(arg) > 2 {
		// Remove --
		arg = arg[2:]
	} else {
		// Remove -
		arg = arg[1:]
	}

	// Handle name=value
	if idx := strings.Index(arg, "="); idx != -1 {
		name = arg[:idx]
		value = arg[idx+1:]
		return 0, name, value
	}

	// Handle name value
	name = arg
	if index+1 < len(args) && !strings.HasPrefix(args[index+1], "-") {
		value = args[index+1]

		// Skip the next arg since we used it
		return 1, name, value
	}

	// Flag without value (true)
	return 0, name, "true"
}
