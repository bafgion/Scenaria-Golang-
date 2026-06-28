//go:build !desktop

package main

import "fmt"

func main() {
	fmt.Println("Scenaria GUI requires build tag desktop. Run: go run -tags desktop ./cmd/scenaria-gui")
}
