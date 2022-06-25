package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Khan/genqlient/graphql"
	"github.com/kachvame/kirechain/chain"
	"github.com/kachvame/kirechain/github"
	"golang.org/x/oauth2"
	"log"
	"os"
)

var emails = [...]string{
	"kiril@senteca.com",
	"51754423+kirilsenteca@users.noreply.github.com",
}

func main() {
	ctx := context.Background()

	token := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(ctx, token)

	client := graphql.NewClient("https://api.github.com/graphql", httpClient)

	var messages []string

	err := github.WalkRepositories(ctx, client, "senteca", func(repository github.Repository) {
		log.Println("querying commits for repo", repository.Name)

		commits, err := github.GetCommitsByAuthor(ctx, client, "senteca", repository.Name, emails[:])
		if err != nil {
			log.Fatalf("failed to get commits for repository %s: %s", repository.Name, err)
		}

		for _, commit := range commits {
			messages = append(messages, commit.Message)
		}

		log.Println("got", len(commits), "commits for repository", repository.Name)
	})

	if err != nil {
		log.Fatalln("failed to fetch repositories:", err)
	}

	file, err := os.Create("entries.json")
	if err != nil {
		log.Fatalln("failed to open entries.json:", err)
	}

	if err = json.NewEncoder(file).Encode(chain.Entries{Commits: messages}); err != nil {
		log.Fatalln("failed to encode commit messages:", err)
	}

	fmt.Println(len(messages))
}
