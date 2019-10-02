package models

import (
	"time"
)

type FeedResponse struct {
	Data DataResponse `json:"data"`
}
type DataResponse struct {
	User UserResponse `json:"user"`
}
type UserResponse struct {
	ID       string           `json:"id"`
	Profile  string           `json:"profile_pic_url"`
	Username string           `json:"username"`
	TimeLine TimeLineResponse `json:"edge_web_feed_timeline"`
}
type TimeLineResponse struct {
	PageInfo PageInfoResponse `json:"page_info"`
	Edges    []EdgeResponse   `json:"edges"`
}
type PageInfoResponse struct {
	HasNextPage bool   `json:"has_next_page"`
	EndCursor   string `json:"end_cursor"`
}
type EdgeResponse struct {
	Node      EdgeNode `json:"node"`
	Link      string
	CreatedAt time.Time
}
type EdgeNode struct {
	ID         string `json:"id"`
	DisplayURL string `json:"display_url"`
	Shortcode  string `json:"shortcode"`
	PicOwner   Owner  `json:"owner"`
}

type Owner struct {
	Fullname  string `json:"full_name"`
	IsPrivate bool   `json:"is_private"`
}

type FeedData struct {
	ID   interface{}
	Link string
}
