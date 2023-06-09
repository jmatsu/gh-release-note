package main

import (
	"errors"
	"fmt"
	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/jmatsu/gh-release-note/cmd"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/exec"
	"time"
)

// https://github.com/c-bata/go-prompt/issues/228#issuecomment-820639887
func handleExit() {
	rawModeOff := exec.Command("/bin/stty", "-raw", "echo")
	rawModeOff.Stdin = os.Stdin
	_ = rawModeOff.Run()
	_ = rawModeOff.Wait()
}

func main() {
	defer handleExit()

	app := &cli.App{
		Name:  "gh release-note",
		Usage: "Generate release notes from pull requests.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "repo",
				Aliases: []string{"R"},
				Usage:   "Select another repository using the [HOST/]OWNER/REPO format",
			},
			&cli.StringFlag{
				Name:     "base",
				Aliases:  []string{"B"},
				Usage:    "Select a base branch name",
				Required: true,
			},
			&cli.IntFlag{
				Name:  "limit",
				Usage: "Select a limit to get pull requests and/or tags at once",
				Value: 100,
			},
			&cli.BoolFlag{
				Name:  "skip-generate",
				Usage: "Specify if you would like not to generate the release note but want pull requests",
				Value: false,
			},
			&cli.StringFlag{
				Name:  "since",
				Usage: "Select a date for the start position of the search range using YYYY-MM-dd[ HH:mm:SSZ]",
			},
			&cli.StringFlag{
				Name:  "since-tag",
				Usage: "Select a tag name for the start position of the search range",
			},
			&cli.StringFlag{
				Name:  "until",
				Usage: "Select a date for the end position of the search range using YYYY-MM-dd[THH:mm:SSZ]",
			},
			&cli.StringFlag{
				Name:  "until-tag",
				Usage: "Select a tag name for the end position of the search range",
			},
		},
		Action: func(c *cli.Context) error {
			return cmd.GenerateReleaseNote(c, func(c *cli.Context) (cmd.ReleaseNoteOption, error) {
				option := cmd.ReleaseNoteOption{
					Base:         c.String("base"),
					Limit:        c.Int("limit"),
					SkipGenerate: c.Bool("skip-generate"),
				}

				if repo, err := getRepo(c, "repo"); err != nil {
					return option, err
				} else {
					option.Repo = repo
				}

				option.Base = c.String("base")

				if c.IsSet("since") {
					t, err := parseDateTime(c, "since")

					if err != nil {
						return option, err
					}

					option.SinceDate = t
				} else if v := c.String("since-tag"); c.IsSet("since-tag") {
					option.SinceTagName = &v
				}

				if c.IsSet("until") {
					t, err := parseDateTime(c, "until")

					if err != nil {
						return option, err
					}

					option.UntilDate = t
				} else if v := c.String("until-tag"); c.IsSet("until-tag") {
					option.UntilTagName = &v
				}

				return option, nil
			})
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func getRepo(c *cli.Context, name string) (repository.Repository, error) {
	var repo repository.Repository

	if c.IsSet(name) {
		if v, err := repository.Parse(c.String(name)); err != nil {
			return repo, errors.Join(errors.New(fmt.Sprintf("%s is not a valid repository", name)), err)
		} else {
			repo = v
		}
	} else if v, err := repository.Current(); err != nil {
		return repo, errors.New("cannot get the current repository")
	} else {
		repo = v
	}

	return repo, nil
}

func parseDateTime(c *cli.Context, name string) (*time.Time, error) {
	if v := c.String(name); c.IsSet(name) {
		if len(v) > 10 {
			t, err := time.Parse(time.DateTime, v)

			if err != nil {
				return nil, errors.New(fmt.Sprintf("%s is not a valid date format", v))
			}

			return &t, nil
		} else {
			t, err := time.Parse(time.DateOnly, v)

			if err != nil {
				return nil, errors.New(fmt.Sprintf("%s is not a valid date format", v))
			}

			return &t, nil
		}
	}

	return nil, nil
}
