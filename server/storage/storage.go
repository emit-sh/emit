package storage

import "io"

type Storage interface {
	Get(path string) (reader io.ReadCloser, err error)
	GetDownloadUrl(path string) (url string, err error)
	Put(path string, reader io.Reader, ttl int) error
	Exists(path string) (exists bool, err error)
	Type() StorageType
}

type StorageType int;
const (
	AWS StorageType = 1 //not created yet
	BackBlaze StorageType = 2 // your only option!
	DigitalOcean StorageType = 3 // not created yet
	LocalFileSystem StorageType = 4
)