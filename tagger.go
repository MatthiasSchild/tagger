package main

import (
	"os"
)

func main() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
