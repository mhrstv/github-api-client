package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github-api-client/internal/github"
	"github-api-client/internal/input"
	"github-api-client/internal/report"
	"github-api-client/internal/stats"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <usernames-file>\n", os.Args[0])
		os.Exit(1)
	}

	usernames, err := input.ReadUsernames(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	client := github.NewClient()
	ctx := context.Background()

	var all []stats.UserStats
	for _, name := range usernames {
		s, err := collect(ctx, client, name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "skipping %s: %v\n", name, err)
			continue
		}
		all = append(all, s)
	}

	if len(all) == 0 {
		log.Fatal("no users could be fetched")
	}
	report.Print(os.Stdout, all)
}

func collect(ctx context.Context, c *github.Client, username string) (stats.UserStats, error) {
	user, err := c.GetUser(ctx, username)
	if err != nil {
		return stats.UserStats{}, err
	}

	repos, err := c.GetRepos(ctx, username)
	if err != nil {
		return stats.UserStats{}, err
	}

	var langs []github.Languages
	for _, r := range repos {
		l, err := c.GetLanguages(ctx, username, r.Name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s/%s languages: %v\n", username, r.Name, err)
			continue
		}
		langs = append(langs, l)
	}

	return stats.Compute(user, repos, langs), nil
}
