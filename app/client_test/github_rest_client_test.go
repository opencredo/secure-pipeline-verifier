package client_test

import (
	"github.com/google/go-github/v38/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"secure-pipeline-poc/app/client"
	"testing"
	"time"
)

var firstCommitDate = time.Date(2021, time.Month(7), 21, 15, 10, 30, 0, time.UTC)
var secondCommitDate = time.Date(2021, time.Month(7), 23, 1, 10, 30, 0, time.UTC)

var keyCreationTs = time.Date(2021, time.Month(8), 03, 14, 50, 23, 0, time.UTC)

// control-1
func TestGetChangesToCiCdReturnsCommits(t *testing.T) {
	assert := assert.New(t)
	mockedHttpClient := createMockedGitHubHttpClientReturnsRepoCommits()

	githubClient := github.NewClient(mockedHttpClient)

	sinceDate := time.Date(2021, time.Month(7), 1, 9, 00, 00, 0, time.UTC)
	cicdChanges, _ := client.GetChangesToCiCd(githubClient, "my-org-123", "my-awesome-app", ".github/workspace", sinceDate)

	firstCommit := cicdChanges[0]
	assert.Equal(firstCommitDate, firstCommit.Date)
	assert.Equal("John White", firstCommit.AuthorName)
	assert.Equal("jwhite", firstCommit.AuthorUsername)
	assert.Equal("jwhite@email.com", firstCommit.AuthorEmail)
	assert.Equal(true, firstCommit.VerifiedSignature)
	assert.Equal("valid", firstCommit.VerificationReason)

	secondCommit := cicdChanges[1]
	assert.Equal(secondCommitDate, secondCommit.Date)
	assert.Equal("Dodgy User", secondCommit.AuthorName)
	assert.Equal("aDodgyOne", secondCommit.AuthorUsername)
	assert.Equal("user12345@email.com", secondCommit.AuthorEmail)
	assert.Equal(false, secondCommit.VerifiedSignature)
	assert.Equal("unsigned", secondCommit.VerificationReason)
}

/*func TestGetChangesToCiCdReturnsError(t *testing.T) {
	assert := assert.New(t)

	mockedHttpClient := createMockedRepositoryCommitsGitHubHttpClientReturnsError()

	githubClient := github.NewClient(mockedHttpClient)

	sinceDate := time.Date(2021, time.Month(7), 1, 9, 00, 00, 0, time.UTC)
	cicdCommits, err := client.GetChangesToCiCd(githubClient, "my-org-123", "my-awesome-app", ".travis-ci", sinceDate)
	assert.Nil( cicdCommits)
	assert.NotNil(err)
	assert.Equal("error retrieving commits - 401 unauthorized", err.Error())
}*/

func createMockedGitHubHttpClientReturnsRepoCommits() *http.Client {
	return mock.NewMockedHTTPClient(
		mock.WithRequestMatch(
			mock.GetReposCommitsByOwnerByRepo,
			[]github.RepositoryCommit{
				{
					HTMLURL: github.String("https://github.com/myorg/myrepo/commit/123456789"),
					Author: &github.User{
						Login: github.String("jwhite"),
					},
					Commit: &github.Commit{
						Author: &github.CommitAuthor{
							Date:  &firstCommitDate,
							Name:  github.String("John White"),

							Email: github.String("jwhite@email.com"),
						},
						Verification: &github.SignatureVerification{
							Verified: github.Bool(true),
							Reason: github.String("valid"),
					},
					},
				},
				{
					HTMLURL: github.String("https://github.com/myorg/myrepo/commit/1987654321"),
					Author: &github.User{
						Login: github.String("aDodgyOne"),
					},
					Commit: &github.Commit{
						Author: &github.CommitAuthor{
							Date:  &secondCommitDate,
							Name:  github.String("Dodgy User"),
							Email: github.String("user12345@email.com"),
						},
						Verification: &github.SignatureVerification{
							Verified: github.Bool(false),
							Reason: github.String("unsigned"),
						},
					},
				},
			},
			[]github.Response{
				{
					NextPage: 0,
				},
			},
		),
	)
}

/*func createMockedRepositoryCommitsGitHubHttpClientReturnsError() *http.Client {
	return mock.NewMockedHTTPClient(
		mock.WithRequestMatch(
			mock.GetReposCommitsByOwnerByRepo,
			nil,
			errors.New("error retrieving commits - 401 unauthorized"),
		),
	)
}*/

