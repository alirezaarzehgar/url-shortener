package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/alirezaarzehgar/url-shortener/cmd"
	"github.com/joho/godotenv"
)

func rootCmdHelp(args []string) {
	fmt.Printf("%s <server>\n", args[0])
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <command>\n", os.Args[0])
		fmt.Println("Commands: ")
		fmt.Println("	server Run shortener service")
		os.Exit(1)
	}

	if err := godotenv.Load(); err != nil {
		slog.Error("failed to load .env", "error", err)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "server":
		cmd.Server(os.Args[2:])
	default:
		rootCmdHelp(os.Args)
		os.Exit(1)
	}
}
