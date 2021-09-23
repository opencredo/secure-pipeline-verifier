package gitlab

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
	"time"
)

type CommitInfo struct {
	Repo               string
	CommitUrl          string
	Date               *time.Time
	AuthorName         string
	AuthorEmail        string
	VerifiedSignature  bool
	VerificationReason string
}

type BranchCommitProtection struct {
	Repo               string
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

func GetChangesToCiCd(client *gitlab.Client, projectPath string, path string, since time.Time) ([]CommitInfo, error) {
	opt := &gitlab.ListCommitsOptions{
		Path:        &path,
		Since:       &since,
		ListOptions: gitlab.ListOptions{PerPage: 20},
	}

	// get all pages of results
	var allCommits []*gitlab.Commit
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
		isVerified, reason := checkCommitSignature(client, projectPath, repoCommit.ID)

		commitsInfo = append(commitsInfo,
			CommitInfo{
				Repo:               projectPath,
				CommitUrl:          repoCommit.WebURL,
				Date:               repoCommit.AuthoredDate,
				AuthorName:         repoCommit.AuthorName,
				AuthorEmail:        repoCommit.AuthorEmail,
				VerifiedSignature:  isVerified,
				VerificationReason: reason,
			},
		)
	}

	return commitsInfo
}

// checkCommitSignature: Checks if a commit has a signature
func checkCommitSignature(client *gitlab.Client, projectPath string, sha string) (bool, string) {
	// For unsigned commits we get a 404 response
	sig, _, _ := client.Commits.GetGPGSiganature(projectPath, sha)
	if sig != nil {
		return true, sig.VerificationStatus
	}
	return false, ""
}

// GetBranchSignatureProtection Control-2
// Client user needs to be an admin of the repo to get this info
// This endpoint returns an error containing string "404 Branch not protected" when the branch is not protected
func GetBranchSignatureProtection(client *gitlab.Client, org string, repo string, branches []string) []BranchCommitProtection {

	var branchesProtection []BranchCommitProtection
	for _, branch := range branches {
		protectedBranch, _, err := client.Branches.GetBranch("da", branch)
		if err != nil {
			branchesProtection = append(branchesProtection,
				BranchCommitProtection{GitHubRepo: org + "/" + repo, BranchName: branch, Error: err.Error()})
			continue
		}

		branchesProtection = append(branchesProtection,
			BranchCommitProtection{
				GitHubRepo:         org + "/" + repo,
				BranchName:         branch,
				SignatureProtected: protectedBranch.GetEnabled(),
			},
		)
	}

	return branchesProtection
}
