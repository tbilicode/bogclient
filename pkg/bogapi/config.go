package bogapi

import "github.com/effective-security/x/configloader"

type Config struct {
	Accounts     []Account `json:"accounts" yaml:"accounts"`
	ClientID     string    `json:"client_id" yaml:"client_id"`
	ClientSecret string    `json:"client_secret" yaml:"client_secret"`
	AuthURL      string    `json:"auth_url" yaml:"auth_url"`
	ApiHost      string    `json:"api_host" yaml:"api_host"`
}

type Account struct {
	ID       string   `json:"id" yaml:"id"`
	Name     string   `json:"name" yaml:"name"`
	Currency []string `json:"currency" yaml:"currency"`
}

func LoadConfig(file string) (*Config, error) {
	cfg := new(Config)
	if file == "" {
		return cfg, nil
	}

	err := configloader.UnmarshalAndExpand(file, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
