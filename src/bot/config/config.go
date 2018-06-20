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

var config TomlConfig

type TomlConfig struct {
	Title        string
	ENV          string
	Ela          elastic            `toml:"elastic"`
	PostConfig   postConifg         `toml:"post_config"`
	IntelBotInfo MastodonClientInfo `toml:"intelbot"`
	FireBotInfo  MastodonClientInfo `toml:"firebot"`
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

type MastodonClientInfo struct {
	ID       string `toml:"client_id"`
	Secret   string `toml:"client_secret"`
	Sever    string `toml:"server"`
	Email    string `toml:"client_email"`
	Password string `toml:"client_password"`
	// differentiate the production and development evn
	Timezone time.Duration `toml:"timezone"`
	FbotName string        `toml:"fbot"`
}

func init() {
	runingEnv := flag.String("env", "development", "running env")
	rootedPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
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
	config.ENV = *runingEnv
}

func GetElastic() elastic {
	return config.Ela
}

func GetPostConfig() postConifg {
	return config.PostConfig
}

func GetRuntimeEnv() string {
	return config.ENV
}

func IntelBotClientInfo() MastodonClientInfo {
	return config.IntelBotInfo
}

func FireBotClientInfo() MastodonClientInfo {
	return config.FireBotInfo
}
