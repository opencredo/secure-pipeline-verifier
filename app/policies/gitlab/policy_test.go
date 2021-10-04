package gitlab

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"secure-pipeline-poc/app/clients/gitlab"
	"secure-pipeline-poc/app/config"
	"secure-pipeline-poc/app/notification"
	"testing"
	"time"
)

func setup() (*http.ServeMux, *httptest.Server) {
	// mux is the HTTP request multiplexer used with the test server.
	mux := http.NewServeMux()

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(mux)

	// Mock slack server
	notification.APIURL = server.URL + "/"

	return mux, server
}

func TestControl1(t *testing.T) {
	mux, server := setup()
	defer teardown(server)

	cfg := &config.Config{
		Project: config.Project{
			Platform: "gitlab",
			Owner:    "myorg",
			Repo:     "myrepo",
		},
		RepoInfoChecks: config.RepoInfoChecks{
			TrustedDataFile:   "./test_data/gitlab-secure-pipeline-example-data.json",
			CiCdPath:          ".travis.yaml",
			ProtectedBranches: []string{"main", "develop"},
		},
	}
	api := gitlab.NewApi("", cfg, server.URL)

	// Mock all responses from the gitlab server.
	prefix := api.Client.BaseURL().Path
	projectPath := "myorg/myrepo"
	endpoint := fmt.Sprintf("%vprojects/%v/repository/commits", prefix, projectPath)
	mux.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `[
		  {
		    "id": "ed899a2f4b50b4370feeea94676502b42383c746",
		    "short_id": "ed899a2f4b5",
		    "title": "Replace sanitize with escape once",
		    "author_name": "Example User",
		    "author_email": "user@example.com",
		    "authored_date": "2021-09-20T11:50:22+03:00",
		    "committer_name": "Administrator",
		    "committer_email": "admin@example.com",
		    "committed_date": "2021-09-20T11:50:22+03:00",
		    "created_at": "2021-09-20T11:50:22+03:00",
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
		    "created_at": "2021-09-20T09:06:12+03:00",
		    "message": "Sanitize for network graph",
		    "parent_ids": [
		      "ae1d9fb46aa2b07ee9836d49862ec4e2c46fbbba"
		    ],
		    "web_url": "https://gitlab.example.com/thedude/gitlab-foss/-/commit/ed899a2f4b50b4370feeea94676502b42383c746"
		  }
		]`)
	})

	endpoint = fmt.Sprintf(
		"%vprojects/%v/repository/commits/ed899a2f4b50b4370feeea94676502b42383c746/signature", prefix, projectPath)
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
		"%vprojects/%v/repository/commits/6104942438c14ec7bd21c6cd5bd995272b3faff6/signature", prefix, projectPath)
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
	// Mock Slack notification
	endpoint = "/chat.postMessage"
	mux.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{
		    "ok": true,
		    "channel": "C1H9RESGL",
		    "ts": "1503435956.000247",
		    "message": {
		        "text": "Here's a message for you",
		        "username": "ecto1",
		        "bot_id": "B19LU7CSY",
		        "attachments": [
		            {
		                "text": "This is an attachment",
		                "id": 1,
		                "fallback": "This is an attachment's fallback"
		            }
		        ],
		        "type": "message",
		        "subtype": "bot_message",
		        "ts": "1503435956.000247"
		    }
		}`)
	})

	sinceDate := time.Date(2021, time.Month(9), 20, 11, 50, 22, 0, time.FixedZone("", 10800))
	ValidateC1(api, cfg, sinceDate)
}

func teardown(server *httptest.Server) {
	server.Close()
}
