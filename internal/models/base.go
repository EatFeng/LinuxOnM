package models

import "time"

type BaseModel struct {
	ID        uint      `grom:"primarykey;AUTO_INCREMENT" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
