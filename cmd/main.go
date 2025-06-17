package main

import (
	configs "go-backend/internal/configs"
)

func main() {
	server := configs.NewServer()
	server.Start()
}
