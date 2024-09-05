package main

import (
	"fmt"
	"os"
)

func main() {
	err := RootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
