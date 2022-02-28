package main

import (
	"fmt"
	"os"
	"github/secure-pipeline-verifier/app/config"
	"github/secure-pipeline-verifier/cmd"
	"time"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Usage:", os.Args[0], "path/to/config/", "YYYY-MM-ddTHH:mm:ss.SSSZ")
		return
	}

	var cfg config.Config
	config.LoadConfig(os.Args[1], &cfg)
	config.LoadTrustedDataToMap(os.Args[1], &cfg)

	sinceDate, err := time.Parse(time.RFC3339, os.Args[2])
	if err != nil {
		fmt.Println("Error " + err.Error() + " occurred while parsing date from " + os.Args[2])
		os.Exit(2)
	}

	branch := ""
	if len(os.Args) == 4 {
		branch = os.Args[3]
	}

	cmd.PerformCheck(&cfg, sinceDate, branch)
}
