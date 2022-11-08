package models

type News struct {
	ID          string `json:"id" db:"id" bson:"_id"`
	Title       string `json:"title" db:"title" bson:"title"`
	Content     string `json:"content" db:"content" bson:"content"`
	AuthorID    string `json:"author_id" db:"author_id" bson:"author_id"`
	AuthorName  string `json:"name" bson:"-"`
	CreatedAt   int64  `json:"created_at" db:"created_at" bson:"created_at"`
	PublishedAt int64  `json:"published_at" db:"published_at" bson:"-"`
}

// News - публикация.
type NewsItem struct {
	ID          int    `json:"id"`
	CommentID   int    `json:"comment_id"`
	Title       string `json:"title"`
	Link        string `json:"link"`
	Description string `json:"description"`
	Timestamp   int    `json:"timestamp"`
}
