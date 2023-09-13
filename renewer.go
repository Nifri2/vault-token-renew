package main

import (
	vault "github.com/hashicorp/vault/api"
	"github.com/rs/zerolog/log"
)

func (token Token) renewToken() error {

	log.Info().Msgf("Starting Cycle - Renewing Token %s", token.Name)

	// Create a default config
	config := vault.DefaultConfig()

	// Set vault address and create a client

	config.Address = token.VaultURL
	client, err := vault.NewClient(config)
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating vault client")
		return err
	}
	defer client.ClearToken()

	// Set token and create a token auth
	client.SetToken(token.Token)
	token_auth := client.Auth().Token()

	// Lookup self to verify that the token is valid and authentication is possible
	token_lookup, err := token_auth.LookupSelf()
	if err != nil {
		log.Err(err).Msg("Error looking up token")
		return err
	}
	log.Info().Msgf("Successfully authenticated to vault %s", token.VaultURL)
	log.Info().Msgf("Sucessfully looked up token")

	// Get current Token TTL and print it for debugging
	token_ttl, err := token_lookup.TokenTTL()
	if err != nil {
		log.Fatal().Err(err).Msg("Error getting token TTL")
		return err
	}
	log.Info().Msgf("Token TTL %s", token_ttl)

	// Renew Token
	renew, err := token_auth.RenewSelf(10)
	if err != nil {
		log.Fatal().Err(err).Msg("Error renewing token")
		return err
	}
	log.Info().Msgf("Successfully Renewed token")

	// Get new Token TTL and print it for debugging
	new_token_ttl, err := renew.TokenTTL()
	if err != nil {
		log.Fatal().Err(err).Msg("Error getting new token TTL")
		return err
	}
	log.Info().Msgf("New token TTL %s", new_token_ttl)

	// Clear Token - Say goodbye
	log.Info().Msg("Ending Cycle - Clearing Token")
	return nil
}
