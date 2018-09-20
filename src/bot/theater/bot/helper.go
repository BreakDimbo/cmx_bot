package bot

import (
	"bot/bredis"
	cons "bot/const"
	gomastodon "bot/go-mastodon"
	"bot/log"
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
)

const (
	LoveYouKey     = "LoveKurisu"
	LoveYouTimeout = 6 * time.Hour
)

func filter(raw string) (polished string) {
	p := bluemonday.StrictPolicy()
	polished = p.Sanitize(raw)
	polished = strings.Replace(polished, "@rintarou", "", -1)
	polished = html.UnescapeString(polished)
	return
}

func isLoveYou(content string) bool {
	return strings.Contains(content, "Love_You") || strings.Contains(content, "love you") ||
		strings.Contains(content, "Love You") || strings.Contains(content, "爱你") || strings.Contains(content, "喜欢你") ||
		strings.Contains(content, "吃了你") || strings.Contains(content, "好き") || strings.Contains(content, "吃掉你") ||
		strings.Contains(content, "梦到你") || strings.Contains(content, "吻你") || strings.Contains(content, "可爱")
}

func isLoved(key string) bool {
	res, err := bredis.Client.Get(key).Result()
	if err == nil && res != "" {
		return true
	}
	return false
}

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
		key := fmt.Sprintf("%s:%s", FoodKey, food)
		err := bredis.Client.Set(key, "true", 1024*24*time.Hour).Err()
		if err != nil {
			log.SLogger.Errorf("save %s to redis error: %v", key, err)
			return
		}
		script := fmt.Sprintf("诶嘿嘿，%s 怎么样？", food)
		AddReply(cons.Itaru, script)

		toot := fmt.Sprintf("@%s %s", ntf.Account.Username, "乙！")
		_, err = self.client.Post(toot)
		if err != nil {
			log.SLogger.Errorf("kurisu reply to error %v", err)
		}
	} else if strings.Contains(content, "桶子") && ntf.Status.Visibility == "public" {
		reply := GetRandomReply(cons.Itaru)
		toot := fmt.Sprintf("@%s %s", ntf.Account.Username, reply)
		_, err := self.client.Post(toot)
		if err != nil {
			log.SLogger.Errorf("kurisu reply to error %v", err)
		}
	}
}
