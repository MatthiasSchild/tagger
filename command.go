package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const rootCmdDescription = `---------------------------------
Tagger for creating a new git tag
---------------------------------

With this command you can automatically your current commit.
Tagger will read all previous tags and will depending on your flags increase the version and tag your commit.
You have multiple tagging strategies:

--major: increase the major version and set minor and patch to 0
--minor: increase the minor version and set patch to 0
--patch: (default): increase the patch version
--datetime: apply the date time strategy (read below)
--hash: add the character commit hash to your version, enter the number of chars (e.g. "8" for v0.1.2-3456abcd)

Date time strategy:
The strategy datetime is more special. It stores the unix timestamp into the version.
The major part will be kept, the minor part will contain the date information and the patch part the time information.
E.g. when you tag v1.0.0 on the 01 Jan 2020 on 09:30:00,
you get the unix timestamp of 1577867400.
and the date will be 1577867400 / (60 * 60 * 24)
  = 18262
The time will be 1577867400 % (60 * 60 * 24)
  = 30600
so the version will result in v1.18262.30600
`

var RootCmd = &cobra.Command{
	Use:   "tagger",
	Short: "Create a new git tag",
	Long:  rootCmdDescription,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := validateFlags()
		if err != nil {
			return err
		}

		tags, err := getAllGitTags()
		if err != nil {
			return fmt.Errorf("failed to fetch git tags: %s", err.Error())
		}

		if len(tags) == 0 {
			return fmt.Errorf("no tags found")
		}

		latestTag := getLatestTag(tags)
		newTag := latestTag.Clone()
		if flagMajor {
			newTag.Major++
			newTag.Minor = 0
			newTag.Patch = 0
		} else if flagMinor {
			newTag.Minor++
			newTag.Patch = 0
		} else {
			newTag.Patch++
		}

		if flagDateTime {
			now := time.Now().Unix()
			newTag.Minor = int(now % (60 * 60 * 24))
			newTag.Patch = int(now / (60 * 60 * 24))
		}

		if flagHash != 0 {
			hash, err := getCurrentGitHash()
			if err != nil {
				return fmt.Errorf("could not get current hash: %s", err.Error())
			}

			newTag.Addition = hash[:flagHash]
		}

		if !flagDry {
			err = createTag(newTag)
			if err != nil {
				fmt.Println("Failed to create tag")
				fmt.Println(err)
				os.Exit(1)
			}
		}

		fmt.Printf("Tagged %s -> %s\n", latestTag, newTag)
		return nil
	},
}

func init() {
	RootCmd.Flags().BoolVar(&flagMajor, "major", false, "Increase major part")
	RootCmd.Flags().BoolVar(&flagMinor, "minor", false, "Increase minor part")
	RootCmd.Flags().BoolVar(&flagPatch, "patch", false, "Increase patch part")
	RootCmd.Flags().BoolVar(&flagDateTime, "datetime", false, "Set minor and patch to date time")
	RootCmd.Flags().IntVar(&flagHash, "hash", 0, "Add commit hash to end")
	RootCmd.Flags().BoolVarP(&flagDry, "dry", "d", false, "Show new tag but don't apply")
}
