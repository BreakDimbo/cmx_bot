package config

import (
	cons "bot/const"
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
	HBotInfo     MastodonClientInfo `toml:"hbot"`
	WikiInfo     MastodonClientInfo `toml:"wikibot"`

	ScriptFile string `toml:"script_file"`

	ActorA MastodonClientInfo `toml:"actor_a"`
	ActorB MastodonClientInfo `toml:"actor_b"`
	ActorC MastodonClientInfo `toml:"actor_c"`
	ActorD MastodonClientInfo `toml:"actor_d"`
	ActorE MastodonClientInfo `toml:"actor_e"`
	ActorF MastodonClientInfo `toml:"actor_f"`
	ActorG MastodonClientInfo `toml:"actor_g"`
	ActorH MastodonClientInfo `toml:"actor_h"`
	ActorI MastodonClientInfo `toml:"actor_i"`
	ActorJ MastodonClientInfo `toml:"actor_j"`
	ActorK MastodonClientInfo `toml:"actor_k"`
	ActorL MastodonClientInfo `toml:"actor_l"`
	ActorM MastodonClientInfo `toml:"actor_m"`
}

type elastic struct {
	Url      string
	Username string
	Password string
}

type postConifg struct {
	DailyTime       string `toml:"daily_cron_time"`
	WeeklyTime      string `toml:"weekly_cron_time"`
	MonthlyTime     string `toml:"monthly_cron_time"`
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

func HBotClientInfo() MastodonClientInfo {
	return config.HBotInfo
}

func WikiBotClientInfo() MastodonClientInfo {
	return config.WikiInfo
}

func ActorBotClientInfo(name string) (MastodonClientInfo, error) {
	switch name {
	case cons.Okabe:
		return config.ActorA, nil
	case cons.Mayuri:
		return config.ActorB, nil
	case cons.Itaru:
		return config.ActorC, nil
	case cons.Kurisu:
		return config.ActorD, nil
	case cons.Moeka:
		return config.ActorE, nil
	case cons.Ruka:
		return config.ActorF, nil
	case cons.NyanNyan:
		return config.ActorG, nil
	case cons.Suzuha:
		return config.ActorH, nil
	case cons.Maho:
		return config.ActorI, nil
	case cons.Kagari:
		return config.ActorJ, nil
	case cons.Yuki:
		return config.ActorK, nil
	case cons.Tennouji:
		return config.ActorL, nil
	case cons.Nae:
		return config.ActorM, nil
	default:
		return MastodonClientInfo{}, fmt.Errorf("no such actor %s", name)
	}
}

func ScriptFilePath() string {
	return config.ScriptFile
}
