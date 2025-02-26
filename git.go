package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"slices"
	"strings"
)

func getAllGitTags() ([]Tag, error) {
	result := make([]Tag, 0)

	cmd := exec.Command("git", "tag", "-l")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	rawTags := strings.Split(string(out), "\n")
	cleanTags := make([]string, 0)

	for _, tag := range rawTags {
		pattern := "^v([0-9]+)\\.([0-9]+)\\.([0-9]+)([+-][a-zA-Z0-9]+)?$"
		matched, err := regexp.MatchString(pattern, tag)
		if err != nil {
			continue
		}
		if matched {
			cleanTag := strings.SplitN(strings.SplitN(tag, "+", 2)[0], "-", 2)[0]
			if !slices.Contains(cleanTags, cleanTag) {
				cleanTags = append(cleanTags, cleanTag)
			}
		}
	}

	for _, tag := range cleanTags {
		var major, minor, patch int

		_, err = fmt.Sscanf(tag, "v%d.%d.%d", &major, &minor, &patch)
		if err != nil {
			continue
		}

		result = append(result, Tag{major, minor, patch, "", ""})
	}

	return result, nil
}

func getCurrentGitHash() (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.Trim(string(out), "\r\n\t "), nil
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

func hasUncommittedChanges() (bool, error) {
	cmd := exec.Command("git", "diff", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return false, err
	}

	trimmedOutput := strings.Trim(string(out), "\r\n\t ")
	return len(trimmedOutput) > 0, nil
}

func commitAll(message string) error {
	cmd1 := exec.Command("git", "add", "--all")
	err := cmd1.Run()
	if err != nil {
		return err
	}
	cmd2 := exec.Command("git", "commit", "-m", message)
	err = cmd2.Run()
	if err != nil {
		return err
	}
	return nil
}
