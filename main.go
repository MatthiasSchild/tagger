package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	optionDry       bool
	optionHash      bool
	optionStrategy  string
	optionNote string
	optionWriteFile bool
	strategy        Strategy
)

func init() {
	flag.BoolVar(&optionDry, "dry", false, "Dry-run, dont tag.")
	flag.BoolVar(&optionHash, "hash", false, "Add the hash as note.")
	flag.StringVar(&optionStrategy, "strategy", "patch", "Version increase strategy.")
	flag.StringVar(&optionNote, "note", "", "Add a note to the tag.")
	flag.BoolVar(&optionWriteFile, "write", false, "Write config to file.")

	flag.Parse()

	strategy = stringToStrategy(optionStrategy)
}

func main() {
	currentVersion := getCurrentVersion()
	increaseVersion(strategy, &currentVersion)

	if optionHash {
		output, err := execute("git", "log", "-n1", "--format=format:%H")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		currentVersion.Note = output[:7]
	}

	if optionNote != "" {
		currentVersion.Note = optionNote
	}

	fmt.Println(currentVersion)
	if !optionDry {
		_, err := execute("git", "tag", currentVersion.String())
		if err != nil {
			fmt.Println(err)
		}
	}
}
