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
	fmt.Println("\t-g, --gdrive")
	fmt.Println("\t\tsearch in google drive")
	os.Exit(exitValue)
}

func contains(s []string, searchterm string) bool {
	for _, str := range s {
		if str == searchterm {
			return true
		}
	}
	return false
}

func main() {
	if len(os.Args) < 2 {
		usage(exitFailure)
	}
	if os.Args[1] == "--help" || os.Args[1] == "-h" {
		usage(exitSuccess)
	}
	//the program turn off verbose mode by default
	isVerbose, isGdrive := false, false
	if contains(os.Args, "-v") || contains(os.Args, "--verbose") {
		fmt.Println("Verbose Enabled")
		isVerbose = true
	}
	if contains(os.Args, "-g") || contains(os.Args, "--gdrive") {
		fmt.Println("Gdrive Enabled")
		isGdrive = true
	}
	folderPath := os.Args[1]
	hunter := hunter.NewHunter(folderPath, isVerbose, isGdrive, []string{"pass", "mot de passe", "password", "@extia."})
	err := hunter.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(exitFailure)
	}
}
