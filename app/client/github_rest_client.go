package client

import (
	"context"
	"fmt"
	"github.com/google/go-github/v38/github"
	_ "github.com/google/go-github/v38/github"
	"golang.org/x/oauth2"
	"time"
)

func NewClient(oauthToken string) *github.Client {
	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: oauthToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

type CommitInfo struct {
	GitHubRepo string
	CommitUrl string
	Date        time.Time
	AuthorName  string
	AuthorUsername  string
	AuthorEmail string
	VerifiedSignature  bool
	VerificationReason string
}

type BranchCommitProtection struct {
	GitHubRepo string
	BranchName string
	SignatureProtected bool
	Error string
}

type AuthKeyPermission struct {
	ID int64
	Title string
	Verified bool
	ReadOnly bool
}

// GetChangesToCiCd Control-1
func GetChangesToCiCd(client *github.Client, org string, repo string, path string, since time.Time) ([]CommitInfo, error) {
	ctx := context.Background()

	opt := &github.CommitsListOptions{
		Path: path, Since: since,
		ListOptions: github.ListOptions{ PerPage: 20 },
	}

	// get all pages of results
	var allCommits []*github.RepositoryCommit
	for {
		commits, resp, err := client.Repositories.ListCommits(ctx, org, repo, opt)
		if err != nil {
			fmt.Printf("Error retrieving changes to CI/CD folder. Error: %s, Response Status: %s", err.Error(), resp.Status)
			return nil, err
		}
		allCommits = append(allCommits, commits...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return getCommitsInfo(org, repo, allCommits), nil
}

// GetBranchSignatureProtection Control-2
// Client user needs to be an admin of the repo to get this info
// This endpoint returns an error containing string "404 Branch not protected" when the branch is not protected
func GetBranchSignatureProtection(client *github.Client, org string, repo string, branches []string) []BranchCommitProtection {
	ctx := context.Background()

	var branchesProtection []BranchCommitProtection
	for _, branch := range branches {
		protectedBranch, _, err := client.Repositories.GetSignaturesProtectedBranch(ctx, org, repo, branch)
		if err != nil {
			branchesProtection = append(branchesProtection, BranchCommitProtection{BranchName: branch, Error: err.Error()})
			continue
		}

		branchesProtection = append(branchesProtection,
			BranchCommitProtection{
				GitHubRepo: org + "/" + repo,
				BranchName: branch,
				SignatureProtected: protectedBranch.GetEnabled(),
			},
		)
	}

	return branchesProtection
}

func getCommitsInfo(org string, repo string, repositoryCommits []*github.RepositoryCommit) []CommitInfo {
	var commitsInfo []CommitInfo
	for _, repoCommit := range repositoryCommits {
		var url = repoCommit.GetHTMLURL()
		var date = repoCommit.Commit.Author.GetDate()
		var authorName = repoCommit.Commit.Author.GetName()
		var authorUsername = repoCommit.Author.GetLogin()
		var authorEmail = repoCommit.Commit.Author.GetEmail()
		var verifiedSignature = repoCommit.Commit.Verification.GetVerified()
		var reason = repoCommit.Commit.Verification.GetReason()

		commitsInfo = append(commitsInfo,
			CommitInfo{
				GitHubRepo: 		org + "/" + repo,
				CommitUrl:          url,
				Date:               date,
				AuthorName:         authorName,
				AuthorUsername:     authorUsername,
				AuthorEmail:        authorEmail,
				VerifiedSignature:  verifiedSignature,
				VerificationReason: reason,
			},
		)
	}

	return commitsInfo
}

// GetAutomationKeysPermissionsInfo Control-4
func GetAutomationKeysPermissionsInfo(client *github.Client, org string, repo string) ([]AuthKeyPermission, error) {
	ctx := context.Background()

	opts := &github.ListOptions{ PerPage: 20 }

	keys, response, err := client.Repositories.ListKeys(ctx, org, repo, opts)
	if err != nil {
		fmt.Printf("Error checking Permissions for authomation keys. Error: %s, Response Status: %s", err.Error(), response.Status)
		return nil, err
	}

	var authKeysPermission []AuthKeyPermission
	for _, key:= range keys{
		authKeysPermission = append(authKeysPermission,
			AuthKeyPermission{
				ID: key.GetID(),
				Title: key.GetTitle(),
				Verified: key.GetVerified(),
				ReadOnly: key.GetReadOnly(),
			},
		)
	}

	return authKeysPermission, nil
}
