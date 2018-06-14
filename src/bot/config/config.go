package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

var config tomlConfig

type tomlConfig struct {
	Title           string
	Ela             elastic            `toml:"elastic"`
	PostConfig      postConifg         `toml:"post_config"`
	MClientInfo     mastodonClientInfo `toml:"mastodon_client_info"`
	FBotMClientInfo mastodonClientInfo `toml:"firebot"`
}

type elastic struct {
	Url      string
	Username string
	Password string
}

type postConifg struct {
	DailyTime       string `toml:"daily_cron_time"`
	WeeklyTime      string `toml:"weekly_cron_time"`
	CleanUnfollower string `toml:"clean_unfollowers_time"`
	Scope           string
}

type mastodonClientInfo struct {
	ID       string        `toml:"client_id"`
	Secret   string        `toml:"client_secret"`
	Sever    string        `toml:"server"`
	Email    string        `toml:"client_email"`
	Password string        `toml:"client_password"`
	Timezone time.Duration `toml:"timezone"`
	Fbot     string        `toml:"fbot"`
}

func init() {
	runingEnv := flag.String("env", "development", "running env")
	rootedPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	flag.Parse()
	fmt.Printf("Running in %s env\n", *runingEnv)
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

func GetPostConfig() postConifg {
	return config.PostConfig
}

func GetMastodonClientInfo() mastodonClientInfo {
	return config.MClientInfo
}

func GetFBotMClientInfo() mastodonClientInfo {
	return config.FBotMClientInfo
}
