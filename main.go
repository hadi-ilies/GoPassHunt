package main

import (
	"fmt"
	"os"
)

const (
	exitFailure = iota
	exitSuccess
)

func usage(exitValue int) {
	fmt.Println("USAGE")
	fmt.Println("\tgo run main.go <folderPath> [options]")
	fmt.Println("DESCRIPTION")
	fmt.Println("\tSearch drives for documents containing passwords")
	fmt.Println("OPTIONS")
	fmt.Println("\t-h, --help")
	fmt.Println("\t\tDisplay the program usage")
	fmt.Println("\t-v, --verbose")
	fmt.Println("\t\tDisplay additonal logs")
	os.Exit(exitValue)
}

func main() {
	fmt.Println("Hello world")
	if os.Args[1] == "--help" || os.Args[1] == "-h" {
		usage(exitSuccess)
	}
}
