package db

import (
	"log"
	"math/rand"
	"time"

	. "github.com/tostyle/igfeed/models"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var db *mgo.Database

const FeedCollection = "feeds"

type FeedDB struct {
	Database string
	Server   string
	Username string
	Password string
}

func (feed FeedDB) Connect() {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{feed.Server},
		Timeout:  30 * time.Second,
		Database: feed.Database,
		Source:   "admin",
		Username: feed.Username,
		Password: feed.Password,
	}
	log.Printf("mongoInfo %#v", mongoDBDialInfo)
	// session, err := mgo.Dial(feed.Server)
	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Panic(err)
	}
	db = session.DB(feed.Database)
}
func (feed FeedDB) AddFeed(feedEdge EdgeResponse) error {

	err := db.C(FeedCollection).Insert(&feedEdge)
	return err
}

func (feed FeedDB) RandomFeed() (error, EdgeResponse) {
	feedEdge := EdgeResponse{}
	skip := rand.Intn(36)
	log.Print(skip)
	err := db.C(FeedCollection).Find(bson.M{"node.picowner.isprivate": false}).Skip(skip).One(&feedEdge)
	if err != nil {
		return err, feedEdge
	}
	return nil, feedEdge
}

func (feed *FeedDB) UpFeed(feedEdge *EdgeResponse) (info *mgo.ChangeInfo, err error) {
	find := bson.M{"node.id": feedEdge.Node.ID}
	feedEdge.CreatedAt = time.Now()
	feedEdge.Link = "http://instagram.com/p/" + feedEdge.Node.Shortcode
	newFeed, err := db.C(FeedCollection).Upsert(find, feedEdge)
	return newFeed, err
}

// func(feed FeedDB) Disconnect () {
// 	if (db != nil) {
// 		db.Close()
// 	}
// }
