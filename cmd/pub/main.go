package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/cheebz/go-pub/cmd/pub/commands"
)

func main() {
	rootCmd := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	var logLevel string
	rootCmd.StringVar(&logLevel, "log-level", "info", "Set the logging level (\"debug\"|\"info\"|\"warn\"|\"error\"|\"fatal\") (default \"info\")")

	rootCmd.Usage = func() {
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintf(os.Stderr, "Usage: %s <command>\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "  serve    Start the server\n")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		rootCmd.PrintDefaults()
		fmt.Fprintln(os.Stderr, "")
	}

	rootCmd.Parse(os.Args[1:])

	command := rootCmd.Arg(0)
	if command == "" {
		rootCmd.Usage()
		os.Exit(1)
	}

	switch command {
	case "serve":
		log.Fatal(commands.Serve())
	default:
		log.Fatalf("Unknown command: %s\n", command)
	}
}
