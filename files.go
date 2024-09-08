package main

import (
	"encoding/json"
	"fmt"
	"os"
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
