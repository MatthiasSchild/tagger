package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flagMajor := flag.Bool("major", false, "Increase major version")
	flagMinor := flag.Bool("minor", false, "Increase minor version")
	flagDry := flag.Bool("dry", false, "Dry run")
	flag.Parse()

	tags, err := getAllGitTags()
	if err != nil {
		fmt.Println("Failed to fetch git tags")
		fmt.Println(err)
		os.Exit(1)
	}

	if len(tags) == 0 {
		fmt.Println("No tags found")
		os.Exit(1)
	}

	latestTag := getLatestTag(tags)
	newTag := latestTag.Clone()
	if *flagMajor {
		newTag.Major++
		newTag.Minor = 0
		newTag.Patch = 0
	} else if *flagMinor {
		newTag.Minor++
		newTag.Patch = 0
	} else {
		newTag.Patch++
	}
	if !*flagDry {
		err = createTag(newTag)
		if err != nil {
			fmt.Println("Failed to create tag")
			fmt.Println(err)
			os.Exit(1)
		}
	}
	fmt.Printf("Tagged %s -> %s\n", latestTag, newTag)
}
