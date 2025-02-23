package dto

import "time"

type ImageInfo struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	IsUsed    bool      `json:"isUsed"`
	Tags      []string  `json:"tags"`
	Size      string    `json:"size"`
}
