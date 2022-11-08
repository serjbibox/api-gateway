package models

type NewsComment struct {
	ID        int    `json:"id"`
	Child     int    `json:"child"`
	Parent    int    `json:"parent"`
	Content   string `json:"content"`
	Timestamp int    `json:"timestamp"`
}
