package client

import (
	"bot/config"
	gomastodon "bot/go-mastodon"
	zlog "bot/log"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/sethgrid/pester"
)

type Bot struct {
	Normal *gomastodon.Client
	WS     *gomastodon.WSClient
}

// New new bot client which should be called only when init
func New(config *config.MastodonClientInfo) (*Bot, error) {
	c := gomastodon.NewClient(&gomastodon.Config{
		Server:       config.Sever,
		ClientID:     config.ID,
		ClientSecret: config.Secret,
	})
	err := c.Authenticate(context.Background(), config.Email, config.Password)
	if err != nil {
		log.Fatalf("[Fatal]: authenticate error of mastodon client: %s\n", err)
		return nil, err
	}
	bc := &Bot{Normal: c, WS: c.NewWSClient()}
	return bc, nil
}

func (bc *Bot) RawPost(toot *gomastodon.Toot) (*gomastodon.Status, error) {
	pc := config.GetPostConfig()
	toot.Visibility = pc.Scope
	status, err := bc.Normal.PostStatus(context.Background(), toot)
	if err != nil {
		zlog.SLogger.Errorf("post toot: %s error: %s", toot, err)
		return nil, err
	}
	return status, nil
}

func (bc *Bot) BlockAccount(accountID string) (gomastodon.ID, error) {
	id := gomastodon.ID(accountID)
	status, err := bc.Normal.AccountBlock(context.Background(), id)
	if err != nil {
		zlog.SLogger.Errorf("block id: %s error: %s", id, err)
		return "", err
	}
	return status.ID, nil
}

func (bc *Bot) UnBlockAccount(accountID string) (gomastodon.ID, error) {
	id := gomastodon.ID(accountID)
	status, err := bc.Normal.AccountUnblock(context.Background(), id)
	if err != nil {
		zlog.SLogger.Errorf("unblock id: %s error: %s", id, err)
		return "", err
	}
	return status.ID, nil
}

func (bc *Bot) Post(toot string) (gomastodon.ID, error) {
	pc := config.GetPostConfig()
	status, err := bc.Normal.PostStatus(context.Background(), &gomastodon.Toot{
		Status:     toot,
		Visibility: pc.Scope,
	})
	if err != nil {
		zlog.SLogger.Errorf("post toot: %s error: %s", toot, err)
		return "", err
	}
	return status.ID, nil
}

func (bc *Bot) PostSpoiler(spolier string, toot string) (gomastodon.ID, error) {
	pc := config.GetPostConfig()
	status, err := bc.Normal.PostStatus(context.Background(), &gomastodon.Toot{
		Status:      toot,
		Visibility:  pc.Scope,
		SpoilerText: spolier,
	})
	if err != nil {
		zlog.SLogger.Errorf("post toot: %s error: %s", toot, err)
		return "", err
	}
	return status.ID, nil
}

func (bc *Bot) PostSensetiveWithPic(spolier string, toot string, sensitive bool, medias []gomastodon.Attachment) (*gomastodon.Status, error) {
	pc := config.GetPostConfig()
	var mediasID []gomastodon.ID

	for _, media := range medias {
		fp := downloadPic(media.URL)
		attachment, err := bc.Normal.UploadMedia(context.Background(), fp)
		if err != nil {
			zlog.SLogger.Errorf("upload pic failed: %s", err)
			continue
		}
		mediasID = append(mediasID, gomastodon.ID(attachment.ID))
		removeFile(fp)
	}

	status, err := bc.Normal.PostStatus(context.Background(), &gomastodon.Toot{
		Status:      toot,
		Visibility:  pc.Scope,
		SpoilerText: spolier,
		Sensitive:   sensitive,
		MediaIDs:    mediasID,
	})
	if err != nil {
		zlog.SLogger.Errorf("post toot: %s error: %s", toot, err)
		return nil, err
	}
	return status, nil
}

func downloadPic(remoteURL string) (localURL string) {
	// Create the file
	filenameSlice := strings.Split(remoteURL, "/")
	filename := filenameSlice[len(filenameSlice)-1]
	filepath := fmt.Sprintf("/tmp/cmx_pic/%s", filename)
	out, err := os.Create(filepath)
	if err != nil {
		zlog.SLogger.Error(err)
	}
	defer out.Close()

	// Get the data

	client := pester.New()
	client.Concurrency = 3
	client.MaxRetries = 5
	client.Backoff = pester.ExponentialBackoff
	client.KeepLog = true

	resp, err := client.Get(remoteURL)
	if err != nil {
		zlog.SLogger.Error(err)
	}

	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		zlog.SLogger.Error(err)
	}

	return filepath
}

func removeFile(filepath string) {
	err := os.Remove(filepath)
	if err != nil {
		zlog.SLogger.Errorf("remove file error: %s", err)
		return
	}
}

func (bc *Bot) PostWithPicture(spolier string, toot string) (gomastodon.ID, error) {
	pc := config.GetPostConfig()
	status, err := bc.Normal.PostStatus(context.Background(), &gomastodon.Toot{
		Status:      toot,
		Visibility:  pc.Scope,
		SpoilerText: spolier,
	})
	if err != nil {
		zlog.SLogger.Errorf("post toot: %s error: %s", toot, err)
		return "", err
	}
	return status.ID, nil
}

func (bc *Bot) DeleteToot(id string) error {
	ctx := context.Background()
	fbotTootID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		zlog.SLogger.Errorf("parse id: %s error: %s", id, err)
		return err
	}
	return bc.Normal.DeleteStatus(ctx, int64(fbotTootID))
}
