package ig

import (
	"encoding/json"
	"fmt"
	"github.com/tostyle/igfeed/models"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type FeedQuery struct {
	CacheFeedItemIds     []string `json:"cached_feed_item_ids"`
	FetchMediaItemCount  int8     `json:"fetch_media_item_count"`
	FetchMediaItemCursor string   `json:"fetch_media_item_cursor"`
	FetchCommentCount    int8     `json:"fetch_comment_count"`
	FetchLike            int8     `json:"fetch_like"`
	HasStories           bool     `json:"has_stories"`
	HasThreadedComment   bool     `json:"has_threaded_comments"`
}

type IgConfig struct {
	Query     FeedQuery
	QueryHash string
	Cookies   []*http.Cookie
}

const APIURL = "https://www.instagram.com/graphql/query/?"

func (ig IgConfig) BuildQuery() *url.URL {
	queryVar, _ := json.Marshal(ig.Query)
	queryString := string(queryVar)
	baseURL, _ := url.Parse(APIURL)

	fmt.Printf("queryString = %v \n", queryString)
	fmt.Printf("queryHash = %v \n", ig.QueryHash)

	params := url.Values{}
	params.Add("query_hash", ig.QueryHash)
	params.Add("variables", queryString)

	baseURL.RawQuery = params.Encode()
	fmt.Printf("Encoded URL is %q\n", baseURL.String())
	return baseURL
}

func MakeIgCookies(mapCookie map[string]string) []*http.Cookie {
	var cookies []*http.Cookie
	for cookieKey, value := range mapCookie {
		// fmt.Printf("key[%s] value[%s]\n", k, v)
		cookie := &http.Cookie{
			Name:  cookieKey,
			Value: value,
		}
		cookies = append(cookies, cookie)
	}
	return cookies
}
func (ig *IgConfig) ChangePageCursor(cursor string) {
	ig.Query.FetchMediaItemCursor = cursor
}

func (ig IgConfig) FetchFeed() (models.FeedResponse, error) {
	url := ig.BuildQuery()
	jar, _ := cookiejar.New(nil)
	jar.SetCookies(url, ig.Cookies)
	feedResponse := models.FeedResponse{}
	client := &http.Client{
		Jar: jar,
	}
	request, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return feedResponse, err
		// log.Println(err)
		// panic(nil)
	}
	res, err := client.Do(request)
	if err != nil {
		return feedResponse, err
		// log.Println(err)
		// panic(nil)
	}
	defer res.Body.Close()
	responseBodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return feedResponse, err
		// log.Println(err)
		// panic(nil)
	}

	err = json.Unmarshal([]byte(responseBodyBytes), &feedResponse)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	return feedResponse, err
}
