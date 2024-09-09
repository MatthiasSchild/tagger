package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"gopkg.in/yaml.v3"
)

func readVersionFromPackageJson() (Tag, error) {
	content, err := os.ReadFile("package.json")
	if err != nil {
		return Tag{}, err
	}

	packageData := &struct {
		Version string `json:"version"`
	}{}

	err = json.Unmarshal(content, packageData)
	if err != nil {
		return Tag{}, err
	}

	groups := versionRegex.FindStringSubmatch(packageData.Version)
	if groups == nil {
		return Tag{}, fmt.Errorf("version in package.json must have format '1.2.3'")
	}

	major, _ := strconv.Atoi(groups[1])
	minor, _ := strconv.Atoi(groups[2])
	patch, _ := strconv.Atoi(groups[3])

	return Tag{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}

func writeVersionToPackageJson(tag Tag) error {
	content, err := os.ReadFile("package.json")
	if err != nil {
		return err
	}

	// Replace the version
	re := regexp.MustCompile(`"version":\s*"[^"]*"`)
	newVersion := fmt.Sprintf("%d.%d.%d", tag.Major, tag.Minor, tag.Patch)
	updatedContent := re.ReplaceAllString(string(content), fmt.Sprintf(`"version": "%s"`, newVersion))

	err = os.WriteFile("package.json", []byte(updatedContent), 0644)
	if err != nil {
		return err
	}
	return nil
}

func readVersionFromPubspecYaml() (Tag, error) {
	content, err := os.ReadFile("pubspec.yaml")
	if err != nil {
		return Tag{}, err
	}

	packageData := &struct {
		Version string `json:"version"`
	}{}

	err = yaml.Unmarshal(content, packageData)
	if err != nil {
		return Tag{}, err
	}

	groups := versionRegex.FindStringSubmatch(packageData.Version)
	if groups == nil {
		return Tag{}, fmt.Errorf("version in pubspec.yaml must have format '1.2.3+4'")
	}

	major, _ := strconv.Atoi(groups[1])
	minor, _ := strconv.Atoi(groups[2])
	patch, _ := strconv.Atoi(groups[3])

	return Tag{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}

func writeVersionToPubspecYaml(tag Tag, increment bool) error {
	content, err := os.ReadFile("pubspec.yaml")
	if err != nil {
		return err
	}

	// Replace the version
	re := regexp.MustCompile(`version:\s*(\d+\.\d+\.\d+)\+(\d+)`)
	newVersion := fmt.Sprintf("%d.%d.%d", tag.Major, tag.Minor, tag.Patch)
	updatedContent := re.ReplaceAllStringFunc(string(content), func(match string) string {
		matches := re.FindStringSubmatch(match)
		if len(matches) > 2 {
			buildNumber, _ := strconv.Atoi(matches[2])
			if increment {
				buildNumber++
			}
			return fmt.Sprintf("version: %s+%d", newVersion, buildNumber)
		}
		return match
	})

	err = os.WriteFile("pubspec.yaml", []byte(updatedContent), 0644)
	if err != nil {
		return err
	}
	return nil
}
