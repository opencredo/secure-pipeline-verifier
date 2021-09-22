package gitlab

import (
	"fmt"
	gitlabx "github.com/xanzy/go-gitlab"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)


func teardown(server *httptest.Server) {
	server.Close()
}

func TestGetChangesToCiCd(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	prefix := client.BaseURL().Path
	endpoint := fmt.Sprintf("%vprojects/myorg/myrepo/repository/commits", prefix)
	mux.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `[
		  {
		    "id": "ed899a2f4b50b4370feeea94676502b42383c746",
		    "short_id": "ed899a2f4b5",
		    "title": "Replace sanitize with escape once",
		    "author_name": "Example User",
		    "author_email": "user@example.com",
		    "authored_date": "2012-09-20T11:50:22+03:00",
		    "committer_name": "Administrator",
		    "committer_email": "admin@example.com",
		    "committed_date": "2012-09-20T11:50:22+03:00",
		    "created_at": "2012-09-20T11:50:22+03:00",
		    "message": "Replace sanitize with escape once",
		    "parent_ids": [
		      "6104942438c14ec7bd21c6cd5bd995272b3faff6"
		    ],
		    "web_url": "https://gitlab.example.com/thedude/gitlab-foss/-/commit/ed899a2f4b50b4370feeea94676502b42383c746"
		  },
		  {
		    "id": "6104942438c14ec7bd21c6cd5bd995272b3faff6",
		    "short_id": "6104942438c",
		    "title": "Sanitize for network graph",
		    "author_name": "randx",
		    "author_email": "user@example.com",
		    "committer_name": "ExampleName",
		    "committer_email": "user@example.com",
		    "created_at": "2012-09-20T09:06:12+03:00",
		    "message": "Sanitize for network graph",
		    "parent_ids": [
		      "ae1d9fb46aa2b07ee9836d49862ec4e2c46fbbba"
		    ],
		    "web_url": "https://gitlab.example.com/thedude/gitlab-foss/-/commit/ed899a2f4b50b4370feeea94676502b42383c746"
		  }
		]`)
	})

	endpoint = fmt.Sprintf(
		"%vprojects/myorg/myrepo/repository/commits/ed899a2f4b50b4370feeea94676502b42383c746/signature", prefix)
	mux.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{
		  "signature_type": "PGP",
		  "verification_status": "verified",
		  "gpg_key_id": 1,
		  "gpg_key_primary_keyid": "8254AAB3FBD54AC9",
		  "gpg_key_user_name": "John Doe",
		  "gpg_key_user_email": "johndoe@example.com",
		  "gpg_key_subkey_id": null,
		  "commit_source": "gitaly"
		}`)
	})

	endpoint = fmt.Sprintf(
		"%vprojects/myorg/myrepo/repository/commits/6104942438c14ec7bd21c6cd5bd995272b3faff6/signature", prefix)
	mux.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{
		  "signature_type": "PGP",
		  "verification_status": "unverified",
		  "gpg_key_id": 1,
		  "gpg_key_primary_keyid": "8254AAB3FBD54AC9",
		  "gpg_key_user_name": "John Doe",
		  "gpg_key_user_email": "johndoe@example.com",
		  "gpg_key_subkey_id": null,
		  "commit_source": "gitaly"
		}`)
	})

	resp, _ := GetChangesToCiCd(client, "myorg", "myrepo", ".github/workflow.yaml", time.Time{})

	fmt.Printf("%v", resp)

}


func setup(t *testing.T) (*http.ServeMux, *httptest.Server, *gitlabx.Client) {
	// mux is the HTTP request multiplexer used with the test server.
	mux := http.NewServeMux()

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(mux)

	// client is the Gitlab client being tested.
	client, err := gitlabx.NewClient("", gitlabx.WithBaseURL(server.URL))
	if err != nil {
		server.Close()
		t.Fatalf("Failed to create client: %v", err)
	}

	return mux, server, client
}

