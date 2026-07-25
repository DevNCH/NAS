package models

import "time"

type File struct {
	ID         int64
	Filename   string
	Filepath   string
	Size       int64
	UploadedBy int64
	CreatedAt  time.Time
}