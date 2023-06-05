package vault

import (
	"github.com/atotto/clipboard"
	"github.com/chewedfeed/automated/internal/config"
	"github.com/hashicorp/vault/api"
)

type Vault struct {
	Address string
	Token   string
}

func NewVault() (*Vault, error) {
	creds, err := config.GetCredentials()
	if err != nil {
		return nil, err
	}

	return &Vault{
		Address: creds.Vault.ServiceAddress,
		Token:   creds.Vault.Credentials.Token,
	}, nil
}

func (v *Vault) GetCreds() ([]string, error) {
	cfg := api.DefaultConfig()
	cfg.Address = v.Address
	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	client.SetToken(v.Token)
	data, err := client.Logical().List("database/config")
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	results := []string{}
	for _, va := range data.Data {
		for _, vv := range va.([]interface{}) {
			results = append(results, vv.(string))
		}
	}

	return results, nil
}

func (v *Vault) CreateDatabaseToken(credsPath, policyName string, read, write bool) (string, error) {
	cfg := api.DefaultConfig()
	cfg.Address = v.Address
	client, err := api.NewClient(cfg)
	if err != nil {
		return "", err
	}
	client.SetToken(v.Token)
	if err := client.Sys().PutPolicy(policyName, policyTemplate(credsPath, read, write)); err != nil {
		return "", err
	}

	secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
		Policies: []string{"default", policyName},
	})
	if err != nil {
		return "", err
	}
	if err := clipboard.WriteAll(secret.Auth.ClientToken); err != nil {
		return secret.Auth.ClientToken, err
	}

	return "Token in clipboard", nil
}

func policyTemplate(credsPath string, read, write bool) string {
	if read && write {
		return `path "database/creds/` + credsPath + `" {
      capabilities = ["read", "list", "create", "update", "delete"]
    }`
	}

	if read && !write {
		return `path "database/creds/` + credsPath + `" {
      capabilities = ["read", "list"]
    }`
	}

	if !read && write {
		return `path "database/creds/` + credsPath + `" {
      capabilities = ["create", "update", "delete"]
    }`
	}

	return ""
}
