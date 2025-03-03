# Sc2Bot

🚧 Starting a StarCraft II bot

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

### Maps

Go to [Map Packs](https://github.com/Blizzard/s2client-proto?tab=readme-ov-file#map-packs), download all the maps and extract them in `$HOME$/.steam/debian-installation/steamapps/compatdata/3430940832/pfx/drive_c/Program Files (x86)/StarCraft II/Maps` using the password `iagreetotheeula`.
