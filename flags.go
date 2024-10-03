package main

import (
	"fmt"
)

var (
	flagMajor    bool
	flagMinor    bool
	flagPatch    bool
	flagDateTime bool
	flagHash     int
	flagDry      bool
	flagWrite    string
	flagBuild    bool
)

func validateFlags() error {
	// When neither major, minor, patch nor datetime is set, set patch
	if !flagMajor && !flagMinor && !flagPatch && !flagDateTime {
		flagPatch = true
	}

	// Only one of those are allowed to be set: major, minor, patch
	if (flagMajor && flagMinor) || (flagMajor && flagPatch) || (flagMinor && flagPatch) {
		return fmt.Errorf("only one of those are allowed: --major, --minor, --patch")
	}

	// When using date time, minor and patch are not allowed
	if flagDateTime && (flagMinor || flagPatch) {
		return fmt.Errorf("when using --datetime, --minor and --patch are not allowed")
	}

	// "hash" should be a number between 1 and (including) 40 (0 equals "disabled")
	if flagHash < 0 || flagHash > 40 {
		return fmt.Errorf("hash must be a number between 1 and including 40")
	}
	if flagHash == 1 {
		fmt.Println("Just one character? This is useless, but here you go...")
	}

	return nil
}
