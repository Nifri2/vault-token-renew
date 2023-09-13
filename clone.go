package main

import (
	"os"

	"github.com/rs/zerolog/log"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func cloneGitRepo(env Environment) (*string, error) {

	// Create a temp directory
	dir, err := os.MkdirTemp("", "tokens")
	if err != nil {
		log.Err(err).Msg("Error creating temp directory")
		os.RemoveAll(dir)
		return nil, err
	}

	log.Info().Msgf("Cloning git repo %s", env.GitURL)

	// CloneOptions - we have to use InsecureSkipTLS because of self signed certs / internal CAs
	options := &git.CloneOptions{
		URL:             env.GitURL,
		InsecureSkipTLS: true,
		Auth: &http.BasicAuth{
			Username: "refresher",  // Doesnt matter
			Password: env.GitToken, // Token Here
		},
	}

	// Clone the repo into the temp directory
	_, err = git.PlainClone(dir, false, options)
	if err != nil {
		log.Err(err).Msg("Error cloning git repo")
		os.RemoveAll(dir)
		return nil, err
	}

	return &dir, nil
}
