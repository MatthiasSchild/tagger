package main

import (
	"log"
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
		log.Fatal(err)
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
