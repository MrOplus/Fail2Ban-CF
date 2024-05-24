package main

import (
	"context"
	"github.com/cloudflare/cloudflare-go"
	"github.com/seculize/islazy/log"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"path/filepath"
)

var api *cloudflare.API

func makeAPIClient() {
	var err error
	api, err = cloudflare.NewWithAPIToken(config.APIKey)
	if err != nil {
		log.Fatal("Error creating API client: ", err)
	}
	log.Info("API client created successfully")
}
func parseConfig() {
	config = &Config{}
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	yamlFile, err := os.ReadFile(path.Join(exPath, "config.yaml"))
	if err != nil {
		log.Fatal("Error reading config file: ", err)
	}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		log.Fatal("Error parsing config file: ", err)
	}
}
func main() {
	parseConfig()
	makeAPIClient()
	args := os.Args
	if len(args) < 3 {
		log.Fatal("Usage: %s action ip <comment>", args[0])
	}
	ctx := context.Background()
	action := args[1]
	ip := args[2]
	comment := ""
	if action == "add" {
		if len(args) < 4 {
			log.Fatal("Usage: %s add ip comment", args[0])
		} else {
			comment = args[3]
		}
	}
	log.Info("Action: %s, IP: %s, Comment: %s", action, ip, comment)

	log.Info("Checking for existing list")
	list, err := getList(ctx, api)
	if err != nil {
		log.Warning("List not found, creating new list")
		list, err = createList(ctx, api)
		if err != nil {
			log.Fatal("Error creating list: ", err)
		}
	}
	log.Info("List found/created successfully, %v", list.ID)
	if action == "add" {
		log.Info("Adding IP to list")
		_, err = addIP(ctx, api, list, ip, comment)
		if err != nil {
			log.Warning("Error adding IP to list: ", err)
		}
		log.Info("IP added successfully")
	} else if action == "delete" {
		log.Info("Finding IP in list")
		item, err := findIp(ctx, api, list, ip)
		if err != nil {
			log.Fatal("Error finding IP in list: ", err)
		}
		log.Info("IP found: %v", item.ID)
		log.Info("Deleting IP from list")
		err = deleteIp(ctx, api, list, item)
		if err != nil {
			log.Warning("Error deleting IP from list: ", err)
		}
	}
}
