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
	optionNote      string
	optionWriteFile bool
	strategy        Strategy
)

var helpText = `Strategies
Following strategies are available: patch, minor, major, datetime (default is patch).

patch, minor and major increases the position by 1.
It also resets the later positions.
E.g. when using "minor" it changes the version v1.2.3 to v1.3.0;
the minor is increased from 2 to 3, and the patch position is reset to 0.

The strategy datetime is more special. It stores the unix timestamp into the version.
The major part will be kept, the minor part will contain the date information and the patch part the time information.
E.g. when you tag v1.0.0 on the 01 Jan 2020 on 00:00:00,
you get the unix timestamp of 1577833212.
The date will be 1577833212 % (60 * 60 * 24)
  = 82812
and the time will be 1577833212 / (60 * 60 * 24)
  = 18261
so the version will result in v1.82812.18261`

func init() {
	flag.Usage = func() {
		flag.PrintDefaults()

		fmt.Println("\n" + helpText)
	}

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
