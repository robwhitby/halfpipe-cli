package project

import (
	"flag"
	"os"
	"path/filepath"
)

type Config struct {
	VaultToken string
	RepoRoot   string
}

func NewConfig() (*Config, error) {
	vaultToken := flag.String("vaultToken", "", "Vault token")
	path := flag.String("path", "", "Path to repo to lint, if not specified pwd will be linted")

	flag.Parse()
	p, err := getPath(*path)
	if err != nil {
		return &Config{}, err
	}

	return &Config{
		VaultToken: *vaultToken,
		RepoRoot:   p,
	}, nil
}

func getPath(path string) (string, error) {
	if path == "." || path == "" {
		dir, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return dir, nil
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return path, nil
}
