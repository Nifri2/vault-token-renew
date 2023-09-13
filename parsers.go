package main

import (
	"flag"
	"io/ioutil"

	"github.com/caarlos0/env"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type TokenList struct {
	Tokens []Token `yaml:"tokens"`
}

type Token struct {
	VaultURL string `yaml:"vault_url"`
	Token    string `yaml:"token"`
	Name     string `yaml:"name"`
}

// Reaf from the Kubernetes Secret
type Environment struct {
	GitURL          string `env:"GIT_URL"`
	GitToken        string `env:"GIT_TOKEN"`
	VaultURL        string `env:"VAULT_URL"`         // for Transit / Sops Decryption
	VaultToken      string `env:"VAULT_TOKEN"`       // for Transit / SOPS Decryption
	VaultTransitKey string `env:"VAULT_TRANSIT_KEY"` // for Transit
}

func setup(config string) (*string, Environment) {
	// Load environment variables
	var environment Environment = Environment{}

	err := env.Parse(&environment)
	if err != nil {
		log.Fatal().Err(err).Msg("Error parsing environment variables - did you provide a kube secret?")
	}
	log.Info().Msg("Environment Loaded successfully")

	// Clone git repo
	dir, err := cloneGitRepo(environment)
	if err != nil {
		log.Fatal().Err(err).Msg("Error cloning git repo")
	}

	config = *dir + "/tokens.yaml"

	decrypt(*dir, config, environment)

	return dir, environment
}

func createTokenList(filepath string) (*TokenList, error) {

	// Genric YAML Parser
	// Parsed into struct defined above

	yfile, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal().Err(err).Msg("Error reading yaml file")
		return &TokenList{}, err
	}

	var data *TokenList

	err = yaml.Unmarshal(yfile, &data)
	if err != nil {
		log.Fatal().Err(err).Msg("Error unmarshalling yaml file")
		return &TokenList{}, err
	}
	return data, nil
}

func parseArgs() string {
	// in case theres a different structure in the future
	config := flag.String("config", "tokens.yaml", "Path to config file inside the git repo")
	flag.Parse()

	return *config

}
