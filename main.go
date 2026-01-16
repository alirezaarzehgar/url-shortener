package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/alirezaarzehgar/url-shortener/cmd"
	"github.com/joho/godotenv"
)

func rootHelp() {
	fmt.Printf("Usage: %s <command>\n", os.Args[0])
	fmt.Println("Commands: ")
	fmt.Println("	server Run shortener service")
	fmt.Println("	migrate Migrate databases")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		rootHelp()
	}

	if err := godotenv.Load(); err != nil {
		slog.Error("failed to load .env", "error", err)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "server":
		cmd.Server(os.Args[1:])
	case "migrate":
		cmd.Migratoin(os.Args[1:])
	default:
		rootHelp()
	}
}
