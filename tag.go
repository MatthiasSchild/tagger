package main

import (
	"fmt"
)

// Tag contains the individual information about a 3-part tag with optional addition
type Tag struct {
	Major    int
	Minor    int
	Patch    int
	Addition string
}

// String builds the string dependent on the values of the structure,
// resulting in a 3-part version (e.g. v1.2.3).
// When addition is set, it will be added with a hyphen (e.g. v1.2.3-1fa342)
func (t Tag) String() string {
	result := fmt.Sprintf("v%d.%d.%d", t.Major, t.Minor, t.Patch)

	if len(t.Addition) > 0 {
		result += "-" + t.Addition
	}

	return result
}

// Equals checks if the major, minor and the patch part of two versions are equal.
// The addition part will be ignored.
func (t Tag) Equals(tag Tag) bool {
	return (t.Major == tag.Major) && (t.Minor == tag.Minor) && (t.Patch == tag.Patch)
}

// Clone clones the Tag structure.
// The addition will be removed from the copy.
func (t Tag) Clone() Tag {
	return Tag{t.Major, t.Minor, t.Patch, ""}
}
