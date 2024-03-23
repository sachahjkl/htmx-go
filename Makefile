server:
	go run cmd/main.go

build:
	go build -o ./tmp/main.exe cmd/main.go

# this requires "air" to be installed
# you can run `go install github.com/cosmtrek/air@latest` to get it
live:
	air 