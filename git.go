package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

func getAllGitTags() ([]Tag, error) {
	result := make([]Tag, 0)

	cmd := exec.Command("git", "tag", "-l")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	tags := strings.Split(string(out), "\n")
	for _, tag := range tags {
		pattern := "^v([0-9]+)\\.([0-9]+)\\.([0-9]+)$"
		matched, err := regexp.MatchString(pattern, tag)
		if err != nil {
			continue
		}
		if matched {
			var major, minor, patch int
			_, err = fmt.Sscanf(tag, "v%d.%d.%d", &major, &minor, &patch)
			if err != nil {
				continue
			}

			result = append(result, Tag{major, minor, patch})
		}
	}

	return result, nil
}

func getLatestTag(tags []Tag) Tag {
	var latest Tag
	for _, tag := range tags {
		if tag.Major > latest.Major {
			latest = tag
		} else if tag.Major == latest.Major {
			if tag.Minor > latest.Minor {
				latest = tag
			} else if tag.Minor == latest.Minor {
				if tag.Patch > latest.Patch {
					latest = tag
				}
			}
		}
	}
	return latest
}

func createTag(tag Tag) error {
	cmd := exec.Command("git", "tag", "-a", tag.String(), "-m", tag.String())
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}