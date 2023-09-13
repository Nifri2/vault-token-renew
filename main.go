package main

// Description: Programm that renews a Token read from an environment variable TOKEN
//              the token is updated by posting it to the vault api endpoint /auth/token/renew-self
//
// You need to providehe following environment variables:
// GIT_URL 				- URL to the git repo containing the token
// GIT_TOKEN 			- Token to access the git repo
// VAULT_URL 			- URL to the vault server
// VAULT_TOKEN 			- Token to access the vault server
// VAULT_TRANSIT_KEY 	- Transit Key to decrypt the token
//
// Provide these Via a Kubernetes Secret

import (
	"os"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	// set zerolog caller to short file name and line number
	// Looks cool and simplyfies debugging
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}

	// set zerolog logger to caller
	log.Logger = log.With().Caller().Logger()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Info().Msg("Starting vault token renewer")

	// for every token in the token list
	// renew the token

	config := parseArgs()
	log.Info().Msgf("Config file: %s", config)

	// Run Setup Routine
	dir, _ := setup(config)

	// Create new path based on temp directory
	config_path := *dir + "/tokens.yaml"

	log.Info().Msgf("Temp directory: %s", *dir)
	log.Info().Msgf("Config file Absolute path: %s", config_path)
	defer os.RemoveAll(*dir)

	// Create token list
	list, err := createTokenList(config_path)
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating token list")
	}

	// print lenght of token list
	log.Info().Msgf("Loaded %d tokens", len(list.Tokens))

	var failed []string

	for _, token := range list.Tokens {
		err := token.renewToken()
		if err != nil {
			log.Err(err).Msgf("Error renewing token %s in child function", token.Name)
			failed = append(failed, token.Name)
		}
	}

	// print failed tokens for debugging
	for _, token := range failed {
		log.Warn().Msgf("Failed to renew token %s", token)
	}

}
