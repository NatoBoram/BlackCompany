sc2path:
	open "${SC2PATH}"

slow:
	clear
	go run ./... -- -realtime true

fast:
	clear
	go run ./... -- -realtime false
