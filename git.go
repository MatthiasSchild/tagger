package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var (
	versionPattern = regexp.MustCompile(`^v(\d+)\.(\d+)\.(\d+)(-(.+))?$`)
)

func getCurrentVersion() Version {
	var versions []Version

	output, err := execute("git", "tag")
	if err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			if exit.ExitCode() == 128 {
				fmt.Println("git exited with status 128. Are you in a git repository?")
				os.Exit(1)
			}
		}

		fmt.Println(err)
		os.Exit(1)
	}

	for _, tagString := range strings.Split(output, "\n") {
		hit := versionPattern.FindStringSubmatch(tagString)
		if len(hit) > 0 {
			major, _ := strconv.Atoi(hit[1])
			minor, _ := strconv.Atoi(hit[2])
			patch, _ := strconv.Atoi(hit[3])
			note := hit[5]

			versions = append(versions, Version{major, minor, patch, note})
		}
	}

	return getMaxVersion(versions)
}
