package gitlab

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/xanzy/go-gitlab"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const projectPath = "myorg/myrepo"

type MockAPIProcessor struct {
	mock.Mock
}

func (m *MockAPIProcessor) GetAutomationKeys(projectPath string) ([]AutomationKey, error) {
	panic("implement me")
}

func (m *MockAPIProcessor) GetChangesToCiCd(path string, projectPath string, since time.Time) ([]CommitInfo, error) {
	panic("implement me")
}

func (m *MockAPIProcessor) GetProjectSignatureProtection(projectPath string) RepoCommitProtection {
	panic("implement me")
}

func (m *MockAPIProcessor) CheckCommitSignature(projectPath string, sha string) (bool, string) {
	args := m.Called(projectPath, sha)
	return args.Bool(0), args.String(1)
}

func (m *MockAPIProcessor) GetCommitsInfo(projectPath string, repositoryCommits []*gitlab.Commit) []CommitInfo {
	args := m.Called(projectPath, repositoryCommits)
	result := args.Get(0)
	return result.([]CommitInfo)
}

func TestGetChangesToCiCd(t *testing.T) {
	// Create a fake server
	mux, server, client := setup(t)
	defer teardown(server)

	// Mock the endpoint
	prefix := client.BaseURL().Path
	endpoint := fmt.Sprintf("%vprojects/%v/repository/commits", prefix, projectPath)
	mux.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `[]`)
	})
	// Mock functions inside the function we're testing
	mockObj := &MockAPIProcessor{}
	p := &Api{
		Client: client,
		Repo:   mockObj,
	}

	// Return an empty array because we want only to test if the function was called
	var commits []*gitlab.Commit
	mockObj.On("GetCommitsInfo", projectPath, commits).Return([]CommitInfo{})

	_, err := p.GetChangesToCiCd(".github/workflow.yaml", projectPath, time.Time{})

	assert.NoError(t, err)

	mockObj.AssertNumberOfCalls(t, "GetCommitsInfo", 1)
}

func TestGetCommitsInfo(t *testing.T) {
	client, _ := gitlab.NewClient("")

	// Mock functions inside the function we're testing
	commitDate := time.Date(2021, time.Month(9), 20, 11, 50, 22, 0, time.FixedZone("", 10800))
	commit := &gitlab.Commit{
		ID:           "ed899a2f4b50b4370feeea94676502b42383c746",
		AuthorName:   "Swaggy Baggins",
		AuthorEmail:  "example@example.com",
		AuthoredDate: &commitDate,
		WebURL:       "https://gitlab.example.com/thedude/gitlab-foss/-/commit/ed899a2f4b50b4370feeea94676502b42383c746",
	}
	mockObj := &MockAPIProcessor{}
	p := &Api{
		Client: client,
		Repo:   mockObj,
	}
	mockObj.On("CheckCommitSignature", projectPath, commit.ID).Return(true, "verified")

	// Call the function we want to test
	resp := p.GetCommitsInfo(projectPath, []*gitlab.Commit{commit})

	mockObj.AssertNumberOfCalls(t, "CheckCommitSignature", 1)

	want := CommitInfo{
		Repo:               projectPath,
		CommitUrl:          commit.WebURL,
		Date:               commit.AuthoredDate,
		AuthorName:         commit.AuthorName,
		AuthorEmail:        commit.AuthorEmail,
		VerifiedSignature:  true,
		VerificationReason: "verified",
	}
	assertResult(t, []CommitInfo{want}, resp)
}

func TestCheckCommitSignature(t *testing.T) {
	// Create a fake server
	mux, server, client := setup(t)
	defer teardown(server)

	sha := "123abc"

	// Mock the endpoint
	prefix := client.BaseURL().Path
	endpoint := fmt.Sprintf("%vprojects/%v/repository/commits/%v/signature", prefix, projectPath, sha)
	mux.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{
          "verification_status": "verified",
          "gpg_key_id": 1,
          "gpg_key_primary_keyid": "8254AAB3FBD54AC9",
          "gpg_key_user_name": "John Doe",
          "gpg_key_user_email": "johndoe@example.com"
        }`)
	})

	// Call the function we want to test
	p := &Api{
		Client: client,
	}
	isVerified, reason := p.CheckCommitSignature(projectPath, sha)
	assertResult(t, isVerified, true)
	assertResult(t, reason, "verified")

}

func TestGetAutomationKeys(t *testing.T) {
	// Create a fake server
	mux, server, client := setup(t)
	defer teardown(server)

	// Mock the endpoint
	prefix := client.BaseURL().Path
	endpoint := fmt.Sprintf("%vprojects/%v/deploy_keys", prefix, projectPath)
	mux.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `[
          {
            "id": 1,
            "title": "Public key",
            "key": "ssh-rsa AAAAB3NzaC1yc2EAAAABJQAAAIEAiPWx6WM4lhHNedGfBpPJNPpZ7yKu+dnn1SJejgt4596k6YjzGGphH2TUxwKzxcKDKKezwkpfnxPkSMkuEspGRt/aZZ9wa++Oi7Qkr8prgHc4soW6NUlfDzpvZK2H5E7eQaSeP3SAwGmQKUFHCddNaP0L+hM7zhFNzjFvpaMgJw0=",
            "created_at": "2021-09-20T11:50:22Z",
            "can_push": false
          }
        ]`)
	})

	// Call the function we want to test
	p := &Api{
		Client: client,
	}
	resp, _ := p.GetAutomationKeys(projectPath)

	createdAt := time.Date(2021, time.Month(9), 20, 11, 50, 22, 0, time.UTC)
	want := AutomationKey{
		ID:           1,
		Title:        "Public key",
		ReadOnly:     true,
		CreationDate: &createdAt,
	}
	assertResult(t, []AutomationKey{want}, resp)
}

// Assert for more complicated data types
func assertResult(t *testing.T, want interface{}, got interface{}) {
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Mismatch (-want +got):\n%s", diff)
	}
}

func setup(t *testing.T) (*http.ServeMux, *httptest.Server, *gitlab.Client) {
	// mux is the HTTP request multiplexer used with the test server.
	mux := http.NewServeMux()

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(mux)

	// client is the Gitlab client being tested.
	client, err := gitlab.NewClient("", gitlab.WithBaseURL(server.URL))
	if err != nil {
		server.Close()
		t.Fatalf("Failed to create client: %v", err)
	}

	return mux, server, client
}

func teardown(server *httptest.Server) {
	server.Close()
}
