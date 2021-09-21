package gitlab

import (
	"github.com/xanzy/go-gitlab"
)

func NewClient(token string) (*gitlab.Client, error) {
	return gitlab.NewClient(token)
}