// control-2 - protected branch
func TestGetBranchSignatureProtection_true(t *testing.T) {
	assert := assert.New(t)

	mockedHttpClient := createMockedGitHubHttpClientReturnsBranchProtected()

	githubClient := github.NewClient(mockedHttpClient)

	branchesProtection := client.GetBranchSignatureProtection(githubClient, "my-org", "my-repo", []string{"master"})
	assert.NotNil(branchesProtection)
	assert.Equal(1, len(branchesProtection))
	branchProtection := branchesProtection[0]
	assert.Equal("my-org/my-repo", branchProtection.GitHubRepo)
	assert.Equal("master", branchProtection.BranchName)
	assert.Equal(true, branchProtection.SignatureProtected)
	assert.Empty(branchProtection.Error)
}

// control-2 - unprotected branch
func TestGetBranchSignatureProtection_false(t *testing.T) {
	assert := assert.New(t)

	mockedHttpClient := createMockedGitHubHttpClientReturnsBranchNotProtected()

	githubClient := github.NewClient(mockedHttpClient)

	branchesProtection := client.GetBranchSignatureProtection(githubClient, "my-org", "my-repo", []string{"develop"})
	assert.NotNil(branchesProtection)
	assert.Equal(1, len(branchesProtection))
	branchProtection := branchesProtection[0]
	assert.Equal("my-org/my-repo", branchProtection.GitHubRepo)
	assert.Equal("develop", branchProtection.BranchName)
	assert.Equal(false, branchProtection.SignatureProtected)
	assert.Equal("Branch not protected", branchProtection.Error)
}

func createMockedGitHubHttpClientReturnsBranchProtected() *http.Client {
	return mock.NewMockedHTTPClient(
		mock.WithRequestMatch(
			mock.GetReposBranchesProtectionRequiredSignaturesByOwnerByRepoByBranch,
			github.SignaturesProtectedBranch{
				URL: github.String("https://api.github.com/repos/my-org/my-repo/branches/master/protection/required_signatures"),
				Enabled: github.Bool(true),
			},
			[]github.Response{
				{
					NextPage: 0,
				},
			},
		),
	)
}

func createMockedGitHubHttpClientReturnsBranchNotProtected() *http.Client {
	return mock.NewMockedHTTPClient(
		mock.WithRequestMatchHandler(
			mock.GetReposBranchesProtectionRequiredSignaturesByOwnerByRepoByBranch,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mock.WriteError(
					w,
					http.StatusNotFound,
					"Branch not protected",
				)
			}),
		),
	)
}

// control-3
func TestGetAutomationKeysExpiryReturnsAKey(t *testing.T) {
	assert := assert.New(t)
	mockedHttpClient := createMockedGitHubHttpClientReturnsRepoDeployKeys()

	githubClient := github.NewClient(mockedHttpClient)

	automationKeys, _ := client.GetAutomationKeysExpiry(githubClient, "my-org-456", "my-other-app")
	assert.NotNil(automationKeys)
	assert.Equal(1, len(automationKeys))

	key := automationKeys[0]

	assert.Equal(int64(1), key.ID )
	assert.Equal("my-deploy-key", key.Title )
	assert.Equal(true, key.ReadOnly )
	assert.Equal(true, key.Verified )
	assert.Equal(keyCreationTs, key.CreationDate )

}

// control-4
func TestGetAutomationKeysPermissions(t *testing.T) {
	assert := assert.New(t)
	mockedHttpClient := createMockedGitHubHttpClientReturnsRepoDeployKeys()

	gitHubClient := github.NewClient(mockedHttpClient)
	automationKeys, _ := client.GetAutomationKeysPermissions(gitHubClient, "my-org", "my-repo")
	assert.NotNil(automationKeys)

	assert.NotNil(automationKeys)
	assert.Equal(1, len(automationKeys))

	key := automationKeys[0]

	assert.Equal(int64(1), key.ID )
	assert.Equal("my-deploy-key", key.Title )
	assert.Equal(true, key.ReadOnly )
	assert.Equal(true, key.Verified )
	assert.Equal(keyCreationTs, key.CreationDate )
}

func createMockedGitHubHttpClientReturnsRepoDeployKeys() *http.Client {
	return mock.NewMockedHTTPClient(
		mock.WithRequestMatch(
			mock.GetReposKeysByOwnerByRepo,
			[]github.Key {
				{
					ID: github.Int64(1),
					Title: github.String("my-deploy-key"),
					ReadOnly: github.Bool(true),
					Verified: github.Bool(true),
					CreatedAt: &github.Timestamp{
						Time: keyCreationTs,
					},
				},
			},
			[]github.Response{
				{
					NextPage: 0,
				},
			},
		),
	)
}