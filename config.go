package main

type Config struct {
	APIKey    string `yaml:"api_key"`
	AccountId string `yaml:"account_id"`
}

var config *Config
