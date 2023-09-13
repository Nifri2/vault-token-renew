package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
)

func decrypt(dir string, config string, environment Environment) error {
	// exectue a shell command to decrypt the token
	log.Info().Msgf("Decryption in directory: %s", dir)

	transit_key := fmt.Sprintf("%s%s", environment.VaultURL, environment.VaultTransitKey)
	log.Info().Msgf("Transit Key: %s", transit_key)

	// check if file in var config exists
	if _, err := os.Stat(config); errors.Is(err, os.ErrNotExist) {
		log.Fatal().Err(err).Msgf("File %s does not exist", config)
	}
	log.Info().Msgf("File %s exists, continue", config)

	// Decrypt the token, we use sops with a vault transit key
	log.Info().Msg("Decrypting token...")

	// The Module Documentation for sops is horrible so we have to use a exec.Command
	cmd := exec.Command("sops", "--verbose", "-d", "--hc-vault-transit", transit_key, "-i", config)
	cmd.Dir = dir
	err := cmd.Run()

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to decrypt token")
	}

	log.Info().Msg("Successfully decrypted token")

	return nil
}
