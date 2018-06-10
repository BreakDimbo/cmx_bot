package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

var config tomlConfig

type tomlConfig struct {
	Title       string
	Ela         elastic            `toml:"elastic"`
	PostCron    postCron           `toml:"post_crontab"`
	MClientInfo mastodonClientInfo `toml:"mastodon_client_info"`
}

type elastic struct {
	Url      string
	Username string
	Password string
}

type postCron struct {
	ConTime string `toml:"cron_time"`
}

type mastodonClientInfo struct {
	ID       string `toml:"client_id"`
	Secret   string `toml:"client_secret"`
	Sever    string `toml:"server"`
	Email    string `toml:"client_email"`
	Password string `toml:"client_password"`
}

func init() {
	runingEnv := flag.String("evn", "production", "running env")
	rootedPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	configPath := fmt.Sprintf("%s/config/%s.toml", rootedPath, *runingEnv)

	dat, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("read config file error: %s", err)
		return
	}

	if err := toml.Unmarshal(dat, &config); err != nil {
		log.Fatal(err)
	}
}

func GetElastic() elastic {
	return config.Ela
}

func GetPostCron() postCron {
	return config.PostCron
}

func GetMastodonClientInfo() mastodonClientInfo {
	return config.MClientInfo
}
