sc2path:
	open "${SC2PATH}"

realtime:
	clear
	go run ./... -- -realtime true
