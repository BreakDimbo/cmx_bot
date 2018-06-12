package mastodon

import (
	"bot/ai_x/const"
	"bot/ai_x/elastics"
	gomastodon "bot/go-mastodon"
	"context"
	"fmt"
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
	Scope           string    `json:"scope"` //curl -XPUT "http://localhost:9200/status/_mapping/status" -H 'Content-Type: application/json' -d '{"properties": {"scope": {"type": "keyword"}}}'
}

func HandleUpdate(e *gomastodon.UpdateEvent, scope string) {
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

	if scope == con.ScopeTypePublic {
		indexS.Scope = scope
	}

	ctx := context.Background()
	p, err := elastics.Client.Index().
		Index("status").
		Type("status").
		Id(indexS.ID).
		BodyJson(indexS).
		Do(ctx)
	if err != nil {
		fmt.Printf("[ERROR] update to es error: %s/n", err)
	}
	fmt.Printf("Indexed status %s to index %s, type %s, scope %s\n", p.Id, p.Index, p.Type, scope)
}

func HandleDelete(e *gomastodon.DeleteEvent) {
	ctx := context.Background()
	_, err := elastics.Client.Delete().Index("status").Type("status").Id(e.ID).Do(ctx)
	if err != nil {
		fmt.Printf("[ERROR] delete from es error: %s\n", err)
		return
	}
	fmt.Printf("delete from es ok with id: %s\n", e.ID)
}

func HandleNotification(e *gomastodon.NotificationEvent) {
	switch e.Notification.Type {
	case "follow":
		ctx := context.Background()
		accountId := e.Notification.Account.ID
		_, err := client.AccountFollow(ctx, accountId)
		if err != nil {
			fmt.Printf("[Error] follow account error: %s", err)
		}
	}
}

func CleanUnfollower() {
	fmt.Printf("Start cleaning unfollowers\n")
	ctx := context.Background()
	pg := &gomastodon.Pagination{Limit: 80}
	followerM := make(map[gomastodon.ID]bool)
	followingM := make(map[gomastodon.ID]bool)

	ca, err := client.GetAccountCurrentUser(ctx)
	checkErr(err)

	followers, err := client.GetAccountFollowers(ctx, ca.ID, pg)
	checkErr(err)
	for _, v := range followers {
		followerM[v.ID] = true
	}

	followings, err := client.GetAccountFollowing(ctx, ca.ID, pg)
	checkErr(err)
	for _, v := range followings {
		followingM[v.ID] = true
	}

	for k, _ := range followingM {
		if _, ok := followerM[k]; !ok {
			_, err := client.AccountUnfollow(ctx, k)
			checkErr(err)
		}
	}
}

func filter(raw string) (polished string) {
	p := bluemonday.StrictPolicy()
	polished = p.Sanitize(raw)
	return
}

func checkErr(err error) {
	if err != nil {
		fmt.Printf("[ERROR] get error: %s\n", err)
	}
}
