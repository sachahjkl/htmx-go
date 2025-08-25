server:
	go run cmd/main.go

build:
	go build -o ./tmp/main.exe cmd/main.go

# this requires "air" to be installed
# you can run `go install github.com/cosmtrek/air@latest` to get it
live:
	air 

install-tailwind:
	bun install tailwindcss @tailwindcss/cli

css: install-tailwind
	bun x @tailwindcss/cli -i ./style/style.css -o ./assets/style.css