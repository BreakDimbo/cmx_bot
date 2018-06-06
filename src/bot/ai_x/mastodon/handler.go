package mastodon

import (
	"bot/ai_x/elastics"
	gomastodon "bot/go-mastodon"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/microcosm-cc/bluemonday"
)

type indexStatus struct {
	ID              string    `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	AccountId       string    `json:"account_id"`
	Content         string    `json:"content"`
	ReblogsCount    int64     `json:"reblogs_count"`
	FavouritesCount int64     `json:"favourites_count"`
	Sensitive       bool      `json:"sensitive"`
}

func HandleUpdate(e *gomastodon.UpdateEvent) {
	polished := filter(e.Status.Content)
	indexS := &indexStatus{
		ID:              string(e.Status.ID),
		CreatedAt:       e.Status.CreatedAt,
		AccountId:       string(e.Status.Account.ID),
		Content:         polished,
		ReblogsCount:    e.Status.ReblogsCount,
		FavouritesCount: e.Status.FavouritesCount,
		Sensitive:       e.Status.Sensitive,
	}
	ctx := context.Background()
	p, err := elastics.Client.Index().
		Index("status").
		Type("status").
		Id(indexS.ID).
		BodyJson(indexS).
		Do(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Indexed status %s to index %s, type %s\n", p.Id, p.Index, p.Type)
}

func HandleDelete(e *gomastodon.DeleteEvent) {
	ctx := context.Background()
	_, err := elastics.Client.Delete().Id(e.ID).Do(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("delete from es ok with id: %s\n", e.ID)
}

func filter(raw string) (polished string) {
	p := bluemonday.StrictPolicy()
	polished = p.Sanitize(raw)
	return
}
