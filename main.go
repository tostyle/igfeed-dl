package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	db "github.com/tostyle/igfeed/database"

	"github.com/tostyle/igfeed/models"
	"github.com/tostyle/igfeed/queue"

	"github.com/tostyle/igfeed/ig"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

// const Database = "tostyle"
// const Server = "mongodb://localhost:30001"
// const QueryHash = "4f0af5ae46172b2f891d5c1b3c6fd6a2"

func handleIgError() {
	if r := recover(); r != nil {
		fmt.Println("recovered from ", r)
	}
}
func app() {
	// workEnv := os.Getenv("WORK_ENV")
	// var fileName string
	// if workEnv == "production" {
	// 	fileName = "./.env.docker"
	// }
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
		log.Fatal("Error loading .env file")
		panic(err)
	}
	Database := os.Getenv("DATABASE_NAME")
	DatabaseUsername := os.Getenv("DATABASE_USERNAME")
	DatabasePassword := os.Getenv("DATABASE_PASSWORD")
	Server := os.Getenv("DATABASE_SERVER")
	QueueConnectString := os.Getenv("QUEUE_CONNECTION")
	QueryHash := os.Getenv("IG_QUERY_HASH")
	Mid := os.Getenv("IG_COOKIE_MID")
	Fid := os.Getenv("IG_COOKIE_FID")
	Domain := os.Getenv("IG_COOKIE_DOMAIN")
	SessionID := os.Getenv("IG_COOKIE_SESSIONID")
	CsrfToken := os.Getenv("IG_CSRF_TOKEN")

	log.Println("Start", Server)
	q := queue.FeedQueue{}
	err = q.ConnectAndPrepareQueue(QueueConnectString)
	if err != nil {
		if q.Connection != nil {
			q.Connection.Close()
		}
		if q.Channel != nil {
			q.Channel.Close()
		}
		panic(err)
	}
	feedDb := db.FeedDB{Database: Database, Server: Server, Username: DatabaseUsername, Password: DatabasePassword}
	feedDb.Connect()

	// defer feedDb.Disconnect()
	feedQuery := ig.FeedQuery{
		CacheFeedItemIds:     make([]string, 0),
		FetchMediaItemCount:  12,
		FetchMediaItemCursor: "",
		FetchCommentCount:    4,
		FetchLike:            3,
		HasStories:           false,
		HasThreadedComment:   true}

	cookieObj := make(map[string]string)
	cookieObj["mid"] = Mid
	cookieObj[Fid] = Domain
	cookieObj["csrftoken"] = CsrfToken
	cookieObj["sessionid"] = SessionID

	cookies := ig.MakeIgCookies(cookieObj)

	igConfig := ig.IgConfig{Query: feedQuery, QueryHash: QueryHash, Cookies: cookies}

	totalPage := 3

	var edges []models.EdgeResponse

	defer handleIgError()
	for index := 0; index < totalPage; index++ {
		feedResponse, err := igConfig.FetchFeed()
		if err != nil {
			panic(err)
		}
		timeline := &feedResponse.Data.User.TimeLine
		// fmt.Printf("%+v\n", res)
		edges = append(edges, timeline.Edges...)
		igConfig.ChangePageCursor(timeline.PageInfo.EndCursor)
	}
	fmt.Printf("total feed = %v \n", len(edges))

	var wg sync.WaitGroup
	for _, feed := range edges {
		wg.Add(1)
		go func(currentFeed models.EdgeResponse) {
			defer wg.Done()
			newFeed, err := feedDb.UpFeed(&currentFeed)
			if err != nil {
				log.Println(err)
				panic(nil)
			}
			// log.Printf("%T\n", newFeed.UpsertedId)
			if newFeed.UpsertedId != nil {
				feedData := models.FeedData{
					ID:   newFeed.UpsertedId,
					Link: currentFeed.Node.DisplayURL,
				}
				feedStr, err := json.Marshal(feedData)
				if err != nil {
					log.Println(err)
				}
				log.Printf("%v\n", string(feedStr))
				err = q.Channel.Publish(
					"igfeed",     // exchange
					"igfeed_pic", // routing key
					false,        // mandatory
					false,        // immediate
					amqp.Publishing{
						DeliveryMode: amqp.Transient,
						ContentType:  "application/json",
						Body:         feedStr,
					})
				if err != nil {
					log.Fatalln(err)
				}
			}

		}(feed)
	}
	wg.Wait()
	fmt.Println("finish job")
}

func main() {
	app()
	// c := cron.New()
	// fmt.Println("startt")
	// c.AddFunc("*/1 * * * * ", func() {
	// 	fmt.Println("run job")
	// 	app()
	// })
	// c.Run()
	// defer c.Stop()
}
