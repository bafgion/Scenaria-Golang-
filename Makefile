.PHONY: test run install-cli gui gui-wails gui-fyne build build-windows

test:
	go test ./...

run:
	go run ./cmd/scenaria --help

install-cli:
	go install ./cmd/scenaria

gui-wails:
	wails dev

gui-fyne:
	go run -tags desktop ./cmd/scenaria-gui

gui: gui-fyne

install-gui:
	go install -tags desktop ./cmd/scenaria-gui

build:
	go build -o bin/scenaria ./cmd/scenaria

build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/scenaria.exe ./cmd/scenaria
	GOOS=windows GOARCH=amd64 go build -tags desktop -o bin/scenaria-gui.exe ./cmd/scenaria-gui

build-portable:
	powershell -ExecutionPolicy Bypass -File scripts/build-portable.ps1
