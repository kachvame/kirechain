package github

import (
	"context"
	"github.com/Khan/genqlient/graphql"
	"log"
)

type Repository = GetRepositoriesOrganizationRepositoriesRepositoryConnectionNodesRepository
type Branch = *GetCommitsOfBranchByAuthorRepositoryRefTargetCommit

var defaultBranches = [...]string{"master", "platform/v1"}

func WalkRepositories(ctx context.Context, client graphql.Client, login string, walkFn func(Repository)) error {
	cursor := ""

	for {
		response, err := GetRepositories(ctx, client, login, cursor)
		if err != nil {
			return err
		}

		repositories := response.Organization.Repositories
		for _, repository := range repositories.Nodes {
			walkFn(repository)
		}

		pageInfo := repositories.PageInfo
		if !pageInfo.HasNextPage {
			break
		}

		cursor = pageInfo.EndCursor
	}

	return nil
}

func GetAllCommitsOfBranchByAuthor(ctx context.Context, client graphql.Client, owner string, repository string, branch string, emails []string) ([]CommitNode, error) {
	branch = "refs/heads/" + branch

	var allCommits []CommitNode
	cursor := ""

	for {
		response, err := GetCommitsOfBranchByAuthor(ctx, client, owner, repository, branch, emails, cursor)
		if err != nil {
			return nil, err
		}

		target := response.Repository.Ref.Target
		if target == nil {
			return allCommits, nil
		}

		branchHistory := target.(Branch).History

		commits := branchHistory.Nodes
		for _, commit := range commits {
			allCommits = append(allCommits, commit.CommitNode)
		}

		pageInfo := branchHistory.PageInfo
		if !pageInfo.HasNextPage {
			break
		}

		cursor = pageInfo.EndCursor
	}

	return allCommits, nil
}

func WalkPRs(ctx context.Context, client graphql.Client, owner string, repository string, walkFn func(number int) error) error {
	cursor := ""

	for {
		response, err := GetPRs(ctx, client, owner, repository, cursor)
		if err != nil {
			return err
		}

		pullRequests := response.Repository.PullRequests
		for _, pullRequest := range pullRequests.Nodes {
			err = walkFn(pullRequest.Number)
			if err != nil {
				return err
			}
		}

		pageInfo := pullRequests.PageInfo
		if !pageInfo.HasNextPage {
			break
		}

		cursor = pageInfo.EndCursor
	}

	return nil
}

func GetCommitsInPRByAuthor(ctx context.Context, client graphql.Client, owner string, repository string, number int, emails []string) ([]CommitNode, error) {
	var allCommits []CommitNode
	cursor := ""

	for {
		response, err := GetCommitsInPR(ctx, client, owner, repository, number, cursor)
		if err != nil {
			return nil, err
		}

		commits := response.Repository.PullRequest.Commits
		for _, commit := range commits.Nodes {
			commitNode := commit.Commit.CommitNode

			isByAuthor := false
			for _, email := range emails {
				if email == commitNode.Author.Email {
					isByAuthor = true
					break
				}
			}

			if !isByAuthor {
				continue
			}

			allCommits = append(allCommits, commitNode)
		}

		pageInfo := commits.PageInfo
		if !pageInfo.HasNextPage {
			break
		}

		cursor = pageInfo.EndCursor
	}

	return allCommits, nil
}

func GetCommitsByAuthor(ctx context.Context, client graphql.Client, owner string, repository string, emails []string) ([]CommitNode, error) {
	var allCommits []CommitNode
	for _, branch := range defaultBranches {
		commits, err := GetAllCommitsOfBranchByAuthor(ctx, client, owner, repository, branch, emails)
		if err != nil {
			return nil, err
		}

		log.Println("got", len(commits), "commits for branch", branch)

		allCommits = append(allCommits, commits...)
	}

	if len(allCommits) == 0 {
		log.Println("found no commits on default branches, skipping repository", repository)

		return allCommits, nil
	}

	err := WalkPRs(ctx, client, owner, repository, func(number int) error {
		commits, err := GetCommitsInPRByAuthor(ctx, client, owner, repository, number, emails)
		if err != nil {
			return err
		}

		if len(commits) > 0 {
			log.Println("got", len(commits), "for PR", number)
		}

		allCommits = append(allCommits, commits...)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return allCommits, nil
}
