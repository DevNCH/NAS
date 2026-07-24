package models

import "time"

type File struct {
	ID         int       `json:"id"`
	Filename   string    `json:"filename"`
	Filepath   string    `json:"filepath"`
	Size       int64     `json:"size"`
	UploadedBy int       `json:"uploaded_by"`
	CreatedAt  time.Time `json:"created_at"`
}
