package main

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	versionRegex = regexp.MustCompile("^v?(\\d+)\\.(\\d+)\\.(\\d+)([+-]([a-zA-Z0-9]+))?$")
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

Writing the new version into file:
With the --write flag, you can tell tagger to write the new version into a file.
Tagger will check, if any uncommitted changes are open.
The process will abort, if any uncommitted changes exist.
Then it will calculate the new tag and write it into file (depending on the target).
Afterwards, tagger will create a new commit with the tag as message and will tag this commit.
You have the following targets as options:
--write=npm for writing the version into the package.json
	For this option, the version will have the format "major.minor.patch"
--write=flutter for writing the version into the pubspec.yaml
	For this option, the version will have the format "major.minor.patch+additional"
	The additional part will be kept from what it was before
--write=flutter+ for the same functionality like just "flutter"
	But here, tagger tries to increment the additional part
`

var RootCmd = &cobra.Command{
	Use:          "tagger",
	Short:        "Create a new git tag",
	Long:         rootCmdDescription,
	SilenceUsage: true,
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
		} else if flagPatch {
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

			newTag.MinusAddition = hash[:flagHash]
		}

		if !flagDry {
			if flagWrite != "" {
				uncommittedChanges, err := hasUncommittedChanges()
				if err != nil {
					return fmt.Errorf("failed to check, if uncommitted changes exist: %v", err)
				}
				if uncommittedChanges {
					return fmt.Errorf("cannot use 'write' flag, because there are uncommitted changes")
				}

				switch flagWrite {
				case "npm":
					err = writeVersionToPackageJson(newTag)
					if err != nil {
						return fmt.Errorf("failed to write package.json: %v", err)
					}
					err = commitAll(newTag.String())
					if err != nil {
						return fmt.Errorf("failed to create commit: %v", err)
					}
				case "flutter":
					err = writeVersionToPubspecYaml(newTag, false)
					if err != nil {
						return fmt.Errorf("failed to write pubspec.yaml: %v", err)
					}
					err = commitAll(newTag.String())
					if err != nil {
						return fmt.Errorf("failed to create commit: %v", err)
					}
				case "flutter+":
					err = writeVersionToPubspecYaml(newTag, true)
					if err != nil {
						return fmt.Errorf("failed to write pubspec.yaml: %v", err)
					}
					err = commitAll(newTag.String())
					if err != nil {
						return fmt.Errorf("failed to create commit: %v", err)
					}
				default:
					return fmt.Errorf("unknown write option: %s", flagWrite)
				}
			}

			err = createTag(newTag)
			if err != nil {
				return fmt.Errorf("failed to create tag: %s", err.Error())
			}
		}

		fmt.Printf("Tagged %s -> %s\n", latestTag, newTag)
		return nil
	},
}

var TagCmd = &cobra.Command{
	Use:          "tag",
	Short:        "Create a tag without looking up current tags",
	Long:         "Create a tag using the user input, instead of searching the latest tag and create a tag depending on it",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return fmt.Errorf("usage: tagger tag [version]")
		}

		var userInput string
		if len(args) == 1 {
			userInput = args[0]
		} else {
			prompt := promptui.Prompt{
				Label: "New version",
				Validate: func(s string) error {
					if !versionRegex.MatchString(s) {
						return fmt.Errorf("the version must be in the format v1.2.3")
					}
					return nil
				},
			}
			result, err := prompt.Run()

			if err != nil {
				return fmt.Errorf("prompt failed %v", err)
			}

			userInput = result
		}

		groups := versionRegex.FindStringSubmatch(userInput)

		if groups == nil {
			return fmt.Errorf("the version must be in the format v1.2.3")
		}

		major, _ := strconv.Atoi(groups[1])
		minor, _ := strconv.Atoi(groups[2])
		patch, _ := strconv.Atoi(groups[3])
		newTag := Tag{
			Major: major,
			Minor: minor,
			Patch: patch,
		}

		tags, err := getAllGitTags()
		if err != nil {
			return fmt.Errorf("failed to fetch git tags for validation: %s", err.Error())
		}

		for _, tag := range tags {
			if tag.Equals(newTag) {
				return fmt.Errorf("version tag already created: %s", tag.String())
			}
		}

		err = createTag(newTag)
		if err != nil {
			return fmt.Errorf("failed to create tag: %s", err.Error())
		}

		fmt.Printf("Tagged %s\n", newTag)
		return nil
	},
}

var ListCmd = &cobra.Command{
	Use:          "list",
	Short:        "List current set version tags",
	Long:         "List current set version tags",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		tags, err := getAllGitTags()
		if err != nil {
			return fmt.Errorf("failed to fetch git tags: %s", err.Error())
		}

		if len(tags) == 0 {
			return fmt.Errorf("no tags found")
		}

		for _, tag := range tags {
			fmt.Println(tag)
		}

		return nil
	},
}

var NpmCmd = &cobra.Command{
	Use:          "npm",
	Short:        "Tag commit using package.json",
	Long:         "Read the version from the package.json file and tag the current commit this version",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		newTag, err := readVersionFromPackageJson()
		if err != nil {
			return err
		}

		tags, err := getAllGitTags()
		if err != nil {
			return fmt.Errorf("failed to fetch git tags for validation: %s", err.Error())
		}

		for _, tag := range tags {
			if tag.Equals(newTag) {
				return fmt.Errorf("version tag already created: %s", tag.String())
			}
		}

		err = createTag(newTag)
		if err != nil {
			return fmt.Errorf("failed to create tag: %s", err.Error())
		}

		fmt.Printf("Tagged %s\n", newTag)
		return nil
	},
}

var FlutterCmd = &cobra.Command{
	Use:          "flutter",
	Short:        "Tag commit using pubspec.yaml",
	Long:         "Read the version from the pubspec.yaml file and tag the current commit this version",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		newTag, build, err := readVersionFromPubspecYaml()
		if err != nil {
			return err
		}

		tags, err := getAllGitTags()
		if err != nil {
			return fmt.Errorf("failed to fetch git tags for validation: %s", err.Error())
		}

		if !flagBuild {
			for _, tag := range tags {
				if tag.Equals(newTag) {
					return fmt.Errorf("version tag already created: %s", tag.String())
				}
			}
		}

		if flagBuild {
			newTag.PlusAddition = strconv.Itoa(build + 1)

			if !flagDry {
				uncommittedChanges, err := hasUncommittedChanges()
				if err != nil {
					return fmt.Errorf("failed to check, if uncommitted changes exist: %v", err)
				}
				if uncommittedChanges {
					return fmt.Errorf("cannot use 'build' flag, because there are uncommitted changes")
				}
				err = writeVersionAndBuildToPubspecYaml(newTag, build+1)
				if err != nil {
					return fmt.Errorf("failed to update build number in pubspec.yaml: %v", err)
				}
				err = commitAll(newTag.String())
				if err != nil {
					return fmt.Errorf("failed to create commit: %v", err)
				}
			}
		}

		if !flagDry {
			err = createTag(newTag)
			if err != nil {
				return fmt.Errorf("failed to create tag: %s", err.Error())
			}
		}

		fmt.Printf("Tagged %s\n", newTag)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(TagCmd, ListCmd, NpmCmd, FlutterCmd)

	RootCmd.Flags().BoolVar(&flagMajor, "major", false, "Increase major part")
	RootCmd.Flags().BoolVar(&flagMinor, "minor", false, "Increase minor part")
	RootCmd.Flags().BoolVar(&flagPatch, "patch", false, "Increase patch part")
	RootCmd.Flags().BoolVar(&flagDateTime, "datetime", false, "Set minor and patch to date time")
	RootCmd.Flags().IntVar(&flagHash, "hash", 0, "Add commit hash to end")
	RootCmd.PersistentFlags().BoolVarP(&flagDry, "dry", "d", false, "Show new tag but don't apply")
	RootCmd.Flags().StringVar(&flagWrite, "write", "", "Write the version into file (see help)")

	FlutterCmd.Flags().BoolVar(&flagBuild, "build", false, "Increase build number and tag with +build")
}
