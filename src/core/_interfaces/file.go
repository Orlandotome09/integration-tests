package _interfaces

import (
	"github.com/google/uuid"
	"io"
)

type FileRepository interface {
	Store(reader io.Reader, contentType string) (fileID *uuid.UUID, err error)
	GetUrl(fileID uuid.UUID) (url string, err error)
	Exists(fileID uuid.UUID) (exists bool, err error)
}

type FileService interface {
	Add(reader io.Reader, contentType string) (fileID *uuid.UUID, err error)
	GetUrl(fileID uuid.UUID) (url string, err error)
	Check(fileID uuid.UUID) (exists bool, err error)
}

type FileAdapter interface {
	GetUrl(fileID uuid.UUID) (url string, err error)
}
