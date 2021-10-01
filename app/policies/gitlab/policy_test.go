package gitlab

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/mock"
	x "github.com/xanzy/go-gitlab"
	"secure-pipeline-poc/app/clients/gitlab"
	"secure-pipeline-poc/app/config"
	"secure-pipeline-poc/app/policies/common"
	"testing"
	"time"
)

type mockobj struct {
	mock.Mock
}

//return args.Get(0).(*MyObject), args.Get(1).(*AnotherObjectOfMine)

func (m *mockobj) GetChangesToCiCd(client *x.Client, projectPath string, path string, since time.Time) ([]gitlab.CommitInfo, error) {
	args := m.Called(client, projectPath, path, since)
	result := args.Get(0)
	err := args.Get(1)
	return result.([]gitlab.CommitInfo), err.(error)
}

func (m *mockobj) userAuthPolicy() *common.Policy {
	args := m.Called()
	result := args.Get(0)
	return result.(*common.Policy)
}

func TestValidateC1(t *testing.T) {
	mockObj := &mockobj{}
	client, _ := x.NewClient("")
	conf := config.Config{}
	projectPath := "org/project"
	filePath := ".travis.yaml"
	sinceDate := time.Date(2021, time.Month(9), 20, 11, 50, 22, 0, time.FixedZone("", 10800))

	commit := gitlab.CommitInfo{
		Repo:               "DUPA!",
		CommitUrl:          "",
		Date:               &sinceDate,
		AuthorName:         "",
		AuthorEmail:        "",
		VerifiedSignature:  true,
		VerificationReason: "",
	}

	policy := common.Policy{}
	mockObj.On("userAuthPolicy").Return(commit)
	mockObj.On("GetChangesToCiCd", client, projectPath, filePath, sinceDate).Return(commit)
	ValidateC1(client, &conf, projectPath, sinceDate)
	fmt.Printf("%v", commit)

}

// Assert for more complicated data types
func assertResult(t *testing.T, want interface{}, got interface{}) {
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Mismatch (-want +got):\n%s", diff)
	}
}
