# gh-release-note

A gh extension to generate release notes from pull requests merged between two git refs.

# THIS COMMAND IS STILL WIP

This command just prints the summary of the pull requests that have beem merged between two git refs.

TODOs

- [ ] Generate the release note from the pull request information
- [ ] Generate the pull request summary for each pull request. (You may have to specify one pull request.)
- [ ] Extract and use additional information from description and/or commit messages for the release note

# Usage

```
Usage:
  gh release-note -h
  gh release-note -R OWNER/REPO -B BRANCH [-S SINCE] [-U UNTIL] [--tag-prefix TAG_PREFIX]

A gh extension to generate release notes from pull requests merged between two git refs.

Available options:
-h, --help
                  Print this help and exit
-R, --repo OWNER/REPO
                  Select another repository using the OWNER/REPO format
-B, --base BRANCH
                  The branch name that the pull requests have been merged into
-S, --since DATE
                  The date from the tag or which the pull request have been merged into the base branch
-U, --until DATE
                  The date until the tag or which the pull requests have been merged into the base branch
--tag-prefix TAG_PREFIX
                  The prefix of the tag name to filter tag names
```