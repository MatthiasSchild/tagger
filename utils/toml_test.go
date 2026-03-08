package utils_test

import (
	"strings"
	"testing"

	"github.com/MatthiasSchild/tagger/utils"
)

const tomlInput = `
# Comment A
version = "0.0.1"

[package]
version = "0.0.2" # Comment B

[[package]]
version = "0.0.3"

[other]
version = "0.0.4"

[[other]]
version = "0.0.5"
# Comment C
`

const tomlExpected = `
# Comment A
version = "0.0.1"

[package]
version = "0.1.0" # Comment B

[[package]]
version = "0.0.3"

[other]
version = "0.0.4"

[[other]]
version = "0.0.5"
# Comment C
`

func TestUpdateVersionInToml(t *testing.T) {
	result := utils.UpdateVersionInToml(tomlInput, "0.1.0")

	resultLines := strings.Split(result, "\n")
	expectLines := strings.Split(tomlExpected, "\n")
	if len(resultLines) != len(expectLines) {
		t.Errorf(
			"line numbers of the result mismatch, result=%d, expect=%d",
			len(resultLines),
			len(expectLines),
		)
		return
	}

	for index := range resultLines {
		if resultLines[index] != expectLines[index] {
			t.Errorf(
				"result mismatches in line %d:\nline: %s\nexpect:%s",
				index+1,
				resultLines[index],
				expectLines[index],
			)
			return
		}
	}
}
