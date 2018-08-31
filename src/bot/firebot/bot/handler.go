package bot

import (
	"bot/bredis"
	con "bot/firebot/const"
	gomastodon "bot/go-mastodon"
	"bot/log"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"os"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

const TtsDir = "/tmp/tts/"

func HandleNotification(e *gomastodon.NotificationEvent) {
	var tootToPost string
	notify := e.Notification

	switch notify.Type {
	case con.NotificationTypeMention:
		fromUser := notify.Account
		toot := notify.Status
		replyToID := notify.Status.InReplyToID
		firstContent := filter(toot.Content)
		isTTs := false
		var filename string

		if fromUser.Username == "xbot" || fromUser.Username == "zbot" {
			index := strings.Index(firstContent, "县民榜：\n")
			tootToPost = firstContent[index+12:]
			timeRange := firstContent[index-15 : index-9]
			lastIndex := strings.Index(tootToPost, "条\n")
			tootToPost = fmt.Sprintf("咳咳...注意！%s最活跃（话唠）县民是:%s。", timeRange, tootToPost[:lastIndex+3])
		} else if strings.Contains(firstContent, "#树洞") || strings.Contains(filter(toot.SpoilerText), "#树洞") {
			tootToPost = firstContent
		} else if strings.Contains(firstContent, "#话唠树洞") || strings.Contains(filter(toot.SpoilerText), "#话唠树洞") {
			isTTs = true
			filename, err = askForTTs(firstContent)
			if err != nil {
				return
			}
		} else {
			content := recurToot(replyToID)
			tootToPost = fmt.Sprintf("@%s:%s// %s", fromUser.Acct, firstContent, content)
			tootToPost = strings.TrimSuffix(tootToPost, "// ")
		}

		var id gomastodon.ID

		if isTTs {
			filepath := TtsDir + filename + ".wav"

			attachment, err := botClient.Normal.UploadMedia(context.Background(), filepath)
			if err != nil {
				log.SLogger.Errorf("upload media error: %v", err)
				return
			}

			toot := &gomastodon.Toot{
				MediaIDs: []gomastodon.ID{attachment.ID},
			}

			status, err := botClient.RawPost(toot)
			if err != nil {
				log.SLogger.Errorf("post to mastodon error: %v", err)
				return
			}

			id = status.ID

			err = os.Remove(filepath)
			if err != nil {
				log.SLogger.Error(err)
			}

		} else {
			id, _ = botClient.PostSensetiveWithPic(filter(toot.SpoilerText), tootToPost, toot.Sensitive, toot.MediaAttachments)
		}

		err := bredis.Client.Set(string(toot.ID), string(id), con.TootIDRedisTimeout).Err()
		if err != nil {
			log.SLogger.Errorf("set id: %s to redis error: %s", id, err)
		}

		log.SLogger.Infof("get toots from user id: %s, first toot: %s, tootToPost: %s\n", fromUser.Username, firstContent, tootToPost)
	}
}

func HandleDelete(e *gomastodon.DeleteEvent) {
	id := e.ID
	fbotTootID, err := bredis.Client.Get(string(id)).Result()
	if err != nil {
		// log.SLogger.Errorf("get toot: %s from redis error: %s", id, err)
		return
	}
	log.SLogger.Infof("start to delete toot with id: %s", fbotTootID)
	botClient.DeleteToot(fbotTootID)
}

func recurToot(tootId interface{}) string {
	if tootId != nil {
		ctx := context.Background()
		originToot, err := botClient.Normal.GetStatus(ctx, tootId.(string))
		if err != nil {
			log.SLogger.Errorf("get status error: %s", err)
			return ""
		}
		polished := filter(originToot.Content)
		return fmt.Sprintf("@%s:%s// %s", originToot.Account.Acct, polished, recurToot(originToot.InReplyToID))
	}
	return ""
}

func askForTTs(s string) (string, error) {
	content := make(map[string]string)
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	filename := hex.EncodeToString(bs)
	content["id"] = filename
	content["txt"] = s
	jsonMap, err := json.Marshal(content)
	if err != nil {
		log.SLogger.Errorf("marshal to json error: %v", err)
		return "", err
	}
	err = bredis.Client.Publish("tts", jsonMap).Err()
	if err != nil {
		log.SLogger.Errorf("pub to redis error: %v", err)
		return "", err
	}
	return filename, nil
}

func filter(raw string) (polished string) {
	p := bluemonday.StrictPolicy()
	p.AllowElements("br")
	polished = p.Sanitize(raw)
	polished = strings.Replace(polished, "@firebot", "", -1)
	polished = strings.Replace(polished, "@fbot", "", -1)
	polished = strings.Replace(polished, "<br/>", "\n", -1)
	polished = html.UnescapeString(polished)
	return
}
