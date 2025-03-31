# BlackCompany

[![Go](https://github.com/NatoBoram/BlackCompany/actions/workflows/go.yaml/badge.svg)](https://github.com/NatoBoram/BlackCompany/actions/workflows/go.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/NatoBoram/BlackCompany)](https://goreportcard.com/report/github.com/NatoBoram/BlackCompany)
[![Go Reference](https://pkg.go.dev/badge/github.com/NatoBoram/BlackCompany.svg)](https://pkg.go.dev/github.com/NatoBoram/BlackCompany)
[![Wakapi](https://wakapi.dev/api/badge/NatoBoram/interval:any/project:BlackCompany)](https://wakapi.dev/summary?interval=any&project=BlackCompany)

BlackCompany is a StarCraft II bot written in Go.

## Dependencies

This bot is developed on Linux.

### StarCraft II

1. Install [Steam](https://store.steampowered.com)
2. Download [`Battle.net-Setup.exe`](https://download.battle.net)
3. Add `Battle.net-Setup.exe` as a non-Steam game
4. In Steam, right click on `Battle.net-Setup.exe`, _Properties..._ and set _Launch options_ to `WINE_SIMULATE_WRITECOPY=1 %command%`
5. Go to _Compatibility_ and check _Force the use of a specific Steam Play compatibility tool_
6. Launch `Battle.net-Setup.exe` via Steam and install it at the default location
7. Go back to `Battle.net-Setup.exe`'s properties and change the _Target_ to its installed location, like `$HOME/.steam/debian-installation/steamapps/compatdata/3430940832/pfx/drive_c/Program Files (x86)/Battle.net/Battle.net Launcher.exe`

You can now launch Battle.net as a non-Steam game and play StarCraft II.

### Environment variables

When `PROTON_PATH` is set, the bot will use Proton to launch the game. Otherwise, it will fallback to [`sl2`](github.com/aiseeq/s2l)'s defaults.

Fill in the `.env.local` file. Here's an example:

```sh
PROTON_PATH="$HOME/.steam/root/steamapps/common/Proton - Experimental/proton"
SC2PATH="$HOME/.steam/debian-installation/steamapps/compatdata/3430940832/pfx/drive_c/Program Files (x86)/StarCraft II"
STEAM_COMPAT_CLIENT_INSTALL_PATH="$HOME/.steam/debian-installation"
STEAM_COMPAT_DATA_PATH="$HOME/.steam/debian-installation/steamapps/compatdata/3430940832"
```

### Maps

Go to [Map Packs](https://github.com/Blizzard/s2client-proto?tab=readme-ov-file#map-packs), download all the maps and extract them in `$SC2PATH/Maps` using the password `iagreetotheeula`.

## Run

There's a `Makefile` to help you run the bot.

```sh
# Runs the game in real-time
make slow

# Runs the game as fast as possible
make fast
```

## Resources

These resources massively helped me kickstart bot development.

- [Install and play StarCraft II using Steam on Linux](https://www.youtube.com/watch?v=HqOEKSR_Eow)
- [Guide to StarCraft II Proto API](https://levelup.gitconnected.com/guide-to-starcraft-ii-proto-api-264811da8a50)
- [VeTerran](https://bitbucket.org/AiSee/VeTerran)
- [s2l](https://pkg.go.dev/github.com/aiseeq/s2l)
- [liquipedia](https://liquipedia.net/starcraft2)
- [StarCraft Wiki](https://starcraft.fandom.com)
