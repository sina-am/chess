build:
	@go build -o bin/game *.go
run: build
	@bin/game