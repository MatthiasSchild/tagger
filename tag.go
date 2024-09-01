package main

import (
	"fmt"
)

type Tag struct {
	Major int
	Minor int
	Patch int
}

func (t Tag) String() string {
	return fmt.Sprintf("v%d.%d.%d", t.Major, t.Minor, t.Patch)
}

func (t Tag) Clone() Tag {
	return Tag{t.Major, t.Minor, t.Patch}
}
