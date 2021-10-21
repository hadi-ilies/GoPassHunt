package main

import (
	"fmt"
	"os"

	"GoPassHunt/hunter"
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
	if len(os.Args) < 2 {
		usage(exitFailure)
	}
	if os.Args[1] == "--help" || os.Args[1] == "-h" {
		usage(exitSuccess)
	}
	//the program turn off verbose mode by default
	isVerbose := false
	folderPath := os.Args[1]
	if len(os.Args) > 2 {
		if os.Args[2] == "--verbose" || os.Args[2] == "-v" {
			fmt.Println("Verbose option is used")
			isVerbose = true
		}
	}
	hunter := hunter.NewHunter(folderPath, isVerbose, []string{"pass", "mot de passe", "password", "@extia."})
	err := hunter.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(exitFailure)
	}
}
