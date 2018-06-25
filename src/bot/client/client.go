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

type BotClient struct {
	Normal *gomastodon.Client
	WS     *gomastodon.WSClient
}

// New new bot client which should be called only when init
func New(config *config.MastodonClientInfo) (*BotClient, error) {
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
	bc := &BotClient{Normal: c, WS: c.NewWSClient()}
	return bc, nil
}

func (bc *BotClient) Post(toot string) (gomastodon.ID, error) {
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

func (bc *BotClient) PostSpoiler(spolier string, toot string) (gomastodon.ID, error) {
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

func (bc *BotClient) PostSensetiveWithPic(spolier string, toot string, sensitive bool, medias []gomastodon.Attachment) (gomastodon.ID, error) {
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
		return "", err
	}
	return status.ID, nil
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

func (bc *BotClient) PostWithPicture(spolier string, toot string) (gomastodon.ID, error) {
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

func (bc *BotClient) DeleteToot(id string) error {
	ctx := context.Background()
	fbotTootID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		zlog.SLogger.Errorf("parse id: %s error: %s", id, err)
		return err
	}
	return bc.Normal.DeleteStatus(ctx, int64(fbotTootID))
}
