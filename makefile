sc2path:
	open "${SC2PATH}"

slow:
	clear
	go run ./... -- -realtime true -windowwidth 1920 -windowheight 1080

fast:
	clear
	go run ./... -- -realtime false -windowwidth 1280 -windowheight 720

clean:
	rm -f __debug_bin* BlackCompany BlackCompany-*.zip BlackCompany.exe BlackCompany.zip vendor

zip:
	make clean
	go build
	zip -9 "BlackCompany.zip" BlackCompany LICENSE.md README.md
