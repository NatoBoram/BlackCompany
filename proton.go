package main

import (
	"fmt"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/NatoBoram/BlackCompany/log"
	"github.com/aiseeq/s2l/protocol/api"
	"github.com/aiseeq/s2l/protocol/client"
	"github.com/shirou/gopsutil/v4/process"
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
		log.Info("Running StarCraft II in real-time")
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

	running, err := runningProton()
	if err != nil {
		return fmt.Errorf("failed to check if Proton is running: %w", err)
	}

	if running != nil {
		log.Info("Proton is already running")
		return nil
	}

	log.Info("Launching StarCraft II using Proton")
	go func() {
		if err := cmd.Run(); err != nil {
			log.Fatal("failed to run StarCraft II: %v", err)
		}
	}()

	return nil
}

func runningProton() (*process.Process, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("failed to get processes: %w", err)
	}

	for _, p := range processes {
		exe, err := p.Exe()
		if err != nil {
			continue
		}

		if strings.HasSuffix(exe, "/files/bin/wineserver") {
			return p, nil
		}
	}

	return nil, nil
}

func protonConfig(bot *api.PlayerSetup, participants ...*api.PlayerSetup) *client.GameConfig {
	mapPath := random1v1Map()

	log.Info("Using map %q\n", mapPath)

	client.SetMap(mapPath)

	participants = append([]*api.PlayerSetup{bot}, participants...)
	config := client.NewGameConfig(participants...)
	config.Connect(8168)
	config.StartGame(mapPath)

	return config
}

func replayConfig(flags Flags) (*client.GameConfig, error) {
	if flags.Replay == "" {
		log.Fatal("no replay file provided")
	}

	if flags.Map == "" {
		log.Fatal("no map provided")
	}

	client.SetMap(flags.Map + ".SC2Map")

	observer := client.NewParticipant(api.Race_NoRace, "Observer")
	config := client.NewGameConfig(observer)
	config.Connect(8168)

	replayPath := path.Join("C:\\Program Files (x86)\\StarCraft II\\Replays", flags.Replay)

	request := &api.Request{Request: &api.Request_StartReplay{
		StartReplay: &api.RequestStartReplay{
			Replay: &api.RequestStartReplay_ReplayPath{
				ReplayPath: replayPath,
			},
			ObservedPlayerId: 1,
			Options: &api.InterfaceOptions{
				Raw:   true,
				Score: true,
			},
			DisableFog: false,
			Realtime:   false,
		},
	}}

	response, err := config.Client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("failed to start replay: %w", err)
	}

	if response.GetStartReplay().GetError() != api.ResponseStartReplay_nil {
		return nil, fmt.Errorf("failed to start replay. Error: %v, Details: %v",
			response.GetStartReplay().GetError(),
			response.GetStartReplay().GetErrorDetails(),
		)
	}

	return config, nil
}
