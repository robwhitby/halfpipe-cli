package linter

import (
	"fmt"
	"regexp"

	"github.com/robwhitby/halfpipe-cli/manifest"
)

type RepoLinter struct{}

func (r RepoLinter) isPrivateRepo(repo string) bool {
	regex := `git@github.com:[a-zA-Z0-9]+\/[a-zA-Z0-9_-]+.git`
	matches, _ := regexp.MatchString(regex, repo)
	return matches
}

func (r RepoLinter) isPublicRepo(repo string) bool {
	regex := `https:\/\/github.com\/[a-zA-Z0-9]+\/[a-zA-Z0-9]+`
	matches, _ := regexp.MatchString(regex, repo)
	return matches
}

func (r RepoLinter) Lint(manifest manifest.Manifest) (Result, error) {
	result := Result{
		Linter: "Repo",
		Errors: []Error{},
	}
	if !r.isPrivateRepo(manifest.Repo.Uri) && !r.isPublicRepo(manifest.Repo.Uri) {
		result.Errors = append(result.Errors, Error{
			Message: fmt.Sprintf("'%s' does not look like a real repo!", manifest.Repo.Uri),
		})
	}

	if r.isPrivateRepo(manifest.Repo.Uri) && manifest.Repo.PrivateKey == "" {
		result.Errors = append(result.Errors, Error{
			Message: "It looks like you are refering to a private repo, but no private key provided in `repo.private_key`",
		})
	}
	return result, nil
}
