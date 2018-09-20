package bot

import (
	"bot/bredis"
	cons "bot/const"
	gomastodon "bot/go-mastodon"
	"bot/log"
	"fmt"
	"strings"
	"time"
)

const (
	LoveYouKey     = "LoveKurisu"
	LoveYouTimeout = 6 * time.Hour
)

func BlockHandler(self *Actor, ntf *gomastodon.Notification, data interface{}) {
	content := filter(ntf.Status.Content)
	log.SLogger.Infof("get notification: %s", content)

	if strings.Contains(content, "EL_PSY_CONGROO") {
		actors, ok := data.(map[string]*Actor)
		if !ok {
			log.SLogger.Errorf("convert data %v to map error", data)
			return
		}

		for _, actor := range actors {
			if actor.Name == self.Name {
				continue
			}
			actor.BlockCh <- string(ntf.Account.ID)
			log.SLogger.Infof("start to block %s", ntf.Account.Username)
		}
	}
}

func UnblockHandler(self *Actor, ntf *gomastodon.Notification, data interface{}) {
	content := filter(ntf.Status.Content)
	log.SLogger.Infof("get notification: %s", content)

	if strings.Contains(content, "Love_You") {
		actors, ok := data.(map[string]*Actor)
		if !ok {
			log.SLogger.Errorf("convert data %v to map error", data)
			return
		}

		for _, actor := range actors {
			if actor.Name == self.Name {
				continue
			}
			actor.UnBlockCh <- string(ntf.Account.ID)
			log.SLogger.Infof("start to unblock %s", ntf.Account.Username)
		}
	}
}

func LoveHandler(self *Actor, ntf *gomastodon.Notification, data interface{}) {
	content := filter(ntf.Status.Content)
	log.SLogger.Infof("get notification: %s", content)

	// if the toot is on public and is love related then will reply he(she) on public line
	if isLoveYou(content) && ntf.Status.Visibility == "public" {
		key := fmt.Sprintf("%s:%s", LoveYouKey, ntf.Account.Username)
		// if loved already, toot hentai and return
		if isLoved(key) {
			toot := fmt.Sprintf("@%s %s", ntf.Account.Username, "够了！变态！")
			_, err := self.client.Post(toot)
			if err != nil {
				log.SLogger.Errorf("kurisu reply to error %v", err)
			}
			return
		}

		// set userID with love timeout in redis
		err := bredis.Client.Set(key, ntf.Account.Username, LoveYouTimeout).Err()
		if err != nil {
			log.SLogger.Errorf("set key to redis error: %v", err)
		}
		reply := GetRandomReply(cons.Kurisu)
		toot := fmt.Sprintf("@%s %s", ntf.Account.Username, reply)
		_, err = self.client.Post(toot)
		if err != nil {
			log.SLogger.Errorf("kurisu reply to error %v", err)
		}
	}
}

func FoodHandler(self *Actor, ntf *gomastodon.Notification, data interface{}) {
	content := filter(ntf.Status.Content)
	log.SLogger.Infof("get notification: %s", content)

	if strings.Contains(content, "#菜谱") {
		// keep diet in redis
		i := strings.Index(content, "#菜谱")
		food := content[i+7:]
		key := fmt.Sprintf("%s:%s", cons.FoodKey, food)
		err := bredis.Client.Set(key, "true", 1024*24*time.Hour).Err()
		if err != nil {
			log.SLogger.Errorf("save %s to redis error: %v", key, err)
			return
		}
		script := fmt.Sprintf("诶嘿嘿，%s 怎么样？", food)
		AddReply(cons.EatSome, script)

		toot := fmt.Sprintf("@%s %s", ntf.Account.Username, "乙！")
		_, err = self.client.Post(toot)
		if err != nil {
			log.SLogger.Errorf("itaru reply to %s error %v", ntf.Account.Username, err)
		}
	} else if strings.Contains(content, "吃啥") || strings.Contains(content, "吃点啥") ||
		strings.Contains(content, "吃什么") {
		reply := GetRandomReply(cons.EatSome)
		toot := fmt.Sprintf("@%s %s", ntf.Account.Username, reply)
		_, err := self.client.Post(toot)
		if err != nil {
			log.SLogger.Errorf("itaru reply to %s error %v", ntf.Account.Username, err)
		}

	} else if strings.Contains(content, "桶子") && ntf.Status.Visibility == "public" {
		reply := GetRandomReply(cons.Itaru)
		toot := fmt.Sprintf("@%s %s", ntf.Account.Username, reply)
		_, err := self.client.Post(toot)
		if err != nil {
			log.SLogger.Errorf("itaru reply to %s error %v", ntf.Account.Username, err)
		}
	}
}
