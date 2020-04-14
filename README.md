# Tagger

Tagger is a command line tool to create git tags.
After installation, execute `tagger` to create a git tag.

It scans the existing tags for the format `v1.2.3` or `v1.2.3-note` and gets the maximum
(e.g. `v2.0.0` is higher than `v1.19.5`).
It has some strategies to increase the version.

## Parameters

You can manipulate the way tagger works with those parameters:

- `dry`: Dryrun, just show the version, don't actually tag
- `hash`: Add the commit hash as note (behind dash)
- `strategy <strategy>`: Chose the version increase strategy (default is `patch`)
- `note <note>`: Set the note part of the version (behind dash)

Example usage: `tagger --strategy datetime --dry`

*Attention*: `note` will overwrite the result of `hash`!

## Strategies

### Patch increase `patch`

This is the default strategy. It just increases the patch part (3rd number) by one.

Example: After the last version `v1.2.3`, the new version will be `v1.2.4`.

### Other version parts `major` `minor`

This is the same like `patch`, but instead the patch part, it will increase the major or minor part.

### Datestamp and Timestamp `datetime`

This method calculates a datestamp and a timestamp.
The datestamp is `unix % (60 * 60 * 24)`, while the timestamp is `unix / (60 * 60 * 24)`.

This strategy embeds the time information into the version.
You can also see the time difference between versions.
It is usefull for automated versioning, because it will use minor and patch at the same time.
The downsite is a bad readability.
When you create many versions automatically, this is less of a matter.

Example: `v0.69900.18366`
