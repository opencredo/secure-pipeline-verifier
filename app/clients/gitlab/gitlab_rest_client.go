package gitlab

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
	"time"
)


type CommitInfo struct {
	Repo         string
	CommitUrl          string
	Date               time.Time
	AuthorName         string
	AuthorEmail        string
	VerifiedSignature  bool
	VerificationReason string
}

type BranchCommitProtection struct {
	Repo         string
	BranchName         string
	SignatureProtected bool
	Error              string
}

type AutomationKey struct {
	ID           int64
	Title        string
	Verified     bool
	ReadOnly     bool
	CreationDate time.Time
}

func GetChangesToCiCd(client *gitlab.Client, org string, repo string, path string, since time.Time) ([]CommitInfo, error) {
	opt := &gitlab.ListCommitsOptions{
		Path: &path,
		Since: &since,
		ListOptions: gitlab.ListOptions{PerPage: 20},
	}

	// get all pages of results
	var allCommits []*gitlab.Commit
	// Project path
	projectPath := fmt.Sprintf("%s/%s", org, repo)
	for {
		commits, resp, err := client.Commits.ListCommits(
			projectPath,
			opt,
		)
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

	return getCommitsInfo(client, projectPath, allCommits), nil
}

func getCommitsInfo(client *gitlab.Client, projectPath string, repositoryCommits []*gitlab.Commit) []CommitInfo {
	var commitsInfo []CommitInfo
	for _, repoCommit := range repositoryCommits {
		isVerified, reason := getCommitSignature(client, projectPath, repoCommit.ID)

		commitsInfo = append(commitsInfo,
			CommitInfo{
				Repo:               projectPath,
				CommitUrl:          repoCommit.WebURL,
				Date:               *repoCommit.AuthoredDate,
				AuthorName:         repoCommit.AuthorName,
				AuthorEmail:        repoCommit.AuthorEmail,
				VerifiedSignature:  isVerified,
				VerificationReason: reason,
			},
		)
	}

	return commitsInfo
}

func getCommitSignature(client *gitlab.Client, projectPath string, sha string) (bool, string) {
	// For unsigned commits we get a 404 response
	sig, _, _ := client.Commits.GetGPGSiganature(projectPath, sha)

	if sig != nil {
		return true, sig.VerificationStatus
	}
	return false, ""
}