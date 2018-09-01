package bot

import (
	"bot/bredis"
	con "bot/firebot/const"
	gomastodon "bot/go-mastodon"
	"bot/log"
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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
		} else {
			content := recurToot(replyToID)
			tootToPost = fmt.Sprintf("@%s:%s// %s", fromUser.Acct, firstContent, content)
			tootToPost = strings.TrimSuffix(tootToPost, "// ")
		}

		var id gomastodon.ID

		if isTTs {
			file, err := askForTTs(firstContent)
			if err != nil {
				log.SLogger.Errorf("ask for tts error: %v", err)
				return
			}
			defer file.Close()

			log.SLogger.Debugf("filename is %s", file.Name())
			attachment, err := botClient.Normal.UploadMedia(context.Background(), file.Name())
			if err != nil {
				log.SLogger.Errorf("upload media error: %v", err)
				return
			}

			toot := &gomastodon.Toot{
				MediaIDs: []gomastodon.ID{attachment.ID},
				Status:   "#话唠树洞",
			}

			status, err := botClient.RawPost(toot)
			if err != nil {
				log.SLogger.Errorf("post to mastodon error: %v", err)
				return
			}

			id = status.ID

			err = os.Remove(file.Name())
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

func askForTTs(s string) (*os.File, error) {
	t := strconv.FormatInt(time.Now().Unix(), 10)
	secrete := calSig(t)
	url := fmt.Sprintf("http://47.93.43.59:5438/mp4_tts?secrete=%s&time=%s", secrete, t)
	body := bytes.NewBuffer([]byte(s))

	res, err := http.Post(url, "text", body)
	if err != nil {
		return nil, err
	}

	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	filename := fmt.Sprintf("/tmp/tts/tts%s.mp4", time.Now())
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	_, err = f.Write(result)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func calSig(t string) string {
	const signature = "a7d040cee7bed322a188c9ec7fd9b8b8b34ac893"
	h := sha1.New()
	h.Write([]byte(signature + string(t)))
	bs := h.Sum(nil)
	bsStr := hex.EncodeToString(bs)
	return bsStr
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
