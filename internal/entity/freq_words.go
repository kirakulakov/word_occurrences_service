package entity

type Comm struct {
	PostId int    `json:"post_id"`
	Word   string `json:"word"`
	Count  int    `json:"int"`
}

type Comment struct {
	PostID int    `json:"postId"`
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}

type Count struct {
	Count int
}
