.PHONY: test run install-cli gui gui-wails build build-windows build-portable

test:
	go test ./...

run:
	go run ./cmd/scenaria --help

install-cli:
	go install ./cmd/scenaria

gui: gui-wails

gui-wails:
	wails dev

build:
	go build -o bin/scenaria ./cmd/scenaria

build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/scenaria.exe ./cmd/scenaria
	cd frontend && npm run build && cd ..
	wails build -platform windows/amd64

build-portable:
	powershell -ExecutionPolicy Bypass -File scripts/build-portable.ps1
