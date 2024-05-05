bootstrap_path=./static/bootstrap/dist

build:
	@go build -o bin/game *.go
	@GOOS=js GOARCH=wasm go build -o static/main.wasm front/*
	@cp /usr/local/go/misc/wasm/wasm_exec.js static/js

download:
	@go get 
	@mkdir -p $(bootstrap_path)/css $(bootstrap_path)/js \
		&& curl -Lo $(bootstrap_path)/css/bootstrap.min.css \
			https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/css/bootstrap.min.css \
		&& curl -Lo $(bootstrap_path)/js/bootstrap.min.js \
			https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/js/bootstrap.bundle.min.js

run: build
	@bin/game