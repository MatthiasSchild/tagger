package utils

import (
	"regexp"
	"strings"
)

var tomlGroupRegex = regexp.MustCompile(`^\s*(\[.*\])\s*$`)
var tomlVersionRegex = regexp.MustCompile(`^(\s*version\s*=\s*")([^"]+)(".*)$`)

func UpdateVersionInToml(content string, version string) string {
	result := make([]string, 0)
	withinPackage := false

	for _, line := range strings.Split(content, "\n") {
		groupMatch := tomlGroupRegex.FindStringSubmatch(line)
		if len(groupMatch) > 0 {
			withinPackage = groupMatch[1] == "[package]"
		} else if withinPackage {
			lineParts := tomlVersionRegex.FindStringSubmatch(line)
			if len(lineParts) != 0 {
				line = lineParts[1] + version + lineParts[3]
			}
		}

		result = append(result, line)
	}

	return strings.Join(result, "\n")
}
