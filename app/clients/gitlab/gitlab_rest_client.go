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

type Api struct {
	Client      *gitlab.Client
	ProjectPath string
	Repo
}

func NewApi(token string, cfg *config.Config, url ...string) *Api {
	var client *gitlab.Client
	if url != nil {
		// Get a client for a specific gitlab server
		client, _ = gitlab.NewClient(token, gitlab.WithBaseURL(url[0]))
	} else {
		client, _ = gitlab.NewClient(token)
	}

	p := &Api{
		Client:      client,
		ProjectPath: fmt.Sprintf("%s/%s", cfg.Project.Owner, cfg.Project.Repo),
	}
	p.Repo = p
	return p
}

// GetChangesToCiCd Control-1
// Returns commits for a specific item since a specific date
func (api *Api) GetChangesToCiCd(path string, since time.Time) ([]CommitInfo, error) {

	opt := &gitlab.ListCommitsOptions{
		Path:        &path,
		Since:       &since,
		ListOptions: gitlab.ListOptions{PerPage: 20},
	}

	// get all pages of results
	var allCommits []*gitlab.Commit
	for {
		commits, resp, err := api.Client.Commits.ListCommits(
			api.ProjectPath,
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

	return api.Repo.GetCommitsInfo(allCommits), nil
}

func (api *Api) GetCommitsInfo(repositoryCommits []*gitlab.Commit) []CommitInfo {
	var commitsInfo []CommitInfo
	for _, repoCommit := range repositoryCommits {
		isVerified, reason := api.Repo.CheckCommitSignature(repoCommit.ID)

		commitsInfo = append(commitsInfo,
			CommitInfo{
				Repo:               api.ProjectPath,
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
func (api *Api) GetProjectSignatureProtection() RepoCommitProtection {

	pushRules, _, _ := api.Client.Projects.GetProjectPushRules(api.ProjectPath)

	repoCommitProtection := RepoCommitProtection{
		Repo:               api.ProjectPath,
		SignatureProtected: pushRules.RejectUnsignedCommits,
	}
	return repoCommitProtection
}

// CheckCommitSignature Checks if a commit has a signature
func (api *Api) CheckCommitSignature(sha string) (bool, string) {
	// For unsigned commits we get api 404 response
	sig, _, _ := api.Client.Commits.GetGPGSiganature(api.ProjectPath, sha)
	if sig != nil {
		return true, sig.VerificationStatus
	}
	return false, ""
}

func (api *Api) GetAutomationKeys() ([]AutomationKey, error) {

	opts := &gitlab.ListProjectDeployKeysOptions{PerPage: 20}
	keys, response, err := api.Client.DeployKeys.ListProjectDeployKeys(api.ProjectPath, opts)
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
