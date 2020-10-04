package main

import (
	"fmt"
	"time"
)

type Version struct {
	Major int
	Minor int
	Patch int
	Note  string
}

func (v Version) String() string {
	s := fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.Note != "" {
		s = fmt.Sprintf("%s-%s", s, v.Note)
	}
	return s
}

func getMaxVersion(versions []Version) Version {
	v := Version{}

	for _, v2 := range versions {
		a := v2.Major > v.Major
		b := v2.Major == v.Major && v2.Minor > v.Minor
		c := v2.Major == v.Major && v2.Minor == v.Minor && v2.Patch > v.Patch

		if a || b || c {
			v = v2
		}
	}

	return v
}

type Strategy int

const (
	StrategyIncreasePatch Strategy = iota
	StrategyIncreaseMinor
	StrategyIncreaseMajor
	StrategyDateTime
)

func stringToStrategy(s string) Strategy {
	strategies := map[string]Strategy{
		"patch":    StrategyIncreasePatch,
		"minor":    StrategyIncreaseMinor,
		"major":    StrategyIncreaseMajor,
		"datetime": StrategyDateTime,
	}

	return strategies[s]
}

func increaseVersion(strategy Strategy, version *Version) {
	switch strategy {
	case StrategyIncreasePatch:
		version.Patch += 1
	case StrategyIncreaseMinor:
		version.Minor += 1
		version.Patch = 0
	case StrategyIncreaseMajor:
		version.Major += 1
		version.Minor = 0
		version.Patch = 0
	case StrategyDateTime:
		unix := time.Now().Unix()
		version.Minor = int(unix % (60 * 60 * 24)) // date
		version.Patch = int(unix / (60 * 60 * 24)) // time
	}
}
