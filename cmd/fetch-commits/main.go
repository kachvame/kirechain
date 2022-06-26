package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Khan/genqlient/graphql"
	"github.com/kachvame/kirechain/chain"
	"github.com/kachvame/kirechain/github"
	"golang.org/x/oauth2"
)

var emails = [...]string{
	"kiril@senteca.com",
	"51754423+kirilsenteca@users.noreply.github.com",
}

var defaultBranches = [...]string{
	"master",
	"platform/v1",
}


func main() {
	if err := run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()
	token := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(ctx, token)

	client := graphql.NewClient("https://api.github.com/graphql", httpClient)

	var messages []string

	err := github.WalkRepositories(ctx, client, "senteca", func(repository github.RepositoryNode) error {
		log.Println("querying commits for repo", repository.Name)

		commits, err := github.GetCommitsByAuthor(ctx, client, "senteca", repository.Name, emails[:], defaultBranches[:])
		if err != nil {
			return fmt.Errorf("failed to get commits for repository %s: %w", repository.Name, err)
		}

		for _, commit := range commits {
			messages = append(messages, commit.Message)
		}

		log.Println("got", len(commits), "commits for repository", repository.Name)
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to fetch repositories: %w", err)
	}

	file, err := os.Create("entries.json")
	if err != nil {
		return fmt.Errorf("failed to open entries.json: %w", err)
	}

	if err = json.NewEncoder(file).Encode(chain.Entries{Commits: messages}); err != nil {
		return fmt.Errorf("failed to encode commit messages: %w", err)
	}

	return nil

}
