// Copyright 2022 Google LLC All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/Aleromerog/artifact-registry-go-tools/pkg/auth"
	"github.com/Aleromerog/artifact-registry-go-tools/pkg/npmrc"
)

const help = `
Update your user .npmrc file to work with Google Cloud Artifact Registry Go Repositories.

Commands:

* refresh, to refresh oauth tokens for Artifact Registry Go endpoints.`

func main() {
	if len(os.Args) < 2 {
		fmt.Println(help)
		return
	}
	switch os.Args[1] {
	case "refresh":
		refreshFlags := flag.NewFlagSet("refresh", flag.ExitOnError)
		var (
			token = refreshFlags.String("token", "", "The oauth token to write to the npmrc file. Most users should not set this field and let the tool find the credentials to use from the environment.")
		)
		refreshFlags.Parse(os.Args[2:])
		refresh(*token)
	case "help", "-help", "--help":
		fmt.Println(help)
	default:
		fmt.Printf("unknown command %q. Please rerun the tool with `--help`\n", os.Args[1])
	}
}

func refresh(token string) {
	ctx := context.Background()
	ctx, cf := context.WithTimeout(ctx, 30*time.Second)
	defer cf()

	// Load per-project npmrc file
	projectConfig, err := npmrc.Load(".npmrc")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if token == "" {
		var err error
		token, err = auth.Token(ctx)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}

	// Add token to npmrc file

	config, err := npmrc.AddTokenToConfigFile(projectConfig, token)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// Getting per-user config file.
	h, err := os.UserHomeDir()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// Save updated npmrc file to per-user config file.
	if err := npmrc.Save(config, path.Join(h, ".npmrc")); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println("Refresh completed.")
}
