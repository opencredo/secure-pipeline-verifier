package gitlab

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
	"secure-pipeline-poc/app/config"
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

type RepoCommitProtection struct {
	Repo               string
	SignatureProtected bool
}

type AutomationKey struct {
	ID           int
	Title        string
	ReadOnly     bool
	CreationDate *time.Time
}

// Repo implements methods for the /project endpoints
type Repo interface {
	GetCommitsInfo(repositoryCommits []*gitlab.Commit) []CommitInfo
	GetProjectSignatureProtection() RepoCommitProtection
	GetChangesToCiCd(path string, since time.Time) ([]CommitInfo, error)
	CheckCommitSignature(sha string) (bool, string)
	GetAutomationKeys() ([]AutomationKey, error)
}

type Api struct{
	Client *gitlab.Client
	ProjectPath string
	Repo
}

func NewApi(token string, cfg *config.Config) *Api {
	client, _ := gitlab.NewClient(token)
	p := &Api{
		Client: client,
		ProjectPath: fmt.Sprintf("%s/%s", cfg.Project.Owner, cfg.Project.Repo),
	}
	p.Repo = p
	return p
}

// GetChangesToCiCd Control-1
// Returns commits for a specific item since a specific date
func (a *Api) GetChangesToCiCd(path string, since time.Time) ([]CommitInfo, error) {

	opt := &gitlab.ListCommitsOptions{
		Path:        &path,
		Since:       &since,
		ListOptions: gitlab.ListOptions{PerPage: 20},
	}

	// get all pages of results
	var allCommits []*gitlab.Commit
	for {
		commits, resp, err := a.Client.Commits.ListCommits(
			a.ProjectPath,
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

	return a.Repo.GetCommitsInfo(allCommits), nil
}

func (a *Api) GetCommitsInfo(repositoryCommits []*gitlab.Commit) []CommitInfo {
	var commitsInfo []CommitInfo
	for _, repoCommit := range repositoryCommits {
		isVerified, reason := a.Repo.CheckCommitSignature(repoCommit.ID)

		commitsInfo = append(commitsInfo,
			CommitInfo{
				Repo:               a.ProjectPath,
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

// GetProjectSignatureProtection Control-2
func (a *Api) GetProjectSignatureProtection() RepoCommitProtection {

	pushRules, _, _ := a.Client.Projects.GetProjectPushRules(a.ProjectPath)

	repoCommitProtection := RepoCommitProtection{
		Repo:               a.ProjectPath,
		SignatureProtected: pushRules.RejectUnsignedCommits,
	}
	return repoCommitProtection

}

// CheckCommitSignature Checks if a commit has a signature
func (a *Api) CheckCommitSignature(sha string) (bool, string) {
	// For unsigned commits we get a 404 response
	sig, _, _ := a.Client.Commits.GetGPGSiganature(a.ProjectPath, sha)
	if sig != nil {
		return true, sig.VerificationStatus
	}
	return false, ""
}

func (a *Api) GetAutomationKeys() ([]AutomationKey, error) {

	opts := &gitlab.ListProjectDeployKeysOptions{PerPage: 20}
	keys, response, err := a.Client.DeployKeys.ListProjectDeployKeys(a.ProjectPath, opts)
	if err != nil {
		fmt.Printf("Error retrieving authomation keys. Error: %s, Response Status: %s", err.Error(), response.Status)
		return nil, err
	}

	var automationKeys []AutomationKey
	for _, key := range keys {
		automationKeys = append(automationKeys,
			AutomationKey{
				ID:           key.ID,
				Title:        key.Title,
				ReadOnly:     !*key.CanPush,
				CreationDate: key.CreatedAt,
			},
		)
	}

	return automationKeys, nil
}
